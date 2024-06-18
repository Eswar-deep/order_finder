package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

var (
	mongoURI     string
	googleAPIKey string
	port         string
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	mongoURI = os.Getenv("MONGO_URI")
	googleAPIKey = os.Getenv("GOOGLE_MAPS_API_KEY")
	port = os.Getenv("PORT")

}

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
			log.Printf("Failed to parse template: %v", err)
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

		http.Redirect(w, r, fmt.Sprintf("/map?lat=%f&lng=%f", coords.Lat, coords.Lng), http.StatusSeeOther)
	})

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "map.html")
	})

	fmt.Printf("Starting server at port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
