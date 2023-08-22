package protocol

import "encoding/json"

type jsonCoder struct {
	CoderType CodeType
}

func (j *jsonCoder) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (j *jsonCoder) Unmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
func (j *jsonCoder) ToString() string {
	return "Json"
}
