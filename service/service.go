package service

import (
	"github.com/AfandyW/motor/models"
	"github.com/AfandyW/motor/repository"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

// create
func (svc *Service) Create(motor models.Motor) error {
	if err := svc.repository.CreateMotors(motor); err != nil {
		return models.NewInternalServerError(err.Error())
	}

	return nil
}

// update
func (svc *Service) Update(newMotor models.Motor) error {
	motor, err := svc.repository.GetMotor(newMotor.ID)
	if err != nil {
		return models.NewInternalServerError(err.Error())
	}

	if err = motor.Exist(); err != nil {
		return err
	}

	motor.Name = newMotor.Name
	motor.Price = newMotor.Price

	if err = svc.repository.UpdateMotors(motor); err != nil {
		return models.NewInternalServerError(err.Error())
	}

	return nil
}

// get all
func (svc *Service) List() ([]models.Motor, error) {
	motors, err := svc.repository.GetMotors()
	if err != nil {
		return nil, models.NewInternalServerError(err.Error())
	}

	return motors, nil
}

// get
func (svc *Service) Get(id int) (models.Motor, error) {
	motor, err := svc.repository.GetMotor(id)
	if err != nil {
		return models.Motor{}, models.NewInternalServerError(err.Error())
	}

	if err = motor.Exist(); err != nil {
		return models.Motor{}, err
	}

	return motor, nil
}

// delete
func (svc *Service) Delete(id int) error {
	motor, err := svc.repository.GetMotor(id)
	if err != nil {
		return models.NewInternalServerError(err.Error())
	}

	if err = motor.Exist(); err != nil {
		return err
	}

	if err = svc.repository.Delete(motor.ID); err != nil {
		return models.NewInternalServerError(err.Error())
	}

	return nil
}
