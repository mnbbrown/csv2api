package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/mnbbrown/csv2api/lib"
	"net/http"
	"strings"
)

var (
	branch string
	commit string
	date   string
)

var (
	dotenv = flag.String("config", ".env", "")
)

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

	SERVE_FROM := conf.Get("SERVE_FROM", "./data")
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
		r.Handle("/api/v1/{filename}", http.HandlerFunc(lib.NewHandler(SERVE_FROM)))
	} else {
		r.Handle("/api/v1/{filename}", apiKeyMiddleware(API_KEY, lib.NewHandler(SERVE_FROM)))
	}
	http.ListenAndServe(":8080", r)

}
