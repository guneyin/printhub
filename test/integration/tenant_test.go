package integration_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/guneyin/printhub/test/integration"
	"github.com/stretchr/testify/require"

	"github.com/guneyin/printhub/model"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"
)

func TestTenantAuthRoute(t *testing.T) {
	app := integration.NewTestApp()

	googleConfigData := map[string]string{
		"apiKey":       os.Getenv("PH_GOOGLE_APIKEY"),
		"clientId":     os.Getenv("PH_GOOGLE_CLIENT_ID"),
		"clientSecret": os.Getenv("PH_GOOGLE_CLIENT_SECRET"),
		"callBackUrl":  "http://localhost:8080/api/tenant/auth/callback",
	}
	googleConfig, _ := json.Marshal(googleConfigData)

	conf := model.ConfigList{model.Config{
		Key:   "google:config",
		Value: string(googleConfig),
	}}

	tests := []integration.Case{{
		Skip:               false,
		Description:        "set google drive config",
		Method:             http.MethodPut,
		Route:              "http://localhost:8081/api/tenant/config",
		Body:               integration.NewBody(conf),
		ExpectedStatusCode: http.StatusOK,
	}, {
		Skip:        false,
		PrintOutput: true,
		Description: "init google auth",
		Method:      http.MethodGet,
		Route:       "http://localhost:8081/api/tenant/disk/auth?provider=google",
		//body:               conf.JSON(),
		ExpectedStatusCode: http.StatusFound,
	},
	}

	for _, test := range tests {
		if test.Skip {
			continue
		}

		req, _ := http.NewRequest(test.Method, test.Route, test.Body.ToReader())
		res, err := app.Test(req, -1)
		require.NoError(t, err)
		assert.Equal(t, test.ExpectedStatusCode, res.StatusCode)

		if test.PrintOutput {
			t.Log(res.Body)
		}
	}
}
