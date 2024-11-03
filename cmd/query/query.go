package main

import (
	"github.com/qxdn/imagesim/global"
	"github.com/qxdn/imagesim/services"
	"github.com/qxdn/imagesim/util"
)

type SimilarImage struct {
	Key        string  `json:"key"`
	Url        string  `json:"url"`
	Other      string  `json:"other"`
	OtherUrl   string  `json:"other_url"`
	Similarity float64 `json:"similarity"`
}

func main() {
	config := global.ReadConfig()
	global.InitGlobal()
	db := global.Db

	images := services.LoadAllImagesFromDB(db)

	var queryResult []SimilarImage

	for _, image := range *images {
		result := services.LoadImagesByPHashDistance(db, image.PHash, config.Query.Distance)
		//global.Logger.Debugf("Image with key %v query finished find %v similar images", image.Key, len(*images))
		for _, img := range *result {
			if img.Key == image.Key {
				continue
			}
			similar := util.ComputeSimilarity(image.PHash, img.PHash)
			queryResult = append(queryResult, SimilarImage{
				Key:        image.Key,
				Url:        image.Url,
				Other:      img.Key,
				OtherUrl:   img.Url,
				Similarity: similar,
			})
		}
	}
	global.Logger.Infof("Query finished find %v similar images", len(queryResult))
	for _, result := range queryResult {
		global.Logger.Info("Key:", result.Key, " , Url:", result.Url, " , Other:", result.Other, " , OtherUrl:", result.OtherUrl, " , Similarity:", result.Similarity)
	}
}
