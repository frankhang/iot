package dto

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Errno  string `json:"errno"`
	Error  string `json:"error"`
}
