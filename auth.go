package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Auth0User struct {
	Sub           string    `json:"sub"`
	GivenName     string    `json:"given_name"`
	FamilyName    string    `json:"family_name"`
	Nickname      string    `json:"nickname"`
	Name          string    `json:"name"`
	Picture       string    `json:"picture"`
	Locale        string    `json:"locale"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
}

func getAuth0UserEmail(token string) string {
	url := config.Auth0Domain + "/userinfo"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// read json body and extract email address in a variable
	body, _ := ioutil.ReadAll(res.Body)
	// read body as json and extract email address
	var user Auth0User
	json.Unmarshal(body, &user)
	email := user.Email
	return email
}

type AuthTestBody struct {
	Role string `json:"role"`
}
type AuthList map[string]string

type RoleMap map[string]string

type DatasetRestrictedMap map[string][]string

func createAuthList(dataMap DataMap) AuthList {
	authList := make(AuthList, len(dataMap))
	for dataSetName, role := range config.RoleMap {
		for _, value := range dataMap[dataSetName] {
			email := value["email"].rawData
			authList[email] = role
		}
	}
	return authList

}

func setupAuthRoute(dataMap DataMap) {
	authList := createAuthList(dataMap)

	http.Handle("/auth", CorsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract token from request
		token := r.Header.Get("Authorization")
		if token != "" {
			// check if token is valid
			email := getAuth0UserEmail(token)
			// if token is valid, check email address
			if email != "" {
				// if email address is valid, check if email address is in the auth list
				if role, ok := authList[email]; ok {
					// if email address is in the auth list, return the role
					// aes encrypt email address and log it
					encryptedRole := encrypt(role)
					// return email address
					w.Write([]byte(encryptedRole))
					return
				} else {
					// if email address is not in the auth list, return "unauthorized"
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})))

	http.Handle("/authRole", CorsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// extract email from body
		role := getUserRole(r.Header)
		// if email is not empty, return email address
		if role != "" {
			// return email address
			w.Write([]byte(role))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})))
}

func getUserRole(header http.Header) string {
	token := header.Get("Authorization")
	if token == "" {
		return ""
	} else {
		role := decrypt(strings.Split(token, " ")[1])
		validatedRole := false
		for _, value := range config.RoleMap {
			if value == role {
				validatedRole = true
			}
		}
		if validatedRole {
			return role
		} else {
			return ""
		}
	}
}
