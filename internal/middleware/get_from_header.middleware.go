package middlewares

import (
	"bmt_order_service/global"
	"bmt_order_service/internal/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetFromHeaderMiddleware struct {
}

func NewGetFromHeaderMiddleware() *GetFromHeaderMiddleware {
	return &GetFromHeaderMiddleware{}
}

func (g *GetFromHeaderMiddleware) GetEmailFromHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetHeader(global.X_USER_EMAIL)
		if email == "" {
			responses.FailureResponse(c, http.StatusBadRequest, "email is not empty")
			c.Abort()
			return
		}

		c.Set(global.X_USER_EMAIL, email)
		c.Next()
	}
}
