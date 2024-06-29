package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prachin77/task7/handlers"
)

func main() {
	r := mux.NewRouter()

	fmt.Println("listening on port 8080")

	r.HandleFunc("/",handlers.DefaultRoute).Methods("GET")
	r.HandleFunc("/app",handlers.GetApp).Methods("GET")

	r.HandleFunc("/getipadd",handlers.GetIpAdd).Methods("GET")
	r.HandleFunc("/getpingform",handlers.GetPingForm).Methods("GET")
	r.HandleFunc("/ping",handlers.PingToUser).Methods("POST")
	r.HandleFunc("/getnslookupform",handlers.GetNsLookUpForm).Methods("GET")
	r.HandleFunc("/nslookup",handlers.NsLookUp).Methods("POST")
	r.HandleFunc("/getcountryform",handlers.GetCountryForm).Methods("GET")
	r.HandleFunc("/getcountryname",handlers.GetCountryName).Methods("POST")
	r.HandleFunc("/getwhoisform",handlers.GetWhoisForm).Methods("GET")
	r.HandleFunc("/getwhoisinfo",handlers.GetWhoisInfo).Methods("POST")
	r.HandleFunc("/getportcheckerform",handlers.GetPortCheckerForm).Methods("GET")
	r.HandleFunc("/checkportstatus",handlers.CheckPortStatus).Methods("POST")
	r.HandleFunc("/getproxycheckform",handlers.GetProxyCheckerForm).Methods("GET")
	r.HandleFunc("/checkproxy",handlers.CheckProxyStatus).Methods("POST")
	r.HandleFunc("/getrevlookupform",handlers.GetRevLookupForm).Methods("GET")
	r.HandleFunc("/reverselookup",handlers.ReverseLookup).Methods("POST")
	r.HandleFunc("/gettracerouteform",handlers.GetTraceRouteForm).Methods("GET")
	r.HandleFunc("/traceroute",handlers.TraceRoute).Methods("POST")
	r.HandleFunc("/getunitconverterform",handlers.	GetUnitConverterForm).Methods("GET")
	r.HandleFunc("/convertunit",handlers.UnitConverter).Methods("POST")
	r.HandleFunc("/getbmform",handlers.GetBandwidthMeterForm).Methods("GET")
	r.HandleFunc("/checkbandwidthmeter",handlers.CheckBandwidth).Methods("POST")
	r.HandleFunc("/getncform",handlers.GetNetworkCalculatorForm).Methods("GET")
	r.HandleFunc("/calculatenetwork",handlers.CalculateNetwork).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080",r))
}
