package global

import (
	"custom-go/generated"
	"custom-go/pkg/plugins"
	"custom-go/pkg/types"
	"math"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

const superRoleCodeKey = "ANT_ADMIN_ROLE_CODE"

func BeforeOriginRequest(hook *types.HttpTransportHookRequest, body *plugins.HttpTransportBody) (*types.WunderGraphRequest, error) {
	requestLogMap.Store(getRequestLogId(hook), &requestLog{
		start: time.Now(),
		ip:    strings.Split(body.Request.Headers["X-Forwarded-For"], ",")[0],
		ua:    body.Request.Headers["User-Agent"],
		body:  body.Request.OriginBody,
	})
	requestURL, _ := url.Parse(body.Request.RequestURI)
	if requestURL != nil && !hasSuperRoleCode(hook.User) {
		matchPath := strings.TrimPrefix(requestURL.Path, "/operations/")
		role2apisInput := generated.Admin__role__api__findManyInternalInput{Path: matchPath, Take: math.MaxInt16}
		role2apisResp, _ := generated.Admin__role__api__findMany.Execute(role2apisInput, hook.InternalClient)
		if len(role2apisResp.Data) > 0 {
			body.Request.Headers[string(types.RbacHeader_x_rbac_requireMatchAny)] = strings.Join(role2apisResp.Data, ",")
		}
	}
	body.Request.Headers["X-Permission"] = "user"
	return body.Request, nil
}

func hasSuperRoleCode(user *types.User) bool {
	superRoleCode := os.Getenv(superRoleCodeKey)
	return len(superRoleCode) > 0 && user != nil && slices.Contains(user.Roles, superRoleCode)
}
