package handlers

import (
	"broker-service/models"
	"errors"
	"net/http"
	"strconv"
)

func (app *Config) HandleBlog(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.createBlog(w, r)
	case http.MethodGet:
		app.getBlog(w, r)
	case http.MethodPut:
		app.updateBlog(w, r)
	case http.MethodDelete:
		app.deleteBlog(w, r)
	default:
		app.AppConfig.ErrorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
	}
}

func (app *Config) HandleListBlogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.AppConfig.ErrorJSON(w, errors.New("method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		app.AppConfig.ErrorJSON(w, errors.New("user_id is required"), http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}
	}

	response, err := app.BlogClient.ListBlogs(userID, page, limit)
	if err != nil {
		app.AppConfig.ErrorJSON(w, err)
		return
	}

	_ = app.AppConfig.WriteJSON(w, http.StatusOK, response)
}

func (app *Config) createBlog(w http.ResponseWriter, r *http.Request) {
	var requestPayload models.CreateBlogRequest
	err := app.AppConfig.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.AppConfig.ErrorJSON(w, err)
		return
	}

	if requestPayload.UserID == "" {
		app.AppConfig.ErrorJSON(w, errors.New("user_id is required"))
		return
	}
	if requestPayload.Title == "" {
		app.AppConfig.ErrorJSON(w, errors.New("title is required"))
		return
	}
	if requestPayload.Text == "" {
		app.AppConfig.ErrorJSON(w, errors.New("text is required"))
		return
	}

	response, err := app.BlogClient.CreateBlog(
		requestPayload.UserID,
		requestPayload.Title,
		requestPayload.Text,
	)
	if err != nil {
		app.AppConfig.ErrorJSON(w, err)
		return
	}

	_ = app.AppConfig.WriteJSON(w, http.StatusCreated, response)
}

func (app *Config) getBlog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		app.AppConfig.ErrorJSON(w, errors.New("id is required"), http.StatusBadRequest)
		return
	}

	response, err := app.BlogClient.GetBlog(id)
	if err != nil {
		app.AppConfig.ErrorJSON(w, err)
		return
	}

	_ = app.AppConfig.WriteJSON(w, http.StatusOK, response)
}

func (app *Config) updateBlog(w http.ResponseWriter, r *http.Request) {
	var requestPayload models.UpdateBlogRequest
	err := app.AppConfig.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.AppConfig.ErrorJSON(w, err)
		return
	}

	if requestPayload.ID == "" {
		app.AppConfig.ErrorJSON(w, errors.New("id is required"))
		return
	}
	if requestPayload.Title == "" && requestPayload.Text == "" {
		app.AppConfig.ErrorJSON(w, errors.New("title or text is required"))
		return
	}

	response, err := app.BlogClient.UpdateBlog(
		requestPayload.ID,
		requestPayload.Title,
		requestPayload.Text,
	)
	if err != nil {
		app.AppConfig.ErrorJSON(w, err)
		return
	}

	_ = app.AppConfig.WriteJSON(w, http.StatusOK, response)
}

func (app *Config) deleteBlog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		app.AppConfig.ErrorJSON(w, errors.New("id is required"), http.StatusBadRequest)
		return
	}

	response, err := app.BlogClient.DeleteBlog(id)
	if err != nil {
		app.AppConfig.ErrorJSON(w, err)
		return
	}

	_ = app.AppConfig.WriteJSON(w, http.StatusOK, response)
}
