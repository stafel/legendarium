package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stafel/legendarium/api"
	"github.com/stafel/legendarium/dao"
	"github.com/stafel/legendarium/provider"
)

// loasds all jsons from folder into db
func migrateFolder(folderPath string) {
	prov := &provider.SqliteProvider{}
	prov.Connect()

	dao := &dao.LegendDao{}
	dao.Connect(prov)

	dao.MigrateFromFolder(folderPath)
}

// adds api routes and starts echo server
func startEchoServer() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello")
	})

	api.AddRoutes(e)

	e.Logger.Fatal(e.Start(":9000"))
}

func main() {

	//migrateFolder("/wherever/you/like")

	startEchoServer()
}
