package delivery

import (
	usecase "TP_DB/internal/interfaces"
	"TP_DB/internal/models"
	errorsMsg "TP_DB/pkg"
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)

type UserDeliveryStruct struct {
	userUseCase usecase.UserUseCaseInterface
}

func CreateUserDelivery(userUseCase usecase.UserUseCaseInterface) *UserDeliveryStruct {
	return &UserDeliveryStruct{userUseCase: userUseCase}
}

func (userDelivery *UserDeliveryStruct) SetHandlers(router *routing.Router) {
	//api := router.Group("/api")
	router.Post("/api/user/<nickname>/create", userDelivery.UserCreate)
	router.Get("/api/user/<nickname>/profile", userDelivery.UserGet)
	router.Post("/api/user/<nickname>/profile", userDelivery.UserChange)

	log.Println("user delivery handlers are set")
}

func (userDelivery *UserDeliveryStruct) UserCreate(ctx *routing.Context) error {
	log.Println("user create request start")

	var user models.User
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		return err
	}

	nickname := ctx.Param("nickname")
	user.Nickname = nickname

	userDB, err := userDelivery.userUseCase.UserGet(nickname)

	if err == nil {
		data, err := json.Marshal(userDB)

		if err != nil {
			return err
		}

		ctx.Response.SetStatusCode(fasthttp.StatusConflict)
		ctx.SetContentType("application/json")
		ctx.Response.SetBody(data)

		log.Println("user already exists")

		return nil
	}

	createdUser, err := userDelivery.userUseCase.UserCreate(&user)

	if err != nil {
		return err
	}

	data, err := json.Marshal(createdUser)

	if err != nil {
		return err
	}

	ctx.Response.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetContentType("application/json")
	ctx.Response.SetBody(data)

	log.Println("user create request finish")

	return nil
}

func (userDelivery *UserDeliveryStruct) UserGet(ctx *routing.Context) error {
	log.Println("user get request start")

	nickname := ctx.Param("nickname")

	user, err := userDelivery.userUseCase.UserGet(nickname)

	if err != nil {
		log.Println(err)

		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetBody(message)
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		return nil
	}

	data, err := json.Marshal(user)

	if err != nil {
		return err
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.Response.SetBody(data)

	log.Println("user get request finish")

	return nil
}

func (userDelivery *UserDeliveryStruct) UserChange(ctx *routing.Context) error {
	log.Println("user change request start")

	var user models.User

	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		return err
	}

	nickname := ctx.Param("nickname")
	user.Nickname = nickname

	userDB, err := userDelivery.userUseCase.UserGet(nickname)

	if err != nil {
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetBody(message)
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		return nil
	}

	usersWithSameParams, err := userDelivery.userUseCase.UsersGet(userDB)

	if err != nil {
		return err
	}

	if len(usersWithSameParams) > 1 {
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})

		ctx.Response.SetStatusCode(fasthttp.StatusConflict)
		ctx.SetContentType("application/json")
		ctx.Response.SetBody(message)

		log.Println("users with same params")

		return nil
	}

	changedUser, err := userDelivery.userUseCase.UserChange(user)

	if err != nil {
		return err
	}

	data, err := json.Marshal(changedUser)

	if err != nil {
		return err
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	ctx.Response.SetBody(data)

	log.Println("user change request finish")

	return nil
}
