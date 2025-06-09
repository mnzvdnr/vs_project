package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		// пока сравнивать не будем, а просто выведем ответы
		// удалите потом этот вывод
		fmt.Println(response.Body.String())
	}
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	url := "/cafe?city=moscow&count="

	requests := []struct {
		count string // передаваемое значение count
		want  int    // ожидаемое количество кафе в ответе
	}{
		{"0", 0},
		{"1", 1},
		{"2", 2},
		{"100", 100},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url+v.count, nil)
		handler.ServeHTTP(response, req)

		//проверка кода ответа
		require.Equal(t, http.StatusOK, response.Code)

		//проверка тела ответа
		var arrayBody []string
		stringBody := response.Body.String()
		if stringBody != "" {
			arrayBody = strings.Split(stringBody, ",")
		} else {
			arrayBody = []string{}
		}

		if len(arrayBody) <= 2 || len(arrayBody) == 100 {
			assert.Equal(t, v.want, len(arrayBody))
		} else {
			assert.Equal(t, min(v.want, len(arrayBody)), len(arrayBody))
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	url := "/cafe?city=moscow&search="

	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе в ответе
	}{
		{"фасоль", 0},
		{"коФе", 2},
		{"вИлка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url+v.search, nil)
		handler.ServeHTTP(response, req)

		//проверка кода ответа
		require.Equal(t, http.StatusOK, response.Code)

		//проверка тела ответа
		var arrayBody []string
		stringBody := response.Body.String()
		//проверка содержимого
		if stringBody != "" {
			arrayBody = strings.Split(stringBody, ",")
			for _, name := range arrayBody {
				assert.True(t, strings.Contains(strings.ToLower(name), strings.ToLower(v.search)))
			}
		} else {
			arrayBody = []string{}
		}
		//проверка кол-ва
		assert.Equal(t, v.wantCount, len(arrayBody))
	}

}
