package httpClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"therebelsource/emulator/logger"
	"time"
)

func SendWithBackoff(request *http.Request, client *http.Client, model interface{}) error {
	var backoffSchedule = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		10 * time.Second,
	}

	var res *http.Response
	var err error

	for _, backoff := range backoffSchedule {
		res, err = Make(request, client)

		if err != nil {
			logger.Warn(fmt.Sprintf("Request backoff failed with error: %s"), err.Error())
			time.Sleep(backoff)

			continue
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			msg := fmt.Sprintf("Request to %s returned status %d", request.URL, res.StatusCode)
			logger.Warn(msg)
			return errors.New(msg)
		}

		return unpack(res, model)
	}

	return errors.New("Unable to send request. All backoffs failed")
}

func unpack(res *http.Response, model interface{}) error {
	body, err := ioutil.ReadAll(res.Body)

	var apiResponse map[string]interface{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Warn(fmt.Sprintf("Could not unpack response body into map[string]interface{}: %s", err.Error()))
		return err
	}

	d := apiResponse["data"]

	b, err := json.Marshal(d)

	if err != nil {
		logger.Warn(fmt.Sprintf("Could not marshal body data: %s", err.Error()))
		return err
	}

	if err := json.Unmarshal(b, model); err != nil {
		logger.Warn(fmt.Sprintf("Could not unmarshal body data into model: %s", err.Error()))
		return err
	}

	return nil
}
