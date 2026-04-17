package http

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
)

type completedBookDTO struct {
	Complete  *bool `json:"complete"`
	Completed *bool `json:"completed"`
}

func (cb completedBookDTO) CompletionValue() (bool, error) {
	if cb.Complete == nil && cb.Completed == nil {
		return false, errors.New("body must contain 'complete' or 'completed' field")
	}

	if cb.Complete != nil && cb.Completed != nil && *cb.Complete != *cb.Completed {
		return false, errors.New("'complete' and 'completed' fields conflict")
	}

	if cb.Completed != nil {
		return *cb.Completed, nil
	}

	return *cb.Complete, nil
}

type BookDTO struct {
	Title      string
	Author     string
	NumOfPages int
}

func (b BookDTO) ValidateForCreate() error { // Валидация только для создания
	if strings.TrimSpace(b.Title) == "" {
		return errors.New("title is empty")
	}

	if strings.TrimSpace(b.Author) == "" {
		return errors.New("author is empty")
	}

	if b.NumOfPages <= 0 {
		return errors.New("the number of pages cannot be less than or equal to zero")
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
