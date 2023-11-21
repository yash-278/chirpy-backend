package main

import (
	"errors"
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

func GetApiKey(r *http.Request) (string, error) {
	token := getHeaderValue("Authorization", r)

	splits := strings.Split(token, " ")

	// Check if the splits slice has at least two elements
	if len(splits) < 2 {
		// Return an error indicating that the Authorization header is not in the expected format
		return "", errors.New("invalid Authorization header format")
	}

	return splits[1], nil
}
