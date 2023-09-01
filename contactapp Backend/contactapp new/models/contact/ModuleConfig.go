package contact

import (
	"sync"

	"github.com/jinzhu/gorm"

	"contactapp/components/log"
)

type ModuleConfig struct {
	DB *gorm.DB
}

// NewContactModuleConfig creates a new instance of ContactModuleConfig
func NewContactModuleConfig(db *gorm.DB) *ModuleConfig {
	return &ModuleConfig{
		DB: db,
	}
}

func (config *ModuleConfig) TableMigration(wg *sync.WaitGroup) {
	//Table list
	var models []interface{} = []interface{}{
		&Contact{},
	}
	//Table migration
	for _, model := range models {
		if err := config.DB.AutoMigrate(model).Error; err != nil {
			log.GetLogger().Print("Auto migration ==> %s", err.Error())
		}
	}
	if err := config.DB.Model(&Contact{}).
		AddForeignKey("user_id", "user(Id)", "CASCADE", "CASCADE").Error; err != nil {
		log.GetLogger().Print("Auto migration ==> %s", err.Error())
	}
	log.GetLogger().Print("auto migration done")
}
