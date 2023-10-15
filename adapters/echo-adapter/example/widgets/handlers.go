package widgets

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"strconv"
)

func ok(c echo.Context) error {
	return c.JSON(200, c.Request().URL.Query())
}

func okPath(c echo.Context) error {
	// crud has already validated that ID is a number, no need to check error
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(200, id)
}

func fakeAuthPreHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authentication") != "password" {
			return echo.NewHTTPError(401, "Authentication header must be 'password'")
		}
		return next(c)
	}
}

func bindAndOk(c echo.Context) error {
	var widget struct {
		json.RawMessage
	}
	if err := c.Bind(&widget); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}
	return c.JSON(200, widget)
}
