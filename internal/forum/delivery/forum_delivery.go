package delivery

import (
	usecase "TP_DB/internal/interfaces"
	"TP_DB/internal/models"
	errorsMsg "TP_DB/pkg"
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type ForumDeliveryStruct struct {
	forumUseCase usecase.ForumUseCaseInterface
}

func CreateForumDelivery(forumUseCase usecase.ForumUseCaseInterface) *ForumDeliveryStruct {
	return &ForumDeliveryStruct{forumUseCase: forumUseCase}
}

func (forumDelivery *ForumDeliveryStruct) SetHandlers(router *routing.Router) {
	router.Get("/api/forum/<slug>/threads", forumDelivery.ForumGetThreads)
	router.Get("/api/forum/<slug>/users", forumDelivery.ForumGetUsers)
	router.Post("/api/forum/create", forumDelivery.ForumCreate)
	router.Get("/api/forum/<slug>/details", forumDelivery.ForumGetBySlug)
	router.Post("/api/forum/<slug>/create", forumDelivery.ForumCreateThread)
}

func (forumDelivery *ForumDeliveryStruct) ForumCreateThread(ctx *routing.Context) error {
	var thread models.Thread
	slug := ctx.Param("slug")

	err := json.Unmarshal(ctx.PostBody(), &thread)
	if err != nil {
		return err
	}

	threadCreated, code := forumDelivery.forumUseCase.ForumCreateThread(&thread, slug)
	data, err := json.Marshal(threadCreated)

	if err != nil {
		return err
	}

	ctx.SetContentType("application/json")
	switch code {
	case 404:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.SetBody(message)
	case 409:
		ctx.Response.SetStatusCode(fasthttp.StatusConflict)
		ctx.SetBody(data)
	case 201:
		ctx.Response.SetStatusCode(fasthttp.StatusCreated)
		ctx.SetBody(data)
	}

	return nil
}

func (forumDelivery *ForumDeliveryStruct) ForumCreate(ctx *routing.Context) error {
	var forum models.Forum

	err := json.Unmarshal(ctx.PostBody(), &forum)
	if err != nil {
		return err
	}

	forumCreated, code := forumDelivery.forumUseCase.ForumCreate(&forum)
	data, err := json.Marshal(forumCreated)

	if err != nil {
		return err
	}

	ctx.SetContentType("application/json")
	switch code {
	case 404:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.SetBody(message)
	case 409:
		ctx.Response.SetStatusCode(fasthttp.StatusConflict)
		ctx.SetBody(data)
	case 201:
		ctx.Response.SetStatusCode(fasthttp.StatusCreated)
		ctx.SetBody(data)
	}

	return nil
}

func (forumDelivery *ForumDeliveryStruct) ForumGetBySlug(ctx *routing.Context) error {
	slug := ctx.Param("slug")

	forum, code := forumDelivery.forumUseCase.ForumGetBySlug(slug)
	data, err := json.Marshal(forum)

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

func (forumDelivery *ForumDeliveryStruct) ForumGetThreads(ctx *routing.Context) error {
	ctx.SetContentType("application/json")

	slug := ctx.Param("slug")
	limit := ctx.QueryArgs().GetUintOrZero("limit")
	if limit == 0 {
		limit = 100
	}
	desc := func() bool {
		if ctx.QueryArgs().Peek("desc") != nil && ctx.QueryArgs().Peek("desc")[0] == 't' {
			return true
		} else {
			return false
		}
	}()
	since := string(ctx.QueryArgs().Peek("since"))

	threads, code := forumDelivery.forumUseCase.ForumGetThreads(slug, limit, desc, since)

	switch code {
	case 404:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.SetBody(message)
	case 200:
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		if len(*threads) > 0 {
			data, err := json.Marshal(threads)
			if err != nil {
				return err
			}
			ctx.SetBody(data)
		} else {
			ctx.SetBodyString("[]")
		}
	}

	return nil
}

func (forumDelivery *ForumDeliveryStruct) ForumGetUsers(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	limit := ctx.QueryArgs().GetUintOrZero("limit")
	if limit == 0 {
		limit = 100
	}
	desc := func() bool {
		if ctx.QueryArgs().Peek("desc") != nil && ctx.QueryArgs().Peek("desc")[0] == 't' {
			return true
		} else {
			return false
		}
	}()
	since := string(ctx.QueryArgs().Peek("since"))

	users, code := forumDelivery.forumUseCase.ForumGetUsers(slug, limit, desc, since)

	ctx.SetContentType("application/json")
	switch code {
	case 404:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.SetBody(message)
	case 200:
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		if len(*users) > 0 {
			data, err := json.Marshal(users)
			if err != nil {
				return err
			}
			ctx.SetBody(data)
		} else {
			ctx.SetBodyString("[]")
		}
	}
	return nil
}
