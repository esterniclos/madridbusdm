package main

import (
	"fmt"
    "encoding/json"
    "net/http"
    "net/url"
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
type GetArriveStopResponse struct {
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

func emtRequest(path string, req_values url.Values) string {
    creds := readCredentials(cred_file)
    api_url := "https://openbus.emtmadrid.es:9443/emt-proxy-server/last/" + path
    req_values.Add("idClient", creds.ClientID)
    req_values.Add("passKey",  creds.Password)

    req, err := http.NewRequest("POST", api_url,
                                bytes.NewBufferString(req_values.Encode()))
    req.Header.Set("X-Custom-Header", "busdm")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify : true}}
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


// Get EMT times
func GetArriveStop(IdStop int) GetArriveStopResponse{

    urlreq := url.Values{}
    urlreq.Add("cultureInfo", "ES")
    urlreq.Add("idStop", strconv.Itoa(IdStop))
    data := emtRequest("geo/GetArriveStop.php", urlreq)
    st := GetArriveStopResponse{}
    json.Unmarshal([]byte(data), &st)
    return st
}



// Response handler
func handler(w http.ResponseWriter, r *http.Request) {
    data := GetArriveStop(608);
    fmt.Fprintf(w, "Hi there, I love %v!", data)
    //getStopArrivalTime (608);
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}

