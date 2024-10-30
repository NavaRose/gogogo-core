package middleware

import (
	"errors"
	"github.com/NavaRose/gogogo-core/exception"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		for _, err := range ctx.Errors {
			var e exception.Http
			switch {
			case errors.As(err.Err, &e):
				ctx.AbortWithStatusJSON(e.StatusCode, e)
				break
			default:
				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					map[string]string{"message": err.Error()},
				)
			}

			return
		}
	}
}
