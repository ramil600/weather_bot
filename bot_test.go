package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) {
	result := []byte(`{
		"ok": true,
		"result": [
			{
				"update_id": 801052952,
				"message": {
					"message_id": 871,
					"from": {
						"id": 5708402489,
						"is_bot": false,
						"first_name": "Ro",
						"last_name": "Mil",
						"language_code": "en"
					},
					"chat": {
						"id": 5708402489,
						"first_name": "Ro",
						"last_name": "Mil",
						"type": "private"
					},
					"date": 1663838199,
					"text": "about"
				}
			},
			{
				"update_id": 801052953,
				"message": {
					"message_id": 877,
					"from": {
						"id": 5708402489,
						"is_bot": false,
						"first_name": "Ro",
						"last_name": "Mil",
						"language_code": "en"
					},
					"chat": {
						"id": 5708402489,
						"first_name": "Ro",
						"last_name": "Mil",
						"type": "private"
					},
					"date": 1663842306,
					"text": "about"
				}
			}
		]
	}`)

	var apiresp APIResponseUpdates
	err := json.Unmarshal(result, &apiresp)

	assert.NoError(t, err)

	fmt.Printf("%+v\n", apiresp)

	var apiresp2 APIResponse
	err = json.Unmarshal(result, &apiresp2)
	assert.NoError(t, err)

	fmt.Printf("%T\n", apiresp2.Result)

	// var updates map[string]interface{}

	// err = json.Unmarshal(data, &updates)
	// assert.NoError(t, err)

	fmt.Printf("%+v\n", apiresp2.Result)


	}

}
