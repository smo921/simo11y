package sources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"simo11y/internal/types"
)

func HTTP(done chan string) <-chan types.StructuredMessage {
	out := make(chan types.StructuredMessage)
	var msgs *types.StructuredMessages

	// setup http listener
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		if msgs, err = types.ReadStructuredMessages(body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Error Reading Messages: %s\n", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Status OK"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)

		// Blocking
		for _, m := range *msgs {
			out <- m
		}

		return
	})

	go func() {
		log.Print("Starting HTTP listener on :8080")
		http.ListenAndServe(":8080", nil)
		defer close(out)

		<-done
	}()

	return out
}
