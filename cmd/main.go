package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/thirdknife/scoutingapp/database"
	base "github.com/thirdknife/scoutingapp/views"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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

func RenderComponent(c echo.Context, status int, cmp templ.Component) error {
	c.Response().WriteHeader(status)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
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

	e.Static("/public", "public")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `${time_custom}: ${method} ${uri} -> status=${status} ${error}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	e.GET("/players", func(c echo.Context) error {

		players, err := database.AllPlayers(db)
		if err != nil {
			return c.HTML(http.StatusInternalServerError, "<p>Error fetching players.</p>")
		}

		return RenderComponent(c, http.StatusOK, base.ListPlayers(players))
	})

	e.GET("/", func(c echo.Context) error {
		return RenderComponent(c, http.StatusOK, base.Home())
	})

	e.Logger.Fatal(e.Start(":42069"))
}
