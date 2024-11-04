package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var statusCases = []struct {
	request string
	code    int
	message string
}{
	{
		request: "/cafe?totalCount=2&city=moscow",
		code:    http.StatusOK,
		message: "Мир кофе,Сладкоежка",
	},
	{
		request: "/cafe?city=moscow",
		code:    http.StatusBadRequest,
		message: "totalCount missing",
	},
	{
		request: "/cafe?totalCount=0x01&city=moscow",
		code:    http.StatusBadRequest,
		message: "wrong totalCount value",
	},
}

func TestMainHandlerStatus(t *testing.T) {
	for i, c := range statusCases {
		t.Run("test_#"+strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", c.request, nil)
			responseRecorder := httptest.NewRecorder()
			handler := http.HandlerFunc(mainHandle)
			handler.ServeHTTP(responseRecorder, req)
			response := responseRecorder.Body.String()
			assert.Equal(t, responseRecorder.Code, c.code)
			assert.NotEmpty(t, response)
			assert.Equal(t, response, c.message)
		})
	}
}

var notFoundCityCases = []string{
	"/cafe?totalCount=3&city=tomsk",
	"/cafe?totalCount=1&city=spb",
	"/cafe?totalCount=0&city=washington",
	"/cafe?totalCount=555&city=rostov",
	"/cafe?totalCount=12&city=someCity",
}

func TestMainHandlerNotFoundCity(t *testing.T) {
	for i, req := range notFoundCityCases {
		t.Run("test_#"+strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", req, nil)
			responseRecorder := httptest.NewRecorder()
			handler := http.HandlerFunc(mainHandle)
			handler.ServeHTTP(responseRecorder, req)
			require.Equal(t, responseRecorder.Code, http.StatusBadRequest)
			require.Equal(t, responseRecorder.Body.String(), "wrong city value")
		})
	}
}

var countMoreThanTotalCases = []struct {
	request    string
	totalCount int
}{
	{
		request:    "/cafe?totalCount=100&city=moscow",
		totalCount: 4,
	},
	{
		request:    "/cafe?totalCount=3&city=moscow",
		totalCount: 3,
	},
	{
		request:    "/cafe?totalCount=6&city=moscow",
		totalCount: 4,
	},
	{
		request:    "/cafe?totalCount=1&city=moscow",
		totalCount: 1,
	},
	{
		request:    "/cafe?totalCount=999999999&city=moscow",
		totalCount: 4,
	},
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	for i, c := range countMoreThanTotalCases {
		t.Run("test_#"+strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("GET", c.request, nil)
			responseRecorder := httptest.NewRecorder()
			handler := http.HandlerFunc(mainHandle)
			handler.ServeHTTP(responseRecorder, req)
			response := responseRecorder.Body.String()
			assert.Equal(t, responseRecorder.Code, http.StatusOK)
			assert.NotEmpty(t, response)
			assert.Len(t, strings.Split(response, ","), c.totalCount)
		})
	}
}
