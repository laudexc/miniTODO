# miniTODO (Bookshelf REST API)

Простой REST API-сервис для управления личной библиотекой.

## Функционал

- Добавление новых книг
- Получение списка всех книг
- Получение книги по названию (`title`)
- Фильтрация книг по автору (`author`)
- Фильтрация по статусу прочтения (`completed=true/false`)
- Отметка книги как прочитанной / непрочитанной
- Удаление книги

Каждая книга содержит:

- `Title` — название
- `Author` — автор
- `NumOfPages` — количество страниц
- `IsRead` — прочитана ли книга
- `AddedToShelf` — время добавления
- `ReadedAt` — время завершения чтения (или `null`)

## Запуск

Из корня проекта:

```powershell
go run .
```

Сервер стартует на:

`http://localhost:9091`

## Тесты

```powershell
go test ./...
```

Если в системе есть проблемы с правами на системный кэш Go:

```powershell
$env:GOCACHE = (Join-Path (Get-Location) '.gocache')
go test ./...
```

## API

### 1) Создать книгу

`POST /books`

```json
{
  "Title": "1984",
  "Author": "George Orwell",
  "NumOfPages": 328
}
```

Пример:

```powershell
curl -X POST "http://localhost:9091/books" `
  -H "Content-Type: application/json" `
  -d "{\"Title\":\"1984\",\"Author\":\"George Orwell\",\"NumOfPages\":328}"
```

### 2) Получить все книги

`GET /books`

Пример:

```powershell
curl "http://localhost:9091/books"
```

### 3) Получить книгу по названию

`GET /books?title=1984`

Пример:

```powershell
curl "http://localhost:9091/books?title=1984"
```

### 4) Получить книги по автору

`GET /books?author=George%20Orwell`

Пример:

```powershell
curl "http://localhost:9091/books?author=George%20Orwell"
```

### 5) Получить прочитанные книги

`GET /books?completed=true`

Пример:

```powershell
curl "http://localhost:9091/books?completed=true"
```

### 6) Получить непрочитанные книги

`GET /books?completed=false`

Пример:

```powershell
curl "http://localhost:9091/books?completed=false"
```

### 7) Отметить книгу как прочитанную / непрочитанную

`PATCH /books/{title}`

Тело запроса (вариант 1):

```json
{
  "complete": true
}
```

Тело запроса (вариант 2):

```json
{
  "completed": false
}
```

Пример:

```powershell
curl -X PATCH "http://localhost:9091/books/1984" `
  -H "Content-Type: application/json" `
  -d "{\"complete\":true}"
```

### 8) Удалить книгу

`DELETE /books/{title}`

Пример:

```powershell
curl -X DELETE "http://localhost:9091/books/1984"
```

## Коды ответов (основные)

- `200 OK` — успешный GET/PATCH
- `201 Created` — книга создана
- `204 No Content` — книга удалена
- `400 Bad Request` — невалидный запрос
- `404 Not Found` — книга не найдена
- `409 Conflict` — книга уже существует
- `500 Internal Server Error` — внутренняя ошибка
