package app

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"contactapp/components/log"
	"contactapp/repository"

)

type App struct {
	sync.Mutex
	Name   string
	Router *mux.Router
	DB     *gorm.DB
	Log    log.Log

	Server *http.Server
	WG     *sync.WaitGroup

	Repository repository.Repository
}
type ModuleConfig interface {
	TableMigration(wg *sync.WaitGroup)
}

// Controller is implemented by the controllers.
type Controller interface {
	RegisterRoutes(router *mux.Router)
}

// NewApp returns app.
func NewApp(name string, db *gorm.DB, log log.Log,
	wg *sync.WaitGroup, repo repository.Repository) *App {
	return &App{
		Name: name,
		DB:   db,
		Log:  log,

		WG: wg,

		Repository: repo,
	}
}
func NewDBConnection(log log.Log) *gorm.DB {
	url := "root:root@tcp(127.0.0.1:3306)/contactappnew?charset=utf8mb4&parseTime=true&loc=Local"

	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Print(err.Error())
		return nil
	}
	// sqlDB is the underlying mysql DB. It is needed to specify connection restrictions.
	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxLifetime(3 * time.Minute)

	db.LogMode(true)
	// utf8_general_ci is the default collate for utf8 and it is okay to not specify it.
	// ci means case insensitive.
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci")

	// gorm logger interface needs to be implemented by the 3rd party logger for decorted output.
	// db.SetLogger(log)
	// blocks update without a where clause.
	db.BlockGlobalUpdate(true)
	return db
}

// MigrateTables will do a table table migration for all modules.
func (app *App) MigrateTables(configs []ModuleConfig) {
	app.WG.Add(len(configs))
	for _, config := range configs {
		// sometimes leads to dead lock(gorm) with go routines
		config.TableMigration(app.WG)
		app.WG.Done()

	}
	app.WG.Wait()
	app.Log.Print("End of Migration")

}
func (app *App) Init() {
	// Router has to initialized first.
	app.initializeRouter()
	// Server will need the router and must be initialized after.
	app.initializeServer()
}
func (app *App) initializeRouter() {
	app.Log.Print(app.Name + " App Route initializing")
	app.Router = mux.NewRouter().StrictSlash(true)
	app.Router = app.Router.PathPrefix("/api/v1/contactapp").Subrouter()
	app.Router.Use(LoggingMiddleware(app))
}

// he LoggingMiddleware function is a middleware that logs information about each incoming HTTP request and the corresponding outgoing response.
func LoggingMiddleware(app *App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app.Log.Print("Received Request: ", r.Method, r.URL.Path, r.RemoteAddr, r.Host)

			next.ServeHTTP(w, r)

			app.Log.Print("Completed Request without any errors...")
		})
	}
}

//CORS (Cross-Origin Resource Sharing) Middleware: Although not explicitly shown in the provided code, the handlers.CORS function from the "github.com/gorilla/handlers" package is commonly used as middleware to handle Cross-Origin Resource Sharing. CORS is an essential security feature in web applications that controls access to resources from different origins (domains).

// initializeServer will initialize server with the given config.
func (app *App) initializeServer() {
	headers := handlers.AllowedHeaders([]string{
		"Content-Type", "X-Total-Count", "token",
	})
	methods := handlers.AllowedMethods([]string{
		http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete, http.MethodOptions,
	})
	originOption := handlers.AllowedOriginValidator(app.checkOrigin)

	app.Server = &http.Server{
		Addr:         "0.0.0.0:3000",
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.CORS(headers, methods, originOption)(app.Router),
	}
	app.Log.Print("Server Exposed On 4000")
}

// StartServer will start the app.
func (app *App) StartServer() error {

	app.Log.Print("Server Time: ", time.Now())
	app.Log.Print("Server Running on port:3000 ")

	if err := app.Server.ListenAndServe(); err != nil {
		app.Log.Print("Listen and serve error: ", err)
		return err
	}
	return nil
}

// RegisterControllerRoutes will register the specified routes in controllers.
func (app *App) RegisterControllerRoutes(controllers []Controller) {
	app.Lock()
	defer app.Unlock()
	// controllers registering routes.
	for _, controller := range controllers {
		// Can't use go routines as gorilla mux doesn't support it.
		controller.RegisterRoutes(app.Router.NewRoute().Subrouter())
	}
}
func (app *App) checkOrigin(origin string) bool {
	// origin will be the actual origin from which the request is made.

	return true
}

// Stop stops the app.
func (app *App) Stop() {
	// Stopping scheduler.
	context, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	// Closing db
	app.DB.Close()
	app.Log.Print("Db closed")

	// Stopping Server.
	err := app.Server.Shutdown(context)
	if err != nil {
		app.Log.Print("Fail to Stop Server...")
		return
	}
	app.Log.Print("Server shutdown gracefully.")
}
