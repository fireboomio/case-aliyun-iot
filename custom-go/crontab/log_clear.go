package crontab

import (
	"custom-go/generated"
	"custom-go/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
)

func init() {
	types.AddRegisteredHook(func(echo.Logger) {
		c := cron.New()
		_, _ = c.AddFunc("1 0 * * *", logClear)
		c.Start()
	})
}

func logClear() {
	client := types.NewEmptyInternalClient()
	logDeleteInput := generated.Admin__log__deleteManyWithScopeInternalInput{}
	_, _ = generated.Admin__log__deleteManyWithScope.Execute(logDeleteInput, client)
}
