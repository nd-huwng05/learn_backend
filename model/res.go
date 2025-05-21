package model

type Response struct {
	StausCode int         `json:"staus_code,omitempty"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"` //giong nhu object trong OOP cac kieu du lieu
}
