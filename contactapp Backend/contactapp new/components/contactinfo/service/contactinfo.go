package service

import (
	"time"

	"github.com/jinzhu/gorm"

	"contactapp/errors"
	"contactapp/models/contactinfo"
	"contactapp/repository"
)

// ContactInfoService is the service for contactinfo
type ContactInfoService struct {
	db           *gorm.DB
	repository   repository.Repository
	associations []string
}

// NewContactInfoService creates a new instance of ContactInfoService
func NewContactInfoService(db *gorm.DB, repo repository.Repository) *ContactInfoService {
	return &ContactInfoService{
		db:           db,
		repository:   repo,
		associations: []string{},
	}
}

// Create ContactInfo
func (service *ContactInfoService) CreateContactInfo(contactinfo *contactinfo.ContactInfo) error {
	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()
	err := service.repository.Add(uow, contactinfo)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}

// does contactinfoexist
func (service *ContactInfoService) doesContactInfoExist(ID uint) error {
	exists, err := repository.DoesContactInfoExist(service.db, ID, contactinfo.ContactInfo{}, repository.Filter("`id` = ?", ID))
	if !exists || err != nil {
		return errors.NewValidationError("ContactInfo not found")
	}
	return nil
}

// Update ContactInfo
func (service *ContactInfoService) UpdateContactInfo(contactinfoToUpdate *contactinfo.ContactInfo) error {
	err := service.doesContactInfoExist(contactinfoToUpdate.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()
	tempContactInfo := contactinfo.ContactInfo{}
	err = service.repository.GetRecordForContactInfo(uow, contactinfoToUpdate.ID, &tempContactInfo, repository.Select("`created_at`"),
		repository.Select("`id` = ?", contactinfoToUpdate.ID))
	if err != nil {
		return err
	}

	err = service.repository.Save(uow, contactinfoToUpdate)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil

}

// Delete ContactInfo
func (service *ContactInfoService) DeleteContactInfo(contactInfoToDelete *contactinfo.ContactInfo) error {
	err := service.doesContactInfoExist(contactInfoToDelete.ID)
	if err != nil {
		return err
	}

	uow := repository.NewUnitOfWork(service.db, false)
	defer uow.RollBack()

	if err := service.repository.UpdateWithMap(uow, contactInfoToDelete, map[string]interface{}{
		"DeletedAt": time.Now(),
	},
		repository.Filter("`id` = ?", contactInfoToDelete.ID)); err != nil {
		uow.RollBack()
		return err

	}

	uow.Commit()
	return nil
}

// Get ContactInfo
func (service *ContactInfoService) GetAllContactInfo(allContactInfo *[]contactinfo.ContactInfo) error {
	uow := repository.NewUnitOfWork(service.db, true)
	defer uow.RollBack()
	err := service.repository.GetAll(uow, allContactInfo)
	if err != nil {
		uow.RollBack()
		return err
	}
	uow.Commit()
	return nil
}
