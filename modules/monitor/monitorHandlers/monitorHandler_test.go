package monitorHandlers

import (
	"github.com/Montheankul-K/assessment-tax/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockConfig struct {
	mockAppConfig *MockAppConfig
	mockDBConfig  *MockDBConfig
	mockAdminAuth *MockAdminAuth
}

type MockAppConfig struct{}
type MockDBConfig struct{}
type MockAdminAuth struct{}

func (m *MockConfig) App() config.IAppConfig {
	return m.mockAppConfig
}

func (m *MockConfig) DB() config.IDBConfig {
	return m.mockDBConfig
}

func (m *MockConfig) AdminAuth() config.IAdminAuth {
	return m.mockAdminAuth
}

func (m *MockAppConfig) Name() string {
	return "name"
}

func (m *MockAppConfig) Port() string {
	return "port"
}

func (m *MockAppConfig) Version() string {
	return "version"
}

func (m *MockDBConfig) Url() string {
	return "url"
}

func (m *MockAdminAuth) Username() string {
	return "username"
}

func (m *MockAdminAuth) Password() string {
	return "password"
}

func setupEchoContext() (ctx echo.Context, recorder *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestHealthCheck(t *testing.T) {
	handler := MonitorHandler(&MockConfig{})

	c, rec := setupEchoContext()
	err := handler.HealthCheck(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	assert.JSONEq(t, `{"name":"name","version":"version"}`, rec.Body.String())
}
