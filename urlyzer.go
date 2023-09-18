package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ANSI escape codes for colors
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func main() {
	var proxyURL string
	flag.StringVar(&proxyURL, "p", "", "Proxy URL (optional)")
	checkFinalURLDestination := flag.Bool("f", false, "Check the final destination of a URL after redirects.")
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

	if *checkFinalURLDestination {
		finalDestination, headers, statusCode, err := getFinalDestination(proxyURL, inputURL)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("%sFinal Destination:%s %s\n", Cyan, Reset, finalDestination)
		fmt.Printf("%sStatus Code:%s %d\n", Yellow, Reset, statusCode)

		fmt.Printf("%sHeaders:%s\n", Green, Reset)
		for key, value := range headers {
			fmt.Printf("  %s%s:%s %s\n", Blue, key, Reset, value)
		}

	} else { //! Regular urlyzer analysis operation
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
}

func getFinalDestination(proxyURL, targetURL string) (string, http.Header, int, error) {
	var client *http.Client

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			return "", nil, 0, err
		}

		transport := &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip certificate check
		}

		client = &http.Client{
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // Prevent redirects
			},
		}
	} else {
		client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // Prevent redirects
			},
		}
	}

	resp, err := client.Get(targetURL) // Use GET instead of HEAD
	if err != nil {
		return "", nil, 0, err
	}

	defer resp.Body.Close()

	finalDestination := targetURL
	statusCode := resp.StatusCode
	if statusCode >= 300 && statusCode <= 399 {
		location, err := resp.Location()
		if err != nil {
			return "", nil, 0, err
		}
		finalDestination, _, statusCode, err = getFinalDestination(proxyURL, location.String())
		if err != nil {
			return "", nil, 0, err
		}
	}

	return finalDestination, resp.Header, statusCode, nil
}
