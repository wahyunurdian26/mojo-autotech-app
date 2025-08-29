package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"mojo-autotech/model"
	"mojo-autotech/utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, model.Response{
				Code: http.StatusUnauthorized,
				Msg:  "Unauthorized",
				Err:  "missing bearer token",
			})
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		secret := os.Getenv("AUTH_JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusUnauthorized, model.Response{
				Code: http.StatusUnauthorized,
				Msg:  "Unauthorized",
				Err:  "AUTH_JWT_SECRET not set",
			})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenStr, &utils.AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, model.Response{
				Code: http.StatusUnauthorized,
				Msg:  "Unauthorized",
				Err:  "invalid token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*utils.AccessClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, model.Response{
				Code: http.StatusUnauthorized,
				Msg:  "Unauthorized",
				Err:  "invalid claims",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
