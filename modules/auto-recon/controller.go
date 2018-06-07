package autoRecon

import (
	"fmt"

	"bytes"
	"encoding/json"
	"net/http"

	"bitbucket.org/go-mis/config"
	"gopkg.in/kataras/iris.v4"
)

func DisbursementDataTransferSave(ctx *iris.Context) {

	payload := struct {
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

// save / update settlement draft
func SaveSettlementDraft(ctx *iris.Context) {

	payload := struct {
		BankExcess          float64
		BranchId            uint64
		MatchedAmount       float64
		MisExcess           float64
		Stage               string
		TransactionDateFrom string
		TransactionDateTo   string
		TransactionType     string
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

	var url string = config.GoFinAutoReconPath + "/api/v1/settlement/draft/save"
	req, err := http.NewRequest("PUT", url, b)
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

// get disbursement data transfer list
func GetDisbursementDataTransferList(ctx *iris.Context) {
	var url string = config.GoFinAutoReconPath + "/api/v1/data-transfer/disbursement/list/" + ctx.Param("branchId") + "/" + ctx.Param("transferDateFrom") + "/" + ctx.Param("transferDateTo") + "/" + ctx.Param("disburseDateFrom") + "/" + ctx.Param("disburseDateTo")
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

// get data transfer list
func GetDataTransferList(ctx *iris.Context) {
	var url string = config.GoFinAutoReconPath + "/api/v1/data-transfer/list/" + ctx.Param("branchId") + "/" + ctx.Param("validationDateFrom") + "/" + ctx.Param("validationDateTo") + "/" + ctx.Param("transferDateFrom") + "/" + ctx.Param("transferDateTo")
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

// get data transfer detail
func GetDataTransferDetail(ctx *iris.Context) {
	var url string = config.GoFinAutoReconPath + "/api/v1/data-transfer/detail/" + ctx.Param("branchId") + "/" + ctx.Param("validationDate") + "/" + "/" + ctx.Param("transferDate")
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

// get settlement detail
func GetSettlementDetail(ctx *iris.Context) {
	var url string = config.GoFinAutoReconPath + "/api/v1/transaction-reconciliation/settlement/" + ctx.Param("id")
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

// get settlement by criteria
func GetSettlementByCriteria(ctx *iris.Context) {
	var url string = config.GoFinAutoReconPath + "/api/v1/transaction-reconciliation/settlement"

	settlementId := ctx.URLParam("settlementId")
	branchId := ctx.URLParam("branchId")
	transactionType := ctx.URLParam("transactionType")
	transactionDateFrom := ctx.URLParam("transactionDateFrom")
	transactionDateTo := ctx.URLParam("transactionDateTo")

	if branchId == "" || transactionDateFrom == "" || transactionDateTo == "" {
		ctx.JSON(iris.StatusBadRequest, iris.Map{
			"status":  "error",
			"message": "branchId, transactionDateFrom, and transactionDateTo are mandatory",
			"data":    iris.Map{},
		})
		return
	}

	req, err := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	q.Add("branchId", branchId)
	q.Add("transactionDateFrom", transactionDateFrom)
	q.Add("transactionDateTo", transactionDateTo)

	if settlementId != "" {
		q.Add("settlementId", settlementId)
	}

	if transactionType != "" {
		q.Add("transactionType", transactionType)
	}

	req.URL.RawQuery = q.Encode()

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

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   body.Data,
	})
	return
}
