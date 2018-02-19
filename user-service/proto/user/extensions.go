package go_micro_srv_user

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func (model *User) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New()
	return scope.SetColumn("Id", uuid.String())
}
