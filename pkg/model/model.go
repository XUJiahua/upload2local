package model

const (
	CodeSuccess = 200
	CodeFailure = 500
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CompleteRequest struct {
	PartFiles []string `json:"part_files"`
	Filename  string   `json:"filename"`
}
