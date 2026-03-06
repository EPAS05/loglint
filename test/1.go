package main

import (
	"log/slog"

	"go.uber.org/zap"
)

func main() {
	// Правильные сообщения (с точки зрения будущих правил)
	slog.Info("starting server")
	slog.Error("failed to connect")
	zap.L().Info("request completed")
	zap.S().Warn("timeout occurred")

	// Неправильные сообщения
	slog.Info("Starting server")       // заглавная буква
	slog.Error("ошибка подключения")   // русский язык
	zap.L().Info("server started!!!")  // лишние восклицательные знаки
	zap.S().Warn("⚠️ warning")         // эмодзи
	slog.Info("user password: secret") // чувствительные данные
	slog.Debug("api_key=" + "12345")   // чувствительные данные в конкатенации
	slog.Info("message" + "." + "rar")
	username := "admin"
	const usern = "admin1"
	slog.Info("message" + ". " + username + " + rar")
	slog.Info("message" + ". " + usern + " + rar")

	// Вызов через переменную-логгер
	logger := slog.Default()
	logger.Info("all good")
}
