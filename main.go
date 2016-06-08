package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blang/semver"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var sr SoffitRequest

		err := json.NewDecoder(r.Body).Decode(&sr)
		if err != nil {
			log.Println("Error decoding JSON", err)
		}

		if bs, err := json.MarshalIndent(sr, "", "  "); err == nil {
			os.Stdout.Write(bs)
			w.Write(bs)
		}

		fmt.Fprintf(w, "<h1>hello world</h1>")
	})

	log.Fatal(http.ListenAndServe(":8089", nil))
}

type SoffitRequest struct {
	Mode        string              `json:"mode"`
	Namespace   string              `json:"namespace"`
	WindowState string              `json:"windowState"`
	Portal      PortalInfo          `json:"portal"`
	Preferences map[string][]string `json:"preferences"`
}

type PortalInfo struct {
	Provider string         `json:"provider"`
	Version  semver.Version `json:"version"`
}
