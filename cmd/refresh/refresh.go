package main

import (
	"sync"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/qxdn/imagesim/dal"
	"github.com/qxdn/imagesim/global"
	"github.com/qxdn/imagesim/model"
	"github.com/qxdn/imagesim/services"
	"gorm.io/gorm"
)

/**
 * @description: 任务函数
 * @param {object} object
 * @param {client} client
 * @param {db} db
 * @param {wg} wg
 * @return {*}
 */
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
	global.InitGlobal()
	db, client := global.Db, global.OSSClient

	db.AutoMigrate(&dal.Image{})

	objects, err := model.OSSListObject(client, config.OSS.BucketName, config.OSS.Directory)

	if err != nil {
		global.Logger.Errorf("OSSListObject error:%v", err)
		return
	}
	glimiter := model.NewGlimiter(config.Refresh.WorkerNum)
	var wg sync.WaitGroup
	for _, obj := range objects {
		wg.Add(1)
		glimiter.Run(func() {
			task(obj, client, db, &wg)
		})
	}
	wg.Wait()
}
