package logger

import (
	"crypto/rand"
	"encoding/hex"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(requestID string) *zap.Logger {

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := config.Build()

	// Add requestID as a field to the logger
	logger = logger.With(zap.String("requestID", requestID))

	return logger

}

func GenerateRequestID() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
