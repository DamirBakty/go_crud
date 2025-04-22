package handlers

import (
	"broker-service/clients"
	"broker-service/config"
	"broker-service/models"
	"errors"
	"net/http"
)

type Config struct {
	AppConfig       *config.AppConfig
	AuthGRPCAddress string
	AuthClient      *clients.AuthClient
	BlogGRPCAddress string
	BlogClient      *clients.BlogClient
}

func (app *Config) HandleRoot(w http.ResponseWriter, r *http.Request) {
	payload := config.JsonResponse{
		Error:   false,
		Message: "Broker service is running",
	}

	_ = app.AppConfig.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.AppConfig.ErrorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var requestPayload models.AuthRequest
	err := app.AppConfig.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.AppConfig.ErrorJSON(w, err)
		return
	}

	var response *models.AuthResponse
	var authErr error

	switch requestPayload.Action {
	case "register":
		response, authErr = app.AuthClient.Register(
			requestPayload.Username,
			requestPayload.Email,
			requestPayload.Password,
		)
	case "login":
		response, authErr = app.AuthClient.Login(
			requestPayload.Email,
			requestPayload.Password,
		)
	case "verify":
		response, authErr = app.AuthClient.VerifyToken(
			requestPayload.Token,
		)
	default:
		app.AppConfig.ErrorJSON(w, errors.New("invalid action"))
		return
	}

	if authErr != nil {
		app.AppConfig.ErrorJSON(w, authErr)
		return
	}

	_ = app.AppConfig.WriteJSON(w, http.StatusOK, response)
}
