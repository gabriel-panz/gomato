package handler

import (
	"net/http"

	view "github.com/gabriel-panz/gomato/view"
	"github.com/labstack/echo/v4"
)

type TimerForm struct {
	Time int `json:"time" form:"time" query:"time"`
}

func HandleShowHome(c echo.Context) error {
	return view.Hello().Render(c.Request().Context(), c.Response())
}

func HandleShowFocus(c echo.Context) error {
	return view.Focus().Render(c.Request().Context(), c.Response())
}

func HandleShowPause(c echo.Context) error {
	form := TimerForm{}
	err := (&echo.DefaultBinder{}).BindBody(c, &form)

	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	t := form.Time / 5
	return view.Pause(t).Render(c.Request().Context(), c.Response())
}
