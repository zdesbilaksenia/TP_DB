package interfaces

import "TP_DB/internal/models"

type ForumRepositoryInterface interface {
	GetForumBySlug(slug string) (models.Forum, error)
	CreateForum(forum models.Forum) (models.Forum, error)
	GetForumThreads(slug string, limit int, desc bool, since string) (models.Threads, error)
	GetForumUsers(slug string, limit int, desc bool, since string) (models.Users, error)
	CreateForumThread(thread models.Thread) (models.Thread, error)
}

type PostRepositoryInterface interface {
	GetPostInfo(id int) (models.Post, error)
	UpdatePostInfo(post models.Post) (models.Post, error)
}

type ServiceRepositoryInterface interface {
	GetStatus() (models.Status, error)
	Clear() error
}

type ThreadRepositoryInterface interface {
	GetThreadBySlug(slug string) (models.Thread, error)
	GetThreadById(id int) (models.Thread, error)
	CreateThreadPosts(posts models.Posts, threadId int, forum string) (models.Posts, int)
	UpdateThreadById(threadUpd models.ThreadUpdate, id int) (models.Thread, error)
	UpdateThreadBySlug(threadUpd models.ThreadUpdate, slug string) (models.Thread, error)
	VoteThread(vote models.Vote, id int) error
	GetThreadPosts(id int, limit int, since string, sort string, desc bool) (models.Posts, error)
}

type UserRepositoryInterface interface {
	UserCreate(user *models.User) error
	UsersGet(user *models.User) (models.Users, error)
	UserGet(nickname string) (models.User, error)
	UserChange(user models.User) (models.User, error)
}
