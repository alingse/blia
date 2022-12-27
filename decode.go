package blia

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/schema"
)

func DecodeBody(r io.ReadCloser, outPtr Validator) (string, error) {
	defer r.Close()
	body, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	bodyStr := string(body)
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(outPtr)
	if err != nil {
		return bodyStr, ErrDecodeBodyFailed
	}
	return bodyStr, outPtr.Validate()
}

var queryDecoder *schema.Decoder

func init() {
	queryDecoder = schema.NewDecoder()
	queryDecoder.IgnoreUnknownKeys(true)
	queryDecoder.SetAliasTag("json")
}

func DecodeQuery(r *http.Request, outPtr Validator) error {
	err := queryDecoder.Decode(outPtr, r.URL.Query())
	if err != nil {
		return ErrDecodeQueryFailed
	}
	return outPtr.Validate()
}
