package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/showwin/speedtest-go/speedtest"
)

// /InternShip/Task7/templates

func DefaultRoute(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("main.html"))
	tmpl.Execute(w, nil)
}

func GetApp(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/homepage.html"))
	tmpl.Execute(w, nil)
}

func GetIpAdd(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("ipconfig")
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Display the IP address in the UI
	fmt.Fprintf(w, "<pre>%s</pre>", output)
}

func GetPingForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/pingform.html"))
	tmpl.Execute(w, nil)
}

func PingToUser(w http.ResponseWriter, r *http.Request) {
	userIpAdd := r.PostFormValue("ipadd")
	fmt.Println("ip address of user : ", userIpAdd)

	cmd := exec.Command("ipconfig")
	// cmd := exec.Command("C:\\Windows\\System32\\ipconfig.exe")

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract IPv4 address from the output using regular expressions
	ipRegex := regexp.MustCompile(`IPv4 Address[\.\s]*:\s*([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)
	matches := ipRegex.FindStringSubmatch(string(output))

	var ipAddress string
	if len(matches) > 1 {
		ipAddress = matches[1]
	} else {
		ipAddress = "IP Address not found"
	}
	fmt.Println("ipv4 address : ", ipAddress)
	pingcmd := exec.Command("C://Windows//System32//ping.exe", userIpAdd)

	pingOutput, pingErr := pingcmd.Output()
	if pingErr != nil {
		log.Fatal(pingErr)
		return
	}
	fmt.Fprintf(w, "<pre>%s</pre>", userIpAdd)
	fmt.Fprintf(w, "<pre>%s</pre>", pingOutput)
}

func GetNsLookUpForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/nslookupform.html"))
	tmpl.Execute(w, nil)
}

func NsLookUp(w http.ResponseWriter, r *http.Request) {
	// Retrieve the value from the POST form parameter "nslookup"
	userNslookupValue := r.PostFormValue("nslookup")
	fmt.Println("User value for nslookup:", userNslookupValue)

	// Execute nslookup command
	nslookupcmd := exec.Command("nslookup", userNslookupValue)

	// Run the command and capture the output
	nslookupOutput, err := nslookupcmd.Output()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Write the original nslookup value to the HTTP response
	fmt.Fprintf(w, "<pre>User value for nslookup: %s</pre>", userNslookupValue)

	// Extract domain name from nslookup output
	lines := strings.Split(string(nslookupOutput), "\n")
	var domain string
	for _, line := range lines {
		if strings.Contains(line, "Name:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				domain = parts[1]
				break
			}
		}
	}

	fmt.Fprintf(w, "<pre>Domain name = %s</pre>", domain, "\n")
	fmt.Fprintf(w, "<pre>nslookup cmd output = %s</pre>", nslookupOutput)
}

func GetCountryForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/getcountryform.html"))
	tmpl.Execute(w, nil)
}

func GetCountryName(w http.ResponseWriter, r *http.Request) {
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Cookie", "PHPSESSID=5953dfa5f54c8e3ba5b0257d4b074592; LG=en; __gads=ID=7279f709b0cd112b:T=1719226442:RT=1719226442:S=ALNI_Mb0ro8J4kGv6TmZjuMxZ0JYvAybYg; __gpi=UID=00000e61d6db7ac6:T=1719226442:RT=1719226442:S=ALNI_MZYPrG36HixPDlouBttqaiR1FLnbA; __eoi=ID=f4272b7168b54539:T=1719226442:RT=1719226442:S=AA-AfjYDPxslAoYgqkcF4WUsACic")
	// req.Header.Set("Cache-Control", "max-age=0")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	// req.Header.Set("Origin", "https://ping.eu")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.199 Safari/537.36")
	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,/;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// req.Header.Set("Sec-Fetch-Site", "same-origin")
	// req.Header.Set("Sec-Fetch-Mode", "navigate")
	// req.Header.Set("Sec-Fetch-User", "?1")
	// req.Header.Set("Sec-Fetch-Dest", "iframe")
	// req.Header.Set("Referer", "https://ping.eu/country-by-ip/")
	// req.Header.Set("Accept-Encoding", "gzip, deflate")
	// req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	// req.Header.Set("Connection", "close")

	userValue := r.PostFormValue("countryname")
	fmt.Println("user value : ", userValue)

	var host string
	fmt.Print("Enter the host name : ")
	fmt.Scanln(&host)

	data := []byte(fmt.Sprintf("host=%s&go=Go", host))

	req, err := http.NewRequest("POST", "https://ping.eu/action.php?atype=12", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.199 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body Length:", len(body))

	// Extract content within <STRONG> tags
	re := regexp.MustCompile(`<STRONG>(.*?)</STRONG>`)
	matches := re.FindAllStringSubmatch(string(body), -1)
	for _, match := range matches {
		fmt.Println("Strong tag content:", match[1])
		fmt.Fprintf(w, "<pre>%s</pre>", match[1])
	}
}

func GetWhoisForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/getwhoisform.html"))
	tmpl.Execute(w, nil)
}

func GetWhoisInfo(w http.ResponseWriter, r *http.Request) {
	userWhoisValue := r.PostFormValue("whois")
	fmt.Println("user value : ", userWhoisValue)

	// Execute whois command with userWhoisValue as argument
	// cmd := exec.Command("whois", userWhoisValue)
	// "C:\Users\Prachin\Downloads\SysinternalsSuite\whois.exe"
	cmd := exec.Command("C:/Users/Prachin/Downloads/SysinternalsSuite/whois.exe", userWhoisValue)

	// Run the command and capture its output
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running whois command:", err)
		http.Error(w, "Failed to fetch whois information", http.StatusInternalServerError)
		return
	}
	fmt.Println("output : ", string(output))
	fmt.Fprintf(w, "<pre>%s</pre>", string(output))
}

func GetPortCheckerForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/getportcheckerform.html"))
	tmpl.Execute(w, nil)
}

func CheckPortStatus(w http.ResponseWriter, r *http.Request) {
	portValue := r.PostFormValue("checkport")
	hostValue := r.PostFormValue("checkhost")
	fmt.Println("User host value:", hostValue)
	fmt.Println("User port value:", portValue)

	// Check if portValue is a valid port number
	_, err := net.LookupPort("tcp", portValue)
	if err != nil {
		fmt.Println("Invalid port number:", err)
		http.Error(w, "Invalid port number", http.StatusBadRequest)
		return
	}

	// Construct the address for the host and port
	address := fmt.Sprintf("%s:%s", hostValue, portValue)

	// Attempt to connect to the address
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Port is closed or unreachable:", err)
		fmt.Fprintf(w, "Port %s on host %s is closed or unreachable\n", portValue, hostValue)
		return
	}
	defer conn.Close()

	// Write success message to response
	fmt.Fprintf(w, "Port %s on host %s is open\n", portValue, hostValue)
}

func GetProxyCheckerForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/getproxycheckerform.html"))
	tmpl.Execute(w, nil)
}

func CheckProxyStatus(w http.ResponseWriter, r *http.Request) {
	hostValue := r.PostFormValue("checkhost")
	portValue := r.PostFormValue("checkport")
	fmt.Println("user host value : ", hostValue)
	fmt.Println("user port value : ", portValue)

	port, err := strconv.Atoi(portValue)
	if err != nil {
		http.Error(w, "Invalid port number", http.StatusBadRequest)
		return
	}

	// Perform a TCP dial to check the proxy status with increased timeout
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", hostValue, port), 10*time.Second)
	if err != nil {
		fmt.Printf("Proxy status check failed: %v\n", err)
		http.Error(w, fmt.Sprintf("Proxy status check failed: %v", err), http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("Proxy status check failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	fmt.Printf("Proxy at %s:%d is reachable\n", hostValue, port)
	fmt.Fprintf(w, "Proxy at %s:%d is reachable", hostValue, port)
}

func GetRevLookupForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/getrevlookupform.html"))
	tmpl.Execute(w, nil)
}

func ReverseLookup(w http.ResponseWriter, r *http.Request) {
	userIP := r.PostFormValue("ipadd")
	fmt.Println("User IP address:", userIP)

	// Perform a reverse DNS lookup
	names, err := net.LookupAddr(userIP)
	if err != nil {
		fmt.Printf("Reverse lookup failed: %v\n", err)
		http.Error(w, fmt.Sprintf("Reverse lookup failed: %v", err), http.StatusInternalServerError)
		fmt.Fprintf(w, "Reverse lookup failed ")
		return
	}

	// Print the results to console (for debugging)
	fmt.Printf("Reverse lookup results for %s:\n", userIP)
	for _, name := range names {
		fmt.Println(name)
	}

	// Send the results back to the client
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Reverse lookup results for %s:\n", userIP)
	for _, name := range names {
		fmt.Fprintf(w, "%s\n", name)
	}
}

func GetTraceRouteForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/gettracerouteform.html"))
	tmpl.Execute(w, nil)
}

func TraceRoute(w http.ResponseWriter, r *http.Request) {
	userValue := r.PostFormValue("traceroute")
	fmt.Println("user value for trace route:", userValue)

	if userValue == "" {
		http.Error(w, "Empty traceroute value", http.StatusBadRequest)
		return
	}

	// Execute the traceroute command
	cmd := exec.Command("tracert", userValue)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running traceroute: %v\n", err)
		http.Error(w, fmt.Sprintf("Error running traceroute: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("Traceroute output:")
	fmt.Println(string(output))

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "<pre>%s</pre>", string(output))
}

func GetUnitConverterForm(w http.ResponseWriter, r *http.Request) {
	// tmpl  := template.Must(template.ParseFiles("/InternShip/Task7/templates/getunitconverterform.html"))
	tmpl, err := template.ParseFiles("/InternShip/Task7/templates/getunitconverterform.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, nil)
}

func UnitConverter(w http.ResponseWriter, r *http.Request) {
	// Get values from the form
	digitValueInt := r.PostFormValue("digitvalue") // Assuming you also have an input for the numeric value
	fromValue := r.PostFormValue("from")
	toValue := r.PostFormValue("to")

	fmt.Println("Digit value : ", digitValueInt)
	fmt.Println("Unit to be converted : ", fromValue)
	fmt.Println("Designated unit : ", toValue)

	// Convert digitValue to float64
	digitValue, err := strconv.ParseFloat(digitValueInt, 64)
	if err != nil {
		http.Error(w, "Invalid digit value", http.StatusBadRequest)
		return
	}

	// Perform unit conversion based on fromValue and toValue
	var result float64
	switch fromValue {
	case "meters":
		switch toValue {
		case "feet":
			result = digitValue * 3.28084
		case "inches":
			result = digitValue * 39.3701
		case "kilometers":
			result = digitValue / 1000
		case "centimeters":
			result = digitValue * 100
		case "yards":
			result = digitValue * 1.09361
		case "meters":
			result = digitValue
		}
	case "feet":
		switch toValue {
		case "meters":
			result = digitValue / 3.28084
		case "inches":
			result = digitValue * 12
		case "kilometers":
			result = digitValue / 3280.84
		case "centimeters":
			result = digitValue * 30.48
		case "yards":
			result = digitValue / 3
		case "feet":
			result = digitValue
		}
	case "inches":
		switch toValue {
		case "meters":
			result = digitValue / 39.3701
		case "feet":
			result = digitValue / 12
		case "kilometers":
			result = digitValue / 39370.1
		case "centimeters":
			result = digitValue * 2.54
		case "yards":
			result = digitValue / 36
		case "inches":
			result = digitValue
		}
	case "kilometers":
		switch toValue {
		case "meters":
			result = digitValue * 1000
		case "feet":
			result = digitValue * 3280.84
		case "inches":
			result = digitValue * 39370.1
		case "centimeters":
			result = digitValue * 100000
		case "yards":
			result = digitValue * 1093.61
		case "kilometers":
			result = digitValue
		}
	case "centimeters":
		switch toValue {
		case "meters":
			result = digitValue / 100
		case "feet":
			result = digitValue / 30.48
		case "inches":
			result = digitValue / 2.54
		case "kilometers":
			result = digitValue / 100000
		case "yards":
			result = digitValue / 91.44
		case "centimeters":
			result = digitValue
		}
	case "yards":
		switch toValue {
		case "meters":
			result = digitValue / 1.09361
		case "feet":
			result = digitValue * 3
		case "inches":
			result = digitValue * 36
		case "kilometers":
			result = digitValue / 1093.61
		case "centimeters":
			result = digitValue * 91.44
		case "yards":
			result = digitValue
		}
	default:
		http.Error(w, "Unsupported 'from' unit", http.StatusBadRequest)
		return
	}

	fmt.Println("Result : ", result)
	fmt.Fprintf(w, "<pre>%f</pre>", result)
}

func GetBandwidthMeterForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/getbandwidthmeterform.html"))
	tmpl.Execute(w, nil)
}

func CheckBandwidth(w http.ResponseWriter, r *http.Request) {
	var speedtestClient = speedtest.New()

	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		// Please make sure your host can access this test server,
		// otherwise you will get an error.
		// It is recommended to replace a server at this time
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()
		// Note: The unit of s.DLSpeed, s.ULSpeed is bytes per second, this is a float64.
		fmt.Printf("Latency: %s, Download: %s, Upload: %s\n", s.Latency, s.DLSpeed, s.ULSpeed)
		fmt.Fprintf(w, "<pre>Latency: %s, Download: %s, Upload: %s\n</pre>", s.Latency, s.DLSpeed, s.ULSpeed)
		s.Context.Reset() // reset counter
	}
}

func GetNetworkCalculatorForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("/InternShip/Task7/templates/getnetworkcalcform.html"))
	tmpl.Execute(w, nil)
}

func CalculateNetwork(w http.ResponseWriter, r *http.Request) {
	userIP := r.PostFormValue("ipadd")
	userMask := r.PostFormValue("mask")
	fmt.Println("User IP Address:", userIP)
	fmt.Println("User Subnet Mask:", userMask)

	// Validate IP address
	ip := net.ParseIP(userIP)
	if ip == nil {
		http.Error(w, "Invalid IP address", http.StatusBadRequest)
		return
	}

	// Check for IPv4 addresses
	if ip.To4() == nil {
		http.Error(w, "IPv6 addresses not currently supported", http.StatusBadRequest)
		return
	}

	// Validate and convert subnet mask
	var ones int
	var err error
	if userMask != "" {
		ones, err = strconv.Atoi(userMask)
		if err != nil || ones < 0 || ones > 32 {
			http.Error(w, "Invalid subnet mask format", http.StatusBadRequest)
			return
		}
	} else {
		ones, _ = ip.DefaultMask().Size() // Extract only the first value (ones)
	}

	// Calculate subnet using IPSubnet function
	originalNetwork := &net.IPNet{
		IP:   ip.Mask(net.CIDRMask(ones, 32)),
		Mask: net.CIDRMask(ones, 32),
	}

	newNet := IPSubnet(originalNetwork, ones, 1) // Example offset of 1

	if newNet == nil {
		fmt.Fprintf(w, "Unable to calculate subnet with size %d and offset 1\n", ones)
		return
	}

	// Format the output similar to ping.eu
	subnetRange := fmt.Sprintf(`
		Address:         %s
		Netmask:         %s = %d
		Wildcard:        %s
		=>
		Network:         %s/%d
		HostMin:         %s
		HostMax:         %s
		Broadcast:       %s
		`,
		userIP,
		net.CIDRMask(ones, 32).String(), ones,
		net.IP(net.CIDRMask(ones, 32)).String(),
		newNet.IP.String(),
		ones,
		newNet.IP.String(),
		newNet.IP.String(),
		net.IP(net.IPv4(0xff, 0xff, 0xff, 0xff).To4()).String(), // Broadcast address for IPv4
	)

	// Write response
	fmt.Fprintf(w, "Subnet Range for %s/%d:\n%s", ip.String(), ones, subnetRange)
}

// IPSubnet calculates a subnet from the given network with the specified size and offset.
func IPSubnet(network *net.IPNet, size int, offset int) *net.IPNet {
	var maskLen int
	ip := network.IP

	if IsIPv4(ip) {
		maskLen = net.IPv4len * 8
		addr := int(ipToI32(ip.To4()))
		ip = i32ToIP(int32(addr + offset*(maskLen-size+1)))
	} else {
		maskLen = net.IPv6len * 8
		a := ipToU64(ip[:net.IPv6len/2])
		b := ipToU64(ip[net.IPv6len/2:])

		if size > 64 {
			b = b + uint64(offset)<<uint64(128-size)
		} else {
			a = a + uint64(offset)<<uint64(64-size)
		}

		ip = make(net.IP, net.IPv6len)
		u64ToIP(ip[:net.IPv6len/2], a)
		u64ToIP(ip[net.IPv6len/2:], b)
	}

	if !network.Contains(ip) {
		return nil
	}

	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(size, maskLen),
	}
}


// IsIPv4 returns true if ip is IPv4 address.
func IsIPv4(ip net.IP) bool {
	return len(ip) == net.IPv4len || (isZeros(ip[0:10]) && ip[10] == 0xff && ip[11] == 0xff)
}

func isZeros(p net.IP) bool {
	for _, b := range p {
		if b != 0 {
			return false
		}
	}
	return true
}

func ipToI32(ip net.IP) uint32 {
	return (uint32(ip[0]) << 24) | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3])
}

func i32ToIP(addr int32) net.IP {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(addr >> 24 & 0xFF)
	ip[1] = byte(addr >> 16 & 0xFF)
	ip[2] = byte(addr >> 8 & 0xFF)
	ip[3] = byte(addr & 0xFF)
	return ip
}

func ipToU64(ip net.IP) uint64 {
	return (uint64(ip[0]) << 56) | (uint64(ip[1]) << 48) | (uint64(ip[2]) << 40) | (uint64(ip[3]) << 32) |
		(uint64(ip[4]) << 24) | (uint64(ip[5]) << 16) | (uint64(ip[6]) << 8) | uint64(ip[7])
}

func u64ToIP(ip []byte, addr uint64) {
	ip[0] = byte(addr >> 56 & 0xFF)
	ip[1] = byte(addr >> 48 & 0xFF)
	ip[2] = byte(addr >> 40 & 0xFF)
	ip[3] = byte(addr >> 32 & 0xFF)
	ip[4] = byte(addr >> 24 & 0xFF)
	ip[5] = byte(addr >> 16 & 0xFF)
	ip[6] = byte(addr >> 8 & 0xFF)
	ip[7] = byte(addr & 0xFF)
}