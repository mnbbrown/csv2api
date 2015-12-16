package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	branch string
	commit string
	date   string
)

var (
	dotenv     = flag.String("config", ".env", "")
	SERVE_FROM = ""
)

func handleAPI(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	filename := vars["filename"]
	filepath := fmt.Sprintf("%s/%s.csv", SERVE_FROM, filename)

	acceptType := req.Header.Get("Accept")

	switch {
	case acceptType == "text/csv":

		message := fmt.Sprintf("Requesting filename %s as CSV. Seriving %s/%s.csv", filename, SERVE_FROM, filename)
		log.Printf(message)
		rw.Header().Set("Content-Type", "text/csv")
		http.ServeFile(rw, req, filepath)

	case acceptType == "application/json" || acceptType == "application/x-json":

		message := fmt.Sprintf("Requesting filename %s as JSON.", filename)
		log.Printf(message)
		fmt.Fprint(rw, message)

	default:

		log.Println(fmt.Sprintf("No Accept header. Defaulting to JSON. Requesting %s as JSON.", filepath))
		csvFile, err := os.Open(filepath)
		if err != nil {
			log.WithField("filepath", filepath).Infof("Error opening CSV: %s", err)
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer csvFile.Close()

		reader := csv.NewReader(csvFile)
		reader.FieldsPerRecord = -1
		if err != nil {
			log.WithField("filepath", filepath).Errorf("Error parsing CSV: %s", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var headers []string
		var allRecords []map[string]string

		row := 0
		for {

			record, err := reader.Read()

			if err == io.EOF {
				break
			} else if err != nil {
				log.WithField("filepath", filepath).Errorf("Error parsing CSV: %s", err)
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if row == 0 {
				headers = make([]string, len(record))
				for i, h := range record {
					headers[i] = strings.ToLower(strings.Replace(h, " ", "_", -1))
				}
				row++
				continue
			}

			oneRecord := make(map[string]string)
			for itemIndex, item := range record {
				oneRecord[headers[itemIndex]] = item
			}
			allRecords = append(allRecords, oneRecord)
			row++
		}

		b, err := json.Marshal(allRecords)
		if err != nil {
			log.WithField("filepath", filepath).Errorf("Error encoding JSON: %s", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(b)
	}
}

func apiKeyMiddleware(key string, h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(rw, "Not Authorized", http.StatusUnauthorized)
			return
		}

		auth := strings.SplitN(authHeader, " ", 2)

		if len(auth) != 2 || auth[0] != "Bearer" || auth[1] != key {
			http.Error(rw, "Not Authorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(rw, req)
	})
}

func main() {

	conf := Load(*dotenv)

	APP_ENV := conf.Get("APP_ENV", "development")
	if APP_ENV == "production" || APP_ENV == "staging" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	API_KEY := conf.Get("API_KEY", "")
	if API_KEY == "" {
		log.Infof("API_KEY not set. Open to public.")
	}

	SERVE_FROM = conf.Get("SERVE_FROM", "./data")
	if SERVE_FROM == "" {
		log.Fatal("SERVE_FROM must be set.")
	}

	log.WithFields(log.Fields{
		"branch":       branch,
		"commit":       commit,
		"date":         date,
		"serving_from": SERVE_FROM,
		"api_key":      API_KEY,
	}).Infoln("Starting CSV2API")

	r := mux.NewRouter()
	if API_KEY == "" {
		r.Handle("/api/v1/{filename}", http.HandlerFunc(handleAPI))
	} else {
		r.Handle("/api/v1/{filename}", apiKeyMiddleware(API_KEY, handleAPI))
	}
	http.ListenAndServe(":8080", r)

}