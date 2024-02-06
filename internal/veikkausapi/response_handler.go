package veikkausapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ResponseCodeIsOk(statusCode int) bool {
	if statusCode >= HttpOk && statusCode < HttpMultipleChoices {
		return true
	}

	return false
}

func HandleResponse(response *http.Response, result interface{}) error {
	body, err := io.ReadAll(response.Body)

	defer response.Body.Close()

	if err != nil {
		return fmt.Errorf("Error reading the response body: %s", err.Error())
	}

	if err = json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("Error unmarshaling response body: %s", err.Error())
	}

	return nil
}
