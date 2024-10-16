package echo

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/UnknownRori/lagra_server/src"
	"github.com/charmbracelet/log"

	"github.com/golang-jwt/jwt/v5"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	App       *echo.Echo
	Db        src.DB
	ConfigJwt echoJwt.Config
}

func CreateServer() Server {
	db, err := src.CreateConn()

	if err != nil {
		panic("Database Connection failed : \n" + err.Error())
	}

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Debug = os.Getenv("APP_DEBUG") == "true"
	e.Use(LoggerMiddleware)
	e.Use(middleware.Recover())

	configJwt :=
		echoJwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(JwtClaims)
			},
			SigningKey: []byte(os.Getenv("ACCESS_TOKEN_KEY")),
		}

	s := Server{
		App:       e,
		Db:        db,
		ConfigJwt: configJwt,
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
