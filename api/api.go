package api

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type APIHandler struct {
	BaseURL  string
	Headers  map[string]string
	Sessions []string
}

func NewAPIHandler(baseURL string, sessions []string, headers map[string]string) (*APIHandler, error) {
	if len(sessions) == 0 {
		return nil, errors.New("at least one session ID is required")
	}
	return &APIHandler{
		BaseURL:  baseURL,
		Headers:  headers,
		Sessions: sessions,
	}, nil
}

func (api *APIHandler) postRequest(endpoint string, data map[string]string) (map[string]interface{}, error) {
	fullURL := api.BaseURL + endpoint

	formData := url.Values{}
	for key, value := range data {
		formData.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", fullURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	for key, value := range api.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Warning: failed to close response body: %v", cerr)
		}
	}()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer func() {
			if cerr := reader.Close(); cerr != nil {
				log.Printf("Warning: failed to close gzip reader: %v", cerr)
			}
		}()
	default:
		reader = resp.Body
	}

	// Чтение тела ответа
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status code " + resp.Status)
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (api *APIHandler) CollectStars() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	for _, session := range api.Sessions {
		data := map[string]string{"session": session}
		response, err := api.postRequest("/galaxy/collect", data)
		if err != nil {
			log.Printf("Error collecting stars for session %s: %v", session, err)
			results = append(results, map[string]interface{}{
				"session": session,
				"error":   err.Error(),
			})
			continue
		}

		response["session"] = session
		results = append(results, response)
	}

	return results, nil
}

func (api *APIHandler) CheckStats() ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0)

	for _, session := range api.Sessions {
		data := map[string]string{"session": session}
		response, err := api.postRequest("/user/info", data)
		if err != nil {
			log.Printf("Error checking stats for session %s: %v", session, err)
			continue
		}
		results = append(results, response)
	}

	return results, nil
}

func (api *APIHandler) BuyStars(galaxyIDs []string, starsCount string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Проверяем, что длина galaxyIDs совпадает с количеством сессий
	if len(galaxyIDs) != len(api.Sessions) {
		return nil, fmt.Errorf("length of galaxyIDs (%d) must match number of sessions (%d)", len(galaxyIDs), len(api.Sessions))
	}

	for i, session := range api.Sessions {
		data := map[string]string{
			"galaxy_id": galaxyIDs[i], // Используем galaxyID для текущей сессии
			"session":   session,
			"stars":     starsCount,
		}

		response, err := api.postRequest("/stars/create", data)
		if err != nil {
			log.Printf("Error buying stars for session %s: %v", session, err)
			results = append(results, map[string]interface{}{
				"session":   session,
				"galaxy_id": galaxyIDs[i],
				"error":     err.Error(),
			})
			continue
		}

		// Добавляем session и galaxyID в ответ
		response["session"] = session
		response["galaxy_id"] = galaxyIDs[i]
		results = append(results, response)
	}

	return results, nil
}
