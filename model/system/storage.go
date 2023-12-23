package system

import (
	"time"

	"gorm.io/gorm"
)

type Storage struct {
	gorm.Model
	StorageName     string `gorm:"unique;not null"`
	StorageRealPath string 
	StorageUrlName  string `gorm:"unique;not null"`
	Status          bool
	DeadLine		time.Time
}
