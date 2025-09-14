package authentication

import (
	"custom-go/generated"
	"custom-go/pkg/plugins"
	"custom-go/pkg/types"
	"errors"
	"fmt"
)

func MutatingPostAuthentication(hook *types.AuthenticationHookRequest) (resp *plugins.AuthenticationResponse, err error) {
	roleCodes, err := GetUserinfo(hook.InternalClient, hook.User.UserId, true)
	if err != nil {
		return
	}
	if len(roleCodes) > 0 {
		hook.User.Roles = roleCodes
	}
	return &plugins.AuthenticationResponse{User: hook.User, Status: "ok"}, nil
}

func GetUserinfo(client *types.InternalClient, userId string, userFoundRequired bool) (roleCodes []string, err error) {
	findInput := generated.Admin__user__findUniqueInternalInput{Id: userId}
	findResp, err := generated.Admin__user__findUnique.Execute(findInput, client)
	if err != nil {
		return
	}
	if findResp.Data.Id == "" && userFoundRequired {
		err = errors.New("用户不存在")
		return
	}
	for _, role := range findResp.Data.Roles {
		roleCodes = append(roleCodes, role.Code)
	}
	return
}

func GetUserinfoByToken(hook *types.BaseRequestContext, accessToken string) (userId, phone string, err error) {
	casdoorUserinfoInput := generated.Casdoor__userinfoInternalInput{
		AccessToken: fmt.Sprintf("Bearer %s", accessToken),
	}
	casdoorUserinfoResp, err := generated.Casdoor__userinfo.Execute(casdoorUserinfoInput, hook.InternalClient)
	if err != nil {
		return
	}
	userinfoData := casdoorUserinfoResp.Data
	if userinfoData.UserId == "" {
		err = errors.New("userinfo未返回userId，请检查OIDC版本")
		return
	}

	userId, phone = userinfoData.UserId, userinfoData.Phone
	return
}

func CreateOneUser(hook *types.BaseRequestContext, userId, phone string) (err error) {
	userIsExistedInput := generated.User__isExistedInternalInput{Id: userId}
	userIsExistedResp, err := generated.User__isExisted.Execute(userIsExistedInput, hook.InternalClient)
	if err != nil || userIsExistedResp.Data.Id != "" {
		return
	}

	userCreateOneInput := generated.User__createOneInternalInput{Id: userId, Phone: phone}
	_, err = generated.User__createOne.Execute(userCreateOneInput, hook.InternalClient)
	return
}
