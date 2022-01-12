package repository

import (
	"TP_DB/internal/models"
	"github.com/jackc/pgx"
)

const (
	CreateUser        = `INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)`
	GetUsers          = `SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2`
	GetUserByNickname = `SELECT nickname, fullname, about, email FROM users WHERE nickname = $1`
	ChangeUser        = `UPDATE users SET about = COALESCE(NULLIF($1, ''), about), email = COALESCE(NULLIF($2, ''), email), 
						fullname = COALESCE(NULLIF($3, ''), fullname) WHERE nickname = $4 RETURNING nickname, fullname, about, email`
)

type UserRepositoryStruct struct {
	DB *pgx.ConnPool
}

func CreateUserRepository(DB *pgx.ConnPool) *UserRepositoryStruct {
	return &UserRepositoryStruct{DB: DB}
}

func (userRepository *UserRepositoryStruct) UserCreate(user *models.User) error {
	_, err := userRepository.DB.Exec(CreateUser, user.Nickname, user.Fullname, user.About, user.Email)

	return err
}

func (userRepository *UserRepositoryStruct) UsersGet(user *models.User) (models.Users, error) {
	var users models.Users

	rows, _ := userRepository.DB.Query(GetUsers, user.Nickname, user.Email)

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)

		if err != nil {
			rows.Close()
			return users, err
		}

		users = append(users, user)
	}

	rows.Close()
	return users, nil
}

func (userRepository *UserRepositoryStruct) UserGet(nickname string) (models.User, error) {
	var user = models.User{Nickname: nickname}

	err := userRepository.DB.QueryRow(GetUserByNickname, nickname).
		Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)

	return user, err
}

func (userRepository *UserRepositoryStruct) UserChange(user models.User) (models.User, error) {
	var changedUser models.User

	err := userRepository.DB.QueryRow(ChangeUser, user.About, user.Email, user.Fullname, user.Nickname).
		Scan(&changedUser.Nickname, &changedUser.Fullname, &changedUser.About, &changedUser.Email)

	return changedUser, err
}
