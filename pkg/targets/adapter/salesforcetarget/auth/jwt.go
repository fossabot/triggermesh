/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

const (
	// grant type for OAuth JWT.
	// See: https://tools.ietf.org/html/rfc7523#page-10
	grantJWT = "urn:ietf:params:oauth:grant-type:jwt-bearer"

	oauthTokenPath = "/services/oauth2/token"
)

// JWTAuthenticator is the JWT OAuth implementation.
// See: https://help.salesforce.com/articleView?id=remoteaccess_oauth_jwt_flow.htm
type JWTAuthenticator struct {
	authURL string
	signKey *rsa.PrivateKey
	claims  *claims

	client *http.Client
	logger *zap.SugaredLogger
}

type claims struct {
	jwt.StandardClaims
}

// Credentials returned from Salesforce Auth.
type Credentials struct {
	Token        string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	InstanceURL  string `json:"instance_url"`
	ID           string `json:"id"`
	Signature    string `json:"signature"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
	CommunityURL string `json:"sfdc_community_url"`
	CommunityID  string `json:"sfdc_community_id"`
}

// NewJWTAuthenticator creates an OAuth JWT authenticator for Salesforce.
func NewJWTAuthenticator(certKey, clientID, user, server string, client *http.Client, logger *zap.SugaredLogger) (*JWTAuthenticator, error) {
	audience := strings.TrimSuffix(server, "/")

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(certKey))
	if err != nil {
		return nil, fmt.Errorf("unable to parse PEM private key: %w", err)
	}

	claims := &claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:   clientID,
			Subject:  user,
			Audience: audience,
		},
	}

	return &JWTAuthenticator{
		authURL: audience + oauthTokenPath,
		claims:  claims,
		signKey: signKey,
		client:  client,
		logger:  logger,
	}, nil
}

// NewCredentials generates a new set of credentials.
func (j *JWTAuthenticator) NewCredentials() (*Credentials, error) {
	// expiry needs to be set to 3 minutes or less
	// See: https://help.salesforce.com/articleView?id=remoteaccess_oauth_jwt_flow.htm
	j.claims.ExpiresAt = time.Now().Add(time.Minute * 3).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, j.claims)
	signedToken, err := token.SignedString(j.signKey)
	if err != nil {
		return nil, fmt.Errorf("could not sign JWT token: %w", err)
	}

	form := url.Values{
		"grant_type": []string{grantJWT},
		"assertion":  []string{signedToken},
	}

	req, err := http.NewRequest("POST", j.authURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("could not build authentication request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := j.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute authentication request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		msg := fmt.Sprintf("received unexpected status code %d from authentication", res.StatusCode)
		if resb, err := ioutil.ReadAll(res.Body); err == nil {
			msg += ": " + string(resb)
		}
		return nil, errors.New(msg)
	}

	c := &Credentials{}
	err = json.NewDecoder(res.Body).Decode(c)
	if err != nil {
		return nil, fmt.Errorf("could not decode authentication response into credentails: %w", err)
	}

	return c, nil
}

// RefreshCredentials renews credentials.
func (j *JWTAuthenticator) RefreshCredentials() (*Credentials, error) {
	return nil, errors.New("not implemented")
}

// CreateOrRenewCredentials will always create a new set of credentials.
func (j *JWTAuthenticator) CreateOrRenewCredentials() (*Credentials, error) {
	return j.NewCredentials()
}
