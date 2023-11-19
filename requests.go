package main

import (
	"net/http"
	"strings"
)

func getHeaderValue(key string, r *http.Request) string {
	value := r.Header.Get(key)
	return value
}

func GetBearerToken(r *http.Request) string {
	token := getHeaderValue("Authorization", r)

	splits := strings.Split(token, " ")

	return splits[1]
}
