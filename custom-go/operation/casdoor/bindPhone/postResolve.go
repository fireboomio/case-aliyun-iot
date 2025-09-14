package bindPhone

import (
	"custom-go/authentication"
	"custom-go/generated"
	"custom-go/pkg/types"
	"strings"
)

func PostResolve(hook *types.HookRequest, body generated.Casdoor__bindPhoneBody) (resp generated.Casdoor__bindPhoneBody, err error) {
	resp = body
	if body.Input.IsSwitchAction {
		return
	}

	accessToken := body.Response.Data.Data.Data.AccessToken
	if accessToken == "" {
		accessToken = strings.TrimPrefix(body.Input.Authorization, "Bearer ")
	}
	userId, _, err := authentication.GetUserinfoByToken(hook, accessToken)
	if err != nil {
		return
	}
	err = authentication.CreateOneUser(hook, userId, body.Input.Phone)
	return
}
