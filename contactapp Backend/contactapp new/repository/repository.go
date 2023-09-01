package repository

import (
	"github.com/jinzhu/gorm"

	"contactapp/errors"
)

type Repository interface {
	GetAll(uow *UnitOfWork, out interface{}, queryProcessor ...QueryProcessor) error
	Add(uow *UnitOfWork, out interface{}) error
	Update(uow *UnitOfWork, out interface{}) error
	UpdateWithMap(uow *UnitOfWork, model interface{}, value map[string]interface{}, queryProcessor ...QueryProcessor) error
	GetRecordForUser(uow *UnitOfWork, userID uint, out interface{}, queryProcessors ...QueryProcessor) error
	Save(uow *UnitOfWork, value interface{}) error
	GetRecord(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error
	GetRecordForContact(uow *UnitOfWork, ID uint, out interface{}, queryProcessors ...QueryProcessor) error
	GetRecordForContactInfo(uow *UnitOfWork, ID uint, out interface{}, queryProcessors ...QueryProcessor) error
}

type GormRepository struct{}

func NewGormRepository() *GormRepository {
	return &GormRepository{}
}

// UnitofWork mamkes  A database connection and a transaction
type UnitOfWork struct {
	DB        *gorm.DB
	Committed bool
	Readonly  bool
}

// NewUnitOfWork creates a new instance of UnitOfWork
func NewUnitOfWork(db *gorm.DB, readonly bool) *UnitOfWork {
	commit := false // what if this commit is true ?
	if readonly {
		return &UnitOfWork{
			DB:        db.New(),
			Committed: commit,
			Readonly:  readonly,
		}
	}

	return &UnitOfWork{
		DB:        db.New().Begin(),
		Committed: commit,
		Readonly:  readonly,
	}
}

// Commit commits the transaction
func (uow *UnitOfWork) Commit() {
	if !uow.Readonly && !uow.Committed {
		uow.Committed = true
		uow.DB.Commit()
	}
}

// Save saves the changes to the database
func (Repository *GormRepository) Save(uow *UnitOfWork, value interface{}) error {
	return uow.DB.Debug().Save(value).Error
}

// RollBack is used to rollback the transaction on faliure
func (uow *UnitOfWork) RollBack() {
	// This condition can be used if Rollback() is defered as soon as UOW is created.
	// So we only rollback if it's not committed.
	if !uow.Readonly && !uow.Committed {
		uow.DB.Rollback()
	}
}

func Filter(condition string, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Debug().Where(condition, args...)
		return db, nil
	}
}

// when creating/updating , specify the fields that you want to save to database.
func Select(query interface{}, args ...interface{}) QueryProcessor {
	return func(db *gorm.DB, out interface{}) (*gorm.DB, error) {
		db = db.Select(query, args...)
		return db, nil
	}
}

func executeQueryProcessors(db *gorm.DB, out interface{}, queryProcessors ...QueryProcessor) (*gorm.DB, error) {
	var err error
	for _, query := range queryProcessors {
		if query != nil {
			db, err = query(db, out)
			if err != nil {
				return db, err
			}
		}
	}
	return db, nil
}

// GetALL returns all the records from the database
func (repository *GormRepository) GetAll(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Find(out).Error
}

// add adds records to table
func (repository *GormRepository) Add(uow *UnitOfWork, out interface{}) error {
	return uow.DB.Create(out).Error
}

// Update updates the record in the database
func (repository *GormRepository) Update(uow *UnitOfWork, out interface{}) error {
	return uow.DB.Model(out).Update(out).Error // why are we using model over here ?
}

func DoesRecordExistForUser(db *gorm.DB, userID uint, out interface{}, queryProcessors ...QueryProcessor) (bool, error) {
	if userID == 0 {
		return false, errors.NewValidationError("User ID cannot be 0 : Invalid tenant ID")
	}
	count := 0

	// if out is nil , then we are not interested in the result , we just want to know if the record exists or not
	// if out is not nil , then we are interested in the result , we want to know if the record exists or not and we want to store the result in out
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return false, err
	}
	if err := db.Debug().Model(out).Where("id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (repository *GormRepository) GetRecordForUser(uow *UnitOfWork, userID uint, out interface{}, queryProcessors ...QueryProcessor) error {
	queryProcessors = append([]QueryProcessor{Filter("id = ?", userID)}, queryProcessors...)
	return repository.GetRecord(uow, out, queryProcessors...)
}

func (repository *GormRepository) GetRecord(uow *UnitOfWork, out interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().First(out).Error
}

func (repository *GormRepository) UpdateWithMap(uow *UnitOfWork, model interface{}, value map[string]interface{}, queryProcessors ...QueryProcessor) error {
	db := uow.DB
	db, err := executeQueryProcessors(db, value, queryProcessors...)
	if err != nil {
		return err
	}
	return db.Debug().Model(model).Updates(value).Error
}

func DoesContactExist(db *gorm.DB, ID uint, out interface{}, queryProcessors ...QueryProcessor) (bool, error) {
	if ID == 0 {
		return false, errors.NewValidationError("ID cannot be 0 : Invalid tenant ID")
	}
	count := 0

	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return false, err
	}
	if err := db.Debug().Model(out).Where("id = ?", ID).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (repository *GormRepository) GetRecordForContact(uow *UnitOfWork, ID uint, out interface{}, queryProcessors ...QueryProcessor) error {
	queryProcessors = append([]QueryProcessor{Filter("id = ?", ID)}, queryProcessors...)
	return repository.GetRecord(uow, out, queryProcessors...)
}

func DoesContactInfoExist(db *gorm.DB, ID uint, out interface{}, queryProcessors ...QueryProcessor) (bool, error) {
	if ID == 0 {
		return false, errors.NewValidationError("ID cannot be 0 : Invalid tenant ID")
	}
	count := 0

	db, err := executeQueryProcessors(db, out, queryProcessors...)
	if err != nil {
		return false, err
	}
	if err := db.Debug().Model(out).Where("id = ?", ID).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (repository *GormRepository) GetRecordForContactInfo(uow *UnitOfWork, ID uint, out interface{}, queryProcessors ...QueryProcessor) error {
	queryProcessors = append([]QueryProcessor{Filter("id = ?", ID)}, queryProcessors...)
	return repository.GetRecord(uow, out, queryProcessors...)
}
