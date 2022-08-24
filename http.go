package main

import (
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow cross domain AJAX requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "true")
		if r.Method == "OPTIONS" { // pre-flight request
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := getBodyString(r)
		if identifyIntrospectionQuery(body) {
			next.ServeHTTP(w, r)
			return
		}
		role := getUserRole(r.Header)
		datasets := extractQueriedDataset(body)
		if role == "" && config.RestrictAllDatasets {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		for _, dataset := range datasets {
			if authorizedRole, ok := config.DatasetRestricted[dataset]; ok {
				if !containsString(authorizedRole, role) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

func httpListen() {
	http.ListenAndServe(config.ServerAddr, nil)
}
