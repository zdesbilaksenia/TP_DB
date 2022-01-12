package usecase

import (
	repository "TP_DB/internal/interfaces"
	"TP_DB/internal/models"
)

type PostUseCaseStruct struct {
	postRepository   repository.PostRepositoryInterface
	userRepository   repository.UserRepositoryInterface
	threadRepository repository.ThreadRepositoryInterface
	forumRepository  repository.ForumRepositoryInterface
}

func CreatePostUseCase(postRepository repository.PostRepositoryInterface, userRepository repository.UserRepositoryInterface, threadRepository repository.ThreadRepositoryInterface, forumRepository repository.ForumRepositoryInterface) *PostUseCaseStruct {
	return &PostUseCaseStruct{postRepository: postRepository, userRepository: userRepository, threadRepository: threadRepository, forumRepository: forumRepository}
}

func (postUseCase *PostUseCaseStruct) PostGetInfo(related []string, id int) (*models.PostFull, int) {
	post, err := postUseCase.postRepository.GetPostInfo(id)

	if err != nil {
		return nil, 404
	}

	var postFull models.PostFull
	postFull.Post = &post

	for _, rel := range related {
		switch rel {
		case "user":
			author, err := postUseCase.userRepository.UserGet(post.Author)
			if err != nil {
				return nil, 404
			}
			postFull.Author = &author
		case "forum":
			forum, err := postUseCase.forumRepository.GetForumBySlug(post.Forum)
			if err != nil {
				return nil, 404
			}
			postFull.Forum = &forum
		case "thread":
			thread, err := postUseCase.threadRepository.GetThreadById(post.Thread)
			if err != nil {
				return nil, 404
			}
			postFull.Thread = &thread
		}
	}

	return &postFull, 200
}

func (postUseCase *PostUseCaseStruct) PostUpdateInfo(id int, postUpdate models.Post) (*models.Post, int) {
	_, err := postUseCase.postRepository.GetPostInfo(id)

	if err != nil {
		return nil, 404
	}

	postEdited, err := postUseCase.postRepository.UpdatePostInfo(postUpdate)
	if err == nil {
		return &postEdited, 200
	}

	return nil, 0
}
