# urlyzer
**`urlyzer`** is a tool to help with analyzing URLs, allowing you to quickly dissect and easily visualize all the parameters offline to keep your URL and Params confidential.

### Why?
There is very often times when you have a really long URL with all kinds of different parameters. When you are trying to analyze the parameters in the query string you may paste it into a text/code editor and break them apart or search for a certain parameter and value. Also, you may often have to decode some of it (if it is URL encoded), which means you have to also copy and paste that to another terminal tool or some other tool like Burp’s Decoder/CyberChef. 

If you are frustrated by this process, often repeated multiple times a day, then this tool is for you!

### Features
- Offline parsing to keep your URL and Params confidential
- Allows you to pipe-in the url from the command line via stdin


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

# Releases
Check out the releases tab to download one of the binaries for your targeted architecture