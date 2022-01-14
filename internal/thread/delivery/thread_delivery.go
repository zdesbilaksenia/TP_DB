package delivery

import (
	usecase "TP_DB/internal/interfaces"
	"TP_DB/internal/models"
	errorsMsg "TP_DB/pkg"
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"strconv"
)

type ThreadDeliveryStruct struct {
	threadUseCase usecase.ThreadUseCaseInterface
}

func CreateThreadDelivery(threadUseCase usecase.ThreadUseCaseInterface) *ThreadDeliveryStruct {
	return &ThreadDeliveryStruct{threadUseCase: threadUseCase}
}

func (threadDelivery *ThreadDeliveryStruct) SetHandlers(router *routing.Router) {
	router.Post("/api/thread/<slug_or_id>/create", threadDelivery.ThreadCreatePosts)
	router.Get("/api/thread/<slug_or_id>/posts", threadDelivery.ThreadGetPosts)
	router.Post("/api/thread/<slug_or_id>/details", threadDelivery.ThreadUpdate)
	router.Get("/api/thread/<slug_or_id>/details", threadDelivery.ThreadGet)
	router.Post("/api/thread/<slug_or_id>/vote", threadDelivery.ThreadVote)
}

func (threadDelivery *ThreadDeliveryStruct) ThreadGet(ctx *routing.Context) error {
	var slug string
	id, err := strconv.Atoi(ctx.Param("slug_or_id"))
	if err != nil {
		slug = ctx.Param("slug_or_id")
		id = -1
	} else {
		slug = ""
	}

	thread, code := threadDelivery.threadUseCase.ThreadGet(slug, id)
	data, err := json.Marshal(thread)
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

func (threadDelivery *ThreadDeliveryStruct) ThreadUpdate(ctx *routing.Context) error {
	var update models.ThreadUpdate
	err := json.Unmarshal(ctx.PostBody(), &update)
	if err != nil {
		return err
	}

	var slug string
	id, err := strconv.Atoi(ctx.Param("slug_or_id"))
	if err != nil {
		slug = ctx.Param("slug_or_id")
		id = -1
	} else {
		slug = ""
	}

	thread, code := threadDelivery.threadUseCase.ThreadUpdate(update, slug, id)
	data, err := json.Marshal(thread)
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

func (threadDelivery *ThreadDeliveryStruct) ThreadCreatePosts(ctx *routing.Context) error {
	var posts models.Posts
	err := json.Unmarshal(ctx.PostBody(), &posts)
	if err != nil {
		return err
	}

	var slug string
	id, err := strconv.Atoi(ctx.Param("slug_or_id"))
	if err != nil {
		slug = ctx.Param("slug_or_id")
		id = -1
	} else {
		slug = ""
	}

	createdPosts, code := threadDelivery.threadUseCase.ThreadCreatePosts(slug, id, posts)

	ctx.SetContentType("application/json")
	switch code {
	case 409:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusConflict)
		ctx.Response.SetBody(message)
	case 404:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.SetBody(message)
	case 201:
		ctx.Response.SetStatusCode(fasthttp.StatusCreated)
		if len(createdPosts) > 0 {
			data, err := json.Marshal(createdPosts)
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

func (threadDelivery *ThreadDeliveryStruct) ThreadVote(ctx *routing.Context) error {
	var vote models.Vote
	err := json.Unmarshal(ctx.PostBody(), &vote)
	if err != nil {
		return err
	}

	var slug string
	id, err := strconv.Atoi(ctx.Param("slug_or_id"))
	if err != nil {
		slug = ctx.Param("slug_or_id")
		id = -1
	} else {
		slug = ""
	}

	thread, code := threadDelivery.threadUseCase.ThreadVote(vote, slug, id)
	data, err := json.Marshal(thread)
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

func (threadDelivery *ThreadDeliveryStruct) ThreadGetPosts(ctx *routing.Context) error {
	var slug string
	id, err := strconv.Atoi(ctx.Param("slug_or_id"))
	if err != nil {
		slug = ctx.Param("slug_or_id")
		id = -1
	} else {
		slug = ""
	}
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
	sort := string(ctx.QueryArgs().Peek("sort"))

	posts, code := threadDelivery.threadUseCase.ThreadGetPosts(slug, id, limit, since, desc, sort)

	ctx.SetContentType("application/json")
	switch code {
	case 404:
		message, _ := json.Marshal(models.Err{Message: errorsMsg.UserNotExist})
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Response.SetBody(message)
	case 200:
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		if len(*posts) > 0 {
			data, err := json.Marshal(posts)
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
