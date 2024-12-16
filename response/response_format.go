package response

type JSONResponse struct {
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}
