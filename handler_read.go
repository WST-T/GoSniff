package main

import "net/http"

func handlerRead(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct{}{})
}
