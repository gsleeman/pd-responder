package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const SREDoNotDisturbPolicyID = "PNCPMTV"
const StatusTriggered = "triggered"

var httpclient *http.Client
var token string

type Event struct {
	Type string `json:"event_type"`
	Data EventData `json:"data"`
}

type EventData struct {
	ID string
	Status string
	Title string
	Reference string `json:"self"`
	EscalationPolicy EscalationPolicy `json:"escalation_policy"`
}

type EscalationPolicy struct {
	ID string
	Summary string
}

type IncidentUpdate struct {
	Type string `json:"type"`
	Status string `json:"status"`
	Resolution string `json:"resolution"`
}

func webhook(w http.ResponseWriter, r *http.Request) {
	len, err := strconv.ParseInt(r.Header.Get("content-length"), 10, 0)
	if err != nil {
		len = 40960
	}
	p := make([]byte,len)
	n, _ := r.Body.Read(p)
	var result map[string]Event
	json.Unmarshal(p[:n], &result)
	event := result["event"]
	w.WriteHeader(202)
	fmt.Println(event.Data.Reference)
	if event.Data.EscalationPolicy.ID == SREDoNotDisturbPolicyID && event.Data.Status == StatusTriggered {
		fmt.Printf("Resolving %s %s: ", event.Data.ID, event.Data.Title)
		body, err := json.Marshal(map[string]IncidentUpdate{
			"incident": {
				Type: "incident_reference",
				Status: "resolved",
				Resolution: "auto-resolved",
			},
		})
		if err != nil {
			log.Println(err)
			return
		}
		req, err := http.NewRequest(http.MethodPut, event.Data.Reference, bytes.NewBuffer(body))
		if err != nil {
			log.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", "Token token=" + token)
		res, err := httpclient.Do(req)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(res.Status)
		}
	}
}

func main() {
	httpclient = &http.Client{}
	token = os.Getenv("PAGERDUTY_TOKEN")
	if len(token) == 0 {
		panic("token variable not set")
	}
	http.ListenAndServe("0.0.0.0:8888", http.HandlerFunc(webhook))
}
