package http

import (
	"encoding/json"
	"testing"
)

func TestResty(t *testing.T) {
	var body map[string]any
	_, err := ShuntHttp.R().SetResult(&body).Get(`http://10.0.1.50:5681/meshes/default/dataplanes+insights`)
	if err != nil {
		t.Log(err)
		return
	}

	data, err := json.MarshalIndent(body, "", "    ")

	t.Log(string(data))
}
