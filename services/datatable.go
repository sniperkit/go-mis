package services

import (
	"net/url"
	"strings"
	"encoding/json"
)

type (
	DataTable struct {
		Columns   []DataColumn `json:"columns"`
		OrderInfo []OrderInfo  `json:"order"`
		Start     int          `json:"start"`
		Length    int          `json:"length"`
		Search    Search       `json:"search"`
		Draw      int          `json:"draw"`
	}

	DataColumn struct {
		Data       string `json:"data"`
		Name       string `json:"name"`
		Searchable bool   `json:"searchable"`
		Orderable  bool   `json:"orderable"`
		Search     `json:"search"`
	}

	Search struct {
		Value string `json:"value"`
		Regex bool   `json:"regex"`
	}

	OrderInfo struct {
		Column uint64 `json:"column"`
		Dir    string `json:"dir"`
	}
)

func ParseDatatableURI(fullURI string) DataTable {
	var dtTables DataTable
	u, _ := url.Parse(fullURI)
	q := u.Query()
	for k, v := range q {
		if len(strings.TrimSpace(v[0])) == 0 {
			err := json.Unmarshal([]byte(k), &dtTables)
			if err == nil {
				return dtTables
			}

		}
	}
	return dtTables
}