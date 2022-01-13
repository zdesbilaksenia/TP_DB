package delivery

import (
	usecase "TP_DB/internal/interfaces"
	"encoding/json"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type ServiceDeliveryStruct struct {
	serviceUseCase usecase.ServiceUseCaseInterface
}

func CreateServiceDelivery(serviceUseCase usecase.ServiceUseCaseInterface) *ServiceDeliveryStruct {
	return &ServiceDeliveryStruct{serviceUseCase: serviceUseCase}
}

func (serviceDelivery *ServiceDeliveryStruct) SetHandlers(router *routing.Router) {
	router.Get("/api/service/status", serviceDelivery.ServiceGetStatus)
	router.Post("/api/service/clear", serviceDelivery.ServiceClear)
}

func (serviceDelivery *ServiceDeliveryStruct) ServiceClear(ctx *routing.Context) error {
	ctx.SetContentType("application/json")
	code := serviceDelivery.serviceUseCase.ServiceClear()
	switch code {
	case 200:
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
	}
	return nil
}

func (serviceDelivery *ServiceDeliveryStruct) ServiceGetStatus(ctx *routing.Context) error {
	ctx.SetContentType("application/json")
	status, code := serviceDelivery.serviceUseCase.ServiceGetStatus()
	data, err := json.Marshal(status)
	if err != nil {
		return err
	}
	switch code {
	case 200:
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(data)
	}
	return nil
}
