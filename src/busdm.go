package main

import (
	"fmt"
    "encoding/json"
 //   "log"
    "net/http"
    "bytes"
	"io/ioutil"
    "os"
    "io"
    "crypto/tls"
    "strconv"
)



// Constant
const cred_file string = "credentials.txt"


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
    ClientID string
    Password string
}



// Opens credentials filename and stores its values.
func readCredentials(filename string) st_credentials{
    var cred st_credentials;

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
   n,err := fmt.Fscanf(fi, "%s %s \n", &cred.ClientID, &cred.Password)
   // Error line:
   if (err != nil && err != io.EOF ) ||(n!=2){
            panic(err)
        }
    return cred;
}

// Get EMT times

type GetStopTimeRequest struct {
    idClient string
    passKey string
    cultureInfo string
    idStop string
}

func getStopTime( IdStop int) string{


    creds := readCredentials(cred_file)
    url := "https://openbus.emtmadrid.es:9443/emt-proxy-server/last/geo/GetArriveStop.php"
    fmt.Println("URL:>", url)
    jsonreq := GetStopTimeRequest{creds.ClientID,
                              creds.Password,
                              "ES",
                              strconv.Itoa(IdStop)}
    fmt.Println(jsonreq)
    jsonStr, err := json.Marshal(jsonreq)

    fmt.Printf("JSON:> %s\n", jsonStr)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("X-Custom-Header", "busdm")
    req.Header.Set("Content-Type", "application/json")



    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
    }
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
    return string(body)
}



// Response handler
func handler(w http.ResponseWriter, r *http.Request) {
    body := getStopTime(608);
    fmt.Fprintf(w, "Hi there, I love %s!", body)
    //getStopArrivalTime (608);
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}

