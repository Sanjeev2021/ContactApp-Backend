package modules

import (
	"contactapp/app"
	"contactapp/repository"
)

func RegisterModuleRoutes(app *app.App, repo *repository.GormRepository) {
	log := app.Log
	log.Print("=======Registering module routes============")
	app.WG.Add(3)
	go registerUserRoutes(app, repo)
	go registerContactRoutes(app, repo)
	go registerContactInfoRoutes(app, repo)

	app.WG.Wait()
}
