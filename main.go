package main

import (
	"sync"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/qxdn/imagesim/dal"
	"github.com/qxdn/imagesim/global"
	"github.com/qxdn/imagesim/model"
	"github.com/qxdn/imagesim/services"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func task(object *model.OSSObject, client *oss.Client, db *gorm.DB, wg *sync.WaitGroup) {
	defer func() {
		if err := recover(); err != nil {
			global.Logger.Errorf("task panic:%v", err)
		}
	}()
	defer wg.Done()
	err := services.ComputeSingle(object, client, db)
	if err != nil {
		global.Logger.Errorf("ComputeSingle error:%v", err)
	}
}

func main() {
	config := global.ReadConfig()
	db, err := gorm.Open(mysql.Open(config.DB.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	err = db.AutoMigrate(&dal.Image{})
	if err != nil {
		panic(err)
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.OSS.AccessKey, config.OSS.SecretKey)).
		WithRegion(config.OSS.Region)

	client := oss.NewClient(cfg)

	objects, err := model.OSSListObject(client, config.OSS.BucketName, config.OSS.Directory)

	if err != nil {
		global.Logger.Errorf("OSSListObject error:%v", err)
		return
	}
	var wg sync.WaitGroup
	for _, obj := range *objects {
		wg.Add(1)
		go task(&obj, client, db, &wg)
	}
	wg.Wait()
}
