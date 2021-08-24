package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func handleVerify(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	if !strings.Contains(r.Header.Get("content-type"), "application/json") {
		response := VerifyResponse{
			VerificationStatus: StatusUnfulfilled,
			Message:            "content type must be JSON",
		}
		resp, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	var input VerifyInput
	body, err := ioutil.ReadAll(r.Body)
	unmarshallErr := json.Unmarshal(body, &input)

	if err != nil || unmarshallErr != nil {
		response := VerifyResponse{
			VerificationStatus: StatusUnfulfilled,
			Message:            "invalid_request_body",
		}

		resp, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	err = VerifyNGO(input)

	if err != nil {
		response := VerifyResponse{
			Message: err.Error(),
		}

		if err.Error() == "unexpected_error" {
			response.VerificationStatus = StatusUnfulfilled
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			response.VerificationStatus = StatusFail
			w.WriteHeader(http.StatusOK)
		}

		resp, _ := json.Marshal(response)
		w.Write(resp)
		return
	}

	response := VerifyResponse{
		VerificationStatus: StatusSuccess,
		Message:            "success",
	}

	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(response)
	w.Write(resp)
}

func StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/verify", handleVerify)
	log.Fatal(http.ListenAndServe(":8080", r))
}
