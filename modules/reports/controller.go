package reports

import (
	"fmt"

	"bitbucket.org/go-mis/config"

	"github.com/kataras/iris"
	"github.com/parnurzeal/gorequest"
)

// AgentRekap - AgentRekap bridge
func AgentRekap(ctx *iris.Context) {
	agentId := ctx.GetString("agentId")
	date := ctx.GetString("date")
	var resReport ResReport

	urlStr := config.UploaderApiPath + "report?agentId=" + agentId + "&date=" + date + "&secretKey=n0de-U>lo4d3r"
	request := gorequest.New()
	_, _, errs := request.Get(urlStr).
		EndStruct(&resReport)

	if len(errs) > 0 {
		fmt.Print(errs)
	}

	ctx.JSON(iris.StatusOK, resReport)

}
