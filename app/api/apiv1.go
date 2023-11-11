package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/stafel/legendarium/dao"
)

// convert string gained from context param to uint
func getIdFromString(idStr string) (uint, error) {
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	if idInt < 0 {
		return 0, errors.New("id can not be negative")
	}

	return uint(idInt), nil
}

func v1GetCharacters(c echo.Context) error {

	d := dao.DefaultConnect()

	return c.JSON(http.StatusOK, d.GetCharacters())
}

func v1GetCharacterById(c echo.Context) error {
	id, err := getIdFromString(c.Param("cid"))
	if err != nil {
		return err
	}

	d := dao.DefaultConnect()

	char, err := d.GetCharacterById(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, char)
}

func v1GetMilestones(c echo.Context) error {
	id, err := getIdFromString(c.Param("cid"))
	if err != nil {
		return err
	}

	d := dao.DefaultConnect()
	ms, err := d.GetMilestonesForCharacterId(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ms)
}

func v1GetLatestMilestone(c echo.Context) error {
	id, err := getIdFromString(c.Param("cid"))
	if err != nil {
		return err
	}

	d := dao.DefaultConnect()
	ms, err := d.GetLatestMilestonesForCharacterId(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ms)
}

func v1GetMilestone(c echo.Context) error {
	cid, err := getIdFromString(c.Param("cid"))
	if err != nil {
		return err
	}

	mid, err := getIdFromString(c.Param("mid"))
	if err != nil {
		return err
	}

	d := dao.DefaultConnect()
	ms, err := d.GetMilestoneForMilestoneId(mid)
	if err != nil {
		return err
	}

	if ms.LegendCharacterID != cid { // prevent accessing a milestone from another character
		return errors.New("no milestone found")
	}

	return c.JSON(http.StatusOK, ms)
}

func AddRoutes(e *echo.Echo) {
	e.GET("/v1/characters", v1GetCharacters)
	e.GET("/v1/characters/:cid", v1GetCharacterById)
	e.GET("/v1/characters/:cid/milestones", v1GetMilestones)
	e.GET("/v1/characters/:cid/milestones/latest", v1GetLatestMilestone)
	e.GET("/v1/characters/:cid/milestones/:mid", v1GetMilestone)
}
