package login

import (
	"custom-go/authentication"
	"custom-go/generated"
	"custom-go/pkg/types"
	"golang.org/x/exp/slices"
)

var userFoundRequiredLoginTypes = []generated.Casdoor_login_post_input_object_loginType_enum{
	generated.Casdoor_login_post_input_object_loginType_enum_sms,
	generated.Casdoor_login_post_input_object_loginType_enum_password,
}

func MutatingPostResolve(hook *types.HookRequest, body generated.Casdoor__loginBody) (resp generated.Casdoor__loginBody, err error) {
	userId, phone, err := authentication.GetUserinfoByToken(hook, body.Response.Data.Data.Data.AccessToken)
	if err != nil {
		return
	}
	if phone != "" {
		syncPhoneInput := generated.User__syncPhoneInternalInput{Id: userId, Phone: phone}
		_, _ = generated.User__syncPhone.Execute(syncPhoneInput, hook.InternalClient)
	}

	if phone != "" || body.Input.LoginType == generated.Casdoor_login_post_input_object_loginType_enum_sms {
		if err = authentication.CreateOneUser(hook, userId, phone); err != nil {
			return
		}
	}

	if _, err = authentication.GetUserinfo(hook.InternalClient, userId,
		slices.Contains(userFoundRequiredLoginTypes, body.Input.LoginType)); err != nil {
		return
	}

	body.Response.Data.Data.PhoneBound = phone != ""
	return body, nil
}
