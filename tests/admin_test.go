package test

import (
	"fmt"
	"github.com/guneyin/printhub/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAdminRoutes(t *testing.T) {
	tests := []testCase{{
		skip:        false,
		description: "create tenant",
		method:      http.MethodPost,
		route:       "http://localhost:8081/api/admin/tenant",
		body: newBody(&model.Tenant{
			Email: "foo@bar.com",
			Name:  "Foo Photo Studio",
		}),
		expectedStatusCode: http.StatusCreated,
	},
	}

	for _, test := range tests {
		if test.skip {
			continue
		}

		data := test.body.toReader()
		req, _ := http.NewRequest(test.method, test.route, data)
		res, err := testApp.Test(req, -1)
		assert.Nil(t, err)
		assert.Equal(t, test.expectedStatusCode, res.StatusCode)

		if test.printOutput {
			fmt.Println(res.Body)
		}
	}
}
