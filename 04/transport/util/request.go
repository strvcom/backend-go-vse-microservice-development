package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func UnmarshalRequest(r *http.Request, b any) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, b); err != nil {
		return err
	}

	if err := validate.Struct(b); err != nil {
		return err
	}

	return nil
}
