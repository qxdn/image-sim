package dal

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Key          string    `gorm:"type:varchar(100);not null;unique"` // object key
	AHash        uint64    `gorm:"type:bigint;not null;index"`        // average hash
	DHash        uint64    `gorm:"type:bigint;not null"`              // difference hash
	PHash        uint64    `gorm:"type:bigint;not null"`              // perceptual hash
	Filename     string    `gorm:"type:varchar(100);not null"`        // object filename
	Url          string    `gorm:"type:varchar(255);not null"`        // object url
	LastModified time.Time `gorm:"type:datetime;not null"`            // object last modified time
}
