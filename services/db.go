package services

import (
	"errors"

	"github.com/qxdn/imagesim/dal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**
 * @Function: LoadImageFromDB
 * @Description: Load image from database
 * @Param: db *gorm.DB, key string
 * @Return: *dal.Image
 **/
func LoadImageFromDB(db *gorm.DB, key string) *dal.Image {
	image := &dal.Image{}
	result := db.Where(&dal.Image{Key: key}).First(image)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	if result.Error != nil {
		panic(result.Error)
	}
	return image
}

/**
 * @Function: LoadImageFromDBWithLock
 * @Description: Load image from database with lock
 * @Param: db *gorm.DB, key string
 * @Return: *dal.Image
 **/
func LoadImageFromDBWithLock(db *gorm.DB, key string) *dal.Image {
	image := &dal.Image{}
	result := db.Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "NOWAIT",
	}).Where(&dal.Image{Key: key}).First(image)
	err := result.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		// 可能是并发
		panic(err)
	}
	return image
}

func SaveDBImage(image *dal.Image, db *gorm.DB) error {
	result := db.Save(image)
	return result.Error
}

/**
 * @Function: LoadAllImagesFromDB
 * @Description: Load all images from database
 * @Param: db *gorm.DB
 * @Return: *[]dal.Image
 **/
func LoadAllImagesFromDB(db *gorm.DB) *[]dal.Image {
	images := &[]dal.Image{}
	result := db.Find(images)
	if result.Error != nil {
		panic(result.Error)
	}
	return images
}

func LoadImagesByPHashDistance(db *gorm.DB, hash uint64, distance int) *[]dal.Image {
	images := &[]dal.Image{}
	result := db.Where("BIT_COUNT(p_hash ^ ?) < ?", hash, distance).Find(images)
	if result.Error != nil {
		panic(result.Error)
	}
	return images
}
