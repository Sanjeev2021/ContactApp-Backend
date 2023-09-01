package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"contactapp/components/contactinfo/service"
	"contactapp/components/log"
	"contactapp/errors"
	"contactapp/models/contactinfo"
	"contactapp/web"
)

// ContactInfoController gives access to CRUD operations for entity
type ContactInfoController struct {
	log     log.Log
	service *service.ContactInfoService
}

// NewContactInfoController returns new instance of ContactInfoController
func NewContactInfoController(contactService *service.ContactInfoService,
	log log.Log) *ContactInfoController {
	return &ContactInfoController{
		service: contactService,
		log:     log,
	}
}
func (controller *ContactInfoController) RegisterRoutes(router *mux.Router) {
	contactInfoRouter := router.PathPrefix("/contactinfo").Subrouter()
	contactInfoRouter.HandleFunc("/", controller.CreateContactInfo).Methods(http.MethodPost)
	contactInfoRouter.HandleFunc("/", controller.GetContactInfo).Methods(http.MethodGet)
	contactInfoRouter.HandleFunc("/{id}", controller.UpdateContactInfo).Methods(http.MethodPut)
	contactInfoRouter.HandleFunc("/{id}", controller.DeleteContactInfo).Methods(http.MethodDelete)

}
func (controller *ContactInfoController) CreateContactInfo(w http.ResponseWriter, r *http.Request) {
	newContactInfo := contactinfo.ContactInfo{}
	//Unmarsh json
	err := web.UnmarshalJSON(r, &newContactInfo)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	//call add test method
	err = controller.service.CreateContactInfo(&newContactInfo)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusInternalServerError))
		return
	}
	//respond with json
	web.RespondJSON(w, http.StatusCreated, newContactInfo)
}

func (controller *ContactInfoController) UpdateContactInfo(w http.ResponseWriter, r *http.Request) {
	updatedContactInfo := contactinfo.ContactInfo{}
	//Unmarshal json
	err := web.UnmarshalJSON(r, &updatedContactInfo)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	updatedContactInfo.ID = uint(id)
	//call add test method
	err = controller.service.UpdateContactInfo(&updatedContactInfo)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusInternalServerError))
		return
	}
	//respond with json
	web.RespondJSON(w, http.StatusOK, updatedContactInfo)

}

func (controller *ContactInfoController) DeleteContactInfo(w http.ResponseWriter, r *http.Request) {
	contactInfoToDelete := contactinfo.ContactInfo{}
	var err error
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	contactInfoToDelete.ID = uint(id)
	//call add test method
	err = controller.service.DeleteContactInfo(&contactInfoToDelete)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusInternalServerError))
		return
	}
	//respond with json
	web.RespondJSON(w, http.StatusOK, contactInfoToDelete)

}

// Getall ContactInfo
func (controller *ContactInfoController) GetContactInfo(w http.ResponseWriter, r *http.Request) {
	allContactInfo := &[]contactinfo.ContactInfo{}
	//call add test method
	err := controller.service.GetAllContactInfo(allContactInfo)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusInternalServerError))
		return
	}
	//respond with json
	web.RespondJSON(w, http.StatusOK, allContactInfo)

}
