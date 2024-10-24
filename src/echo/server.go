package echo

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/UnknownRori/lagra_server/src"
	"github.com/charmbracelet/log"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/golang-jwt/jwt/v5"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	App            *echo.Echo
	Db             src.DB
	StorageService *StorageService
	ConfigJwt      echoJwt.Config
}

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func CreateServer() Server {
	db, err := src.CreateConn()

	if err != nil {
		panic("Database Connection failed : \n" + err.Error())
	}

	storage := NewStorageService("storage")

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Debug = os.Getenv("APP_DEBUG") == "true"
	e.Use(LoggerMiddleware)
	e.Use(middleware.Recover())
	validate = validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)
	e.Validator = &CustomValidator{validator: validate}

	configJwt :=
		echoJwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(JwtClaims)
			},
			SigningKey: []byte(os.Getenv("ACCESS_TOKEN_KEY")),
		}

	s := Server{
		App:            e,
		Db:             db,
		ConfigJwt:      configJwt,
		StorageService: storage,
	}

	RegisterRouter(&s)

	return s
}

func (s *Server) Start(port string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		log.Infof("Starting server at : http://localhost:%s", port)
		if err := s.App.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatal("Shutting down server : Cannot start server")
			s.App.Logger.Fatal("shutting down the server : cannot start server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.App.Shutdown(ctx); err != nil {
		log.Fatal(err)
		s.App.Logger.Fatal(err)
	}

	if err := s.Db.Close(); err != nil {
		log.Fatal(err)
	}
	log.Info("Server shutdown completed...")
}
