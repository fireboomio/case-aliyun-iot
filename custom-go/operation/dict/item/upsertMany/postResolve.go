package upsertMany

import (
	"custom-go/generated"
	"custom-go/pkg/types"
	"sync"
)

func PostResolve(hook *types.HookRequest, body generated.Dict__item__upsertManyBody) (resp generated.Dict__item__upsertManyBody, err error) {
	invoke(body.Response.Data.Data.Code, hook.InternalClient)
	return body, nil
}

var updateDictHooks sync.Map

func Subscribe(code string, hook func(*types.InternalClient)) {
	var hooks []func(*types.InternalClient)
	if v, ok := updateDictHooks.Load(code); ok {
		hooks = append(v.([]func(*types.InternalClient)), hook)
	} else {
		hooks = []func(*types.InternalClient){hook}
	}
	updateDictHooks.Store(code, hooks)
}

func invoke(code string, client *types.InternalClient) {
	v, ok := updateDictHooks.Load(code)
	if !ok {
		return
	}
	for _, hook := range v.([]func(*types.InternalClient)) {
		hook(client)
	}
}
