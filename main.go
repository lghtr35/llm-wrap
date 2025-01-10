package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lghtr35/llm-wrapping/models"
)

func main() {
	vendorConfigs := readConfiguration()
	handler := NewCommandHandler(vendorConfigs)
	s := gin.New()
	api := s.Group("/v1/api")
	{
		api.POST("/command", handler.Handle)
	}

	s.Run(":11242")
}

func readConfiguration() map[string]models.VendorConfig {
	f, err := os.Open("config.vendor.json")
	if err != nil {
		panic(err)
	}
	strContent, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	f.Close()
	var vendorConfigs []models.VendorConfig
	err = json.Unmarshal(strContent, &vendorConfigs)
	if err != nil {
		panic(err)
	}

	res := make(map[string]models.VendorConfig)

	for _, config := range vendorConfigs {
		res[config.Name] = config
	}

	return res
}
