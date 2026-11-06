# miniTODO

REST API для списка книг. Данные хранятся в памяти и сбрасываются при перезапуске сервера.

## Запуск

Нужен Go **1.25.7 или новее**. В папке с `go.mod` выполните:

```sh
go run .
```

Откройте [http://localhost:9091/books](http://localhost:9091/books) в браузере. Если список пуст, сервер готов принимать запросы. Оставьте его работать и отправляйте запросы из другого терминала или из Postman. Для остановки нажмите **Ctrl+C**.

## Запросы к API

К каждому пути в таблице добавляйте адрес `http://localhost:9091`. Для `POST` и `PATCH` передавайте JSON с заголовком `Content-Type: application/json`.

| Запрос | Действие | Пример тела |
| --- | --- | --- |
| `POST /books` | Добавить книгу | `{"Title":"1984","Author":"George Orwell","NumOfPages":328}` |
| `GET /books` | Получить весь список | |
| `GET /books?title=1984` | Найти книгу по названию | |
| `GET /books?author=George%20Orwell` | Найти книги автора | |
| `GET /books?completed=true` | Показать прочитанные книги. Для непрочитанных укажите `false` | |
| `PATCH /books/1984` | Отметить книгу как прочитанную. Для обратного действия укажите `false` | `{"completed":true}` |
| `DELETE /books/1984` | Удалить книгу по названию | |

### Bash, Zsh или Git Bash

```bash
url=http://localhost:9091/books

curl -fsS -H 'Content-Type: application/json' \
  -d '{"Title":"1984","Author":"George Orwell","NumOfPages":328}' "$url"
curl -fsS "$url?author=George%20Orwell"
curl -fsS -X PATCH -H 'Content-Type: application/json' \
  -d '{"completed":true}' "$url/1984"
curl -fsS "$url?completed=true"
curl -fsS -X DELETE "$url/1984"
```

### PowerShell

```powershell
$url = 'http://localhost:9091/books'

Invoke-RestMethod $url -Method Post -ContentType 'application/json' -Body '{"Title":"1984","Author":"George Orwell","NumOfPages":328}'
Invoke-RestMethod "${url}?author=George%20Orwell"
Invoke-RestMethod "$url/1984" -Method Patch -ContentType 'application/json' -Body '{"completed":true}'
Invoke-RestMethod "${url}?completed=true"
Invoke-RestMethod "$url/1984" -Method Delete
```

### Postman

Выберите метод из таблицы и введите полный URL. Для `POST` и `PATCH` в разделе **Body** выберите **raw → JSON** и вставьте тело из таблицы. Нажмите **Send**.

При успешном добавлении сервер возвращает `201 Created`, при удалении — `204 No Content` с пустым телом. Если добавить книгу с тем же названием повторно, сервер вернёт `409 Conflict`.

## Тесты

```sh
go test ./...
```
