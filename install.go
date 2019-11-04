package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sync"
)

func getStringFromReq(url string) *string {
	resp, err := client.Get(url)

	handleErr(err)

	defer resp.Body.Close()

	bodyB, err := ioutil.ReadAll(resp.Body)

	handleErr(err)

	body := string(bodyB)

	return &body
}

func getCurrentVersion() string {
	url := "https://ungoogled-software.github.io/ungoogled-chromium-binaries/releases/debian/buster_amd64/"

	body := getStringFromReq(url)

	re, err := regexp.Compile(`[a-zA-Z\-\.0-9]+\.buster1`)
	handleErr(err)

	return re.FindString(*body)
}

func getPacketsLinks(packetVersion string) []string {
	url := "https://ungoogled-software.github.io/" +
		"ungoogled-chromium-binaries/releases/debian/buster_amd64/" +
		packetVersion

	body := getStringFromReq(url)

	re, err := regexp.Compile(`https://github.com/Eloston/.*?\.deb`)
	handleErr(err)

	return re.FindAllString(*body, -1)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

var client *http.Client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		DisableKeepAlives:   false,
	},
}

func main() {
	var installedVersion string
	installedVersionChan := make(chan string)

	go func() {
		re, _ := regexp.Compile(`[a-zA-Z\-\.0-9]+\.buster1`)
		cmd := exec.Command("apt", "list", "ungoogled-chromium")
		stdout, err := cmd.Output()
		handleErr(err)

		versionB := re.Find(stdout)

		installedVersionChan <- string(versionB)
	}()

	currentVersion := getCurrentVersion()
	installedVersion = <-installedVersionChan

	installedVersionChan = nil

	if currentVersion == "" {
		fmt.Println("Failed to retrieve current version")
		return
	}

	if installedVersion == currentVersion {
		fmt.Println("Ungoogled chromium is already up to date")
		return
	}

	const downloadPath string = "binaries/"

	if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
		os.Mkdir(downloadPath, 0755)
	}

	var wg sync.WaitGroup

	for _, url := range getPacketsLinks(currentVersion) {
		wg.Add(1)

		go func() {
			fmt.Println("[-] Downloading", url)

			defer wg.Done()

			resp, err := client.Get(url)
			handleErr(err)
			defer resp.Body.Close()

			filepath := downloadPath + path.Base(url)
			fd, err := os.Create(filepath)
			handleErr(err)
			defer fd.Close()

			_, err = io.Copy(fd, resp.Body)
			handleErr(err)

			fmt.Println("[+] Downloaded", url)
		}()
	}

	wg.Wait()
}
