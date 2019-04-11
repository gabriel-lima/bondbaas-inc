package handler

import (
	"bondbaas/model"
	"bondbaas/presenter"
	"encoding/json"
	"io"
	"io/ioutil"
)

func SerializePayload(presenter presenter.GenericPresenter, body io.ReadCloser) (payload map[string]interface{}, err error) {
	raw, err := ioutil.ReadAll(body)
	if err != nil {
		presenter.Fail(422, err.Error())
		return
	}

	err = json.Unmarshal([]byte(raw), &payload)
	if err != nil {
		presenter.Fail(422, err.Error())
		return
	}

	return payload, err
}

func SerializeTableModel(presenter presenter.GenericPresenter, body io.ReadCloser) (payload model.TableModel, err error) {
	payload = model.TableModel{}
	err = json.NewDecoder(body).Decode(&payload)
	if err != nil {
		presenter.Fail(422, err.Error())
		return
	}

	return payload, err
}
