package db

import (
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
    "os"
)

var giveDirectlyURL = "https://partners.every.org/v0.2/nonprofit/givedirectly/fundraiser/nethackathon-vi-fundraising/raised?apiKey="


type Fundraiser struct {
	Currency   string     `json: "currency"`
	Raised     string     `json: "raised"`
    GoalAmount string     `json: "goalAmount"`
    GoalType     string   `json: "goalType"`
    Supporters int        `json: "supporters"`
}

func GetFundraiserData() Fundraiser {
    giveDirectlyAPIKey := os.Getenv("GIVEDIRECTLYAPI")
	resp, err := http.Get(giveDirectlyURL + giveDirectlyAPIKey)
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
