package veikkausapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ResponseCodeIsOk(response *http.Response) bool {
	if c := response.StatusCode; c >= http.StatusOK && c < http.StatusMultipleChoices {
		return true
	} else {
		return false
	}
}

func HandleError(response *http.Response) error {
	if statusCode := response.StatusCode; statusCode < http.StatusBadRequest {
		return &UnsupportedStatusCodeError{Code: statusCode}
	}

	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return fmt.Errorf("could not convert response to a byte-stream: %v", err)
	}

	return ParseAPIError(body)
}

func HandleResponse(response *http.Response, responseInterface interface{}) error {
	body, err := io.ReadAll(response.Body)

	defer response.Body.Close()

	if err != nil {
		return fmt.Errorf("error reading the response body: %s", err.Error())
	}

	if err = json.Unmarshal(body, &responseInterface); err != nil {
		return fmt.Errorf("error unmarshaling response body: %s", err.Error())
	}

	return nil
}
