package reports

type ResReport struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Path    string `json:"path"`
}
