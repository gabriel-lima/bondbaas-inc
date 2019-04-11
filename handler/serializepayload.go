package handler

import (
	"bondbaas/model"
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

func getPayloadSchema(request *http.Request) (model.TableModel, error) {
	table := model.TableModel{}
	err := json.NewDecoder(request.Body).Decode(&table)
	return table, err
}
