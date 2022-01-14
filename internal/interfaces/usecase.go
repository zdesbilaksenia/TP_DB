package interfaces

import "TP_DB/internal/models"

type ForumUseCaseInterface interface {
	ForumCreate(forum *models.Forum) (*models.Forum, int)
	ForumGetBySlug(slug string) (*models.Forum, int)
	ForumGetThreads(slug string, limit int, desc bool, since string) (*models.Threads, int)
	ForumGetUsers(slug string, limit int, desc bool, since string) (*models.Users, int)
	ForumCreateThread(thread *models.Thread, forumSlug string) (*models.Thread, int)
}

type PostUseCaseInterface interface {
	PostGetInfo(related []string, id int) (*models.PostFull, int)
	PostUpdateInfo(id int, postUpdate models.Post) (*models.Post, int)
}

type ServiceUseCaseInterface interface {
	ServiceClear() int
	ServiceGetStatus() (*models.Status, int)
}

type ThreadUseCaseInterface interface {
	ThreadCreatePosts(slug string, id int, posts models.Posts) (models.Posts, int)
	ThreadUpdate(threadUpd models.ThreadUpdate, slug string, id int) (*models.Thread, int)
	ThreadGet(slug string, id int) (*models.Thread, int)
	ThreadVote(vote models.Vote, slug string, id int) (*models.Thread, int)
	ThreadGetPosts(slug string, id int, limit int, since string, desc bool, sort string) (*models.Posts, int)
}

type UserUseCaseInterface interface {
	UserCreate(user *models.User) (models.User, models.Users, int)
	UserGet(nickname string) (models.User, error)
	UserChange(user models.User) (models.User, models.Users, int)
	UsersGet(user models.User) (models.Users, error)
}
