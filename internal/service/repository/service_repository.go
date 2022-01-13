package repository

import (
	"TP_DB/internal/models"
	"github.com/jackc/pgx"
)

const (
	Clean       = `TRUNCATE thread, forum, users, votes, post, nickname_forum CASCADE`
	CountPost   = `SELECT COUNT(*) FROM post`
	CountUser   = `SELECT COUNT(*) FROM users`
	CountForum  = `SELECT COUNT(*) FROM forum`
	CountThread = `SELECT COUNT(*) FROM thread`
)

type ServiceRepositoryStruct struct {
	DB *pgx.ConnPool
}

func CreateServiceRepository(DB *pgx.ConnPool) *ServiceRepositoryStruct {
	return &ServiceRepositoryStruct{DB: DB}
}

func (serviceRepository *ServiceRepositoryStruct) GetStatus() (models.Status, error) {
	var status models.Status

	err := serviceRepository.DB.QueryRow(CountPost).Scan(&status.Post)
	err = serviceRepository.DB.QueryRow(CountUser).Scan(&status.User)
	err = serviceRepository.DB.QueryRow(CountForum).Scan(&status.Forum)
	err = serviceRepository.DB.QueryRow(CountThread).Scan(&status.Thread)

	return status, err
}

func (serviceRepository *ServiceRepositoryStruct) Clear() error {
	_, err := serviceRepository.DB.Exec(Clean)

	return err
}
