package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	clearScreen()
	blue := "\033[34m"
	green := "\033[32m"
	yellow := "\033[33m"
	red := "\033[31m"
	reset := "\033[0m"

	fmt.Printf("%s[*] Enter the URL: %s", blue, reset)
	reader := bufio.NewReader(os.Stdin)
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	if url == "" {
		fmt.Printf("%s[-] Please enter a valid URL%s\n", red, reset)
		return
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	checkURL(url, green, yellow, blue, red, reset)
}

func checkURL(url, green, yellow, blue, red, reset string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s[-] Could not connect to %s%s\n", red, url, reset)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		fmt.Printf("%s[+] %s is up%s\n", green, url, reset)
	case 404:
		fmt.Printf("%s[-] %s is down%s\n", red, url, reset)
	case 403:
		fmt.Printf("%s[!] %s is blocked%s\n", yellow, url, reset)
	case 500, 501, 502, 503, 504:
		fmt.Printf("%s[*] %s server error%s\n", blue, url, reset)
	case 400:
		fmt.Printf("%s[-] %s bad request%s\n", red, url, reset)
	case 401:
		fmt.Printf("%s[!] %s unauthorized%s\n", yellow, url, reset)
	case 429:
		fmt.Printf("%s[!] %s too many requests%s\n", yellow, url, reset)
	case 301:
		fmt.Printf("%s[+] %s moved permanently%s\n", green, url, reset)
	case 302:
		fmt.Printf("%s[+] %s temporary redirect%s\n", green, url, reset)
	default:
		fmt.Printf("%s[*] %s returned status code %d%s\n", blue, url, resp.StatusCode, reset)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
