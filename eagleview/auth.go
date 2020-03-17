package eagleview

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type InitialRequest struct {
	Username     string
	Password     string
	SourceId     string
	ClientSecret string
}

// type RefreshRequest struct {
// 	RefreshToken string
// }

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ClientId     string `json:"as:client_id"`
	IssuedTime   string `json:".issued"`
	ExpiresTime  string `json:".expires"`
	ExpiresIn    int    `json:"expires_in"`
}

type Token struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ClientId     string
	IssuedAt     time.Time
	ExpiresAt    time.Time
}

type ErrorResponse struct {
	Message     string `json:"error"`
	Description string `json:"error_description"`
}

// type ErrorData struct {
// 	Code        string      `json:"code"`
// 	Description string      `json:"description"`
// 	Field       string      `json:"field"`
// 	Instance    interface{} `json:"instance"`
// }

func GetInitVars() ([]string, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../") // This is needed to run the tests
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return []string{}, err
	}

	username := viper.GetString("EAGLEVIEW_API_USERNAME")
	password := viper.GetString("EAGLEVIEW_API_PASSWORD")
	source_id := viper.GetString("EAGLEVIEW_API_SOURCE_ID")
	client_secret := viper.GetString("EAGLEVIEW_API_CLIENT_SECRET")

	// These should always start out empty
	// access_token := viper.GetString("EAGLEVIEW_API_ACCESS_TOKEN")
	// access_token_expires := viper.GetString("EAGLEVIEW_API_ACCESS_TOKEN_EXPIRES")
	// refresh_token := viper.GetString("EAGLEVIEW_API_REFRESH_TOKEN")

	result := make([]string, 4)
	result[0] = username
	result[1] = password
	result[2] = source_id
	result[3] = client_secret
	// result[4] = access_token
	// result[5] = access_token_expires
	// result[6] = refresh_token

	return result, nil
}
