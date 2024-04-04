package db

import (
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
)

var progressURL = "https://api.nethackathon.org/charity/progress"


type Fundraiser struct {
	Currency   string     `json: "currency"`
	Raised     string     `json: "raised"`
    GoalAmount string     `json: "goalAmount"`
    GoalType     string   `json: "goalType"`
    Supporters int        `json: "supporters"`
}

func GetFundraiserData() Fundraiser {
	resp, err := http.Get(progressURL)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//bodyString := string(body)
    fundraiser := Fundraiser{}
    err = json.Unmarshal(body, &fundraiser) 
    if err != nil {
        log.Fatalln(err)
    }
	return fundraiser
}
