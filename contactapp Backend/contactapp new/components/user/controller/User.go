package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	// "contactapp/components/contact/controller"
	"contactapp/components/log"
	"contactapp/components/user/service"
	"contactapp/errors"
	"contactapp/models/user"
	"contactapp/web"

)

// UserController gives access to CRUD operations for entity
type UserController struct {
	log     log.Log
	service *service.UserService
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewUserController returns new instance of UserController
func NewUserController(userService *service.UserService, log log.Log) *UserController {
	return &UserController{
		service: userService,
		log:     log,
	}
}

func (controller *UserController) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", controller.RegisterUser).Methods(http.MethodPost)
	userRouter.HandleFunc("/", controller.GetAllUsers).Methods(http.MethodGet)
	userRouter.HandleFunc("/{id}", controller.UpdateUser).Methods(http.MethodPut)
	userRouter.HandleFunc("/{id}", controller.DeleteUser).Methods(http.MethodDelete)
	userRouter.HandleFunc("/login", controller.Login).Methods(http.MethodPost)
	fmt.Println("==============================userRegisterRoutes==========================")
}

// Login handles the login request
func (controller *UserController) Login(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		fmt.Println("Error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Credentials:", credentials)
	user, err := controller.service.GetUserByUsername(credentials.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Password != credentials.Password {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&user.ID)
}

func (controller *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	newUser := user.User{}
	// Unmarshal json.
	err := web.UnmarshalJSON(r, &newUser)
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	// Call add test method.
	err = controller.service.CreateUser(&newUser)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResposeWriter.
	web.RespondJSON(w, http.StatusCreated, newUser)
}
func (controller *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers := &[]user.User{}
	var totalCount int
	err := controller.service.GetAllUsers(allUsers, &totalCount)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	// Writing Response with OK Status to ResonseWriter,
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allUsers)
}
func (controller *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==============================userToUpdate==========================")
	userToUpdate := user.User{}

	// Unmarshal JSON.
	fmt.Println(r.Body)
	err := web.UnmarshalJSON(r, &userToUpdate)
	if err != nil {
		fmt.Println("==============================err from UnmarshalJSON==========================")
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	vars := mux.Vars(r)

	intID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	userToUpdate.ID = uint(intID)
	fmt.Println("==============================userToUpdate==========================")
	fmt.Println(&userToUpdate)
	// Call update test method.
	err = controller.service.UpdateUser(&userToUpdate)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}

	web.RespondJSON(w, http.StatusOK, userToUpdate)
}
func (controller *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	controller.log.Print("********************************DeleteTest call**************************************")
	usetToDelete := user.User{}
	var err error
	vars := mux.Vars(r)
	intID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err)
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	usetToDelete.ID = uint(intID)
	err = controller.service.DeleteUser(&usetToDelete)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, err)
		return
	}
	web.RespondJSON(w, http.StatusOK, "Delete User successfull.")
}
