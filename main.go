package main

import (
	"fmt"

	"myapi/routes"

	"github.com/pilinux/gorest/config"
	"github.com/pilinux/gorest/database"
	"github.com/pilinux/gorest/lib/middleware"
)

var configure = config.Config()

func main() {
	if configure.Database.RDBMS.Activate == "yes" {
		// Initialize RDBMS client
		if err := database.InitDB().Error; err != nil {
			fmt.Println(err)
			return
		}
	}

	// JWT
	middleware.AccessKey = []byte(configure.Security.JWT.AccessKey)
	middleware.AccessKeyTTL = configure.Security.JWT.AccessKeyTTL
	middleware.RefreshKey = []byte(configure.Security.JWT.RefreshKey)
	middleware.RefreshKeyTTL = configure.Security.JWT.RefreshKeyTTL
	router, err := routes.SetupRouter(configure)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = router.Run(":" + configure.Server.ServerPort)
	if err != nil {
		fmt.Println(err)
		return
	}
}
