package reports

import (
	"fmt"

	"bitbucket.org/go-mis/config"

	"github.com/parnurzeal/gorequest"
	iris "gopkg.in/kataras/iris.v4"
)

// AgentRekap - AgentRekap bridge
func AgentRekap(ctx *iris.Context) {
	agentID := ctx.URLParam("agentId")
	date := ctx.URLParam("date")
	//fmt.Print(agentID)
	var resReport ResReport

	urlStr := config.UploaderApiPath + "report/agent?agentId=" + agentID + "&installmentDate=" + date + "&secretKey=n0de-U>lo4d3r"
	request := gorequest.New()
	_, _, errs := request.Get(urlStr).
		EndStruct(&resReport)

	if len(errs) > 0 {
		fmt.Print(errs)
	}

	if (resReport == ResReport{}) {
		resReport.Status = "error"
		resReport.Message = "error bridge to uploader"
	}

	if resReport.Path != "" {
		resReport.Path = config.UploaderApiPath + resReport.Path

	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   resReport,
	})

}
