package slack

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"therebelsource/emulator/appErrors"
	"therebelsource/emulator/httpClient"
)

var channels = map[string]map[string]string{
	"test": {
		"dev_critical_log": "https://hooks.slack.com/services/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",
		"dev_deploy_log":   "https://hooks.slack.com/services/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",
	},
	"dev": {
		"dev_critical_log": "https://hooks.slack.com/services/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",
		"dev_deploy_log":   "https://hooks.slack.com/services/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",
	},
	"prod": {
		"prod_critical_log": "https://hooks.slack.com/services/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",
		"prod_deploy_log":   "https://hooks.slack.com/services/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",
	},
	"staging": {
		"staging_critical_log": "https://hooks.slack.com/services/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",
		"staging_deploy_log":   "https://hooks.slack.com/services/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",
	},
}

func getUrl(channel string) string {
	env := os.Getenv("APP_ENV")
	if os.Getenv("APP_ENV") == "test" {
		env = "dev"
	}

	return channels[os.Getenv("APP_ENV")][fmt.Sprintf("%s_%s", env, channel)]
}

func createSimpleLogBody(title string, log string) map[string]interface{} {
	return map[string]interface{}{
		"text": "An event happened\n",
		"blocks": []map[string]interface{}{
			{
				"type": "header",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": title,
				},
			},
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": log,
				},
			},
		},
	}
}

func createErrorLogBody(err *appErrors.Error) map[string]interface{} {
	return map[string]interface{}{
		"text": "An error occurred\n",
		"blocks": []map[string]interface{}{
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("Master Code: %d", err.MasterCode),
				},
			},
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("Code: %d", err.Code),
				},
			},
			{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("Message: %s", err.Error()),
				},
			},
		},
	}
}

func SendLog(title string, log string, channel string) *appErrors.Error {
	client, err := httpClient.NewHttpClient(&tls.Config{})

	if err != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Failed writing to slack")
	}

	b, err := json.Marshal(createSimpleLogBody(title, log))

	if err != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Failed writing to slack")
	}

	response, clientErr := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    getUrl(channel),
		Method: "POST",
		Body:   b,
	})

	if clientErr != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Failed writing to slack")
	}

	if response.Status != 200 {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Failed writing to slack")
	}

	return nil
}

func SendErrorLog(error *appErrors.Error, channel string) *appErrors.Error {
	client, err := httpClient.NewHttpClient(&tls.Config{})

	if err != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Failed writing to slack")
	}

	b, err := json.Marshal(createErrorLogBody(error))

	if err != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Failed writing to slack")
	}

	response, clientErr := client.MakeJsonRequest(&httpClient.JsonRequest{
		Url:    getUrl(channel),
		Method: "POST",
		Body:   b,
	})

	if clientErr != nil {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Failed writing to slack")
	}

	if response.Status != 200 {
		return appErrors.New(appErrors.ApplicationError, appErrors.ApplicationRuntimeError, "Failed writing to slack")
	}

	return nil
}
