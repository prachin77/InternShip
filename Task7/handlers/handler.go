package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

// /InternShip/Task7/

type GeoIP struct {
	Ip          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	Zipcode     string  `json:"zipcode"`
	Lat         float64 `json:"latitude"`
	Lon         float64 `json:"longitude"`
}

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
