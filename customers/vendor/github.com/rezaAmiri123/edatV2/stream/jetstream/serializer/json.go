package amserializer

import (
	"encoding/json"

	"github.com/rezaAmiri123/edatV2/stream/jetstream"
)

type JsonSerializer struct {}

var _ jetstream.MessageSerializer = (*JsonSerializer)(nil)

func NewJsonSerializer() jetstream.MessageSerializer {
	return JsonSerializer{}
}

func (s JsonSerializer) Serialize(message jetstream.MessageSerializerData) ([]byte, error){
	return json.Marshal(message)
}

func (s JsonSerializer) Deserialize(message []byte) (jetstream.MessageSerializerData, error){
	var mesageData jetstream.MessageSerializerData
	err := json.Unmarshal(message, &mesageData)

	return mesageData, err
}
