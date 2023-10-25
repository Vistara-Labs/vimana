package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const APIEndpoint = "http://localhost:8080/save-vimana-usage"

type APIVimanaUsage struct {
	EthAddress string  `json:"eth_address,omitempty"`
	Kvm        bool    `json:"kvm"`
	CpuCount   int     `json:"cpu_count"`
	RamSize    float64 `json:"ram_size"`
	DiskSize   float64 `json:"disk_size"`
	InitDate   string  `json:"init_date"`
	SpaceCore  string  `json:"space_core"`
}

func SendAnonymousData(conf *InitConfig) {
	data := APIVimanaUsage{
		EthAddress: conf.EthAddress,
		Kvm:        conf.Kvm,
		CpuCount:   conf.CpuCount,
		RamSize:    conf.RamSize,
		DiskSize:   conf.DiskSize,
		InitDate:   conf.InitDate,
		SpaceCore:  conf.SpaceCore,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal config data:", err)
		return
	}

	req, err := http.NewRequest("POST", APIEndpoint, bytes.NewBuffer(jsonData))
	fmt.Println("req:", req)
	if err != nil {
		log.Println("Failed to create request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send data:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to send analytics data, status:", resp.Status)
	}
}
