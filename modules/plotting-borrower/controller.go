package plottingBorrower

import (
	"errors"
	"strconv"
	"strings"

	investor "bitbucket.org/go-mis/modules/investor"
	"bitbucket.org/go-mis/services"

	"gopkg.in/kataras/iris.v4"
)

// This function saves potting paramaters as borrower criteria
// into investor table. The data would be saved in json format
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

// This function get all investors which their borrowerCriteria is not null
func ListPlottingParams(ctx *iris.Context) {

	totalRows := 0
	investors := []struct {
		ID                       uint64 `gorm:"column:id" json:"id"`
		InvestorNo               uint64 `gorm:"column:investorNo" json:"investorNo"`
		InvestorName             string `gorm:"column:name" json:"investorName"`
		IsBorrowerCriteriaActive *bool  `gorm:"column:isBorrowerCriteriaActive" json:"isBorrowerCriteriaActive"`
	}{}

	query := `select investor.id, investor."investorNo", cif."name", investor."isBorrowerCriteriaActive" 
	from investor
	join r_cif_investor rci on rci."investorId" = investor.Id
	join cif on cif.id = rci."cifId"
	where "borrowerCriteria" <> '{}' and investor."deletedAt" is null`
	services.DBCPsql.Raw(query).Scan(&investors)
	for _, _ = range investors {
		totalRows += 1
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"status":    "success",
		"data":      investors,
		"totalRows": totalRows,
	})

}
