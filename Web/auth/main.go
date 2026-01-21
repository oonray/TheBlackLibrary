package main

import (
	"net/http"
	"time"
	"log"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
  "github.com/gorilla/mux"
)

var (
	auth *Auth = NewAuth()
	router *mux.Router = mux.NewRouter()
)

// Example_multiFactorRegisterAndLogin demonstrates handling Multi Factor registration and Logins. This uses the higher level APIs to
// perform all of the various requirements. The Crude and Abstract examples are purely domain logic and will often
// describe aspects that should be considered during their implementation if they are important; these aspects
// are not strictly concerns related to the library as there are too many logical implementations to count.
func main() {
	w, err := webauthn.New(auth.Config)
	if err != nil {
		log.Fatal(fmt.Errorf("Could not create webauthn %s",err))
	}

	// Register the handlers. The second component describes the action (i.e. register/login), the final component
	// describes the step (i.e. start/finish).
	router.HandleFunc("/webauthn/register/start", handlerExampleMultiFactorCreateChallenge(w))
	router.HandleFunc("/webauthn/register/finish", handlerExampleMultiFactorValidateCreateChallengeResponse(w))

	router.HandleFunc("/webauthn/login/start", handlerExampleMultiFactorLoginChallenge(w))
	router.HandleFunc("/webauthn/login/finish", handlerExampleMultiFactorLoginChallengeResponse(w))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	if err = server.ListenAndServe(); err != nil {
		log.Fatal(fmt.Sprintf("Could not listen on :8080 | %s",err))
	}
}

