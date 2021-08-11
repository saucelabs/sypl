package shared

import (
	"bytes"
	"encoding/json"
	"log"
)

// Prettify encodes data returning its JSON-stringified version.
//
// Note: Only exported fields of the data structure will be printed.
func Prettify(data interface{}) string {
	buf := new(bytes.Buffer)

	enc := json.NewEncoder(buf)
	enc.SetIndent("", "\t")

	if err := enc.Encode(data); err != nil {
		log.Println(ErrorPrefix, "prettify: Failed to encode data.", err)

		return ""
	}

	return buf.String()
}
