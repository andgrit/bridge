package misc

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
	"io"
)

type UnexpectedError struct {
	ErrorMessage string
}

const MarshalError = `{"errorMessage":"marshal error"}`
func UnexpectedErrorToJson(err error) string {
	unexpectedError := UnexpectedError{ErrorMessage: err.Error()}
	ret, err := json.Marshal(unexpectedError)
	if err != nil {
		return MarshalError
	}
	return string(ret)
}

// WriteResponse will write the error or the response
func WriteResponse(responseWriter http.ResponseWriter, response interface{}, err error) {
	var sendString string
	if err != nil {
		sendString = UnexpectedErrorToJson(err)
	} else {
		ret, err := json.Marshal(&response)
		if err != nil {
			sendString = UnexpectedErrorToJson(err)
		} else {
			sendString = string(ret) // yeah!
		}
	}
	fmt.Fprintf(responseWriter, sendString)
}

func RequestBodyUnmarshal(request *http.Request, result interface{}) error {
	return ReadCloserUnmarshal(request.Body, result)
}

func ResponseBodyUnmarshal(response *http.Response, result interface{}) error {
	return ReadCloserUnmarshal(response.Body, result)
}

func ReadCloserUnmarshal(body io.ReadCloser, result interface{}) error {
	s, err := ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		return err
	}

	return json.Unmarshal(s, result)
}

