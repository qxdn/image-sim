package main

import (
	"fmt"

	"github.com/qxdn/imagesim/model"
	"github.com/qxdn/imagesim/util"
)

func main() {
	// 读图
	demo1, _ := util.ReadImage("./demo/demo1.jpg")
	demo2, _ := util.ReadImage("./demo/demo2.jpg")
	// 算hash
	hash1, _ := model.ComputeHashes(&demo1)
	hash2, _ := model.ComputeHashes(&demo2)

	fmt.Printf("demo1: AHash:%v,DHash:%v,PHash:%v\n", hash1.AHash, hash1.DHash, hash1.PHash)
	fmt.Printf("demo2: AHash:%v,DHash:%v,PHash:%v\n", hash2.AHash, hash2.DHash, hash2.PHash)
	// 比较hash
	fmt.Println("AHash Distance:", util.HashDistance(hash1.AHash, hash2.AHash))
	fmt.Println("DHash Distance:", util.HashDistance(hash1.DHash, hash2.DHash))
	fmt.Println("PHash Distance:", util.HashDistance(hash1.PHash, hash2.PHash))
}
