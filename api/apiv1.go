package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/stafel/legendarium/dao"
)

func v1GetCharacters(c echo.Context) error {

	d := dao.DefaultConnect()

	return c.JSON(http.StatusOK, d.GetCharacters())
}

func v1GetCharacterById(c echo.Context) error {
	id_str := c.Param("cid")

	id_int, err := strconv.Atoi(id_str)
	if err != nil {
		log.Fatal(err)
	}

	id := uint(id_int)

	d := dao.DefaultConnect()

	return c.JSON(http.StatusOK, d.GetCharacterById(id))
}

func AddRoutes(e *echo.Echo) {
	e.GET("/v1/characters", v1GetCharacters)
	e.GET("/v1/characters/:cid", v1GetCharacterById)
	//e.GET("/v1/characters/:cid/milestones", nil)
	//e.GET("/v1/characters/:cid/milestones/latest", nil)
	//e.GET("/v1/characters/:cid/milestones/:mid", nil)
}
