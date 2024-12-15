package admin

import (
	"context"
	"github.com/guneyin/printhub/market"
	"github.com/guneyin/printhub/model"
	"github.com/guneyin/printhub/service/tenant"
	"github.com/guneyin/printhub/service/user"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var (
	svc       *Service
	userSvc   *user.Service
	tenantSvc *tenant.Service
)

const (
	tenantEmail = "<EMAIL>"
	tenantName  = "<NAME>"
)

func init() {
	market.InitTestMarket()
	svc = GetService()
	userSvc = user.GetService()
	tenantSvc = tenant.GetService()
}

func TestService_TenantCreate(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	ctx := context.Background()
	Convey("TestService_TenantCreate", t, func() {
		created := &model.Tenant{Email: tenantEmail, Name: tenantName}
		err := svc.TenantCreate(ctx, created)
		So(err, ShouldBeNil)
		So(created.UUID, ShouldNotBeEmpty)

		u, err := userSvc.GetByEmail(ctx, tenantEmail, model.UserRoleTenant)
		So(err, ShouldBeNil)
		So(u.UUID, ShouldNotBeEmpty)

		found, err := tenantSvc.GetByUUID(ctx, created.UUID)
		So(err, ShouldBeNil)
		So(found, ShouldNotBeNil)
		So(found.UUID, ShouldEqual, created.UUID)
	})
}
