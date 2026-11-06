# miniTODO

REST API для списка книг. Книги хранятся в памяти и исчезают при перезапуске сервера.

## Быстрый запуск

1. Установите Go **1.25.7 или новее**.
2. Откройте терминал в папке с файлом `go.mod` и запустите:

   ```powershell
   go run .
   ```

3. Откройте в браузере [http://localhost:9091/books](http://localhost:9091/books). При первом запуске список книг пуст.

Чтобы остановить сервер, нажмите **Ctrl+C** в терминале.

## Добавить первую книгу

Выполните в другом окне PowerShell, пока сервер работает:

```powershell
Invoke-RestMethod http://localhost:9091/books -Method Post -ContentType 'application/json' -Body '{"Title":"1984","Author":"George Orwell","NumOfPages":328}'
```

Обновите страницу `/books`, чтобы увидеть книгу в списке.

## Остальные запросы

| Запрос | Что делает |
| --- | --- |
| `GET /books?title=1984` | Находит книгу по названию |
| `GET /books?author=George%20Orwell` | Показывает книги автора |
| `GET /books?completed=true` | Показывает прочитанные книги (`false` — непрочитанные) |
| `PATCH /books/1984` | Меняет статус чтения. Тело запроса — `{"completed":true}` или `{"completed":false}` |
| `DELETE /books/1984` | Удаляет книгу по названию |

Тесты можно запустить из той же папки командой `go test ./...`.
