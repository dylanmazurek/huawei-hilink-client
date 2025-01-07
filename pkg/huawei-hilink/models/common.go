package models

type ErrorResponse struct {
	Code    string `xml:"code"`
	Message string `xml:"message"`
}
