package http

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
)

type completeTaskDTO struct {
	Complete bool
}

type TaskDTO struct {
	Title       string
	Description string
}

func (t TaskDTO) ValidateForCreate() error { // Валидация только для создания
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title is empty")
	}

	if strings.TrimSpace(t.Description) == "" {
		return errors.New("description is empty")
	}

	return nil
}

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		log.Fatalln("Impossible error in MarshalIndent", err)
	}

	return string(b)
}

func NewErrorDTO(err error) ErrorDTO {
	errDTO := ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}

	return errDTO
}
