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
	cookies := flag.Bool("c", false, "Parse cookies from a cookie header.")
	sas := flag.Bool("sas", false, "Parse SAS URI & identify it's type.")
	flag.Parse()

	// Check if input is being piped
	// Check if input is being piped
	var inputURL string
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Read from standard input
		inputBytes, _ := io.ReadAll(os.Stdin)
		inputURL = strings.TrimSpace(string(inputBytes))
	} else if len(flag.Args()) > 0 {
		// Read from command line argument
		inputURL = flag.Arg(0)
	} else {
		fmt.Println("Please provide a URL or a cookie header as a command line argument or through standard input.")
		os.Exit(1)
	}

	if *cookies {
		// Parse the cookie header from the input
		cookiesMap, err := parseCookies(inputURL)
		if err != nil {
			fmt.Println("Error parsing cookies:", err)
			os.Exit(1)
		}

		// Print the parsed cookies
		fmt.Printf("%sCookies:%s\n", Green, Reset)
		for name, value := range cookiesMap {
			fmt.Printf("  %s%s:%s %s\n", Blue, name, Reset, value)
		}

	} else if *sas {
		parsedURL, err := url.Parse(inputURL)
		if err != nil {
			log.Fatalf("Error parsing URL: %v", err)
		}

		values := parsedURL.Query()
		// Identify the type of SAS URI
		sasType := identifySASURIType(values)
		fmt.Printf("\033[31mSAS URI Type: %s\033[0m\n", sasType)

		// Print the parsed query parameters
		for key, value := range values {
			longFormName, longFormValue := getLongForm(sasType, key, value[0]) // Assuming each key has only one value
			fmt.Printf("%s=%s\t||\t%s = %s\n", key, value[0], longFormName, longFormValue)
		}
	} else if *checkFinalURLDestination {
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

func parseCookies(cookieHeader string) (map[string]string, error) {
	cookies := make(map[string]string)
	pairs := strings.Split(cookieHeader, ";")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		equalIndex := strings.Index(pair, "=")
		if equalIndex < 0 {
			return nil, fmt.Errorf("invalid cookie pair: %s", pair)
		}

		name := pair[:equalIndex]
		value := pair[equalIndex+1:]
		cookies[name] = value
	}

	return cookies, nil
}

// SAS URI Code:
func identifySASURIType(queryParams url.Values) string {
	if _, ok := queryParams["ss"]; ok {
		return "Account SAS URI"
	}
	if _, ok := queryParams["sktid"]; ok {
		return "Delegation SAS URI"
	}
	return "Service SAS URI"
}

func getLongForm(sasType, key, value string) (string, string) {
	switch sasType {
	case "Account SAS URI":
		return getLongFormAccountSAS(key, value)
	case "Service SAS URI":
		return getLongFormServiceSAS(key, value)
	case "Delegation SAS URI":
		return getLongFormDelegationSAS(key, value)
	default:
		return key, value // If the SAS type is unknown, return the original key and value
	}
}

func getLongFormAccountSAS(key, value string) (string, string) {
	// Mapping from short form to long form for keys
	keyMap := map[string]string{
		"sv":  "signedVersion",
		"ss":  "signedServices",
		"srt": "signedResourceTypes",
		"sp":  "signedPermissions",
		"st":  "signedStart Time",
		"se":  "signedExpiry Time",
		"sip": "signedIP",
		"spr": "signedProtocol",
		"ses": "signedEncriptionScope",
		"sig": "signature",
	}

	// Mapping from short form to long form for values
	valueMap := map[string]string{
		"r": "Read",
		"w": "Write",
		"d": "Delete",
		"y": "Permanent Delete",
		"l": "List",
		"a": "Add",
		"c": "Create or Container[srt]",
		"u": "Update",
		"p": "Process",
		"t": "Tag or Table[ss]",
		"f": "Filter or File[ss]",
		"i": "Set Immutability Policy",
		"b": "Blob",
		"q": "Queue",
		"s": "Service",
		"o": "Object",
	}

	longFormKey, ok := keyMap[key]
	if !ok {
		longFormKey = key // If no mapping found, use the original key
	}

	longFormValue := ""
	//for _, char := range value {
	//	longForm, ok := valueMap[string(char)]
	//	if ok {
	//		if longFormValue != "" {
	//			longFormValue += ", "
	//		}
	//		longFormValue += longForm
	//	}
	//}
	//if longFormValue == "" {
	//	longFormValue = value // If no mapping found, use the original value
	//}

	switch key {
	case "ss", "srt", "sp":
		longFormValue = getCombinedLongForm(value, valueMap)
	default:
		longFormValue, ok = valueMap[value]
		if !ok {
			longFormValue = value // If no mapping found, use the original value
		}
	}

	return longFormKey, longFormValue
}

func getLongFormServiceSAS(key, value string) (string, string) {
	// Mapping from short form to long form for keys
	keyMap := map[string]string{
		"sv":  "signedVersion",
		"ss":  "signedServices",
		"srt": "signedResourceTypes",
		"sp":  "signedPermissions",
		"st":  "signedStart Time",
		"se":  "signedExpiry Time",
		"sip": "signedIP",
		"spr": "signedProtocol",
		"ses": "signedEncriptionScope",
		"sig": "signature",
	}

	// Mapping from short form to long form for values
	valueMap := map[string]string{
		"r": "Read",
		"w": "Write",
		"d": "Delete",
		"y": "Permanent Delete",
		"l": "List",
		"a": "Add",
		"c": "Create or Container[srt]",
		"u": "Update",
		"p": "Process",
		"t": "Tag or Table[ss]",
		"f": "Filter or File[ss]",
		"i": "Set Immutability Policy",
		"b": "Blob",
		"q": "Queue",
		"s": "Service",
		"o": "Object",
	}

	longFormKey, ok := keyMap[key]
	if !ok {
		longFormKey = key // If no mapping found, use the original key
	}

	longFormValue := ""
	for _, char := range value {
		longForm, ok := valueMap[string(char)]
		if ok {
			if longFormValue != "" {
				longFormValue += ", "
			}
			longFormValue += longForm
		}
	}
	if longFormValue == "" {
		longFormValue = value // If no mapping found, use the original value
	}

	return longFormKey, longFormValue
}

func getLongFormDelegationSAS(key, value string) (string, string) {
	// Mapping from short form to long form for keys
	keyMap := map[string]string{
		"sv":  "signedVersion",
		"ss":  "signedServices",
		"srt": "signedResourceTypes",
		"sp":  "signedPermissions",
		"st":  "signedStart Time",
		"se":  "signedExpiry Time",
		"sip": "signedIP",
		"spr": "signedProtocol",
		"ses": "signedEncriptionScope",
		"sig": "signature",
	}

	// Mapping from short form to long form for values
	valueMap := map[string]string{
		"r": "Read",
		"w": "Write",
		"d": "Delete",
		"y": "Permanent Delete",
		"l": "List",
		"a": "Add",
		"c": "Create or Container[srt]",
		"u": "Update",
		"p": "Process",
		"t": "Tag or Table[ss]",
		"f": "Filter or File[ss]",
		"i": "Set Immutability Policy",
		"b": "Blob",
		"q": "Queue",
		"s": "Service",
		"o": "Object",
	}

	longFormKey, ok := keyMap[key]
	if !ok {
		longFormKey = key // If no mapping found, use the original key
	}

	longFormValue := ""
	for _, char := range value {
		longForm, ok := valueMap[string(char)]
		if ok {
			if longFormValue != "" {
				longFormValue += ", "
			}
			longFormValue += longForm
		}
	}
	if longFormValue == "" {
		longFormValue = value // If no mapping found, use the original value
	}

	return longFormKey, longFormValue
}

func getCombinedLongForm(value string, valueMap map[string]string) string {
	longFormValue := ""
	for _, char := range value {
		longForm, ok := valueMap[string(char)]
		if ok {
			if longFormValue != "" {
				longFormValue += ", "
			}
			longFormValue += longForm
		}
	}
	if longFormValue == "" {
		longFormValue = value // If no mapping found, use the original value
	}
	return longFormValue
}
