package main

import (
	"github.com/thirdknife/scoutingapp/database"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRegistry struct {
	templates map[string]*template.Template
	Reload    bool
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates[name].Execute(w, data)
}

func main() {

	e := echo.New()

	// Register templates
	templates := make(map[string]*template.Template)
	templates["home"] = template.Must(template.ParseFiles("views/layouts/base.html", "views/pages/home.html"))
	templates["about"] = template.Must(template.ParseFiles("views/layouts/base.html", "views/pages/about.html"))
	templates["signup"] = template.Must(template.ParseFiles("views/layouts/base.html", "views/pages/signup.html"))
	templates["players"] = template.Must(template.ParseFiles("views/layouts/base.html", "views/pages/players.html"))
	templates["dashboard"] = template.Must(template.ParseFiles("views/layouts/base.html", "views/pages/dashboard.html"))

	e.Static("/public", "public")
	e.Use(middleware.Logger())

	t := &TemplateRegistry{
		templates: templates,
		Reload:    false, // Enable template caching
	}

	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home", nil)
	})

	e.GET("/about", func(c echo.Context) error {
		return c.Render(http.StatusOK, "about", nil)
	})

	e.GET("/signup", func(c echo.Context) error {
		return c.Render(http.StatusOK, "signup", nil)
	})

	e.GET("/dashboard", func(c echo.Context) error {
		return c.Render(http.StatusOK, "dashboard", nil)
	})

	e.GET("/players", func(c echo.Context) error {
		type PlayersWrapper struct {
			Players []*database.Player
		}
		fakePlayers := []*database.Player{
			{
				Name:  "Foo",
				Score: 0,
			},
			{
				Name:  "Bar",
				Score: 1,
			},
		}
		return c.Render(http.StatusOK, "players", PlayersWrapper{fakePlayers})
	})

	e.Logger.Fatal(e.Start(":42069"))

}
