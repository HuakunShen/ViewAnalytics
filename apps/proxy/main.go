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
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir()) || strings.Contains(os.Args[0], "JetBrains")
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		// Dir:         migrationDir,
		Automigrate: isGoRun,
	})

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		proxyRoute := "/api/proxy"
		e.Router.GET(proxyRoute+"/:url", func(c echo.Context) error {
			// path := c.PathParam("path")
			fmt.Println(c.Request().URL.Path)
			// this url path starts with /api/proxy/, now we need to remove /api/proxy/ and get the rest of the path
			url := c.Request().URL.Path[len(proxyRoute)+1:]
			fmt.Println(url)
			// redirect to the url
			return c.JSON(http.StatusOK, map[string]interface{}{url: url})
			return c.Redirect(http.StatusTemporaryRedirect, url)
		})
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
