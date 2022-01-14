package usecase

import (
	repository "TP_DB/internal/interfaces"
	"TP_DB/internal/models"
)

type ThreadUseCaseStruct struct {
	threadRepository repository.ThreadRepositoryInterface
	userRepository   repository.UserRepositoryInterface
	postRepository   repository.PostRepositoryInterface
}

func CreateThreadUseCase(threadRepository repository.ThreadRepositoryInterface, userRepository repository.UserRepositoryInterface, postRepository repository.PostRepositoryInterface) *ThreadUseCaseStruct {
	return &ThreadUseCaseStruct{threadRepository: threadRepository, userRepository: userRepository, postRepository: postRepository}
}

func (threadUseCase *ThreadUseCaseStruct) ThreadGet(slug string, id int) (*models.Thread, int) {
	var err error
	var thread models.Thread
	if slug != "" {
		thread, err = threadUseCase.threadRepository.GetThreadBySlug(slug)
		if err != nil {
			return nil, 404
		}
	} else if id != -1 {
		thread, err = threadUseCase.threadRepository.GetThreadById(id)
		if err != nil {
			return nil, 404
		}
	}

	if err != nil {
		return nil, 0
	}

	return &thread, 200
}

func (threadUseCase *ThreadUseCaseStruct) ThreadUpdate(threadUpd models.ThreadUpdate, slug string, id int) (*models.Thread, int) {
	var err error
	var thread models.Thread
	if slug != "" {
		_, err = threadUseCase.threadRepository.GetThreadBySlug(slug)
		if err != nil {
			return nil, 404
		}
		thread, err = threadUseCase.threadRepository.UpdateThreadBySlug(threadUpd, slug)
	} else if id != -1 {
		_, err = threadUseCase.threadRepository.GetThreadById(id)
		if err != nil {
			return nil, 404
		}
		thread, err = threadUseCase.threadRepository.UpdateThreadById(threadUpd, id)
	}

	if err != nil {
		return nil, 0
	}

	return &thread, 200
}

func (threadUseCase *ThreadUseCaseStruct) ThreadCreatePosts(slug string, id int, posts models.Posts) (models.Posts, int) {
	var thread models.Thread
	var err error
	if slug != "" {
		thread, err = threadUseCase.threadRepository.GetThreadBySlug(slug)
	} else if id != -1 {
		thread, err = threadUseCase.threadRepository.GetThreadById(id)
	}
	if err != nil {
		return nil, 404
	}
	posts, errInt := threadUseCase.threadRepository.CreateThreadPosts(posts, thread.Id, thread.Forum)
	if errInt != 0 {
		return nil, errInt
	}
	return posts, 201
}

func (threadUseCase *ThreadUseCaseStruct) ThreadGetPosts(slug string, id int, limit int, since string, desc bool, sort string) (*models.Posts, int) {
	var thread models.Thread
	var err error
	if slug != "" {
		thread, err = threadUseCase.threadRepository.GetThreadBySlug(slug)
	} else if id != -1 {
		thread, err = threadUseCase.threadRepository.GetThreadById(id)
	}

	if err != nil {
		return nil, 404
	}

	posts, err := threadUseCase.threadRepository.GetThreadPosts(thread.Id, limit, since, sort, desc)
	if err == nil {
		return &posts, 200
	}

	return nil, 0
}

func (threadUseCase *ThreadUseCaseStruct) ThreadVote(vote models.Vote, slug string, id int) (*models.Thread, int) {
	var err error
	var thread models.Thread
	if slug != "" {
		thread, err = threadUseCase.threadRepository.GetThreadBySlug(slug)
		if err != nil {
			return nil, 404
		}
	} else if id != -1 {
		thread, err = threadUseCase.threadRepository.GetThreadById(id)
		if err != nil {
			return nil, 404
		}
	}

	user, err := threadUseCase.userRepository.UserGet(vote.Nickname)
	if err != nil {
		return nil, 404
	}
	vote.Nickname = user.Nickname

	err = threadUseCase.threadRepository.VoteThread(vote, thread.Id)
	if err == nil {
		thread, err := threadUseCase.threadRepository.GetThreadById(thread.Id)
		if err == nil {
			return &thread, 200
		}
		return nil, 0
	}

	return nil, 0
}
