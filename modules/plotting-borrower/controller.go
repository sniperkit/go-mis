package plottingBorrower

import (
	"errors"
	"strconv"
	"strings"

	investor "bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/services"

	"gopkg.in/kataras/iris.v4"
)

func SavePlottingParams(ctx *iris.Context) {
	// convert requestbody to string
	pp := string(ctx.Request.Body())
	if pp == "" {
		err := errors.New("No Plotting Params were found in the request body.")
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// get investor id
	var invId uint64
	s := pp[1 : len(pp)-1]

	ppArr := []string{}

	if strings.Contains(s, ",") {
		ppArr = strings.Split(s, ",")
	} else {
		ppArr = append(ppArr, s)
	}

	for _, val := range ppArr {
		if strings.Contains(val, "investorId") {
			v := strings.Split(val, ":")
			if strings.Contains(v[1], "\"") || strings.Contains(v[1], "'") {
				v[1] = v[1][1 : len(v[1])-1]
			}
			v[1] = strings.TrimSpace(v[1])
			v[1] = strings.TrimSuffix(v[1], "\n")
			v[1] = strings.TrimSuffix(v[1], "\"")
			id, err := strconv.Atoi(v[1])
			if err != nil {
				err := errors.New("Error converting investorId.")
				ctx.JSON(iris.StatusInternalServerError, iris.Map{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}
			invId = uint64(id)
			break
		}
	}

	// set plottingParrams as borrowerCrtiteria on investor
	inv := &investor.Investor{ID: invId}
	services.DBCPsql.Model(&inv).Update("borrowerCriteria", pp)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
	})

}
