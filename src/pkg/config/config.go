package config

import (
	db "JwtTestTask/src/pkg/database"
	"JwtTestTask/src/pkg/logger"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ServerParams struct {
	ServerHost string
}

type JwtParams struct {
	SigningKey      string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

type SmtParams struct {
	Host     string
	Port     string
	Username string
	Password string
}

func Init() {
	rootDir, err := os.Getwd()
	if err != nil {
		logger.Log.Fatal("Ошибка при получении текущей директории: %v", err)
	}

	envFilePath := filepath.Join(filepath.Dir(filepath.Dir(rootDir)), ".env")
	envErr := godotenv.Load(envFilePath)
	if envErr != nil {
		logger.Log.Fatal("Ошибка при загрузке .env файла: %v", envErr)
	}
}

func GetDbParams() db.DbInitModel {
	dbInitModel := db.DbInitModel{
		DbHost:     os.Getenv("DB_HOST"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
		DbPort:     os.Getenv("DB_PORT"),
	}

	if dbInitModel.DbHost == "" || dbInitModel.DbUser == "" || dbInitModel.DbPassword == "" || dbInitModel.DbName == "" || dbInitModel.DbPort == "" {
		logger.Log.Fatal("Ошибка: не все параметры базы данных были получены. Проверьте .env файл.")
	}

	return dbInitModel
}

func GetServerParams() ServerParams {
	serverHost := os.Getenv("SERVER_PORT")
	if serverHost == "" {
		logger.Log.Fatal("Ошибка: параметр SERVER_PORT не был получен. Проверьте .env файл.")
	}

	return ServerParams{ServerHost: serverHost}
}

func GetJwtParams() JwtParams {
	AccessDurationStr := os.Getenv("JWT_ACCESS_DURATION")
	RefreshDurationStr := os.Getenv("JWT_REFRESH_DURATION")
	AccessDurationInt, err := strconv.Atoi(AccessDurationStr)
	if err != nil {
		AccessDurationInt = 15
		logger.Log.Printf("Ошибка при преобразовании JWT_ACCESS_DURATION: %v. Используется значение по умолчанию: %d минут.", err, AccessDurationInt)
	}

	RefreshDurationInt, err := strconv.Atoi(RefreshDurationStr)
	if err != nil {
		RefreshDurationInt = 30
		logger.Log.Printf("Ошибка при преобразовании JWT_REFRESH_DURATIOn: %v. Используется значение по умолчанию: %d дней.", err, AccessDurationInt)
	}

	AccessDuration := time.Duration(AccessDurationInt) * time.Minute
	RefreshDuration := time.Duration(RefreshDurationInt) * 24 * time.Hour
	signingKey := os.Getenv("JWT_SIGNING_KEY")
	if signingKey == "" {
		logger.Log.Fatal("Ошибка: параметр JWT_SIGNING_KEY не был получен. Проверьте .env файл.")
	}

	return JwtParams{SigningKey: signingKey, AccessDuration: AccessDuration, RefreshDuration: RefreshDuration}
}

func GetSmtpParams() SmtParams {
	smtParams := SmtParams{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	if smtParams.Host == "" || smtParams.Port == "" || smtParams.Username == "" || smtParams.Password == "" {
		logger.Log.Fatal("Ошибка: не все параметры SMTP были получены. Проверьте .env файл.")
	}

	return smtParams
}
