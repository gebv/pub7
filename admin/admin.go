package admin

import (
	"html/template"
	"io"

	"github.com/tarantool/go-tarantool"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger
var TarantoolConnection *tarantool.Connection

func SetLogger(l *zap.SugaredLogger) {
	log = l
}

func Setup(e *echo.Echo, pathTpls string) error {
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("10M"))
	e.Use(middleware.Secure())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Renderer = &Template{
		// "public/views/*.html"
		templates: template.Must(template.ParseGlob(pathTpls)),
	}

	e.GET("/dashbaord", DashboardHandler)
	e.GET("/api/v1/chats/states", ListChatsHandler)
	return nil
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
