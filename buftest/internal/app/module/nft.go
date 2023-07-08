package module

type CreateNftRequest struct {
	ClassId string `json:"classId"`
	Uri     string `json:"uri"`
	UriHash string `json:"uriHash"`
}

type CreateNftResponse struct {
	Data string `json:"data"`
}

type ClassByIDRequest struct {
	Name   string `json:"name"`
	TxHash string `json:"txHash"`
	Owner  string `json:"owner"`
	ID     string `json:"id"`
}

type ClassByIDResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Uri  string `json:"uri"`
}
