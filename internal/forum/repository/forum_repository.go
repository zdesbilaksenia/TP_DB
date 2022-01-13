package repository

import (
	"TP_DB/internal/models"
	"github.com/jackc/pgx"
)

const (
	CreateThread            = `INSERT INTO thread (title, author, message, created, forum, slug) VALUES ($1, $2, $3, $4, $5, COALESCE(NULLIF($6, ''), $6)) RETURNING id, title, author, forum, message, votes, slug, created`
	CreateForum             = `INSERT INTO forum(title, "user", slug) VALUES ($1,  $2, $3) RETURNING slug, title, "user", posts, threads`
	GetForumBySlug          = `SELECT title, "user", slug, posts, threads FROM forum WHERE slug=$1`
	GetThreads              = `SELECT * FROM thread WHERE forum = $1 AND created >= $2 ORDER BY created LIMIT $3`
	GetThreadsDesc          = `SELECT * FROM thread WHERE forum = $1 AND created <= $2 ORDER BY created DESC LIMIT $3`
	GetThreadsNoSince       = `SELECT * FROM thread WHERE forum = $1 ORDER BY created DESC LIMIT $2`
	GetThreadsNoSinceNoDesc = `SELECT * FROM thread WHERE forum = $1 ORDER BY created LIMIT $2`
	GetUsers                = `SELECT nickname, fullname, about, email FROM nickname_forum WHERE forum=$1`
)

type ForumRepositoryStruct struct {
	DB *pgx.ConnPool
}

func CreateForumRepository(DB *pgx.ConnPool) *ForumRepositoryStruct {
	return &ForumRepositoryStruct{DB: DB}
}

func (forumRepository *ForumRepositoryStruct) CreateForum(forum models.Forum) (models.Forum, error) {
	err := forumRepository.DB.QueryRow(CreateForum, forum.Title, forum.User, forum.Slug).
		Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads)

	return forum, err
}

func (forumRepository *ForumRepositoryStruct) GetForumBySlug(slug string) (models.Forum, error) {
	var forum models.Forum

	err := forumRepository.DB.QueryRow(GetForumBySlug, slug).
		Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)

	return forum, err
}

func (forumRepository *ForumRepositoryStruct) CreateForumThread(thread models.Thread) (models.Thread, error) {
	err := forumRepository.DB.QueryRow(CreateThread, thread.Title, thread.Author, thread.Message, thread.Created, thread.Forum, thread.Slug).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)

	return thread, err
}

func (forumRepository *ForumRepositoryStruct) GetForumThreads(slug string, limit int, desc bool, since string) (models.Threads, error) {
	var rows *pgx.Rows
	var err error
	switch desc {
	case true:
		switch since {
		case "":
			rows, err = forumRepository.DB.Query(GetThreadsNoSince, slug, limit)
		default:
			rows, err = forumRepository.DB.Query(GetThreadsDesc, slug, since, limit)
		}
	case false:
		switch since {
		case "":
			rows, err = forumRepository.DB.Query(GetThreadsNoSinceNoDesc, slug, limit)
		default:
			rows, err = forumRepository.DB.Query(GetThreads, slug, since, limit)
		}
	}

	if err != nil {
		return nil, err
	}

	var threads models.Threads
	for rows.Next() {
		var thread models.Thread
		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author,
			&thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
		if err != nil {
			return nil, err
		}
		threads = append(threads, &thread)
	}

	rows.Close()
	return threads, nil
}

func (forumRepository *ForumRepositoryStruct) GetForumUsers(slug string, limit int, desc bool, since string) (models.Users, error) {
	var rows *pgx.Rows
	var err error
	var query string

	if since == "" {
		if desc {
			query = GetUsers + ` ORDER BY nickname DESC LIMIT NULLIF($2, 0)`
		} else {
			query = GetUsers + ` ORDER BY nickname ASC LIMIT NULLIF($2, 0)`
		}

		rows, err = forumRepository.DB.Query(query, slug, limit)
	} else {
		if desc {
			query = GetUsers + ` AND nickname < $2 ORDER BY nickname DESC LIMIT NULLIF($3, 0)`
		} else {
			query = GetUsers + ` AND nickname > $2 ORDER BY nickname ASC LIMIT NULLIF($3, 0)`
		}

		rows, err = forumRepository.DB.Query(query, slug, since, limit)
	}

	if err != nil {
		return nil, err
	}

	var users models.Users
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	rows.Close()
	return users, nil
}
