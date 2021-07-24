package dto

type Transfer struct {
	Timestamp int64  `json:"time,string"`
	Cmd       string `json:"cmd"`
	Package   string `json:"package"`
}
