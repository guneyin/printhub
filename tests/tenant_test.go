package test

import (
	"encoding/json"
	"fmt"
	"github.com/guneyin/printhub/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestTenantAuthRoute(t *testing.T) {
	googleConfigData := map[string]string{
		"apiKey":       "AIzaSyBE_hb2H1ax1SCzXFPrbJ5ouIz87mqc_5c",
		"clientId":     "1070061319069-aoa9jikah08q3a7gjoqtv5d092du24rc.apps.googleusercontent.com",
		"clientSecret": "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",
		"callBackUrl":  "http://localhost:8080/api/tenant/auth/callback",
	}
	googleConfig, _ := json.Marshal(googleConfigData)

	conf := model.ConfigList{model.Config{
		Key:   "google:config",
		Value: string(googleConfig),
	}}

	tests := []testCase{{
		skip:               false,
		description:        "set google drive config",
		method:             http.MethodPut,
		route:              "http://localhost:8081/api/tenant/config",
		body:               newBody(conf),
		expectedStatusCode: http.StatusOK,
	}, {
		skip:        false,
		printOutput: true,
		description: "init google auth",
		method:      http.MethodGet,
		route:       "http://localhost:8081/api/tenant/disk/auth?provider=google",
		//body:               conf.JSON(),
		expectedStatusCode: http.StatusFound,
	},
	}

	for _, test := range tests {
		if test.skip {
			continue
		}

		req, _ := http.NewRequest(test.method, test.route, test.body.toReader())
		res, err := testApp.Test(req, -1)
		assert.Nil(t, err)
		assert.Equal(t, test.expectedStatusCode, res.StatusCode)

		if test.printOutput {
			fmt.Println(res.Body)
		}
	}
}
