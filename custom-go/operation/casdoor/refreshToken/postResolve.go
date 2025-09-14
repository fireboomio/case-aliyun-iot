package refreshToken

import (
	"custom-go/authentication"
	"custom-go/generated"
	"custom-go/pkg/types"
)

func PostResolve(hook *types.HookRequest, body generated.Casdoor__refreshTokenBody) (resp generated.Casdoor__refreshTokenBody, err error) {
	userId, phone, err := authentication.GetUserinfoByToken(hook, body.Response.Data.Casdoor_refreshToken_post.Data.AccessToken)
	if err != nil {
		return
	}
	if phone != "" {
		syncPhoneInput := generated.User__syncPhoneInternalInput{Id: userId, Phone: phone}
		_, _ = generated.User__syncPhone.Execute(syncPhoneInput, hook.InternalClient)
	}
	return
}
