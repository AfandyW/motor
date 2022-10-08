package repository

import (
	"database/sql"

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
		return err
	}

	return nil
}

// lists motor
func (r *Repository) GetMotors() (motors []models.Motor, err error) {
	rows, err := r.db.Query("select id, name, price from motors")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var motor models.Motor

		err = rows.Scan(
			&motor.ID,
			&motor.Name,
			&motor.Price,
		)
		if err != nil {
			return nil, err
		}
		motors = append(motors, motor)
	}

	return motors, nil
}

// get motor
func (r *Repository) GetMotor(id int) (motor models.Motor, err error) {
	rows, err := r.db.Query("select id, name, price from motors where id = $1", id)
	if err != nil {
		return models.Motor{}, err
	}

	if rows.Next() {

		err = rows.Scan(
			&motor.ID,
			&motor.Name,
			&motor.Price,
		)
		if err != nil {
			return models.Motor{}, err
		}
	}

	return motor, nil
}

// delete motor
func (r *Repository) Delete(id int) (err error) {
	_, err = r.db.Exec("delete from motors where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// update
func (r *Repository) UpdateMotors(motor models.Motor) (err error) {
	_, err = r.db.Exec("update motors set name = $2, price = $3 where id = $1", motor.ID, motor.Name, motor.Price)
	if err != nil {
		return err
	}

	return nil
}
