package testpkg

import (
	"log/slog"
)

func main() {
	slog.Info("starting server")
	slog.Info("request completed 123")

	slog.Info("Starting server")       // want "log message should start with a lowercase letter"
	slog.Info("запуск сервера")        // want "log message should contain only English characters"
	slog.Info("server started!!!")     // want "log message should not contain special characters or emojis"
	slog.Info("⚠️ warning")            // want "log message should contain only English characters" "log message should not contain special characters or emojis"
	slog.Info("user password: secret") // want `log message may contain sensitive data \(word: password\)` `log message should not contain special characters or emojis`
}
