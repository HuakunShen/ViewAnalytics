package main

import (
	"fmt"
	"io"
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

func getOriginalUrl(c echo.Context, route string) string {
	url := c.Request().URL.Path[len(route)+1:]
	rawQuery := c.Request().URL.RawQuery
	if rawQuery != "" {
		url += "?" + rawQuery
	}
	return url
}

func logToDatabaseProxyRecord(app *pocketbase.PocketBase, url string, referer string, ip string) error {
	collection, err := app.Dao().FindCollectionByNameOrId("proxy_records")
	if err != nil {
		return err
	}
	record := models.NewRecord(collection)
	record.Set("url", url)
	record.Set("ip", ip)
	record.Set("referer", referer)
	if err := app.Dao().SaveRecord(record); err != nil {
		return err
	}
	return nil
}

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
			url := getOriginalUrl(c, proxyRoute)
			// log to database
			logToDatabaseProxyRecord(app, url, c.Request().Header.Get("Referer"), c.RealIP())
			// redirect to the url
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Copy headers from the original response
			for key, values := range resp.Header {
				for _, value := range values {
					c.Response().Header().Add(key, value)
				}
			}

			// Set the status code
			c.Response().WriteHeader(resp.StatusCode)
			// Stream the response body
			_, err = io.Copy(c.Response(), resp.Body)
			return err
		})

		e.Router.GET("/api/health", func(c echo.Context) error {
			return c.String(http.StatusOK, "OK")
		})

		redirectRoute := "/api/redirect"
		e.Router.GET(redirectRoute+"/:url", func(c echo.Context) error {
			url := getOriginalUrl(c, redirectRoute)
			fmt.Println("url", url)
			// log to database
			logToDatabaseProxyRecord(app, url, c.Request().Header.Get("Referer"), c.RealIP())
			// redirect to the url
			return c.Redirect(http.StatusTemporaryRedirect, url)
		})
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
