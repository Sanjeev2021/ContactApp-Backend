package controller

import (
	//"net/http"
	"net/http"
	"strconv"

	//"github.com/gorilla/mux"
	"github.com/gorilla/mux"

	"contactapp/components/contact/service"
	"contactapp/components/log"
	"contactapp/errors"
	"contactapp/models/contact"
	"contactapp/web"
	//"contactapp/errors"
	//"contactapp/web"

)

// ContactController handles requests related to contact
type ContactController struct {
	log     *log.Log
	service *service.ContactService
}

// NewContactController creates a new instance of ContactController
func NewContactController(contactservice *service.ContactService, log *log.Log) *ContactController {
	return &ContactController{
		log:     log,
		service: contactservice,
	}
}

func (controller *ContactController) RegisterRoutes(router *mux.Router) {
	contactRouter := router.PathPrefix("/contact").Subrouter()
	contactRouter.HandleFunc("/", controller.GetAllContacts).Methods(http.MethodGet)
	contactRouter.HandleFunc("/", controller.RegisterContact).Methods(http.MethodPost)
	//	contactRouter.HandleFunc("/{id}", controller.GetContact).Methods(http.MethodGet)
	contactRouter.HandleFunc("/{id}", controller.UpdateContact).Methods(http.MethodPut)
	contactRouter.HandleFunc("/{id}", controller.DeleteContact).Methods(http.MethodDelete)
	controller.log.Print("===========Contact routes registered===========")
}

func (controller *ContactController) RegisterContact(w http.ResponseWriter, r *http.Request) {
	newContact := contact.Contact{}
	//Unmarshal json
	err := web.UnmarshalJSON(r, &newContact)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	//call add test method
	err = controller.service.CreateContact(&newContact)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusInternalServerError))
		return
	}
	//respond with json
	web.RespondJSON(w, http.StatusCreated, newContact)
}

func (controller *ContactController) UpdateContact(w http.ResponseWriter, r *http.Request) {
	contactToUpdate := contact.Contact{}
	//Unmarshal json
	err := web.UnmarshalJSON(r, &contactToUpdate)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	vars := mux.Vars(r)

	//get id from url
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	contactToUpdate.ID = uint(id)
	//call update method
	err = controller.service.UpdateContact(&contactToUpdate)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusInternalServerError))
		return
	}
	//respond with json
	web.RespondJSON(w, http.StatusOK, contactToUpdate)
}

func (controller *ContactController) DeleteContact(w http.ResponseWriter, r *http.Request) {
	contactToDelete := contact.Contact{}
	var err error
	vars := mux.Vars(r)
	//get id from url
	intID, err := strconv.Atoi(vars["id"])
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
		return
	}
	contactToDelete.ID = uint(intID)
	err = controller.service.DeleteContact(&contactToDelete)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusInternalServerError))
		return
	}
	//respond with json
	web.RespondJSON(w, http.StatusOK, contactToDelete)
}

// GetAllContacts gets all contacts
func (controller *ContactController) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	allContacts := &[]contact.Contact{}
	var totalCount int
	err := controller.service.GetAllContacts(allContacts, totalCount)
	if err != nil {
		controller.log.Print(err.Error())
		web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusInternalServerError))
		return
	}
	//respond with json
	web.RespondJSONWithXTotalCount(w, http.StatusOK, totalCount, allContacts)
}
