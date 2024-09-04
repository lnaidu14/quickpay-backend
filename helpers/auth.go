package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Auth0ManagementApiResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Auth0UserProfileIdentities struct {
	Provider    string `json:"provider"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	UserId      string `json:"user_id"`
	Connection  string `json:"connection"`
	IsSocial    bool   `json:"is_social"`
}

type Auth0UserProfile struct {
	CreatedAt     string                       `json:"created_at"`
	Email         string                       `json:"email"`
	EmailVerified bool                         `json:"email_verified"`
	FamilyName    string                       `json:"family_name"`
	GivenName     string                       `json:"given_name"`
	Identities    []Auth0UserProfileIdentities `json:"identities"`
	Name          string                       `json:"name"`
	Nickname      string                       `json:"nickname"`
	Picture       string                       `json:"picture"`
	UpdatedAt     string                       `json:"updated_at"`
	UserId        string                       `json:"user_id"`
	LastIp        string                       `json:"last_ip"`
	LastLogin     string                       `json:"last_login"`
	LoginsCount   uint                         `json:"logins_count"`
}

func FetchManagementApiToken() (Auth0ManagementApiResponse, error) {
	url := os.Getenv("AUTH0_MANAGEMENT_API_TOKEN")

	payload := strings.NewReader("{\"client_id\":\"Zu8Mu31v9Fn5fseMoUv0qG0fxYH3dJMP\",\"client_secret\":\"nlYDK37JWs3j4pUZRdgJwc8u2eDqjIsgETfoD2jw4_TekXsAYZev3x_i_NhdGPPe\",\"audience\":\"https://dev-quickpay.us.auth0.com/api/v2/\",\"grant_type\":\"client_credentials\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var parsedBody Auth0ManagementApiResponse

	err := json.Unmarshal(body, &parsedBody)

	if err != nil {
		fmt.Println(err)
	}

	return parsedBody, nil
}

func FetchUserProfile(accessToken string, userId string) (Auth0UserProfile, error) {
	url := "https://dev-quickpay.us.auth0.com/api/v2/users/" + userId

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var parsedBody Auth0UserProfile

	err := json.Unmarshal(body, &parsedBody)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("parsedBody: ", parsedBody)

	return parsedBody, nil
}
