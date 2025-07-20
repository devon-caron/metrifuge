package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
	"net/http"
)

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: "https://dev-i7nwtfj2nh8lqm6v.us.auth0.com/.well-known/jwks.json"}, nil)
	configuration := auth0.NewConfiguration(client, []string{"http://localhost:8000"}, "https://dev-i7nwtfj2nh8lqm6v.us.auth0.com/", jose.RS256)
	validator := auth0.NewValidator(configuration, auth0.RequestTokenExtractorFunc(auth0.FromHeader))

	myJwt, err := validator.ValidateRequest(r)
	myJwtStr, _ := json.Marshal(myJwt)

	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	fmt.Printf("Validated token: \n%v\n", myJwtStr)

	// At this point, the token is valid. You can now create and return your custom token if needed.
	customToken := map[string]string{
		"access_token": "your_custom_token_here",
		"token_type":   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(customToken)
	if err != nil {
		fmt.Println("error: " + err.Error())
		return
	}
}
