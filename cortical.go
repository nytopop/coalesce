// coalesce/cortical.go

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func GetKeywordsForText(key, text string) []string {
	//key := "255035d0-8ce2-11e6-a057-97f4c970893c"
	url := "http://api.cortical.io:80/rest/text/keywords?retina_name=en_associative"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(text)))
	if err != nil {
		return []string{}
	}

	// add api key to header
	req.Header = map[string][]string{
		"api-key": {key},
	}

	// create client
	client := &http.Client{}

	// send POST
	resp, err := client.Do(req)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	// decode response to []string
	var keywords []string
	if err := json.NewDecoder(resp.Body).Decode(&keywords); err != nil {
		return []string{}
	}

	return keywords
}
