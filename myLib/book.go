package myLib

import "time"

type Book struct {
	Title      string
	Author     string
	NumOfPages int
	IsRead     bool

	AddedToShelf time.Time
	ReadedAt     *time.Time
}

func NewBook(title string, author string, numPgs int) Book {
	return Book{
		Title:      title,
		Author:     author,
		NumOfPages: numPgs,
		IsRead:     false,

		AddedToShelf: time.Now(),
		ReadedAt:     nil,
	}
}

func (b *Book) Complete() {
	completeTime := time.Now()

	b.IsRead = true
	b.ReadedAt = &completeTime
}

func (b *Book) Uncomplete() {
	b.IsRead = false
	b.ReadedAt = nil
}
