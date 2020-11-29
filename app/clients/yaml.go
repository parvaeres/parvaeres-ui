package clients

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

func Yaml2Json(in string) (out string, err error) {
	var body interface{}
	if err := yaml.Unmarshal([]byte(in), &body); err != nil {
		panic(err)
	}

	body = convert(body)

	if b, err := json.Marshal(body); err != nil {
		panic(err)
	} else {
		out = string(b)
	}
	return
}