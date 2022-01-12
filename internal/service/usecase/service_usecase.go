package usecase

import (
	repository "TP_DB/internal/interfaces"
	"TP_DB/internal/models"
)

type ServiceUseCaseStruct struct {
	serviceRepository repository.ServiceRepositoryInterface
}

func CreateServiceUseCase(serviceRepository repository.ServiceRepositoryInterface) *ServiceUseCaseStruct {
	return &ServiceUseCaseStruct{serviceRepository: serviceRepository}
}

func (serviceUseCase *ServiceUseCaseStruct) ServiceClear() int {
	err := serviceUseCase.serviceRepository.Clear()
	if err == nil {
		return 200
	}
	return 0
}

func (serviceUseCase *ServiceUseCaseStruct) ServiceGetStatus() (*models.Status, int) {
	status, err := serviceUseCase.serviceRepository.GetStatus()
	if err == nil {
		return &status, 200
	}
	return nil, 0
}
