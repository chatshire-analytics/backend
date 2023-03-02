package setup

import (
	"chatgpt-service/internal/config"
	"chatgpt-service/internal/pkg/client"
	"chatgpt-service/internal/pkg/constants"
	"chatgpt-service/internal/pkg/store"
	"github.com/go-pg/pg"
	"net/http"
	"time"
)

// InitializeDatabase TODO: refactor to configuration environment variable
func InitializeDatabase(_ *config.GlobalConfig) *store.Database {
	endpoint := "pg-ffssb.vpc-pub-cdb-kr.ntruss.com"
	port := "5432"
	user := "chatshire"
	password := "ethdenver1!"
	databaseName := "postgres"

	db := pg.Connect(&pg.Options{
		Addr:     endpoint + ":" + port,
		User:     user,
		Password: password,
		Database: databaseName,
	})

	return &store.Database{DB: db}
}

func NewOpenAIClient(cfg *config.GlobalConfig, options ...client.ClientOption) (*client.OpenAIClient, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(constants.DefaultTimeoutSeconds * time.Second),
	}
	cl := &client.OpenAIClient{
		UserAgent:     constants.DefaultUserAgent,
		ApiKey:        cfg.OpenAIEnv.API_KEY,
		AccessToken:   cfg.OpenAIEnv.ACCESS_TOKEN,
		BaseURL:       constants.DefaultBaseURL,
		HttpClient:    httpClient,
		DefaultEngine: constants.DefaultEngine,
		IdOrg:         constants.DefaultClientName,
	}
	for _, clientOption := range options {
		err := clientOption(cl)
		if err != nil {
			return nil, err
		}
	}
	return cl, nil
}

func NewFlipsideClient(cfg *config.GlobalConfig) (*client.FlipsideClient, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(constants.DefaultTimeoutSeconds * time.Second),
	}
	cl := &client.FlipsideClient{
		UserAgent:  constants.DefaultUserAgent,
		HttpClient: httpClient,
		ApiKey:     cfg.FlipsideEnv.API_KEY,
		BaseURL:    constants.DefaultFlipsideBaseURL,
	}
	return cl, nil
}
