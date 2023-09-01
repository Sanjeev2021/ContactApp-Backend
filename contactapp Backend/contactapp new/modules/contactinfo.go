package modules

import (
	"contactapp/app"
	"contactapp/components/contactinfo/controller"
	"contactapp/components/contactinfo/service"
	"contactapp/repository"

)

func registerContactInfoRoutes(appObj *app.App, repository repository.Repository) {
	defer appObj.WG.Done()
	contactInfoService := service.NewContactInfoService(appObj.DB, repository)

	contactInfoController := controller.NewContactInfoController(contactInfoService, appObj.Log)

	appObj.RegisterControllerRoutes([]app.Controller{
		contactInfoController,
	})
}
