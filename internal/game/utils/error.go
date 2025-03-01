package utils

import "github.com/labstack/echo/v4"

type APIError struct {
	Error   string `json:"error"`
	Details string `json:"details, omitempty"`
}

func ErrorResponse(c echo.Context, statusCode int, message string, details ...string) error {
	errorResponse := APIError{
		Error: message,
	}

	if len(details) > 0 {
		errorResponse.Details = details[0]
	}

	return c.JSON(statusCode, errorResponse)
}
