package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"contactapp/app"
	"contactapp/components/log"
	module "contactapp/modules"
	"contactapp/repository"
)

func main() {
	// Create New Instace of DB
	log := log.GetLogger()
	db := app.NewDBConnection(*log)

	if db == nil {
		log.Print("Db connection failed.")
	}
	defer func() {
		db.Close()
		log.Print("Db closed")
	}()
	var wg sync.WaitGroup
	var repository = repository.NewGormRepository()
	app := app.NewApp("Contact App", db, *log,
		&wg, repository)
	// Initialize app components.
	fmt.Println("1>>>>>>>>>>>>>>>>>>>>>>>>>>")

	app.Init()
	module.RegisterModuleRoutes(app, repository)

	fmt.Println("2>>>>>>>>>>>>>>>>>>>>>>>>>>")

	fmt.Println("3>>>>>>>>>>>>>>>>>>>>>>>>>>")

	// Need to make sure app starts within 60 seconds of deployment so heroku is able to find port.
	go func() {
		err := app.StartServer()
		if err != nil {
			stopApp(app)
		}
	}()
	module.Configure(app)

	// Stop Server On System Call or Interrupt.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	stopApp(app)

	fmt.Print("sub barabar hai")
}
func stopApp(app *app.App) {
	app.Stop()
	app.WG.Wait()
	log.GetLogger().Print("App stopped.")
	os.Exit(0)
}
