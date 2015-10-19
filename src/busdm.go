package main

import (
	"fmt"
    // "encoding/json"
 //   "log"
    "net/http"
    "bytes"
	"io/ioutil"
)



// Constant
cred_file string = "credentials.txt"


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



type st_walkingdistance struct {
    StopId int
    SecondsWalking int
}

type st_credentials struct {
    ClientID String
    Password String
}  st_credentials;



// Opens credentials filename and stores its values.
func readCredentials{

	

    fmt.Println ("Reading Credentials", filename)
    // open input file
    fi, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    // close fi on exit and check for its returned error
    defer func() {
        if err := fi.Close(); err != nil {
            panic(err)
        }
    }()


   // read Credentials
   n,err := fmt.Fscanf(fi, "%d %d \n", &cred.ClientID, &cred.Password) 
   // Error line:
   if (err != nil && err != io.EOF ) ||(n!=2){
            panic(err)            
        }
   
   

}

// Get EMT times
func getStopTime( IdStop int){


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

