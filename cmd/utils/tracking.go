package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// record usage data if user agrees to contribute
const APIEndpoint = "https://api-api-dev.bk7bbm.oss-acorn.io/save-vimana-usage"

type APIVimanaUsage struct {
	EthAddress string  `json:"eth_address,omitempty"`
	Kvm        bool    `json:"kvm"`
	CpuCount   int     `json:"cpu_count"`
	RamSize    float64 `json:"ram_size"`
	DiskSize   float64 `json:"disk_size"`
	InitDate   string  `json:"init_date"`
	SpaceCore  string  `json:"space_core"`
	IpAddress  string  `json:"ip_address"`
}

func SaveAnalyticsData(conf *InitConfig) {
	ip, _ := GetExternalIP()
	data := APIVimanaUsage{
		EthAddress: conf.EthAddress,
		Kvm:        conf.Kvm,
		CpuCount:   conf.CpuCount,
		RamSize:    conf.RamSize,
		DiskSize:   conf.DiskSize,
		InitDate:   conf.InitDate,
		SpaceCore:  conf.SpaceCore,
		IpAddress:  ip,
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

func GetExternalIP() (string, error) {
	response, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	ip, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}
