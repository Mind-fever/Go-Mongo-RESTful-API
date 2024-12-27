package middlewares

import (
	"errors"
	errors2 "github.com/Programacion-2-Trabajo-Practico-Integrador-Asselborn-Martinez/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last().Err
		var customErr *errors2.CustomError
		if ok := errors.As(err, &customErr); ok {
			c.JSON(httpStatusFromCode(customErr.Code), gin.H{"error": customErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}

func httpStatusFromCode(code string) int {
	switch code {
	case "ERR_6", "ERR_17", "ERR_19", "ERR_21", "ERR_22":
		// 404 Not Found
		// Errores relacionados con recursos no encontrados
		return http.StatusNotFound

	case "ERR_1", "ERR_8", "ERR_9", "ERR_10", "ERR_11", "ERR_12", "ERR_13", "ERR_14", "ERR_15", "ERR_16", "ERR_18", "ERR_20":
		// 400 Bad Request
		// Errores relacionados con validación o conflictos en los datos
		return http.StatusBadRequest

	case "ERR_23":
		// 401 Unauthorized
		// Error de autenticación
		return http.StatusUnauthorized

	case "ERR_7":
		// 403 Forbidden
		// Error de acceso prohibido
		return http.StatusForbidden

	case "ERR_2", "ERR_3", "ERR_4", "ERR_5", "ERR_24", "ERR_25", "ERR_26":
		// 500 Internal Server Error
		// Errores relacionados con fallas internas del sistema o la base de datos
		return http.StatusInternalServerError

	default:
		// 500 Internal Server Error por defecto
		return http.StatusInternalServerError
	}
}
