package repository

import (
	"database/sql"
	"errors"

	"github.com/AfandyW/motor/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	// validate
	return &Repository{
		db: db,
	}
}

// create
func (r *Repository) CreateMotors(motor models.Motor) (err error) {
	_, err = r.db.Exec("insert into motors(name,price) values($1,$2)", motor.Name, motor.Price)
	if err != nil {
		return errors.New("create motors return error : " + err.Error())
	}

	return nil
}

// lists motor
func (r *Repository) GetMotors() (motors []models.Motor, err error) {
	rows, err := r.db.Query("select id, name, price from motors")
	if err != nil {
		return nil, errors.New("get motors return error : " + err.Error())
	}

	for rows.Next() {
		var motor models.Motor

		err = rows.Scan(
			&motor.ID,
			&motor.Name,
			&motor.Price,
		)
		if err != nil {
			return nil, errors.New("get motors, scan data return error : " + err.Error())
		}
		motors = append(motors, motor)
	}

	return motors, nil
}

// get motor
func (r *Repository) GetMotor(id int) (motor models.Motor, err error) {
	rows, err := r.db.Query("select id, name, price from motors where id = $1", id)
	if err != nil {
		return models.Motor{}, errors.New("get motor return error : " + err.Error())
	}

	if rows.Next() {

		err = rows.Scan(
			&motor.ID,
			&motor.Name,
			&motor.Price,
		)
		if err != nil {
			return models.Motor{}, errors.New("get motor, scan data return error : " + err.Error())
		}
	}

	return motor, nil
}

// delete motor
func (r *Repository) Delete(id int) (err error) {
	_, err = r.db.Exec("delete from motors where id = $1", id)
	if err != nil {
		return errors.New("delete motor return error : " + err.Error())
	}
	return nil
}

// update
func (r *Repository) UpdateMotors(motor models.Motor) (err error) {
	_, err = r.db.Exec("update motors set name = $2, price = $3 where id = $1", motor.ID, motor.Name, motor.Price)
	if err != nil {
		return errors.New("update motor return error : " + err.Error())
	}

	return nil
}
