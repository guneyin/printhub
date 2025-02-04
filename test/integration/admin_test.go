package integration_test

import (
	"net/http"
	"testing"

	"github.com/guneyin/printhub/test/integration"
	"github.com/stretchr/testify/require"

	"github.com/guneyin/printhub/model"
	"github.com/stretchr/testify/assert"
)

func TestAdminRoutes(t *testing.T) {
	app := integration.NewTestApp()

	tests := []integration.Case{{
		Skip:        false,
		Description: "create tenant",
		Method:      http.MethodPost,
		Route:       "http://localhost:8081/api/admin/tenant",
		Body: integration.NewBody(&model.Tenant{
			Email: "foo@bar.com",
			Name:  "Foo Photo Studio",
		}),
		ExpectedStatusCode: http.StatusCreated,
	},
	}

	for _, test := range tests {
		if test.Skip {
			continue
		}

		data := test.Body.ToReader()
		req, _ := http.NewRequest(test.Method, test.Route, data)
		res, err := app.Test(req, -1)
		require.NoError(t, err)
		assert.Equal(t, test.ExpectedStatusCode, res.StatusCode)

		if test.PrintOutput {
			t.Log(res.Body)
		}
	}
}
