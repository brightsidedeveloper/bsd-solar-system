package util

import (
	"encoding/json"
	"net/http"
)

type JSON struct{}

func (j *JSON) Write(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload == nil {
		http.Error(w, "No payload provided", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (j *JSON) Success(w http.ResponseWriter, payload interface{}) {
	j.Write(w, http.StatusOK, payload)
}

func (j *JSON) Error(w http.ResponseWriter, status int, message string) {
	j.Write(w, status, map[string]string{"error": message})
}

func (j *JSON) ValidationError(w http.ResponseWriter, message string) {
	j.Error(w, http.StatusUnprocessableEntity, message)
}
