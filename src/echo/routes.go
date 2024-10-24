package echo

import (
	"fmt"
	"net/http"
	"strings"

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

	server.App.GET("/storage/:filename", func(c echo.Context) error {
		filename := c.Param("filename")
		filePath := server.StorageService.path + "/" + filename
		strings.ReplaceAll(filePath, "..", "")

		return c.File(filePath)
	})

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

		file, err := c.FormFile("img")
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Must be a file",
				"status":  "fail",
			})
		}

		filename, err := server.StorageService.store(file)
		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "File failed to save",
				"status":  "fail",
			})
		}

		item.ImgUrl = filename
		fmt.Println(item.CategoryId)

		uuid, err := models.CreateItem(&server.Db, item)

		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Items already exist or invalid payload",
				"status":  "fail",
			})
		}

		uuidString := uuid.String()
		itemModel, err := models.FetchItemByUuid(&server.Db, uuidString)

		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Database error",
				"status":  "fail",
			})
		}

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

	authenticatedOnly.POST("carts", func(c echo.Context) error {
		cc := c.(*AuthenticatedContext)
		user := cc.User
		var cart models.NewCart
		err := c.Bind(&cart)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Required name, price, category",
				"status":  "fail",
			})
		}

		uuid, err := models.CreateCart(&server.Db, cart, user)

		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Items already exist or invalid payload",
				"status":  "fail",
			})
		}

		uuidString := uuid.String()
		carts, err := models.FetchCartsByUuid(&server.Db, uuidString, user)

		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Database error",
				"status":  "fail",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"data": echo.Map{
				"carts": carts,
			},
			"status": "success",
		})
	})

	authenticatedOnly.GET("carts", func(c echo.Context) error {
		cc := c.(*AuthenticatedContext)
		user := cc.User
		carts, err := models.FetchCarts(&server.Db, user)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Database offline",
				"status":  "fail",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"data": echo.Map{
				"carts": carts,
			},
			"status": "success",
		})
	})

	authenticatedOnly.GET("transactions", func(c echo.Context) error {
		cc := c.(*AuthenticatedContext)
		user := cc.User
		carts, err := models.FetchTransactions(&server.Db, user)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Database offline",
				"status":  "fail",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"data": echo.Map{
				"transactions": carts,
			},
			"status": "success",
		})
	})

	authenticatedOnly.POST("transactions", func(c echo.Context) error {
		cc := c.(*AuthenticatedContext)
		user := cc.User
		var newTrans models.NewTransaction
		err := c.Bind(&newTrans)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Invalid request",
				"status":  "fail",
			})
		}

		carts, err := models.FetchCarts(&server.Db, user)
		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Database error",
				"status":  "fail",
			})
		} else if carts == nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Empty carts",
				"status":  "fail",
			})
		}
		err = models.CleanCarts(&server.Db, user)
		if err != nil {
			log.Error(err.Error())
		}
		uuid, err := models.CreateTransaction(&server.Db, newTrans, carts, user)

		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Items already exist or invalid payload",
				"status":  "fail",
			})
		}

		uuidString := uuid.String()
		transaction, err := models.FetchTransactionByUuid(&server.Db, uuidString, user)

		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Database error",
				"status":  "fail",
			})
		}

		err = models.CreateTransactionItemsFromCarts(&server.Db, transaction, carts)
		if err != nil {
			log.Error(err.Error())

			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Invalid state",
				"status":  "fail",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"data": echo.Map{
				"transaction": transaction,
			},
			"status": "success",
		})
	})

	authenticatedOnly.GET("transactions/:uuid", func(c echo.Context) error {
		transactionUuid := c.Param("uuid")
		cc := c.(*AuthenticatedContext)
		user := cc.User

		transaction, err := models.FetchDetailTransactions(&server.Db, transactionUuid, user)

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Invalid request",
				"status":  "fail",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"data": echo.Map{
				"transaction": transaction,
			},
			"status": "success",
		})
	})
}
