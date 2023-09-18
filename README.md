# urlyzer
**`urlyzer`** is a tool to help with analyzing URLs, allowing you to quickly dissect and easily visualize all the parameters offline to keep your URL and Params confidential.

### Why?
There is very often times when you have a really long URL with all kinds of different parameters. When you are trying to analyze the parameters in the query string you may paste it into a text/code editor and break them apart or search for a certain parameter and value. Also, you may often have to decode some of it (if it is URL encoded), which means you have to also copy and paste that to another terminal tool or some other tool like Burp’s Decoder/CyberChef. 

If you are frustrated by this process, often repeated multiple times a day, then this tool is for you!

---    
### Features
- Offline parsing to keep your URL and Params confidential
- Allows you to pipe-in the url from the command line via stdin
- [*New Feature] Check the final destination (after redirects) of a URL
- [*New Feature] **Proxy** (`-p`) to allow to forward proxy the traffic for further inspection to Burp (or your proxy of choice) 

---    
## Example use cases:
### Running the tool as a script:
```Shell
go run urlyzer.go "https://www.example.com/path/towin?param1=value1&param2=value%202"
```
### Running the tool from one of the binaries:
**Simple Tool Use:**
```Shell
urlyzer "https://www.example.com/path/towin?param1=value1&param2=value%202#MyFragment"
```
**Piped Input:**
```Shell
echo "https://www.example.com/path/towin?param1=value1&param2=value%202#MyFragment" | urlyzer 
```

### Example Output
```YAML
Scheme: https
Host: www.example.com
Port: 
Path: /path/towin
Query String: param1=value1&param2=value%202
Query Parameters:
  param2: value 2
  param1: value1
Fragment: MyFragment
```

### Running final destination analysis of a URL
This is a case where the URL may have one or multiple redirects and you wish to check it's final destination. 
To do the analysis, pass in the `-f` flag after `urlyzer`:
```Shell
urlyzer -f "https://aka.ms/powershell-release?tag=lts"
```
**Example Output**
```YAML
Final Destination: https://github.com/PowerShell/PowerShell/releases/tag/v7.2.13
Status Code: 200
Headers:
  Location: [https://powershell-release-redirect.azurewebsites.net/api/powershell-release-redirect?code=GyM2EB/tN8KmH2swJp/o/TdJ72z9CLP2/g3lBa9gnPWDZOAbDrD5EA==&tag=lts]
  X-Response-Cache-Status: [True]
  Pragma: [no-cache]
  Cache-Control: [max-age=0, no-cache, no-store]
  Date: [Sun, 24 Sep 2023 20:21:23 GMT]
  Strict-Transport-Security: [max-age=31536000 ; includeSubDomains]
  Content-Length: [0]
  Server: [Kestrel]
  Request-Context: [appId=cid-v1:9b037ab9-fa5a-4c09-81bd-41ffa859f01e]
  Expires: [Sun, 24 Sep 2023 20:21:23 GMT]
```
#### Using the Proxy
If you want to forward proxy the traffic from the analysis for further analysis you can use the `-p` flag follwed by the proxy address:
```Terminal
echo " https://aka.ms/powershell-release?tag=lts" | ./urlyzer -f -p "http://127.0.0.1:8080"
```

---    
# Releases
Check out the releases tab to download one of the binaries for your targeted architecture.