package client

import (
	"context"
	"encoding/json"
	"time"

	"github.com/monzo/terrors"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"swallowtail/libraries/util"
	"swallowtail/s.googlesheets/templates"
)

var (
	// TODO refactor
	client      GooglesheetsClient
	tok         *oauth2.Token
	credentials *ClientCredentials
)

type ClientCredentials struct {
	Installed *Credentials `json:"installed"`
}

type Credentials struct {
	ClientID            string   `json:"client_id"`
	ProjectID           string   `json:"project_id"`
	AuthURI             string   `json:"auth_id"`
	TokenURI            string   `json:"token_uri"`
	AuthProviderCertURL string   `json:"auth_provider_x509_cert_url"`
	ClientSecret        string   `json:"client_secret"`
	RedirectURLs        []string `json:"redirect_uris"`
}

func init() {
	// Credentials
	clientID := util.SetEnv("GOOGLESHEETS_CLIENT_ID")
	projectID := util.SetEnv("GOOGLESHEETS_PROJECT_ID")
	authURI := util.SetEnv("GOOGLESHEETS_AUTH_URI")
	tokenURI := util.SetEnv("GOOGLESHEETS_TOKEN_URI")
	authProviderCertURL := util.SetEnv("GOOGLESHEETS_AUTH_PROVIDER_CERT_URL")
	clientSecret := util.SetEnv("GOOGLESHEETS_CLIENT_SECRET")
	redirectURLA := util.SetEnv("GOOGLESHEETS_REDIRECT_URL_1")
	redirectURLB := util.SetEnv("GOOGLESHEETS_REDIRECT_URL_2")

	credentials = &ClientCredentials{
		Installed: &Credentials{
			ClientID:            clientID,
			ProjectID:           projectID,
			AuthURI:             authURI,
			TokenURI:            tokenURI,
			AuthProviderCertURL: authProviderCertURL,
			ClientSecret:        clientSecret,
			RedirectURLs:        []string{redirectURLA, redirectURLB},
		},
	}

	// Token
	accessToken := util.SetEnv("GOOGLESHEETS_ACCESS_TOKEN")
	tokenType := util.SetEnv("GOOGLESHEETS_TOKEN_TYPE")
	refreshToken := util.SetEnv("GOOGLESHEETS_REFRESH_TOKEN")
	expiry := util.SetEnv("GOOGLESHEETS_EXPIRY")

	t, _ := time.Parse(time.RFC3339, expiry)
	tok = &oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		RefreshToken: refreshToken,
		Expiry:       t,
	}
}

// GooglesheetsClient ...
type GooglesheetsClient interface {
	Ping(ctx context.Context) bool
	Values(ctx context.Context, sheetID string, rowsRange string) ([][]interface{}, error)
	UpdateRows(ctx context.Context, sheetID, rowsRange string, values [][]interface{}) error
	CreateSheet(ctx context.Context, sheetType templates.SheetType) (string, error)
}

// Init sets a new google sheets client for interacting with googlesheets.
func Init(ctx context.Context) error {
	b, err := json.Marshal(credentials)
	if err != nil {
		return terrors.Augment(err, "Failed to create credentials for googlesheets client", nil)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return terrors.Augment(err, "Failed to load config", nil)
	}

	c := config.Client(ctx, tok)

	srv, err := sheets.New(c)
	if err != nil {
		return terrors.Augment(err, "Failed to create sheets client", nil)
	}

	client = &googlesheetsClient{
		s: srv,
	}

	return nil
}

// Values ...
func Values(ctx context.Context, sheetID string, rowsRange string) ([][]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Retrieve values from googlesheets spreadsheet")
	defer span.Finish()
	return client.Values(ctx, sheetID, rowsRange)
}

// UpdateRows ...
func UpdateRows(ctx context.Context, sheetID string, rowsRange string, values [][]interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Update values from googlesheets spreadsheet")
	defer span.Finish()
	return client.UpdateRows(ctx, sheetID, rowsRange, values)
}

// CreateSheet ...
func CreateSheet(ctx context.Context, sheetType templates.SheetType) (string, error) {
	// TODO: metrics & opnetracing.
	return client.CreateSheet(ctx, sheetType)
}

func isValidHTTPStatusCode(c int) bool {
	if c < 200 || c > 299 {
		return false
	}
	return true
}
