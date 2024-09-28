package api

import "net/http"

// Broker api Handler
func Test(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Massage: "Hello from Auth!",
	}

	_ = writeJSON(w, http.StatusAccepted, payload)
}
