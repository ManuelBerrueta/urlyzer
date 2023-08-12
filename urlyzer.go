package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {
	urlFlag := flag.String("u", "", "URL to parse")
	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("Please provide a URL using the -u flag.")
		os.Exit(1)
	}

	parsedURL, err := url.Parse(*urlFlag)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	// Scheme (Protocol)
	fmt.Printf("\033[1;36mScheme:\033[0m %s\n", parsedURL.Scheme)

	// Host
	fmt.Printf("\033[1;36mHost:\033[0m %s\n", parsedURL.Hostname())

	// Port
	fmt.Printf("\033[1;36mPort:\033[0m %s\n", parsedURL.Port())

	// Path
	fmt.Printf("\033[1;36mPath:\033[0m %s\n", parsedURL.Path)

	// Query String
	fmt.Printf("\033[1;36mQuery String:\033[0m %s\n", parsedURL.RawQuery)

	// Parse query parameters
	queryParams, _ := url.ParseQuery(parsedURL.RawQuery)
	if len(queryParams) > 0 {
		fmt.Println("\033[1;36mQuery Parameters:\033[0m")
		for key, values := range queryParams {
			fmt.Printf("  %s: %s\n", key, strings.Join(values, ", "))
		}
	}
}
