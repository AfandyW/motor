package repository

import (
	"database/sql"

	"github.com/AfandyW/motor/models"
)

// create
func CreateMotors(db *sql.DB, motor models.Motor) (err error) {
	_, err = db.Exec("insert into motors(name,price) values($1,$2)", motor.Name, motor.Price)
	if err != nil {
		return err
	}

	return nil
}

// lists motor
func GetMotors(db *sql.DB) (motors []models.Motor, err error) {
	rows, err := db.Query("select id, name, price from motors")
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
func GetMotor(db *sql.DB, id string) (motor models.Motor, err error) {
	rows, err := db.Query("select id, name, price from motors where id = $1", id)
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
func Delete(db *sql.DB, id string) (err error) {
	_, err = db.Exec("delete from motors where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// update
func UpdateMotors(db *sql.DB, motor models.Motor) (err error) {
	_, err = db.Exec("update motors set name = $2, price = $3 where id = $1", motor.ID ,motor.Name, motor.Price)
	if err != nil {
		return err
	}

	return nil
}