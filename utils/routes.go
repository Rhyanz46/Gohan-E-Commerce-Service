package utils

import "gorm.io/gorm"

type RouterConfig struct {
	PrimaryDB *gorm.DB
}
