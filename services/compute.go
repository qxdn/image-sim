package services

import (
	"errors"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/qxdn/imagesim/dal"
	"github.com/qxdn/imagesim/global"
	"github.com/qxdn/imagesim/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**
 * @Function: ComputeSingle
 * @Description: Compute single image
 * @Param: object *model.OSSObject, client *oss.Client, db *gorm.DB
 * @Return: error
 **/
func ComputeSingle(object *model.OSSObject, client *oss.Client, db *gorm.DB) error {
	dbImage := LoadImageFromDB(db, object.Key)
	needSave := false
	if dbImage == nil {
		global.Logger.Info("Image not found with key:", object.Key)
		needSave = true
	} else {
		global.Logger.Info("Image found with key:", object.Key)
		if dbImage.LastModified.Before(object.LastModified) {
			// 需要更新数据库
			global.Logger.Infof("Image with key %v LastModified has changed need update", object.Key)
			needSave = true
		}
	}
	// 如果需要更新或者创建
	if needSave {
		imageHash, err := ComputeOSSHash(object, client)
		if err != nil {
			global.Logger.Error("compute hash or image fail", err)
			return err
		}
		// 开启事务
		err = db.Transaction(func(tx *gorm.DB) error {
			// 加锁读取
			dbImage := LoadImageFromDBWithLock(db, object.Key)
			dbImage = CreateDBImage(dbImage, object, imageHash)
			return SaveDBImage(dbImage, tx)
		})
		return err
	}
	return nil
}

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

func CreateDBImage(image *dal.Image, object *model.OSSObject, hash *model.ImageHash) *dal.Image {
	if image == nil {
		image = &dal.Image{}
	}
	image.Key = object.Key
	image.AHash = hash.AHash
	image.DHash = hash.DHash
	image.PHash = hash.PHash
	image.Filename = object.Filename
	image.Url = object.Url
	image.LastModified = object.LastModified
	return image
}

func SaveDBImage(image *dal.Image, db *gorm.DB) error {
	result := db.Save(image)
	return result.Error
}

func ComputeOSSHash(object *model.OSSObject, client *oss.Client) (*model.ImageHash, error) {
	imageOSS, err := object.ReadImage(client)
	if err != nil {
		global.Logger.Error("read image from oss fail", err)
		return nil, err
	}
	return model.ComputeHashes(imageOSS)
}
