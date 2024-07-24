package faucet

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed dist/*
var dist embed.FS

func ServeFaucet(c *Config) {
	distFS, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(distFS)))

	http.HandleFunc("/credit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err = r.ParseMultipartForm(10 << 20); err != nil {
				http.Error(w, "Failed to parse form data", http.StatusBadRequest)
				return
			}

			address := r.FormValue("address")
			if err = c.CreditTETH(address); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "TETH credited to %s", address)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	fmt.Println("Server running at 8580")
	http.ListenAndServe(":8580", nil)
}
