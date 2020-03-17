package eagleview

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	// "net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	// "github.com/spf13/viper"
	"githost.in/peartree/pkg/log"
)

const (
	baseUrl    = "https://webservices-integrations.eagleview.com"
	timeLayout = "Mon, 02 Jan 2006 15:04:05 MST"
)

func (client *Client) tokenize() error {
	// log.Info().Interface("client", client).Msg("tokenize")

	// Check our current token situation
	accessToken := os.Getenv("EAGLEVIEW_API_ACCESS_TOKEN")
	refreshToken := os.Getenv("EAGLEVIEW_API_REFRESH_TOKEN")
	accessTokenExpires := os.Getenv("EAGLEVIEW_API_ACCESS_TOKEN_EXPIRES") // unix timestamp

	// Check ENV for any tokens
	if accessToken != "" {
		client.Token.AccessToken = accessToken
	}

	if refreshToken != "" {
		client.Token.RefreshToken = refreshToken
	}

	if accessTokenExpires != "" {
		expiresTimestamp, err := strconv.Atoi(accessTokenExpires)
		if err != nil {
			return err
		}
		client.Token.ExpiresAt = time.Unix(int64(expiresTimestamp), 0)
	}

	// Check the client pieces

	if client.Token.AccessToken != "" {
		rightNow := time.Now()

		if client.Token.ExpiresAt.Unix() < rightNow.Unix() {
			// We're still valid, move along
			return nil
		}
	}

	// No valid access token

	// Prepare to get a new token

	reqUrl := fmt.Sprintf(baseUrl+"%s", `/Token`)
	base64Auth := base64.StdEncoding.EncodeToString([]byte(client.InitialRequest.SourceId + ":" + client.InitialRequest.ClientSecret))

	// Check for refresh token
	res := &http.Response{}

	if client.Token.RefreshToken != "" {
		reqBody := url.Values{}
		reqBody.Set("grant_type", `refresh_token`)
		reqBody.Add("refresh_token", client.Token.RefreshToken)

		httpClient := &http.Client{}
		req, err := http.NewRequest("POST", reqUrl, bytes.NewBufferString(reqBody.Encode()))
		if err != nil {

			return err
		}

		req.ContentLength = int64(len(reqBody.Encode()))
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64Auth))
		req.Header.Add("Content-Type", `application/x-www-form-urlencoded`)

		res, err = httpClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		log.Debug().Interface("response Refresh Status:", res.Status).Msg("tokenize")
	} else {
		// No tokens, need to get some
		if client.InitialRequest.Username == "" {
			return errors.New("Username missing")
		}

		if client.InitialRequest.Password == "" {
			return errors.New("Password missing")
		}

		reqBody := url.Values{}
		reqBody.Set("grant_type", `password`)
		reqBody.Add("username", client.InitialRequest.Username)
		reqBody.Add("password", client.InitialRequest.Password)

		httpClient := &http.Client{}
		req, err := http.NewRequest("POST", reqUrl, bytes.NewBufferString(reqBody.Encode()))
		if err != nil {
			return err
		}

		req.ContentLength = int64(len(reqBody.Encode()))
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64Auth))
		req.Header.Add("Content-Type", `application/x-www-form-urlencoded`)

		res, err = httpClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		log.Debug().Interface("response Initial Status:", res.Status).Msg("tokenize")
	}

	if res.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		// fmt.Println("[***** response data:", string(data))
		// fmt.Println("*****]")

		// cannot unmarshal string into Go struct field ErrorResponse.error of type eagleview.ErrorData

		eagleViewErrors := ErrorResponse{}
		// eagleData := ErrorData{}
		// log.Debug().Interface("data", data).Interface("eagleViewErrors", eagleViewErrors).Msg("tokenize")

		err = json.Unmarshal(data, &eagleViewErrors)
		if err != nil {
			return err
		}

		// log.Debug().Interface("eagleViewErrors", eagleViewErrors).Msg("tokenize")

		// return errors.New(eagleViewErrors.Error.Field + ": " + eagleViewErrors.Error.Description)
		errMesg := eagleViewErrors.Message
		if eagleViewErrors.Description != "" {
			errMesg = errMesg + ": " + eagleViewErrors.Description
		}

		return errors.New(errMesg)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	tokenResponse := TokenResponse{}

	err = json.Unmarshal(data, &tokenResponse)
	if err != nil {
		return err
	}

	if tokenResponse.ClientId != client.InitialRequest.SourceId {
		return errors.New("Client ID does not match Source ID")
	}

	token := Token{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		TokenType:    tokenResponse.TokenType,
		ClientId:     tokenResponse.ClientId,
	}

	// "Mon, 09 Mar 2020 22:28:49 GMT"
	issuedTime, err := time.Parse(timeLayout, tokenResponse.IssuedTime)
	if err != nil {
		return err
	}

	token.IssuedAt = issuedTime

	expiresTime, err := time.Parse(timeLayout, tokenResponse.ExpiresTime)
	if err != nil {
		return err
	}

	token.ExpiresAt = expiresTime

	// log.Info().Interface("issuedTime", issuedTime).Msg("tokenize")

	client.Token = token

	// TODO: Add new token values to ENV
	os.Setenv("EAGLEVIEW_API_ACCESS_TOKEN", client.Token.AccessToken)
	os.Setenv("EAGLEVIEW_API_REFRESH_TOKEN", client.Token.RefreshToken)
	os.Setenv("EAGLEVIEW_API_ACCESS_TOKEN_EXPIRES", strconv.Itoa(int(client.Token.ExpiresAt.Unix())))

	// fmt.Println("[***** made it to the bottom! *****]")
	// fmt.Println("client.Token.AccessToken", client.Token.AccessToken)
	// fmt.Println("client.Token.RefreshToken", client.Token.RefreshToken)
	// fmt.Println("client.Token.AccessTokenExpiresAt", client.Token.ExpiresAt)

	return nil
}

func (client *Client) Request(method string, endpoint, contentType string, params map[string]interface{}, response interface{}) error {

	err := client.tokenize()
	if err != nil {
		return err
	}

	httpClient := &http.Client{}
	req := &http.Request{}

	paramsAsBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reqUrl := fmt.Sprintf(baseUrl+"%s", endpoint)

	if len(params) > 0 {
		req, err = http.NewRequest(method, reqUrl, bytes.NewBuffer(paramsAsBytes))
		if err != nil {
			return err
		}
		req.ContentLength = int64(len(paramsAsBytes))
	} else {
		req, err = http.NewRequest(method, reqUrl, nil)
		if err != nil {
			return err
		}
	}

	// log.Debug().Interface("req", req).Msg("Request")

	// fmt.Println("paramsAsBytes", string(paramsAsBytes))
	// 	req, err := http.NewRequest(method, url, bytes.NewBuffer(paramsAsBytes))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.AuthOptions.ApiKey))

	// req, err := http.NewRequest(method, reqUrl, nil)	// orig non-params working
	// req, err := http.NewRequest(method, reqUrl, bytes.NewBuffer(paramsAsBytes))	// works with POST body
	// req, err := http.NewRequest(method, reqUrl, paramsAsBytes) // URL-encoded payload
	// if err != nil {
	// 	return err
	// }

	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", client.Token.AccessToken))

	if contentType == "" {
		req.Header.Add("Content-Type", `application/json`)
	} else {
		req.Header.Add("Content-Type", contentType)
	}

	// requestDump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(requestDump))
	// fmt.Println("requestDump", string(requestDump))

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// DEBUG start
	// data, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("***** status:", res.Status)
	// fmt.Println("***** actual data:", string(data))
	// DEBUG end

	if res.StatusCode != http.StatusOK {

		decoder := json.NewDecoder(res.Body)
		var errs map[string]interface{}
		err := decoder.Decode(&errs)
		if err != nil {
			return err
		}

		log.Debug().Interface("errs", errs).Msg("Request")
		eagleViewErrors := ErrorResponse{}

		if errs["Message"] != nil {
			eagleViewErrors.Message = errs["Message"].(string)
		}

		if errs["ModelState"] != nil {
			modelState := errs["ModelState"].(map[string]interface{})
			descriptions := []string{}

			for k, _ := range modelState {
				if modelState[k] != nil {
					details := modelState[k].([]interface{})

					if details != nil {
						for j, _ := range details {
							if details[j] != nil {
								descriptions = append(descriptions, details[j].(string))
							}

						}

						eagleViewErrors.Description = strings.Join(descriptions, "\n")
					}

				}
			}
		}

		if eagleViewErrors.Description == "" {
			return errors.New(eagleViewErrors.Message)
		}

		return errors.New(eagleViewErrors.Message + ": " + eagleViewErrors.Description)

		// --
		// data, err := ioutil.ReadAll(res.Body)
		// if err != nil {
		// 	return err
		// }

		// eagleViewErrors := ErrorResponse{}

		// err = json.Unmarshal(data, &eagleViewErrors)
		// if err != nil {
		// 	return err
		// }

		// errMesg := eagleViewErrors.Message
		// if eagleViewErrors.Description != "" {
		// 	errMesg = errMesg + ": " + eagleViewErrors.Description
		// }

		// return errors.New(errMesg)
	}

	json.NewDecoder(res.Body).Decode(response)
	return nil

	// availableProducts := []Product{}

	// err = json.Unmarshal(data, &availableProducts)
	// if err != nil {
	// 	return err
	// }

	// log.Debug().Interface("availableProducts", availableProducts).Msg("Request")

	// --
	// url := fmt.Sprintf(baseUrl+"%s", endpoint)

	// httpClient := &http.Client{}
	// reqBody := url.Values{}
	// reqBody.Set("grant_type", `password`)
	// reqBody.Add("username", client.InitialRequest.Username)
	// reqBody.Add("password", client.InitialRequest.Password)
	// // paramsAsBytes, err := json.Marshal(params)
	// // if err != nil {
	// // 	return err
	// // }
	// req, err := http.NewRequest(method, url, bytes.NewBuffer(data.Encode(reqBody)))
	// if err != nil {
	// 	return err
	// }
	// base64Auth := base64.StdEncoding.EncodeToString([]byte(client.InitialRequest.SourceId + ":" + client.InitialRequest.ClientSecret))
	// req.ContentLength = int64(len(reqBody))
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", base64Auth))
	// res, err := httpClient.Do(req)
	// if err != nil {
	// 	return err
	// }
	// defer res.Body.Close()
	// if res.StatusCode != http.StatusOK {
	// 	data, err := ioutil.ReadAll(res.Body)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	eagleViewErrors := ErrorResponse{}
	// 	err = json.Unmarshal(data, &eagleViewErrors)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return errors.New(eagleViewErrors.Error.Field + ": " + eagleViewErrors.Error.Description)

	// }
	// json.NewDecoder(res.Body).Decode(response)
	// return nil
}

// func (client *Client) request(method string, endpoint string, params map[string]interface{}, response interface{}) error {
// 	url := fmt.Sprintf(baseUrl+"%s", endpoint)

// 	httpClient := &http.Client{}
// 	paramsAsBytes, err := json.Marshal(params)
// 	if err != nil {
// 		return err
// 	}
// 	req, err := http.NewRequest(method, url, bytes.NewBuffer(paramsAsBytes))
// 	if err != nil {
// 		return err
// 	}
// 	req.ContentLength = int64(len(paramsAsBytes))
// 	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.AuthOptions.ApiKey))
// 	res, err := httpClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()
// 	if res.StatusCode != http.StatusOK {
// 		data, err := ioutil.ReadAll(res.Body)
// 		if err != nil {
// 			return err
// 		}

// 		yelpErrors := ErrorResponse{}
// 		err = json.Unmarshal(data, &yelpErrors)
// 		if err != nil {
// 			return err
// 		}

// 		return errors.New(yelpErrors.Error.Field + ": " + yelpErrors.Error.Description)

// 	}
// 	json.NewDecoder(res.Body).Decode(response)
// 	return nil
// }
