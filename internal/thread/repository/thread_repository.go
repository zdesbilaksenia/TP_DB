package repository

import (
	"TP_DB/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	"strings"
	"time"
)

const (
	GetThreadBySlug = `SELECT id, title, author, forum, message, votes, slug, created FROM thread WHERE slug = $1`
	GetThreadById   = `SELECT id, title, author, forum, message, votes, slug, created FROM thread WHERE id = $1`
	CreatePost      = `INSERT INTO post (parent, author, message, thread, forum, created) VALUES ($1, $2, $3, $4, $5, $6)
					RETURNING id, parent, author, message, isedited, forum, thread, created`
	UpdateThreadBySlug = `UPDATE thread SET title = COALESCE(NULLIF($1, ''), title), message = COALESCE(NULLIF($2, ''), message) WHERE slug = $3 RETURNING *`
	UpdateThreadById   = `UPDATE thread SET title = COALESCE(NULLIF($1, ''), title), message = COALESCE(NULLIF($2, ''), message) WHERE id = $3 RETURNING *`
	CreateVote         = `INSERT INTO votes (thread, nickname, voice) VALUES ($1, $2, $3)`
	UpdateVote         = `UPDATE votes SET voice = $1 WHERE thread = $2 AND nickname = $3`
	GetPosts           = `SELECT id, parent, author, message, isedited, forum, thread, created FROM post WHERE thread = $1`
	GetPostsTree       = `SELECT id, parent, author, message, isedited, forum, thread, created FROM post WHERE path[1] in`
)

type ThreadRepositoryStruct struct {
	DB *pgx.ConnPool
}

func CreateThreadRepository(DB *pgx.ConnPool) *ThreadRepositoryStruct {
	return &ThreadRepositoryStruct{DB: DB}
}

func (threadRepository *ThreadRepositoryStruct) GetThreadBySlug(slug string) (models.Thread, error) {
	var thread models.Thread
	var created sql.NullTime
	var slugThr sql.NullString

	err := threadRepository.DB.QueryRow(GetThreadBySlug, slug).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &slugThr, &created)

	if err == nil {
		if created.Valid {
			thread.Created = created.Time
		}

		if slugThr.Valid {
			thread.Slug = slugThr.String
		} else {
			thread.Slug = ""
		}
	}

	return thread, err
}

func (threadRepository *ThreadRepositoryStruct) GetThreadById(id int) (models.Thread, error) {
	var thread models.Thread
	var created sql.NullTime
	var slugThr sql.NullString

	err := threadRepository.DB.QueryRow(GetThreadById, id).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &slugThr, &created)

	if err == nil {
		if created.Valid {
			thread.Created = created.Time
		}

		if slugThr.Valid {
			thread.Slug = slugThr.String
		} else {
			thread.Slug = ""
		}
	}

	return thread, err
}

func (threadRepository *ThreadRepositoryStruct) UpdateThreadBySlug(threadUpd models.ThreadUpdate, slug string) (models.Thread, error) {
	var thread models.Thread

	err := threadRepository.DB.QueryRow(UpdateThreadBySlug, threadUpd.Title, threadUpd.Message, slug).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)

	return thread, err
}

func (threadRepository *ThreadRepositoryStruct) UpdateThreadById(threadUpd models.ThreadUpdate, id int) (models.Thread, error) {
	var thread models.Thread

	err := threadRepository.DB.QueryRow(UpdateThreadById, threadUpd.Title, threadUpd.Message, id).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)

	return thread, err
}

func (threadRepository *ThreadRepositoryStruct) CreateThreadPosts(posts models.Posts, threadId int, forum string) (models.Posts, int) {
	var insertedPosts models.Posts
	var sqlValues []interface{}

	sqlQuery := `INSERT INTO post (parent, author, message, forum, thread, created) VALUES `
	if len(posts) == 0 {
		return models.Posts{}, 0
	}
	created := time.Now()
	for i, post := range posts {
		author := ""
		err := threadRepository.DB.QueryRow("select nickname from users where nickname = $1", post.Author).Scan(&author)
		if err == pgx.ErrNoRows {
			return nil, 404
		}
		if post.Parent != 0 {
			id := -1
			err := threadRepository.DB.QueryRow("select id from post where thread = $1 and id = $2", threadId, post.Parent).Scan(&id)
			if err == pgx.ErrNoRows {
				return nil, 409
			}
		}
		sqlValuesString := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6)
		sqlQuery += sqlValuesString
		sqlValues = append(sqlValues, post.Parent, post.Author, post.Message, forum, threadId, created)
	}
	sqlQuery = strings.TrimSuffix(sqlQuery, ",")
	sqlQuery += ` RETURNING id, parent, author, message, isedited, forum, thread, created;`
	rows, err := threadRepository.DB.Query(sqlQuery, sqlValues...)
	if err != nil {
		return nil, 500
	}
	defer rows.Close()
	for rows.Next() {
		post := models.Post{}
		err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil || post.Author == "" {
			return nil, 500
		}
		insertedPosts = append(insertedPosts, &post)
	}
	if len(insertedPosts) == 0 {
		return nil, 0
	}
	return insertedPosts, 0
}

func (threadRepository *ThreadRepositoryStruct) getThreadPostsFlat(id int, limit int, since string, desc bool) (models.Posts, error) {
	var posts models.Posts
	var rows *pgx.Rows
	var err error

	query := GetPosts

	if desc {
		if since != "" {
			query += ` and id < $2 order by id desc limit $3`
			rows, err = threadRepository.DB.Query(query, id, since, limit)
		} else {
			query += ` order by id desc limit $2`
			rows, err = threadRepository.DB.Query(query, id, limit)
		}
	} else {
		if since != "" {
			query += ` and id > $2 order by id limit $3`
			rows, err = threadRepository.DB.Query(query, id, since, limit)
		} else {
			query += ` order by id limit $2`
			rows, err = threadRepository.DB.Query(query, id, limit)
		}
	}

	if err != nil {
		rows.Close()
		return nil, err
	}

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			rows.Close()
			return nil, err
		}
		posts = append(posts, &post)
	}

	rows.Close()
	return posts, err
}

func (threadRepository *ThreadRepositoryStruct) getThreadPostsTree(id int, limit int, since string, desc bool) (models.Posts, error) {
	var posts models.Posts
	var rows *pgx.Rows
	var err error

	query := GetPosts

	if desc {
		if since != "" {
			query += ` and path < (select path FROM post where id = $2) order by path desc, id desc limit $3`
			rows, err = threadRepository.DB.Query(query, id, since, limit)
		} else {
			query += ` order by path desc, id desc limit $2`
			rows, err = threadRepository.DB.Query(query, id, limit)
		}
	} else {
		if since != "" {
			query += ` and path > (select path FROM post where id = $2) order by path, id limit $3`
			rows, err = threadRepository.DB.Query(query, id, since, limit)
		} else {
			query += ` order by path, id limit $2`
			rows, err = threadRepository.DB.Query(query, id, limit)
		}
	}

	if err != nil {
		rows.Close()
		return nil, err
	}

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			rows.Close()
			return nil, err
		}
		posts = append(posts, &post)
	}

	rows.Close()
	return posts, nil
}

func (threadRepository *ThreadRepositoryStruct) getThreadPostsParentTree(id int, limit int, since string, desc bool) (models.Posts, error) {
	var posts models.Posts
	var rows *pgx.Rows
	var err error

	query := GetPostsTree

	if desc {
		if since != "" {
			query += ` (select id from post where thread = $1 and parent = 0 and path[1] < 
						(select path[1] from post where id = $2) order by id desc limit $3)
						order by path[1] desc, path, id`
			rows, err = threadRepository.DB.Query(query, id, since, limit)
		} else {
			query += ` (select id from post where thread = $1 and parent = 0 order by id desc limit $2)
						order by path[1] desc, path, id`
			rows, err = threadRepository.DB.Query(query, id, limit)
		}
	} else {
		if since != "" {
			query += ` (select id from post where thread = $1 and parent = 0 and path[1] > 
						(select path[1] from post where id = $2) order by id limit $3)
						order by path, id`
			rows, err = threadRepository.DB.Query(query, id, since, limit)
		} else {
			query += ` (select id from post where thread = $1 and parent = 0 order by id limit $2)
						order by path, id`
			rows, err = threadRepository.DB.Query(query, id, limit)
		}
	}

	if err != nil {
		rows.Close()
		return nil, err
	}

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)
		if err != nil {
			rows.Close()
			return nil, err
		}
		posts = append(posts, &post)
	}

	rows.Close()
	return posts, nil
}

func (threadRepository *ThreadRepositoryStruct) GetThreadPosts(id int, limit int, since string, sort string, desc bool) (models.Posts, error) {
	if sort == "" || sort == "flat" {
		return threadRepository.getThreadPostsFlat(id, limit, since, desc)
	} else if sort == "tree" {
		return threadRepository.getThreadPostsTree(id, limit, since, desc)
	} else if sort == "parent_tree" {
		return threadRepository.getThreadPostsParentTree(id, limit, since, desc)
	}

	return nil, errors.New("err")
}

func (threadRepository *ThreadRepositoryStruct) VoteThread(vote models.Vote, id int) error {
	_, err := threadRepository.DB.Exec(CreateVote, id, vote.Nickname, vote.Voice)

	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "duplicate") {
		_, err = threadRepository.DB.Exec(UpdateVote, vote.Voice, id, vote.Nickname)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
