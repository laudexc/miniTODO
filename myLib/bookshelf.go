package myLib

import "sync"

type Bookshelf struct {
	books map[string]Book
	mtx   sync.RWMutex
}

func NewBookshelf() *Bookshelf {
	return &Bookshelf{
		books: make(map[string]Book),
	}
}

// добавить книгу
func (b *Bookshelf) AddBook(book Book) error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if _, ok := b.books[book.Title]; ok {
		return ErrBookAlreadyExists
	}

	b.books[book.Title] = book
	return nil
}

// Получение списка всех книг
func (b *Bookshelf) ListBooks() map[string]Book {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	tmp := make(map[string]Book, len(b.books))

	for k, v := range b.books {
		tmp[k] = v
	}

	return tmp
}

// Получение информации о конкретной книге
func (b *Bookshelf) GetBook(title string) (Book, error) {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	task, ok := b.books[title]
	if !ok {
		return Book{}, ErrBookNotFound
	}

	return task, nil
}

// Получение списка всех книг, с учетом возможной фильтрации по автору
func (b *Bookshelf) ListByAuthorBook(author string) map[string]Book {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	byAuthorBooks := make(map[string]Book)
	for title, book := range b.books {
		if book.Author == author {
			byAuthorBooks[title] = book
		}
	}

	return byAuthorBooks
}

// Получение списка всех книг, с учетом возможной фильтрации по прочитано?
func (b *Bookshelf) ListUncompletedBooks() map[string]Book {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	uncompletedBooks := make(map[string]Book)
	for title, book := range b.books {
		if !book.IsRead {
			uncompletedBooks[title] = book
		}
	}

	return uncompletedBooks
}

// Получение списка всех книг, с учетом возможной фильтрации по не прочитано?
func (b *Bookshelf) ListCompletedBooks() map[string]Book {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	completedBooks := make(map[string]Book)
	for title, book := range b.books {
		if book.IsRead {
			completedBooks[title] = book
		}
	}

	return completedBooks
}

// Отмечать отдельные книги как прочитанные
func (b *Bookshelf) ReadBook(title string) (Book, error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	book, ok := b.books[title]
	if !ok {
		return Book{}, ErrBookNotFound
	}

	book.Complete()
	b.books[title] = book

	return book, nil
}

// Отмечать отдельные книги как непрочитанные
func (b *Bookshelf) UnreadBook(title string) (Book, error) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	book, ok := b.books[title]
	if !ok {
		return Book{}, ErrBookNotFound
	}

	book.Uncomplete()

	b.books[title] = book

	return book, nil
}

// Удаление книги из библиотеки
func (b *Bookshelf) DeleteBook(title string) error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	_, ok := b.books[title]
	if !ok {
		return ErrBookNotFound
	}

	delete(b.books, title)

	return nil
}
