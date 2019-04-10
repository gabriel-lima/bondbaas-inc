package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func getPayload(request *http.Request) (payload map[string]interface{}, err error) {
	var raw []byte
	raw, _ = ioutil.ReadAll(request.Body)

	err = json.Unmarshal([]byte(raw), &payload)

	return payload, err
}
