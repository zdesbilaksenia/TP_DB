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

func (userUseCase *UserUseCaseStruct) UserCreate(user *models.User) (models.User, error) {
	err := userUseCase.userRepository.UserCreate(user)

	if err != nil {
		return models.User{}, err
	}

	userResult, err := userUseCase.userRepository.UserGet(user.Nickname)

	if err != nil {
		return models.User{}, err
	}

	return userResult, nil
}

func (userUseCase *UserUseCaseStruct) UserGet(nickname string) (models.User, error) {
	user, err := userUseCase.userRepository.UserGet(nickname)

	return user, err
}

func (userUseCase *UserUseCaseStruct) UserChange(user models.User) (models.User, error) {
	user, err := userUseCase.userRepository.UserChange(user)

	return user, err
}

func (userUseCase *UserUseCaseStruct) UsersGet(user models.User) (models.Users, error) {
	users, err := userUseCase.userRepository.UsersGet(&user)

	return users, err
}
