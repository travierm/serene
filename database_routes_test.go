package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetDatabases(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/databases", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetDatabases(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var responseUsers []Database
		err := json.Unmarshal(rec.Body.Bytes(), &responseUsers)
		assert.NoError(t, err)
		assert.Equal(t, databases, responseUsers)
	}
}

func TestCreateDatabase(t *testing.T) {
	e := echo.New()
	userJSON := `{"name":"ds_main"}`
	req := httptest.NewRequest(http.MethodPost, "/databases", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, CreateDatabase(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var createdUser Database
		err := json.Unmarshal(rec.Body.Bytes(), &createdUser)
		assert.NoError(t, err)
		assert.Equal(t, "ds_main", createdUser.Name)
		assert.Equal(t, len(databases), createdUser.ID)
	}
}
