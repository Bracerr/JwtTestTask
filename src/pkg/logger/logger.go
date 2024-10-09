package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	Log  *logrus.Logger
	once sync.Once
)

func Init() {
	once.Do(func() {
		Log = logrus.New()

		Log.SetLevel(logrus.InfoLevel)

		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		logFilePath := filepath.Join("..", "..", "logrus_example.log")
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			Log.Fatal("Ошибка при открытии файла для логирования")
		}

		Log.SetOutput(io.MultiWriter(os.Stdout, file))
	})
}
