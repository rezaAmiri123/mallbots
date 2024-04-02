package amserializer

import (
	"encoding/json"

	"github.com/rezaAmiri123/edatV2/am"
)

type JsonSerializer struct {}

var _ am.MessageSerializer = (*JsonSerializer)(nil)

func NewJsonSerializer() am.MessageSerializer {
	return JsonSerializer{}
}

func (s JsonSerializer) Serialize(message am.SerializerMessageData) ([]byte, error) {
	return json.Marshal(message)
}

func (s JsonSerializer) Deserialize(message []byte) (am.SerializerMessageData, error) {
	var mesageData am.SerializerMessageData
	err := json.Unmarshal(message, &mesageData)

	return mesageData, err

}
