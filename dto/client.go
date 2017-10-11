package dto

type Client struct {
	ClientNo int    `json:"clientNo"`
	Name     string `json:"name"`
	Action   bool   `json:"action"` //
	Self     bool   `json:"self"`
}
