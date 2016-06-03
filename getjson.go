package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
)

func GetJSON(URL string, result interface{}) (err error) {
	var resp *http.Response
	resp, err = http.Get(URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
		return
	}

	mediatype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return
	}
	if mediatype != "application/json" {
		err = fmt.Errorf("Unexpected Content-Type %s", mediatype)
		return
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(result)

	return
}
