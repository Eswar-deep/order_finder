package main

import (
	"fmt"
	"net/http"
	"text/template"
)

const mongoURI = "mongodb://mongo:eSFmYQgVMdmYjtwAzXgLKiLEVRDRWmCX@viaduct.proxy.rlwy.net:54410"

const tomtomAPIKey = "N2NWaw1sogQ3oT2Rhn2GBTIuWnwIEckT"

// super
func main() {

	err := connectToMongoDB(mongoURI, "nokasa", "orders")
	if err != nil {
		fmt.Printf("Error connecting to MongoDB: %v\n", err)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		name := r.Form.Get("name")
		phoneno := r.Form.Get("phoneno")
		address := r.Form.Get("address")
		preferredtime := r.Form.Get("preferredtime")

		if preferredtime == "" {
			preferredtime = "Not specified" // Default value if preferredtime is empty
		}

		err = insertFormData(name, phoneno, address, preferredtime)
		if err != nil {
			http.Error(w, "Failed to store data", http.StatusInternalServerError)
			fmt.Printf("Error inserting data into MongoDB: %v\n", err)
			return
		}

		coords, err := geocodeAddress(address)
		if err != nil {
			http.Error(w, "Failed to geocode address", http.StatusInternalServerError)
			fmt.Printf("Error geocoding address: %v\n", err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/map?lat=%f&lon=%f", coords.Lat, coords.Lon), http.StatusSeeOther)
	})

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "map.html")
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
		fmt.Println("Started server at port 8080")
	}
}
