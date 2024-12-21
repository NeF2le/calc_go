package logging

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriterWrapper) Write(body []byte) (int, error) {
	rw.body = body
	return rw.ResponseWriter.Write(body)
}

func SetupLogger() *zap.Logger {
	config := zap.NewProductionConfig()

	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		fmt.Println("Error initialization logger", err)
	}

	return logger
}

func LoggingMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapper := &ResponseWriterWrapper{
				ResponseWriter: w,
				statusCode:     http.StatusOK, 
			}

			next.ServeHTTP(wrapper, r)

			duration := time.Since(start)

			var responseMap map[string]interface{}
			if err := json.Unmarshal(wrapper.body, &responseMap); err != nil {
				logger.Error("Ошибка при десериализации ответа", zap.Error(err))
			}

			logger.Info("HTTP запрос",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
				zap.Int("status_code", wrapper.statusCode),
				zap.Any("response", responseMap),
			)
		})
	}
}