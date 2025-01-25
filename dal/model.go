package dal

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type baseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Image struct {
	baseModel
	Key          string                `gorm:"type:varchar(100);not null;uniqueIndex:udx_key"` // object key
	AHash        uint64                `gorm:"type:bigint unsigned;not null;index"`            // average hash
	DHash        uint64                `gorm:"type:bigint unsigned;not null;index"`            // difference hash
	PHash        uint64                `gorm:"type:bigint unsigned;not null;index"`            // perceptual hash
	Filename     string                `gorm:"type:varchar(100);not null"`                     // object filename
	Url          string                `gorm:"type:varchar(255);not null"`                     // object url
	LastModified time.Time             `gorm:"type:datetime;not null"`                         // object last modified time
	DeletedAt    soft_delete.DeletedAt `gorm:"uniqueIndex:udx_key"`                            // soft delete
}
