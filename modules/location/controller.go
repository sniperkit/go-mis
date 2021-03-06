package location

import (
	"fmt"

	"bitbucket.org/go-mis/services"
	"gopkg.in/kataras/iris.v4"
)

// GetLocation - get location
func GetLocation(ctx *iris.Context) {
	query := "SELECT * FROM inf_location "

	if ctx.URLParam("province") != "" && ctx.URLParam("city") != "" && ctx.URLParam("kecamatan") != "" {
		query += "WHERE province = '" + ctx.URLParam("province") + "' AND city = '" + ctx.URLParam("city") + "' AND kecamatan = '" + ctx.URLParam("kecamatan") + "' AND kelurahan != '0' "
	} else if ctx.URLParam("province") != "" && ctx.URLParam("city") != "" {
		query += "WHERE province = '" + ctx.URLParam("province") + "' AND city = '" + ctx.URLParam("city") + "' AND kecamatan != '0' AND kelurahan = '0'"
	} else if ctx.URLParam("province") != "" {
		query += "WHERE province = '" + ctx.URLParam("province") + "' AND city != '0' AND kecamatan = '0' AND kelurahan = '0'"
	} else {
		query += "WHERE province != '0' AND city = '0' AND kecamatan = '0' AND kelurahan = '0'"
	}

	locationSchema := []Location{}
	services.DBCPsql.Raw(query).Scan(&locationSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   locationSchema,
	})
}

// GetLocationById - get location by id
func GetLocationById(ctx *iris.Context) {
	locationCode := ctx.Param("location_code")
	locationSchema := Location{}

	services.DBCPsql.Table("inf_location").Where("\"locationCode\" = ?", locationCode).Scan(&locationSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   locationSchema,
	})
}

func TranslateLoc(ctx *iris.Context) {
	payload := struct {
		LocCodes interface{}
	}{}

	if err := ctx.ReadJSON(&payload); err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err, lt := TranslateLocation(payload.LocCodes)
	if err != nil {
		fmt.Println(err)
	}
	lt.GetLocation()

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   lt,
	})
}
