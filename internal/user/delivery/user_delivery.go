package delivery

import (
	usecase "TP_DB/internal/interfaces"
	"TP_DB/internal/models"
	errorsMsg "TP_DB/pkg"
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
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
}

func (userDelivery *UserDeliveryStruct) UserCreate(ctx *routing.Context) error {
	ctx.SetContentType("application/json")

	var user models.User
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		return err
	}

	nickname := ctx.Param("nickname")
	user.Nickname = nickname

	createdUser, createdUsers, code := userDelivery.userUseCase.UserCreate(&user)
	dataUser, err := json.Marshal(createdUser)
	if err != nil {
		return err
	}

	dataUsers, err := json.Marshal(createdUsers)
	if err != nil {
		return err
	}

	switch code {
	case 409:
		ctx.Response.SetStatusCode(fasthttp.StatusConflict)
		ctx.SetBody(dataUsers)
	case 201:
		ctx.Response.SetStatusCode(fasthttp.StatusCreated)
		ctx.SetBody(dataUser)
	}

	return nil
}

func (userDelivery *UserDeliveryStruct) UserGet(ctx *routing.Context) error {
	ctx.SetContentType("application/json")

	nickname := ctx.Param("nickname")

	user, err := userDelivery.userUseCase.UserGet(nickname)

	if err != nil {
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
	ctx.Response.SetBody(data)

	return nil
}

func (userDelivery *UserDeliveryStruct) UserChange(ctx *routing.Context) error {
	ctx.SetContentType("application/json")

	var user models.User

	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		return err
	}

	nickname := ctx.Param("nickname")
	user.Nickname = nickname

	changedUser, _, code := userDelivery.userUseCase.UserChange(user)
	data, err := json.Marshal(changedUser)

	if err != nil {
		return err
	}

	switch code {
	case 404:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.SetBody(message)
	case 409:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusConflict)
		ctx.Response.SetBody(message)
	case 200:
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(data)
	}

	return nil
}
