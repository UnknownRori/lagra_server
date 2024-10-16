package echo

import (
	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := c.Request()
		response := c.Response()

		// Proceed to the next handler
		if err := next(c); err != nil {
			c.Error(err)
		}

		log.Infof("[%d] %s - %s : %s", response.Status, request.Method, request.RequestURI, request.RemoteAddr)

		return nil
	}
}
