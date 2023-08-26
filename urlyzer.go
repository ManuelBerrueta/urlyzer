package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {
	flag.Parse()

	var inputURL string

	// Check if input is being piped
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Read URL from standard input
		inputBytes, _ := io.ReadAll(os.Stdin)
		inputURL = strings.TrimSpace(string(inputBytes))
	} else if len(flag.Args()) > 0 {
		inputURL = flag.Arg(0)
	} else {
		fmt.Println("Please provide a URL as a command line argument or through standard input.")
		os.Exit(1)
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	// Scheme (Protocol)
	fmt.Printf("\033[1;36mScheme:\033[0m %s\n", parsedURL.Scheme)

	// UserInfo
	if parsedURL.User.String() != "" {
		fmt.Printf("\033[1;36mUserInfo:\033[0m %s\n", parsedURL.User)
	}

	// Host
	fmt.Printf("\033[1;36mHost:\033[0m %s\n", parsedURL.Hostname())

	// Port
	if parsedURL.Port() != "" {
		fmt.Printf("\033[1;36mPort:\033[0m %s\n", parsedURL.Port())
	}

	// Path
	fmt.Printf("\033[1;36mPath:\033[0m %s\n", parsedURL.Path)

	// Query String
	if parsedURL.RawQuery != "" {
		fmt.Printf("\033[1;36mQuery String:\033[0m %s\n", parsedURL.RawQuery)

		// Parse query parameters
		queryParams, _ := url.ParseQuery(parsedURL.RawQuery)
		if len(queryParams) > 0 {
			fmt.Println("\033[1;36mQuery Parameters:\033[0m")
			for key, values := range queryParams {
				fmt.Printf("  \033[1;32m%s:\033[0m %s\n", key, strings.Join(values, ", "))
			}
		}
	}

	// Fragment (Hash)
	if parsedURL.Fragment != "" {
		fmt.Printf("\033[1;36mFragment:\033[0m %s\n", parsedURL.Fragment)
	}
}
