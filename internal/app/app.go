package app

import (
	forumDelivery "TP_DB/internal/forum/delivery"
	forumRep "TP_DB/internal/forum/repository"
	forumUC "TP_DB/internal/forum/usecase"
	postDelivery "TP_DB/internal/post/delivery"
	postRep "TP_DB/internal/post/repository"
	postUC "TP_DB/internal/post/usecase"
	serviceDelivery "TP_DB/internal/service/delivery"
	serviceRep "TP_DB/internal/service/repository"
	serviceUC "TP_DB/internal/service/usecase"

	threadDelivery "TP_DB/internal/thread/delivery"
	threadRep "TP_DB/internal/thread/repository"
	threadUC "TP_DB/internal/thread/usecase"
	userDelivery "TP_DB/internal/user/delivery"
	userRep "TP_DB/internal/user/repository"
	userUC "TP_DB/internal/user/usecase"
	"github.com/jackc/pgx"
	routing "github.com/qiangxue/fasthttp-routing"
	"log"
)

type App struct {
	// options
	userDelivery    *userDelivery.UserDeliveryStruct
	forumDelivery   *forumDelivery.ForumDeliveryStruct
	postDelivery    *postDelivery.PostDeliveryStruct
	threadDelivery  *threadDelivery.ThreadDeliveryStruct
	serviceDelivery *serviceDelivery.ServiceDeliveryStruct
	//db           *sql.DB
}

const dbConfig = "host=127.0.0.1 port=5432 user=ksenia dbname=postgres password=password sslmode=disable"

func ConnectDatabase(config string) *pgx.ConnPool {
	pgxConnectionConfig, err := pgx.ParseConnectionString(config)
	if err != nil {
		log.Println(err)
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgxConnectionConfig,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		log.Println(err)
	}

	return pool
}

func (app *App) CreateRouter() *routing.Router {
	router := routing.New()

	app.userDelivery.SetHandlers(router)
	app.forumDelivery.SetHandlers(router)
	app.postDelivery.SetHandlers(router)
	app.threadDelivery.SetHandlers(router)
	app.serviceDelivery.SetHandlers(router)

	return router
}

func NewApp() (*App, error) {
	db := ConnectDatabase(dbConfig)

	userRepository := userRep.CreateUserRepository(db)
	postRepository := postRep.CreatePostRepository(db)
	threadRepository := threadRep.CreateThreadRepository(db)
	forumRepository := forumRep.CreateForumRepository(db)
	serviceRepository := serviceRep.CreateServiceRepository(db)

	userUseCase := userUC.CreateUserUseCase(userRepository)
	postUseCase := postUC.CreatePostUseCase(postRepository, userRepository, threadRepository, forumRepository)
	threadUseCase := threadUC.CreateThreadUseCase(threadRepository, userRepository, postRepository)
	forumUseCase := forumUC.CreateForumUseCase(forumRepository, userRepository, threadRepository)
	serviceUseCase := serviceUC.CreateServiceUseCase(serviceRepository)

	return &App{
		userDelivery:    userDelivery.CreateUserDelivery(userUseCase),
		forumDelivery:   forumDelivery.CreateForumDelivery(forumUseCase),
		postDelivery:    postDelivery.CreatePostDelivery(postUseCase),
		threadDelivery:  threadDelivery.CreateThreadDelivery(threadUseCase),
		serviceDelivery: serviceDelivery.CreateServiceDelivery(serviceUseCase),
	}, nil
}
