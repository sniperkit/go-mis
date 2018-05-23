package autoRecon

import (
	"fmt"

	"bitbucket.org/go-mis/config"
	"bytes"
	"encoding/json"
	"gopkg.in/kataras/iris.v4"
	"net/http"
)

func DisbursementDataTransferSave(ctx *iris.Context) {

	payload := struct {
		ID               uint64
		TransferDate     string
		DisburseDateFrom string
		DisburseDateTo   string
		Amount           float64
		ReferenceCode    string
		SettlementId     uint64
		BranchId         uint64
	}{}

	err := ctx.ReadJSON(&payload)
	if err != nil {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"message":      "Bad Request",
			"errorMessage": err.Error(),
		})
		return
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(payload)

	var url string = config.GoFinAutoReconPath + "/api/v1/data-transfer/disbursement/save"
	req, err := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    iris.Map{},
		})
		return
	}
	defer resp.Body.Close()

	var body struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Success bool        `json:"success"`
	}

	json.NewDecoder(resp.Body).Decode(&body)
	fmt.Println(body)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   body.Message,
	})
	return
}

// get unmatched data transfer
func GetUnmatchedDataTransfer(ctx *iris.Context) {
	var url string = config.GoFinAutoReconPath + "/api/v1/data-transfer/unmatched/" + ctx.Param("branchId") + "/" + ctx.Param("transacType") + "/" + ctx.Param("transferDateFrom") + "/" + ctx.Param("transferDateTo")
	resp, err := http.Get(url)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    iris.Map{},
		})
		return
	}
	defer resp.Body.Close()

	var body struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Success bool        `json:"success"`
	}

	json.NewDecoder(resp.Body).Decode(&body)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   body.Data,
	})
	return
}

// get unmatched bank statement
func GetUnmatchedBankStatement(ctx *iris.Context) {
	var url string = config.GoFinAutoReconPath + "/api/v1/bank-statement/unmatched/" + ctx.Param("transferDateFrom") + "/" + ctx.Param("transferDateTo")
	resp, err := http.Get(url)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    iris.Map{},
		})
		return
	}
	defer resp.Body.Close()

	var body struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Success bool        `json:"success"`
	}

	json.NewDecoder(resp.Body).Decode(&body)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   body.Data,
	})
	return
}

// get settlement draft
func GetSettlementDraft(ctx *iris.Context) {
	var url string = config.GoFinAutoReconPath + "/api/v1/settlement/draft/" + ctx.Param("transactionType") + "/" + ctx.Param("branchId") + "/" + ctx.Param("transferDateFrom") + "/" + ctx.Param("transferDateTo")
	resp, err := http.Get(url)
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    iris.Map{},
		})
		return
	}
	defer resp.Body.Close()

	var body struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Success bool        `json:"success"`
	}

	json.NewDecoder(resp.Body).Decode(&body)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   body.Data,
	})
	return
}
