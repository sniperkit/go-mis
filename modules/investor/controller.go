package investor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"bitbucket.org/go-mis/config"
	"bitbucket.org/go-mis/modules/r"
	"bitbucket.org/go-mis/modules/virtual-account"
	"bitbucket.org/go-mis/services"
	iris "gopkg.in/kataras/iris.v4"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

func Init() {
	services.DBCPsql.AutoMigrate(&Investor{})
	services.BaseCrudInit(Investor{}, []Investor{})
}

func GetInvestorDetailById(id uint64) (investor Investor, err error) {
	investor = Investor{}
	queryDetailInvestorByID := `select * from investor `
	queryDetailInvestorByID += strings.Replace("WHERE id = ?", "?", strconv.Itoa(int(id)), -1)
	if err := services.DBCPsql.Raw(queryDetailInvestorByID).Scan(&investor).Error; err != nil {
		return investor, err
	}
	return investor, nil
}

func CheckInvestorSwiftCode(ctx *iris.Context) {
	swiftCode := ctx.Param("swiftCode")

	type Body struct {
		SwiftCode string `json:"username"`
	}

	body := Body{
		SwiftCode: swiftCode,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(body)

	res, _ := http.Post(config.GoWithdrawalPath+"/api/v1/cashout/swiftcode-check", "application/json; charset=utf-8", b)

	resp := Response{}

	json.NewDecoder(res.Body).Decode(&resp)

	if resp.Status == "404 Not Found." {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   false,
		})
	} else {
		ctx.JSON(iris.StatusOK, iris.Map{
			"status": "success",
			"data":   true,
		})
	}
}

func FetchDetail(ctx *iris.Context) {
	stage := ctx.URLParam("stage")
	request := ctx.URLParam("request")

	type CashoutInvestorID struct {
		Id uint64 `gorm:"column:id" json:"id"`
	}

	listIdCashoutInvestors := []CashoutInvestorID{}
	listInvestor := []Investor{}
	queryPendingInvestorID := ""

	if request == "" {
		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": "Can't process your request further, expected request query params value not empty"})
		return
	} else {

		if request == "ALL" {
			if stage == "" {
				queryPendingInvestorID = `
				select i.id  from investor i join r_investor_cashout ric on i.id = ric."investorId" 
				join cashout on cashout.id = ric."cashoutId"
				where cashout.stage not like 'SUCCESS'; 
				`
			} else {
				queryPendingInvestorID = `
					select i.id  from investor i join r_investor_cashout ric on i.id = ric."investorId" 
					join cashout on cashout.id = ric."cashoutId" 
				`
				queryPendingInvestorID += strings.Replace("where cashout.stage like '?'", "?", stage, -1)
			}
			services.DBCPsql.Raw(queryPendingInvestorID).Find(&listIdCashoutInvestors)
			for _, v := range listIdCashoutInvestors {
				if investor, err := GetInvestorDetailById(v.Id); err != nil {
					ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
					return
				} else {
					listInvestor = append(listInvestor, investor)
				}
			}
		} else {
			type Payload struct {
				Id []uint64 `json:"id"`
			}

			listStringId := Payload{}
			if err := ctx.ReadJSON(&listStringId); err != nil {
				ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
				return
			} else if len(listStringId.Id) == 0 {
				ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": "Request payload can't be empty"})
				return
			} else {
				for _, v := range listStringId.Id {
					if investor, err := GetInvestorDetailById(v); err != nil {
						ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
						return
					} else {
						listInvestor = append(listInvestor, investor)
					}
				}
			}
		}
	}

	// 	if request == "ALL" && stage == ""{
	// 		queryPendingInvestorID = `
	// 			select i.id  from investor i join r_investor_cashout ric on i.id = ric."investorId"
	// 			join cashout on cashout.id = ric."cashoutId"
	// 			where cashout.stage not like 'SUCCESS';
	// 		`
	// 	} else if request == "ALL" && stage != "" {
	// 		queryPendingInvestorID = `
	// 			select i.id  from investor i join r_investor_cashout ric on i.id = ric."investorId"
	// 			join cashout on cashout.id = ric."cashoutId"
	// 		`
	// 		queryPendingInvestorID += strings.Replace("where cashout.stage like '?'","?",stage,-1)
	// 	} else if request == "SPECIFIC" && stage == "" {

	// 	} else if request == "SPECIFIC" && stage != "" {

	// 	}

	// 	if err != nil {
	// 		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
	// 		return
	// 	} else {
	// 		for _, v := range listStringId.Id {
	// 			investor := GetInvestorDetailById(v)
	// 			listInvestor = append(listInvestor, investor)
	// 		}
	// 	}

	// 	var queryPendingInvestorID string;
	// 	if stage == "ALL" {
	// 		queryPendingInvestorID = `
	// 			select i.id  from investor i join r_investor_cashout ric on i.id = ric."investorId"
	// 			join cashout on cashout.id = ric."cashoutId"
	// 			where cashout.stage not like 'SUCCESS';
	// 		`
	// 	} else {
	// 		queryPendingInvestorID = `
	// 			select i.id  from investor i join r_investor_cashout ric on i.id = ric."investorId"
	// 			join cashout on cashout.id = ric."cashoutId"
	// 		`
	// 		queryPendingInvestorID += strings.Replace("where cashout.stage like '?'","?",stage,-1)
	// 	}

	// 	services.DBCPsql.Raw(queryPendingInvestorID).Find(&listIdCashoutInvestors)

	// 	for _, v := range listIdCashoutInvestors {
	// 		investor := GetInvestorDetailById(v.Id)
	// 		listInvestor = append(listInvestor, investor)
	// 	}

	// } else {
	// 	type Payload struct {
	// 		Id []uint64 `json:"id"`
	// 	}

	// 	listStringId := Payload{}
	// 	err := ctx.ReadJSON(&listStringId)

	// 	if err != nil {
	// 		ctx.JSON(iris.StatusBadRequest, iris.Map{"status": "error", "message": err.Error()})
	// 		return
	// 	} else {
	// 		for _, v := range listStringId.Id {
	// 			investor := GetInvestorDetailById(v)
	// 			listInvestor = append(listInvestor, investor)
	// 		}
	// 	}
	// }

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   listInvestor,
	})
}

// InvestorWithoutVA - Retrieve list of investor without VA
func InvestorWithoutVA(ctx *iris.Context) {
	query := "SELECT cif.\"name\", investor.id AS \"investorId\", investor.\"investorNo\", virtual_account.\"bankName\", virtual_account.\"virtualAccountNo\", virtual_account.\"virtualAccountName\" "
	query += "FROM investor "
	query += "LEFT OUTER JOIN r_investor_virtual_account ON r_investor_virtual_account.\"investorId\" = investor.id  "
	query += "LEFT OUTER JOIN virtual_account ON virtual_account.id = r_investor_virtual_account.\"investorId\" "
	query += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.id "
	query += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "
	query += "WHERE virtual_account.\"virtualAccountNo\" IS NULL  "
	query += "AND cif.\"deletedAt\" IS NULL "
	query += "AND cif.\"isActivated\" = true AND cif.\"isValidated\" = true AND cif.\"isVerified\" = false "
	query += "AND virtual_account.\"deletedAt\" IS NULL "

	investorWithoutVaSchema := []InvestorWithoutVaSchema{}
	services.DBCPsql.Raw(query).Scan(&investorWithoutVaSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data":   investorWithoutVaSchema,
	})
}

type InvestorVASchema struct {
	InvestorID   uint64 `json:"investorId"`
	InvestorNo   uint64 `json:"investorNo"`
	InvestorName string `json:"investorName"`
	VaBri        string `json:"vaBri"`
	VaBca        string `json:"vaBca"`
}

// InvestorRegisterVA - register VA to investor
func InvestorRegisterVA(ctx *iris.Context) {
	investorVASchema := InvestorVASchema{}

	if err := ctx.ReadJSON(&investorVASchema); err != nil {
		fmt.Println(investorVASchema)
		ctx.JSON(iris.StatusInternalServerError, iris.Map{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	vaSchemaBCA := &virtualAccount.VirtualAccount{BankName: "BCA", VirtualAccountCode: "04435", VirtualAccountName: investorVASchema.InvestorName, VirtualAccountNo: investorVASchema.VaBca}
	services.DBCPsql.Create(vaSchemaBCA)

	rInvestorVaBca := &r.RInvestorVirtualAccount{InvestorId: investorVASchema.InvestorID, VirtualAccountId: vaSchemaBCA.ID}
	services.DBCPsql.Create(rInvestorVaBca)

	vaSchemaBRI := &virtualAccount.VirtualAccount{BankName: "BRI", VirtualAccountCode: "99959", VirtualAccountName: investorVASchema.InvestorName, VirtualAccountNo: investorVASchema.VaBri}
	services.DBCPsql.Create(vaSchemaBRI)

	rInvestorVaBri := &r.RInvestorVirtualAccount{InvestorId: investorVASchema.InvestorID, VirtualAccountId: vaSchemaBRI.ID}
	services.DBCPsql.Create(rInvestorVaBri)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status": "success",
		"data": iris.Map{
			"bca": vaSchemaBCA,
			"bri": vaSchemaBRI,
		},
	})
}

type TotalData struct {
	TotalRows int64 `gorm:"column:totalRows" json:"totalRows"`
}

type InvestorSchema struct {
	AccountID  uint64 `gorm:"column:accountId" json:"accountId"`
	CifID      uint64 `gorm:"column:cifId" json:"cifId"`
	InvestorID uint64 `gorm:"column:investorId" json:"investorId"`
	Fullname   string `gorm:"column:fullname" json:"fullname"`
	Username   string `gorm:"column:username" json:"username"`
	PhoneNo    string `gorm:"column:phoneNo" json:"phoneNo"`
}

func GetInvestorForTopup(ctx *iris.Context) {
	queryCount := "SELECT count(*) as \"totalRows\" "
	queryCount += "FROM investor "
	queryCount += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.id "
	queryCount += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "
	queryCount += "JOIN r_account_investor ON r_account_investor.\"investorId\" = investor.id "
	queryCount += "WHERE cif.\"deletedAt\" IS NULL AND investor.\"deletedAt\" IS NULL "

	queryGetInvestor := "SELECT r_account_investor.id AS \"accountId\", cif.id AS \"cifId\", investor.id AS \"investorId\", cif.\"name\" AS \"fullname\", cif.\"username\" AS \"username\", cif.\"phoneNo\"   "
	queryGetInvestor += "FROM investor "
	queryGetInvestor += "JOIN r_cif_investor ON r_cif_investor.\"investorId\" = investor.id "
	queryGetInvestor += "JOIN cif ON cif.id = r_cif_investor.\"cifId\" "
	queryGetInvestor += "JOIN r_account_investor ON r_account_investor.\"investorId\" = investor.id "
	queryGetInvestor += "WHERE cif.\"deletedAt\" IS NULL AND investor.\"deletedAt\" IS NULL "

	if ctx.URLParam("search") != "" {
		queryCount += "AND cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
		queryGetInvestor += "AND cif.\"name\" ~* '" + ctx.URLParam("search") + "' "
	}

	var limitPagination int64 = 10
	var offset int64 = 0

	if ctx.URLParam("limit") != "" {
		queryGetInvestor += "LIMIT " + ctx.URLParam("limit") + " "
	} else {
		queryGetInvestor += "LIMIT " + strconv.FormatInt(limitPagination, 10) + " "
	}

	if ctx.URLParam("page") != "" {
		offset, _ = strconv.ParseInt(ctx.URLParam("page"), 10, 64)
		queryGetInvestor += "OFFSET " + strconv.FormatInt(offset, 10)
	} else {
		queryGetInvestor += "OFFSET 0"
	}

	totalData := TotalData{}
	services.DBCPsql.Raw(queryCount).Find(&totalData)

	investorSchema := []InvestorSchema{}
	services.DBCPsql.Raw(queryGetInvestor).Find(&investorSchema)

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"totalRows": totalData.TotalRows,
		"data":      investorSchema,
	})
}
