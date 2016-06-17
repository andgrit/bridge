package controllers

import (
	"encoding/base64"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/kr/pretty"
	"bytes"
	"encoding/gob"
)

// decode decodes a cookie using base64.
func decode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}

func deserialize(src []byte, dst interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(src))
	err := dec.Decode(dst)
	return err
}

func TestDecode(t *testing.T) {
	// This DATA was noticed in the browser looking at the cookie:
	const DATA = "MTQ2NTk0MDYzOXxEdi1CQkFFQ180SUFBUkFCRUFBQU1QLUNBQUlHYzNSeWFXNW5EQVVBQTJadmJ3WnpkSEpwYm1jTUJRQURZbUZ5QTJsdWRBUUNBRlFEYVc1MEJBSUFWZz09fCBvLlCKgoLSqOkzaWps_XWqkjtSjvAMfjkMQ8NAdg66"

	b, err := decode([]byte(DATA)) // pull out the base64 encoding to get a 3 part string
	pretty.Println(len(DATA), DATA)
	assert.NoError(t, err)
	parts := bytes.SplitN(b, []byte("|"), 3)

	// parts: value
	// 0: timestamp
	// 1: session data
	// 2: hash code of the session data

	// session data is base64 encoded as well
	b, err = decode(parts[1])
	pretty.Println("decoded:", len(string(b)), string(b))

	// session data is v map marshaled via gob
	var values map[interface{}]interface{}
	deserialize(b, &values)
	pretty.Println("values:", values)
}
