package model_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/qxdn/imagesim/model"
	"github.com/qxdn/imagesim/util"
)

func TestHash(t *testing.T) {
	// 读图
	demo1, err1 := util.ReadImageFromPath("../demo/demo1.jpg")
	demo2, err2 := util.ReadImageFromPath("../demo/demo2.jpg")
	if err := errors.Join(err1, err2); err != nil {
		t.Errorf("ReadImage error:%v", err)
		return
	}
	// 算hash
	hash1, err1 := model.ComputeHashes(demo1)
	hash2, err2 := model.ComputeHashes(demo2)
	if err := errors.Join(err1, err2); err != nil {
		t.Errorf("ComputeHashes error:%v", err)
		return
	}

	fmt.Printf("demo1: AHash:%v,DHash:%v,PHash:%v\n", hash1.AHash, hash1.DHash, hash1.PHash)
	fmt.Printf("demo2: AHash:%v,DHash:%v,PHash:%v\n", hash2.AHash, hash2.DHash, hash2.PHash)
	// 比较hash
	fmt.Println("AHash Distance:", util.HashDistance(hash1.AHash, hash2.AHash))
	fmt.Println("DHash Distance:", util.HashDistance(hash1.DHash, hash2.DHash))
	fmt.Println("PHash Distance:", util.HashDistance(hash1.PHash, hash2.PHash))
}
