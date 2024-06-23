package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"vimana/cmd/utils"

	"github.com/asmcos/requests"

	"github.com/shirou/gopsutil/host"
	"github.com/spf13/cobra"
)

func migrateCommand() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Upgrade vimana to latest version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("vimana version: ", Version)
			upgradeVimana()
		},
	}
	return migrateCmd
}

func upgradeVimana() {
	resp, err := requests.Get("https://api.github.com/repos/Vistara-Labs/vimana/releases")
	if err != nil {
		return
	}
	// Status code
	if resp.R.StatusCode != 200 {
		fmt.Println("Downloading vimana error: status code =", resp.R.StatusCode)
		return
	}
	// reponse
	// println(resp.Text())

	type person struct {
		Login               string `json:"login"`
		Id                  int    `json:"id"`
		Node_id             string `json:"node_id"`
		Avatar_url          string `json:"avatar_url"`
		Gravatar_id         string `json:"gravatar_id"`
		Url                 string `json:"url"`
		Html_url            string `json:"html_url"`
		Followers_url       string `json:"followers_url"`
		Following_url       string `json:"following_url"`
		Gists_url           string `json:"gists_url"`
		Starred_url         string `json:"starred_url"`
		Subscriptions_url   string `json:"subscriptions_url"`
		Organizations_url   string `json:"organizations_url"`
		Repos_url           string `json:"repos_url"`
		Events_url          string `json:"events_url"`
		Received_events_url string `json:"received_events_url"`
		Type                string `json:"type"`
		Site_admin          bool   `json:"site_admin"`
	}

	type asset struct {
		Url                  string `json:"url"`
		Id                   int    `json:"id"`
		Node_id              string `json:"node_id"`
		Name                 string `json:"name"`
		Label                string `json:"label"`
		Uploader             person
		Content_type         string `json:"content_type"`
		State                string `json:"state"`
		Size                 int    `json:"size"`
		Download_count       int    `json:"download_count"`
		Created_at           string `json:"created_at"`
		Updated_at           string `json:"updated_at"`
		Browser_download_url string `json:"browser_download_url"`
	}

	type ReleaseItem struct {
		Url              string `json:"url"`
		Assets_url       string `json:"asset_url"`
		Upload_url       string `json:"upload_url"`
		Html_url         string `json:"html_url"`
		Id               int    `json:"id"`
		Author           person
		Node_id          string `json:"node_id"`
		Tag_name         string `json:"tag_name"`
		Target_commitish string `json:"target_commitish"`
		Name             string `json:"name"`
		Draft            bool   `json:"draft"`
		Prerelease       bool   `json:"prerelease"`
		Created_at       string `json:"created_at"`
		Published_at     string `json:"published_at"`
		Assets           []asset
		Tarball_url      string `json:"tarball_url"`
		Zipball_url      string `json:"zipball_url"`
		Body             string `json:"body"`
	}
	var Releases []ReleaseItem
	resp.Json(&Releases)
	//for i, s := range Releases {
	//	fmt.Println(i, s.Tag_name)
	//}
	var current_version = strings.Split(Version, "v")[1]
	var latest_version = strings.Split(Releases[0].Tag_name, "-v")[1]
	//fmt.Println(latest_version)
	if current_version != latest_version {
		fmt.Println("No need to upgrade vimana")
		return
	}

	v, _ := host.Info()
	// convert to JSON. String() is also implemented
	OS := v.OS
	ARCH := v.KernelArch
	//fmt.Println(OS, ARCH)
	if ARCH == "x86_64" {
		ARCH = "amd64"
	} else if ARCH == "arm64" || ARCH == "aarch64" {
		ARCH = "arm64"
	}

	// https://github.com/Vistara-Labs/vimana/releases/download/vimana-v0.0.151/vimana-darwin-arm64.tar.gz
	file_name := "vimana-" + OS + "-" + ARCH + ".tar.gz"
	url := "https://github.com/Vistara-Labs/vimana/releases/download/vimana-v" + latest_version + "/" + file_name

	err = utils.DownloadTarGzFile(file_name, url)
	if err == nil {
		fmt.Println("Download vimana via " + url + " successfully")
	} else {
		fmt.Println(err)
		return
	}

	gzipStream, err := os.Open(file_name)
	if err != nil {
		fmt.Println("error")
		return
	}
	if err := utils.ExtractTarGz(gzipStream); err != nil {
		fmt.Println("ExtractTarGz failed: %w", err)
		return
	}

	//
	VIMANA_BIN_PATH := "/usr/local/bin/vimana"
	vimana := "./vimana-" + OS + "-" + ARCH + "/vimana"

	err = os.Rename(vimana, VIMANA_BIN_PATH)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Upgrade vimana succesfully")
	}
	os.Chmod(VIMANA_BIN_PATH, os.FileMode(0755))
	os.RemoveAll("./vimana-" + OS + "-" + ARCH + "/")
	//

}
