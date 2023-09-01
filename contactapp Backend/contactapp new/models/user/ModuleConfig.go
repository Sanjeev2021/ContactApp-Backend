// table migration is happening here
package user

import (
	//	"fmt"
	"sync"

	"github.com/jinzhu/gorm"

	"contactapp/components/log"
)

type ModuleConfig struct {
	DB *gorm.DB
}

// NewModuleConfig creates a new instance of ModuleConfig
// meaning of this ?
func NewUserModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {
	// table list
	var models []interface{} = []interface{}{
		&User{},
	}
	//Table migrant
	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.GetLogger().Print("Auto Migration ==>", err) // it is used to log an error, it is a string that defines a log message
		}
	}

	log.GetLogger().Print("Table migration completed")

}
