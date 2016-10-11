// coalesce/cortical.go

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Access cortical.io api to generate tags
func GetKeywordsForText(key, text string) ([]string, error) {
	url := "http://api.cortical.io:80/rest/text/keywords?retina_name=en_associative"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(text)))
	if err != nil {
		return []string{}, err
	}

	req.Header = map[string][]string{
		"api-key": {key},
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	var keywords []string
	if err := json.NewDecoder(resp.Body).Decode(&keywords); err != nil {
		return []string{}, err
	}

	return keywords, nil
}
