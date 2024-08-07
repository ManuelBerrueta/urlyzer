# urlyzer
**`urlyzer`** is a tool to help with analyzing URLs, allowing you to quickly dissect and easily visualize all the parameters offline to keep your URL and Params confidential. 

If you find it useful, consider giving `urlyzer` a ⭐ star 👆 😉

### Why?
Frequently, you may encounter lengthy URLs containing various parameters. If you are trying to analyze these parameters in the query string or in a fragment, you may paste it into a text/code editor and break them apart or search for a certain parameter and value. Also, you may often have to decode some of it (if it is URL encoded), which means you have to also copy and paste that to another terminal tool or some other tool like Burp’s Decoder/CyberChef.

> It turns out that also URL Fragments (`#`) may also contain query strings. One such case is with **OAuth 2.0** + **OpenID Connect (OIDC)**, and `urlyzer` also parses these!

If you are frustrated by this process, often repeated multiple times a day, then this 🛠 tool is for you!

---    
### Features
- Offline parsing to keep your URL and Query String Parameters confidential.
- Allows you to pipe-in (`|`) the URL from the command line via `stdin`.
- Check the final destination (`-f`) after redirects of a URL.
- **Proxy** (`-p`) to allow to forward proxy the traffic for further inspection to Burp (or your proxy of choice)  . 
- Check Azure Blob Storage SAS URIs (`sas`) for type & parse the parameters with details of what each of those parameters mean in long form.
  -  More support coming to this soon!
- Parse cookies 🍪 with `-c`.
- Modify query string key:value pairs with `-qr`
- Extract just the query string parameters you are interested in with `qs`

---    
## Example use cases:
### Running the tool as a script:
```shell
go run urlyzer.go "https://www.example.com/path/towin?param1=value1&param2=value%202"
```
### Running the tool from one of the release binaries:
**Simple Tool Use:**
```shell
urlyzer "https://www.example.com/path/towin?param1=value1&param2=value%202#MyFragment"
```
**Piped Input:**
```shell
echo "https://www.example.com/path/towin?param1=value1&param2=value%202#MyFragment" | urlyzer 
```

### Example Output
```yaml
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

### Parsing cookies
Sometimes looking through cookies in the requests can be messy, especially when there are a bunch and really long ones.
You can use urlyzer like this to parse them:
```shell
urlyzer -c "key1=value1; key2=value2; key3=value3"
```
**🔥 TIP:** To copy the cookies from the browser, use the Developer Tools Console and run this:
```js
let cookies = document.cookie;
copy(cookies);
```

### Modify query string key:value pairs
With this example, we can change the `response_type` parameter value to `token` instead of `code`:
```shell
urlyzer -qr response_type=token "https://www.example.com/path/towin?param1=value1&response_type=code#MyFragment"
```
You can also modify more than one parameter by using a comma && as a bonus you can verify the modified URL by parsing it with a pipe through urlyzer again:
```shell
urlyzer -qr response_type=token,param1=test "https://www.example.com/path/towin?param1=value1&response_type=code#MyFragment" | urlyzer
```

### Extract just the query string parameter(s) you are interested in analyzing:
If you have a long URL with lots of parameters, you might just be interested in the values of one or two of those. 
You can query the values of just the parameters you are interested in:
```shell
urlyzer -qs response_type "https://www.example.com/path/towin?param1=value1&response_type=code#MyFragment"
```
Query for more than one parameter:
```shell
urlyzer -qs response_type,param1 "https://www.example.com/path/towin?param1=value1&response_type=code#MyFragment"
```

### Running final destination analysis of a URL
This is a case where the URL may have one or multiple redirects and you wish to check it's final destination. 
To do the analysis, pass in the `-f` flag after `urlyzer`:
```shell
urlyzer -f "https://aka.ms/powershell-release?tag=lts"
```
**Example Output**
```yaml
Final Destination: https://github.com/Powershell/Powershell/releases/tag/v7.2.13
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
If you want to forward proxy the traffic of the final destination analysis for further investigations you can use the `-p` flag followed by the proxy address:
```Terminal
echo " https://aka.ms/powershell-release?tag=lts" | ./urlyzer -f -p "http://127.0.0.1:8080"
```

---    
# Releases
Check out the releases link on the right to download one of the binaries for your targeted architecture.


---    
# Running in a container
Running in a container is easy with the provided Dockerfile and it doesn't any other dependencies, well other than having Docker!
1. Clone the code
2. Build the image by running `docker build --tag urlyzer .`
3. Run the container with `urlyzer`: `docker run urlyzer [flag] "URL"`

## Examples
### Running regular `urlyzer`
`docker run urlyzer "https://login.microsoftonline.com/common/oauth2/authorize?client_id=1234#revx0r.com"`

### Checking final destination/redirects
`docker run urlyzer -f "https://aka.ms/powershell-release?tag=lts"`
