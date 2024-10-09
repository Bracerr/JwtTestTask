# JwtTokens REST API (Go & Echo Framework)

## Технологии
- ### **Go**: 1.23.1

## P.S
### При операции обновления токенов их передача вынесена в body, а не в headers, как положено. Это сделано для удобства проверки работоспособности приложения (описано в Swagger)  
### Access token шифруется с помощью алгоритма HS512(SHA512 + ключ для подписи)

## Запуск приложения
### Docker
```bash
docker-compose up --build
```

### Local
- ### Создайте в основной директории файл .env на основе файла .envExample (Example уже заполнен для подключения к docker-compose DB)
```bash
go mod tidy
```
```bash
docker-compose up -d db
```
```bash
cd src/cmd
```
```bash
go run main.go
```

# Docs(Swagger)
- ### http://localhost:8080/swagger/index.html#/



