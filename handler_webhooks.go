package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	type response struct {
		User
	}

	polkaKey, err := GetApiKey(r)
	if polkaKey != cfg.polkaKey || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusOK, "OK")
		return
	}

	_, err = cfg.DB.GetUser(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	_, err = cfg.DB.UpdateUser(params.Data.UserID, "", "", true)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User{},
	})
}
