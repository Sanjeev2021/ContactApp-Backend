package modules

import (
	"fmt"

	"contactapp/app"
	"contactapp/models/user"
)

// Configure will config all modules
func Configure(appObj *app.App) {
	userModule := user.NewUserModuleConfig(appObj.DB)
	fmt.Println("hgegfe>>>>>>>>>>>>>>>>>>>>>>>>>>")

	appObj.MigrateTables([]app.ModuleConfig{userModule})
}
