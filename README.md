# tages-grpc Тестовое задание

Хранилище сервера находится по адресу: ./server/files. Все загружаемые на сервер файлы помещаются туда.
Хранилище клиента находится по адресу: ./client/files. Файлы, которые будут отправляться на сервер, должны находиться там. Так же файлы скачанные с сервера будут сохраняться туда.

## Запуск сервера

```bash
go run .\server\server.go
```
## Запуск клиента

```bash
go run .\client\client.go
```
После запуска встречает меню с выбором нужного сервиса. Выбор цифрами!
```bash
Menu: 
  1 » Upload file to server
  2 » Download file from server
  3 » File list
  4 » Exit
Choose your option »
```
Введите имя файла с расширением
```bash
Enter file name » test.txt
```
