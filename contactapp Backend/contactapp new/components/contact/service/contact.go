package service

import (
	"time"

	"github.com/jinzhu/gorm"

	"contactapp/errors"
	"contactapp/models/contact"
	"contactapp/repository"
)

type ContactService struct {
	db           *gorm.DB
	repository   repository.Repository
	associations []string
}

func NewContactService(db *gorm.DB, repo repository.Repository) *ContactService {
	return &ContactService{
		db:           db,
		repository:   repo,
		associations: []string{},
	}
}

func (service *ContactService) CreateContact(contact *contact.Contact) error {
	//creating unit of work
	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()
	//creating contact
	err := service.repository.Add(uow, contact)
	if err != nil {
		uow.RollBack()
		return err
	}
	//commiting transaction
	uow.Commit()
	return nil
}

func (service *ContactService) doesContactExist(ID uint) error {
	exists, err := repository.DoesContactExist(service.db, ID, contact.Contact{}, repository.Filter("`id` = ?", ID))
	if !exists || err != nil {
		return errors.NewValidationError("Contact not found")
	}
	return nil
}

func (service *ContactService) UpdateContact(contactToUpdate *contact.Contact) error {
	err := service.doesContactExist(contactToUpdate.ID)
	if err != nil {
		return err
	}

	// do we need to create this function in repository

	uow := repository.NewUnitOfWork(service.db, false) // transaction should be readonly or not so false
	defer uow.RollBack()
	tempContact := contact.Contact{}
	err = service.repository.GetRecordForContact(uow, contactToUpdate.ID, &tempContact, repository.Select("`created_at`"),
		repository.Filter("`id` = ?", contactToUpdate.ID))
	if err != nil {
		return err
	}

	err = service.repository.Save(uow, &contactToUpdate)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

func (service *ContactService) DeleteContact(contactToDelete *contact.Contact) error {
	err := service.doesContactExist(contactToDelete.ID)
	if err != nil {
		return err
	}
	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	if err := service.repository.UpdateWithMap(uow, contactToDelete, map[string]interface{}{
		"DeletedAt": time.Now(),
	},
		repository.Filter("`id` = ?", contactToDelete.ID)); err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

func (service *ContactService) GetAllContacts(allContacts *[]contact.Contact, totalCount int) error {
	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()
	err := service.repository.GetAll(uow, allContacts)
	if err != nil {
		return err
	}
	uow.Commit()
	return nil
}
