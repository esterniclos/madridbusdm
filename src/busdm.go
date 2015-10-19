package main

import (
	"fmt"
    // "encoding/json"
 //   "log"
    "net/http"
    "bytes"
	"io/ioutil"
)




// Struct to store EMT  response 
type jsonEmt struct {
	Arrives []struct {
		BusDistance     int     `json:"busDistance"`
		BusID           string  `json:"busId"`
		BusPositionType int     `json:"busPositionType"`
		BusTimeLeft     int     `json:"busTimeLeft"`
		Destination     string  `json:"destination"`
		IsHead          string  `json:"isHead"`
		Latitude        float64 `json:"latitude"`
		LineID          string  `json:"lineId"`
		Longitude       float64 `json:"longitude"`
		StopID          int     `json:"stopId"`
	} `json:"arrives"`
}



type walkingdistance struct {
    StopId int
    SecondsWalking int
}

// Get EMT times
func getStopArrivalTime( IdStop int){


    var url string ; 
    url = "https://openbus.emtmadrid.es:9443/emt-proxy-server/last/geo/GetArriveStop.php"
  
  
    fmt.Println("URL:>", url)

     var jsonStr = []byte(`
{
    "idClient": "to be read", 
    "passKey": "to be read",
    "idStop": "608"

}`)


    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "busdm")
    req.Header.Set("Content-Type", "application/json")



    cfg := &tls.Config{InsecureSkipVerify: true}
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
    
}



// Response handler
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
    getStopArrivalTime (608);
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}

