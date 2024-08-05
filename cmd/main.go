package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/thirdknife/scoutingapp/database"

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

const databaseDir = "/tmp/"

func createFakeDatabaseFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// This is good. we don't want the fake database to already exist.
		} else {
			// try to clean up from prior runs.
			err := os.Remove(path)
			if err != nil {
				return fmt.Errorf("failed to clean up pre-existing fake database: %w", err)
			}
		}
	}
	// Create a new database.
	db, err := database.Load(path)
	if err != nil {
		return fmt.Errorf("error creating database: %v", err)
	}
	db.Create(&database.Player{
		Name:  "Foo",
		Score: 0,
	})
	db.Create(&database.Player{
		Name:  "Bar",
		Score: 1,
	})
	if err := database.SaveToFile(db); err != nil {
		return fmt.Errorf("error saving fake scout database: %v", err)
	}
	return nil
}

func main() {
	userHash := "FAKE_SCOUT_HASH"
	dbPath := filepath.Join(databaseDir, userHash+".db")
	if err := createFakeDatabaseFile(dbPath); err != nil {
		fmt.Printf("error creating fake database: %v", err)
		os.Exit(1)
	}

	db, err := database.Load(dbPath)
	if err != nil {
		fmt.Printf("Error loading database: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Loaded database from %s\n", dbPath)

	e := echo.New()

	// Register templates
	templates := make(map[string]*template.Template)
	templates["home"] = template.Must(template.ParseFiles("views/layouts/base.html", "views/pages/home.html"))
	templates["about"] = template.Must(template.ParseFiles("views/layouts/base.html", "views/pages/about.html"))
	templates["signup"] = template.Must(template.ParseFiles("views/layouts/base.html", "views/pages/signup.html"))
	templates["players"] = template.Must(template.ParseFiles("views/layouts/base.html", "views/pages/players.html"))
	templates["dashboard"] = template.Must(template.ParseFiles("views/layouts/base.html", "views/pages/dashboard.html"))

	e.Static("/public", "public")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `${time_custom}: ${method} ${uri} -> status=${status} ${error}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

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
		players, err := database.AllPlayers(db)
		if err != nil {
			return c.HTML(http.StatusInternalServerError, "<p>Error fetching players.</p>")
		}
		return c.Render(http.StatusOK, "players", players)
	})

	e.POST("/players", func(c echo.Context) error {
		name := c.FormValue("name")
		// birthdate := c.FormValue("birthdate")
		// height := c.FormValue("height")
		// weight := c.FormValue("weight")
		// club := c.FormValue("club")
		// position := c.FormValue("position")
		// managerName := c.FormValue("manager_name")
		// telephone := c.FormValue("telephone")

		player := &database.Player{
			Name:  name,
			Score: 2,
		}

		if result := db.Debug().Create(player); result.Error != nil {
			fmt.Printf("pack %v: %v\n", name, err)
			return c.HTML(http.StatusInternalServerError, "<p>Error adding player.</p>")
		}

		players, err := database.AllPlayers(db)
		if err != nil {
			fmt.Println(err)
			return c.HTML(http.StatusInternalServerError, "<p>Error fetching players.</p>")
		}

		return c.Render(http.StatusOK, "players", players)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
