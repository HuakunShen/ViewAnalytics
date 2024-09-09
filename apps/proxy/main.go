package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/HuakunShen/view-stats-proxy/apps/proxy/migrations"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()
	isDevMode := strings.HasPrefix(os.Args[0], os.TempDir()) || strings.Contains(os.Args[0], "JetBrains") || strings.Contains(os.Args[0], "debug") || strings.Contains(os.Args[0], "pb_dev")
	fmt.Println("isDevMode", isDevMode)
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isDevMode check is to enable it only during development)
		// Dir:         migrationDir,
		Automigrate: isDevMode,
	})

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		proxyRoute := "/api/proxy"
		e.Router.GET(proxyRoute+"/:url", func(c echo.Context) error {
			// this url path starts with /api/proxy/, now we need to remove /api/proxy/ and get the rest of the path
			url := c.Request().URL.Path[len(proxyRoute)+1:]
			// log to database
			collection, err := app.Dao().FindCollectionByNameOrId("proxy_records")
			if err != nil {
				return err
			}
			record := models.NewRecord(collection)
			record.Set("url", url)
			record.Set("ip", c.RealIP())
			if err := app.Dao().SaveRecord(record); err != nil {
				return err
			}

			// redirect to the url
			return c.Redirect(http.StatusTemporaryRedirect, url)
		})
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
