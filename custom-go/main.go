package main

import (
	_ "custom-go/crontab"
	"custom-go/server"
	_ "github.com/joho/godotenv"
	_ "github.com/shopspring/decimal"
	// 根据需求，开启注释
	_ "custom-go/customize"
	_ "custom-go/function/admin/device"
	_ "custom-go/function/assistant/device"
	_ "custom-go/iot/defaulted"
	_ "custom-go/proxy/admin/user"
)

func main() {
	server.Execute()
}
