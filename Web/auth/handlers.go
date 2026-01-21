package main

import (
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type Auth struct {
	Users Users
	Config *webauthn.Config
}

func NewAuth() *Auth {
	return &Auth{
		Config: &webauthn.Config{
			RPDisplayName: "Go WebAuthn",
			RPID:          "app.awesome-go-webauthn.com",
			RPOrigins:     []string{"https://app.awesome-go-webauthn.com"},
		},
	}
}

func (a *Auth)CreateChallenge(w *webauthn.WebAuthn) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Crude / Abstract  of retrieving the user this registration will belong to. The user must be logged in
		// for this step unless you plan to register the user and the credential at the same time i.e. usernameless.
		// The user should have a unique and stable value returned from WebAuthnID that can be used to retrieve the
		// account details for the user.
		email := r.FormValue("email")
		user, err := a.Users.GetUser(email)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		var (
			creation *protocol.CredentialCreation
			s        *webauthn.SessionData
		)

		opts := []webauthn.RegistrationOption{
			webauthn.WithExclusions(webauthn.Credentials(user.WebAuthnCredentials()).CredentialDescriptors()),
			webauthn.WithExtensions(map[string]any{"credProps": true}),
		}

		if creation, s, err = w.BeginMediatedRegistration(user, protocol.MediationDefault, opts...); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		// Crude  saving the session data securely to be loaded in the finish step of the register action. This
		// should be stored in such a way that the user and user agent has no access to it. For  using an opaque
		// session cookie.
		a.Session.Save(s)

		encoder := json.NewEncoder(rw)

		if err = encoder.Encode(creation); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
	}
}

func (a *Auth)ValidateCreateChallengeResponse(w *webauthn.WebAuthn) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Crude / Abstract  of retrieving the user performing the multi-factor authentication. The user must be
		// logged in for this step.
		user, err := LoadUser()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		// Crude  loading the session data securely from the start step for the register action. This should be
		// loaded from a place the user and user agent has no access to it. For  using an opaque session cookie.
		s, err := a.session.Load()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		credential, err := w.FinishRegistration(user, *s, r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		// Crude / Abstract  of adding the credential to the list of credentials for the user. This is critical
		// for performing future logins.
		user.credentials = append(user.credentials, *credential)

		// Crude / Abstract  of saving the updated user. This is critical for performing future logins.
		if err = SaveUser(user); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}

func (a *Auth)LoginChallenge(w *webauthn.WebAuthn) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Crude / Abstract  of retrieving the user for this multi-factor authentication. Because this is a
		// multi-factor authentication the user MUST be logged in at this stage and the returned struct/interface must
		// be deterministically matched to their account.
		user, err := LoadUser()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		assertion, s, err := w.BeginMediatedLogin(user, protocol.MediationDefault)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		// Crude  saving the session data securely to be loaded in the finish step of the login action. This
		// should be stored in such a way that the user and user agent has no access to it. For  using an opaque
		// session cookie.
		a.session.Save(s)

		encoder := json.NewEncoder(rw)

		if err = encoder.Encode(assertion); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
	}
}

func (a *Auth)LoginChallengeResponse(w *webauthn.WebAuthn) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Crude / Abstract  of retrieving the user performing the multi-factor authentication. The user must be
		// logged in for this step.
		user, err := LoadUser()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		// Crude  loading the session data securely from the start step for the login action. This should be
		// loaded from a place the user and user agent has no access to it. For  using an opaque session cookie.
		s, err := a.Session.Load()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		validatedCredential, err := w.FinishLogin(user, *s, r)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		var found bool

		// Modify the matching credential in the user struct which is critical for proper future validations as the
		// metadata for this credential has been updated. No type assertion is required here since the LoadUser function
		// returns the concrete implementation, you may have to adjust this if you return the abstract implementation
		// instead.
		for i, credential := range user.credentials {
			if bytes.Equal(validatedCredential.ID, credential.ID) {
				user.credentials[i] = *validatedCredential

				// Crude / Abstract  of saving the user with their updated credentials. This is critical for
				// proper future validations.
				if err = SaveUser(user); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)

					return
				}

				found = true

				break
			}
		}

		// Should error if we can't update the credentials for the user.
		if !found {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}

