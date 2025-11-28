package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Logger returns the fiber logger middleware
func Logger() any {
	return logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	})
}