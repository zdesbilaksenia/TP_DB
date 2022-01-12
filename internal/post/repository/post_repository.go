package repository

import (
	"TP_DB/internal/models"
	"github.com/jackc/pgx"
)

const (
	GetPostInfo    = `SELECT id, parent, author, message, isedited, forum, thread, created FROM post WHERE id = $1`
	UpdatePostInfo = `UPDATE post SET message  = COALESCE(NULLIF($1, ''), message), isedited = CASE WHEN $1 = '' OR message = $1 THEN isedited else true end
					WHERE id = $2 RETURNING id, parent, author, message, isedited, forum, thread, created`
)

type PostRepositoryStruct struct {
	DB *pgx.ConnPool
}

func CreatePostRepository(DB *pgx.ConnPool) *PostRepositoryStruct {
	return &PostRepositoryStruct{DB: DB}
}

func (postRepository *PostRepositoryStruct) GetPostInfo(id int) (models.Post, error) {
	var post models.Post
	err := postRepository.DB.QueryRow(GetPostInfo, id).
		Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created)

	return post, err
}

func (postRepository *PostRepositoryStruct) UpdatePostInfo(post models.Post) (models.Post, error) {
	var postEdited models.Post
	err := postRepository.DB.QueryRow(UpdatePostInfo, post.Message, post.Id).
		Scan(&postEdited.Id, &postEdited.Parent, &postEdited.Author, &postEdited.Message, &postEdited.IsEdited, &postEdited.Forum, &postEdited.Thread, &postEdited.Created)

	return postEdited, err
}
