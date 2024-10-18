package echo

import (
	"net/http"
	"github.com/UnknownRori/lagra_server/src/models"

	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthenticatedContext struct {
	echo.Context
	models.User
}

func AuthenticatedContextMiddleware(s *Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*JwtClaims)
			uuid := claims.Uuid

			db := &s.Db
			userModel, err := models.FetchUserByUuid(db, uuid)
			if err != nil {
				log.Fatal(err)
				c.Error(err)
			}

			return next(&AuthenticatedContext{
				Context: c,
				User:    userModel,
			})
		}
	}
}

func AdminContextMiddleware(s *Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(*AuthenticatedContext)
			user := cc.User

			if user.Role != "ADMIN" {
				log.Fatal("Invalid access")
				return c.JSON(http.StatusNotFound, echo.Map{
					"message": "Data not found!",
					"status":  "fail",
				})
			}

			return next(cc)
		}
	}
}
