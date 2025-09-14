package updateOne

import (
	"custom-go/generated"
	"custom-go/pkg/types"
)

func MutatingPreResolve(_ *types.HookRequest, body generated.Admin__role__updateOneBody) (resp generated.Admin__role__updateOneBody, err error) {
	if linkMenuIds := body.Input.LinkMenuIds; len(linkMenuIds) > 0 {
		var createDatas []*generated.Admin_adminRole2MenuCreateManyAdminRoleInput
		for _, item := range linkMenuIds {
			createDatas = append(createDatas, &generated.Admin_adminRole2MenuCreateManyAdminRoleInput{MenuId: item})
		}
		body.Input.CreateManyMenuRole = &generated.Admin_adminRole2MenuCreateManyAdminRoleInputEnvelope{Data: createDatas}
	}
	return body, nil
}
