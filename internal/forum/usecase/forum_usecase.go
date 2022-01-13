package usecase

import (
	repository "TP_DB/internal/interfaces"
	"TP_DB/internal/models"
)

type ForumUseCaseStruct struct {
	forumRepository  repository.ForumRepositoryInterface
	userRepository   repository.UserRepositoryInterface
	threadRepository repository.ThreadRepositoryInterface
}

func CreateForumUseCase(forumRepository repository.ForumRepositoryInterface, userRepository repository.UserRepositoryInterface, threadRepository repository.ThreadRepositoryInterface) *ForumUseCaseStruct {
	return &ForumUseCaseStruct{forumRepository: forumRepository, userRepository: userRepository, threadRepository: threadRepository}
}

func (forumUseCase *ForumUseCaseStruct) ForumCreate(forum *models.Forum) (*models.Forum, int) {
	var user models.User
	user.Nickname = forum.User

	user, err := forumUseCase.userRepository.UserGet(user.Nickname)
	if err != nil {
		return nil, 404
	}
	forum.User = user.Nickname

	var createdForum models.Forum
	createdForum, err = forumUseCase.forumRepository.GetForumBySlug(forum.Slug)
	if err == nil {
		return &createdForum, 409
	}

	createdForum, err = forumUseCase.forumRepository.CreateForum(*forum)
	if err == nil {
		return &createdForum, 201
	}

	return nil, 0
}

func (forumUseCase *ForumUseCaseStruct) ForumCreateThread(thread *models.Thread, forumSlug string) (*models.Thread, int) {
	forum, code := forumUseCase.ForumGetBySlug(forumSlug)
	if code == 404 {
		return nil, 404
	}
	thread.Forum = forum.Slug

	_, err := forumUseCase.userRepository.UserGet(thread.Author)
	if err != nil {
		return nil, 404
	}

	threadDB, err := forumUseCase.threadRepository.GetThreadBySlug(thread.Slug)
	if err == nil && threadDB.Slug != "" {
		return &threadDB, 409
	}

	threadCreated, err := forumUseCase.forumRepository.CreateForumThread(*thread)

	return &threadCreated, 201
}

func (forumUseCase *ForumUseCaseStruct) ForumGetBySlug(slug string) (*models.Forum, int) {
	forum, err := forumUseCase.forumRepository.GetForumBySlug(slug)
	if err != nil {
		return nil, 404
	}

	return &forum, 200
}

func (forumUseCase *ForumUseCaseStruct) ForumGetThreads(slug string, limit int, desc bool, since string) (*models.Threads, int) {
	_, err := forumUseCase.forumRepository.GetForumBySlug(slug)
	if err != nil {
		return nil, 404
	}

	threads, err := forumUseCase.forumRepository.GetForumThreads(slug, limit, desc, since)
	if err == nil {
		return &threads, 200
	}

	return nil, 0
}

func (forumUseCase *ForumUseCaseStruct) ForumGetUsers(slug string, limit int, desc bool, since string) (*models.Users, int) {
	forum, err := forumUseCase.forumRepository.GetForumBySlug(slug)
	if err != nil {
		return nil, 404
	}

	users, err := forumUseCase.forumRepository.GetForumUsers(forum.Slug, limit, desc, since)
	if err == nil {
		return &users, 200
	}

	return nil, 0
}
