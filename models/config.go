package models

type VendorConfig struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	ApiKey string `json:"api_key"`
	Model  string `json:"model"`
}
