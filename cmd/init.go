package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"vimana/cmd/utils"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var initFilePath = os.Getenv("HOME") + "/.vimana/init.toml"

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func InitializeSystem(force bool, noTrack bool) error {
	if err := fileExists(initFilePath); err && !force {
		fmt.Print("Initialization has already been done. Found init.toml.\n")
		return nil
	}
	// config := InitConfig{}
	// Determine the path to the init.toml file
	initPath := filepath.Join(os.Getenv("HOME"), ".vimana", "init.toml")

	config, err := utils.LoadVimanaConfig(initPath)
	// config, err := LoadVimanaConfig(initPath)
	if err != nil {
		return err
	}

	if noTrack {
		config.Analytics.Enabled = false
	}
	// Only prompt if initializing for the first time
	if config.Analytics.Enabled {
		fmt.Println("Vimana collects usage data to improve our software and show our network growth.")
		var response string
		fmt.Println("Would you like to contribute? [Y/n]")
		fmt.Scanln(&response)
		response = strings.ToLower(response)
		if response == "n" {
			config.Analytics.Enabled = false
		}
	}

	fmt.Println("Do you want to create a new Ethereum address? (Y/n)")
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "n" {
		address, privateKey, err := createEthAddress()
		if err != nil {
			log.Fatalf("Failed to create Ethereum address: %s", err)
		}
		config.EthAddress = address
		config.EthPrivateKey = privateKey
		log.Printf("Ethereum address: %s", address)
	}

	cpuCount, _ := cpu.Counts(true)
	config.CpuCount = cpuCount

	memInfo, _ := mem.VirtualMemory()
	config.RamSize = float64(memInfo.Total) / (1 << 30) // Convert to GB

	diskInfo, _ := disk.Usage("/")
	config.DiskSize = float64(diskInfo.Total) / (1 << 30) // Convert to GB

	config.InitDate = time.Now().Format(time.RFC1123)

	kvmSupport, err := checkKvmSupport()
	if err != nil {
		return err
	}
	log.Printf("CPU Count: %d", config.CpuCount)
	log.Printf("Total RAM: %v GB\n", float64(memInfo.Total)/(1<<30))   // Convert to GB
	log.Printf("Total Disk: %v GB\n", float64(diskInfo.Total)/(1<<30)) // Convert to GB
	log.Printf("KVM Support: %v\n", kvmSupport)

	err = utils.SaveConfig(config, initPath)
	if err != nil {
		return err
	}
	return nil
}

func createEthAddress() (address string, privateKey string, err error) {
	// Generate a new private key
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	privateKey = hex.EncodeToString(crypto.FromECDSA(key))
	return address, privateKey, nil
}

func checkKvmSupport() (bool, error) {
	// Execute the kvm-ok command
	out, err := exec.Command("sh", "-c", "kvm-ok 2>&1 | grep -o 'KVM acceleration can be used'").Output()
	if err != nil {
		// log.Printf("kvm-ok tool might be missing or another error occurred: %w", err)
		return false, nil
	}
	// If the string "KVM acceleration can be used" is found, KVM is supported.
	return string(out) == "KVM acceleration can be used", nil
}
