package models

type Motor struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (m Motor) Validate() error {
	cErr := []string{}
	if m.Name == "" {
		cErr = append(cErr, "name cannot null")
	}

	if m.Price == 0.0 {
		cErr = append(cErr, "price cannot null")
	}

	if len(cErr) > 0 {
		return NewBadRequest(cErr)
	}

	return nil
}

func (m Motor) Exist() error {
	if m.ID == 0 {
		return NewNotFound("motor not found")
	}
	return nil
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
