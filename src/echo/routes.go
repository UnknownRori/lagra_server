package echo

import (
	"net/http"

	"github.com/UnknownRori/lagra_server/src/models"

	"github.com/charmbracelet/log"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRouter(server *Server) {
	api := server.App.Group("/api")
	api.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}),
	)

	v1 := api.Group("/v1")

	v1.GET("/ping", func(c echo.Context) error {
		return c.String(200, "Pong!")
	})

	v1.POST("/users", func(c echo.Context) error {
		var user models.NewUser
		err := c.Bind(&user)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Username, Password, and displayName should be filled!",
				"status":  "fail",
			})
		}

		uuid, err := models.CreateUser(&server.Db, user)

		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Username already exist or invalid payload",
				"status":  "fail",
			})
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"data": echo.Map{
				"user": &models.ReturnUser{
					Uuid:        uuid.String(),
					Username:    user.Username,
					DisplayName: user.DisplayName,
					Role:        user.Role,
				},
			},
			"status": "success",
		})
	})

	v1.POST("/auth/login", func(c echo.Context) error {
		var user models.LoginUser
		err := c.Bind(&user)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Username, Password, and displayName should be filled!",
				"status":  "fail",
			})
		}

		fetchUser, err := models.FetchUserByUsername(&server.Db, user.Username)

		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "Data not found!",
				"status":  "fail",
			})
		}

		token, err := CreateClaim(fetchUser.Uuid)

		return c.JSON(http.StatusOK, echo.Map{
			"data": echo.Map{
				"token": token,
			},
			"status": "success",
		})
	})

	authenticatedOnly := v1.Group("/")
	authenticatedOnly.Use(echoJwt.WithConfig(server.ConfigJwt))
	authenticatedOnly.Use(AuthenticatedContextMiddleware(server))
	authenticatedOnly.GET("auth/me", func(c echo.Context) error {
		cc := c.(*AuthenticatedContext)
		user := cc.User

		return c.JSON(http.StatusOK, echo.Map{
			"data": echo.Map{
				"user": user,
			},
			"status": "success",
		})
	})
}
