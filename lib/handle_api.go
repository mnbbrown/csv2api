package lib

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strings"
)

type Header struct {
	Key   int
	Value string
}

func sendJSON(fields []string, rw http.ResponseWriter, filepath string) {

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

	var headers []*Header
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

			// parse headers
			for i, h := range record {

				// convert header to lower case and remove all spaces.
				header := strings.ToLower(strings.Replace(h, " ", "_", -1))

				// filter by field
				if len(fields) > 0 {
					// if header is included in fields filter
					for _, field := range fields {
						if header == field {
							headers = append(headers, &Header{i, header})
						}
					}

				} else {
					// no field filter - add all headers.
					headers = append(headers, &Header{i, header})
				}
			}

			row++
			continue
		}

		oneRecord := make(map[string]string)
		for _, header := range headers {
			oneRecord[header.Value] = record[header.Key]
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
	rw.Write(b)
}

func NewHandler(serve_root string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		filename := vars["filename"]
		filepath := fmt.Sprintf("%s/%s.csv", serve_root, filename)

		acceptType := req.Header.Get("Accept")

		fieldList := strings.ToLower(strings.Replace(req.URL.Query().Get("fields"), " ", "", -1))
		var fields []string
		if fieldList != "" {
			fields = strings.Split(fieldList, ",")
		}

		switch {
		case acceptType == "text/csv":

			message := fmt.Sprintf("Requesting filename %s as CSV. Serving %s/%s.csv", filename, serve_root, filename)
			log.Printf(message)
			rw.Header().Set("Content-Type", "text/csv")
			http.ServeFile(rw, req, filepath)

		default:
			log.Println(fmt.Sprintf("No Accept header. Defaulting to JSON. Requesting %s as JSON.", filepath))
			rw.Header().Set("Content-Type", "application/json")
			sendJSON(fields, rw, filepath)
		}
	}
}
