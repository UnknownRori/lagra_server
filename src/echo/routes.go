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
				"message": "Username, Password, and email should be filled!",
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
					Uuid:     uuid.String(),
					Username: user.Username,
					Role:     user.Role,
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

	admin := authenticatedOnly.Group("")
	admin.Use(AdminContextMiddleware(server))
	admin.POST("categories", func(c echo.Context) error {
		var category models.NewCategory
		err := c.Bind(&category)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Required name",
				"status":  "fail",
			})
		}

		uuid, err := models.CreateCategory(&server.Db, category)

		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Category already exist or invalid payload",
				"status":  "fail",
			})
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"data": echo.Map{
				"category": &models.Category{
					Uuid: uuid.String(),
					Name: category.Name,
				},
			},
			"status": "success",
		})
	})

	v1.GET("/categories", func(c echo.Context) error {
		categories, err := models.FetchCategories(&server.Db)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Database offline",
				"status":  "fail",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"data": echo.Map{
				"categories": categories,
			},
			"status": "success",
		})
	})

	admin.POST("items", func(c echo.Context) error {
		var item models.NewItem
		err := c.Bind(&item)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Required name, price, category",
				"status":  "fail",
			})
		}

		uuid, err := models.CreateItem(&server.Db, item)

		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Items already exist or invalid payload",
				"status":  "fail",
			})
		}

		uuidString := uuid.String()
		log.Info(uuidString)
		itemModel, err := models.FetchItemByUuid(&server.Db, uuidString)
		log.Info(err)

		return c.JSON(http.StatusCreated, echo.Map{
			"data": echo.Map{
				"item": itemModel,
			},
			"status": "success",
		})
	})

	v1.GET("/items", func(c echo.Context) error {
		items, err := models.FetchItems(&server.Db)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Database offline",
				"status":  "fail",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"data": echo.Map{
				"items": items,
			},
			"status": "success",
		})
	})
}
