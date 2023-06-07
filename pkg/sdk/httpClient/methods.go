package httpClient

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Response struct {
	Response *http.Response
	Error    error
}

func JSONTransformer(response Response, model interface{}) error {
	if response.Error != nil {
		return response.Error
	}

	b, err := io.ReadAll(response.Response.Body)
	defer response.Response.Body.Close()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, model); err != nil {
		return err
	}

	return nil
}

func Get(url string, headers ...header) Response {
	client := NewClient(WithTimeout(10 * time.Second))
	request, err := NewRequest(request{
		Headers: headers,
		Url:     url,
		Method:  "get",
		Body:    nil,
	})

	if err != nil {
		return Response{Error: err}
	}

	res, err := Make(request, client)
	return Response{
		Response: res,
		Error:    err,
	}
}

func Post(url string, body []byte, headers ...header) Response {
	client := NewClient(WithTimeout(10 * time.Second))
	request, err := NewRequest(request{
		Headers: headers,
		Url:     url,
		Method:  "post",
		Body:    body,
	})

	if err != nil {
		return Response{Error: err}
	}

	res, err := Make(request, client)
	return Response{
		Response: res,
		Error:    err,
	}
}

func Put(url string, body []byte, headers ...header) Response {
	client := NewClient(WithTimeout(10 * time.Second))
	request, err := NewRequest(request{
		Headers: headers,
		Url:     url,
		Method:  "put",
		Body:    body,
	})

	if err != nil {
		return Response{Error: err}
	}

	res, err := Make(request, client)
	return Response{
		Response: res,
		Error:    err,
	}
}

func Patch(url string, body []byte, headers ...header) Response {
	client := NewClient(WithTimeout(10 * time.Second))
	request, err := NewRequest(request{
		Headers: headers,
		Url:     url,
		Method:  "patch",
		Body:    body,
	})

	if err != nil {
		return Response{Error: err}
	}

	res, err := Make(request, client)
	return Response{
		Response: res,
		Error:    err,
	}
}

func Head(url string, headers ...header) Response {
	client := NewClient(WithTimeout(10 * time.Second))
	request, err := NewRequest(request{
		Headers: headers,
		Url:     url,
		Method:  "head",
		Body:    nil,
	})

	if err != nil {
		return Response{Error: err}
	}

	res, err := Make(request, client)
	return Response{
		Response: res,
		Error:    err,
	}
}

func Delete(url string, headers ...header) Response {
	client := NewClient(WithTimeout(10 * time.Second))
	request, err := NewRequest(request{
		Headers: headers,
		Url:     url,
		Method:  "delete",
		Body:    nil,
	})

	if err != nil {
		return Response{Error: err}
	}

	res, err := Make(request, client)
	return Response{
		Response: res,
		Error:    err,
	}
}
