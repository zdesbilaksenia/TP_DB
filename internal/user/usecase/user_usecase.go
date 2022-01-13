package usecase

import (
	repository "TP_DB/internal/interfaces"
	"TP_DB/internal/models"
)

type UserUseCaseStruct struct {
	userRepository repository.UserRepositoryInterface
}

func CreateUserUseCase(userRepository repository.UserRepositoryInterface) *UserUseCaseStruct {
	return &UserUseCaseStruct{userRepository: userRepository}
}

func (userUseCase *UserUseCaseStruct) UserCreate(user *models.User) (models.User, models.Users, int) {
	users, err := userUseCase.userRepository.UsersGet(user)
	if err == nil && len(users) != 0 {
		return models.User{}, users, 409
	}

	err = userUseCase.userRepository.UserCreate(user)

	if err != nil {
		return models.User{}, nil, 0
	}

	userResult, err := userUseCase.userRepository.UserGet(user.Nickname)

	if err != nil {
		return models.User{}, nil, 0
	}

	return userResult, nil, 201
}

func (userUseCase *UserUseCaseStruct) UserGet(nickname string) (models.User, error) {
	user, err := userUseCase.userRepository.UserGet(nickname)

	return user, err
}

func (userUseCase *UserUseCaseStruct) UserChange(user models.User) (models.User, models.Users, int) {
	users, err := userUseCase.userRepository.UsersGet(&user)
	if err == nil && len(users) == 0 {
		return models.User{}, users, 404
	}
	if err == nil && len(users) > 1 {
		return models.User{}, users, 409
	}

	user, err = userUseCase.userRepository.UserChange(user)
	if err != nil {
		return models.User{}, nil, 0
	}
	return user, nil, 200
}

func (userUseCase *UserUseCaseStruct) UsersGet(user models.User) (models.Users, error) {
	users, err := userUseCase.userRepository.UsersGet(&user)

	return users, err
}
