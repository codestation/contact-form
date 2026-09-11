package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"megpoid.dev/go/contact-form/oapi"
)

func TestRootHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := rootHandler(ctx)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))

	var response appInfo
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&response))
	assert.Equal(t, appInfo{Name: "forms", Version: "v1"}, response)
}

func TestRootIsNotInSwagger(t *testing.T) {
	spec, err := oapi.GetSwagger()
	require.NoError(t, err)

	assert.Nil(t, spec.Paths.Find("/"))
}
