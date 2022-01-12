package delivery

import (
	usecase "TP_DB/internal/interfaces"
	"TP_DB/internal/models"
	errorsMsg "TP_DB/pkg"
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type PostDeliveryStruct struct {
	postUseCase usecase.PostUseCaseInterface
}

func CreatePostDelivery(postUseCase usecase.PostUseCaseInterface) *PostDeliveryStruct {
	return &PostDeliveryStruct{postUseCase: postUseCase}
}

func (postDelivery *PostDeliveryStruct) SetHandlers(router *routing.Router) {
	router.Get("/api/post/<id>/details", postDelivery.PostGetInfo)
	router.Post("/api/post/<id>/details", postDelivery.PostUpdateInfo)
}

func (postDelivery *PostDeliveryStruct) PostGetInfo(ctx *routing.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	related := ctx.QueryArgs().Peek("related")

	postFull, code := postDelivery.postUseCase.PostGetInfo(strings.Split(string(related), ","), id)
	data, err := json.Marshal(postFull)
	if err != nil {
		return err
	}

	ctx.SetContentType("application/json")
	switch code {
	case 404:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.SetBody(message)
	case 200:
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(data)
	}

	return nil
}

func (postDelivery *PostDeliveryStruct) PostUpdateInfo(ctx *routing.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	var postUpdate models.Post
	err := json.Unmarshal(ctx.PostBody(), &postUpdate)
	if err != nil {
		return err
	}
	postUpdate.Id = id

	post, code := postDelivery.postUseCase.PostUpdateInfo(id, postUpdate)
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	ctx.SetContentType("application/json")
	switch code {
	case 404:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.SetBody(message)
	case 200:
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(data)
	}

	return nil
}
