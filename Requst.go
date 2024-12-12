package main

import (
	"fmt"
	"os"
	"os/exec"
	"net/http"
	"strings"
	"regexp"
	"bufio"

)

// Function to print colored messages
func printColored(color, message string) {
	fmt.Printf("%s%s%s\n", color, message, "\033[0m")
}

func requ(url string) {
	if url == "" {
		printColored("\033[31m", "[-] Please enter a valid URL")
		os.Exit(1)
	}

	// Ensure the URL starts with "http://" or "https://"
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	// Make an HTTP request to check the status of the URL
	resp, err := http.Get(url)
	if err != nil {
		printColored("\033[31m", "[-] Could not connect to "+url)
		return
	}
	defer resp.Body.Close()

	// Handle different status codes
	switch resp.StatusCode {
	case 200:
		printColored("\033[32m", "[+] "+url+" is up")

		// Extract domain from the URL
		re := regexp.MustCompile(`(?:http[s]?://)?([^/]+)`)
		matches := re.FindStringSubmatch(url)
		if len(matches) < 2 {
			printColored("\033[31m", "[-] Invalid URL format")
			return
		}
		domain := matches[1]

		// Run the dig command to get DNS records
		cmd := exec.Command("dig", domain, "ANY", "+short")
		output, err := cmd.CombinedOutput()
		if err != nil {
			printColored("\033[31m", "Error occurred while fetching DNS records: "+err.Error())
			return
		}

		if len(output) > 0 {
			printColored("\033[34m", "\nDNS Records for "+domain+":")
			fmt.Println(string(output))
		} else {
			printColored("\033[33m", "[!] No DNS records found for "+domain)
		}

	case 404:
		printColored("\033[31m", "[-] "+url+" is down")
	case 403:
		printColored("\033[33m", "[!] "+url+" is blocked")
	case 500, 503, 504, 502, 501:
		printColored("\033[34m", "[*] "+url+" is not found")
	case 400:
		printColored("\033[31m", "[-] "+url+" bad request")
	case 401:
		printColored("\033[33m", "[!] "+url+" unauthorized")
	case 429:
		printColored("\033[33m", "[!] "+url+" too many requests")
	case 301:
		printColored("\033[32m", "[+] "+url+" moved permanently")
	case 302:
		printColored("\033[32m", "[+] "+url+" temporary redirect")
	default:
		printColored("\033[31m", "[-] Unexpected status code for "+url)
	}
}

func main() {
	// Read URL from user input
	fmt.Print("\033[34m[*] Enter the url: \033[0m")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	url := scanner.Text()

	requ(url)
}
