package service

import (
	"errors"

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
	err := svc.repository.CreateMotors(motor)
	if err != nil {
		return err
	}

	return nil
}

// update
func (svc *Service) Update(newMotor models.Motor) error {
	motor, err := svc.repository.GetMotor(newMotor.ID)
	if err != nil {
		return err
	}

	if motor.ID == 0 {
		return errors.New("data motors not found")
	}

	motor.Name = newMotor.Name
	motor.Price = newMotor.Price

	err = svc.repository.UpdateMotors(motor)
	if err != nil {
		return err
	}

	return nil
}

// get all
func (svc *Service) List() ([]models.Motor, error) {
	motors, err := svc.repository.GetMotors()
	if err != nil {
		return nil, err
	}

	return motors, nil
}

// get
func (svc *Service) Get(id int) (models.Motor, error) {
	motor, err := svc.repository.GetMotor(id)
	if err != nil {
		return models.Motor{}, err
	}

	return motor, nil
}

// delete
func (svc *Service) Delete(id int) error {
	motor, err := svc.repository.GetMotor(id)
	if err != nil {
		return err
	}

	if motor.ID == 0 {
		return errors.New("data motors not found")
	}

	err = svc.repository.Delete(motor.ID)
	if err != nil {
		return err
	}

	return nil
}
