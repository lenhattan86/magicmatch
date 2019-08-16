package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.paypal.com/PaaS/MagicMatch/scheduler"
)

// index shows the home page of magicmatch
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "MagicMatch!")
}

// matchesIndex returns all matches for Aurora to read
func matchesIndex(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(scheduler.MatchesMap); err != nil {
		log.Errorf(err.Error())
	}
}

// matchTaskIdIndex return a match for a given TaskId
func matchTaskIdIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId := vars["taskId"]
	if err := json.NewEncoder(w).Encode(scheduler.MatchesMap[taskId]); err != nil {
		log.Errorf(err.Error())
	}
}

// MesosIndex returns the data from Mesos
func mesosIndex(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(capacitiesMap); err != nil {
		log.Errorf(err.Error())
	}
}

// AuroraIndex returns tasks from Aurora
func auroraIndex(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		panic(err)
	}
}

// HostIndex returns Host Info
func hostIndex(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(hostMapCache); err != nil {
		panic(err)
	}
}
