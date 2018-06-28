package feature_flag

import (
	"bitbucket.org/Amartha/go-control"
	"bitbucket.org/go-mis/config"
	"gopkg.in/kataras/iris.v4"
	"strconv"
)

var Control go_control.Control

func Init() {
	Control = go_control.Init("MIS", config.FlagServerPath)
}

func GetStatusForFlag(ctx *iris.Context) {
	flagName := ctx.Param("flagName")
	bid := ctx.Param("branchID")

	branchID, err := strconv.ParseUint(bid, 10, 64)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message": "Bad Request",
			"errorMessage": "branchId not integer",
		})
		return
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"enabled": Control.IsEnabledForBranchID(flagName, branchID),
		},
	})
}