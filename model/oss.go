package model

import (
	"context"
	"errors"
	"image"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/qxdn/imagesim/global"
	"github.com/qxdn/imagesim/util"
)

/**
 * @Struct: OSSObject
 * @Description: Object in OSS
 **/
type OSSObject struct {
	BuckName     string    // Bucket name
	Key          string    // Object key
	Filename     string    // Object filename
	LastModified time.Time // Object last modified time
}

/**
 * @Function: ReadImage
 * @Description: Read image from OSSObject
 * @Param: client *oss.Client
 * @Return: *image.Image, error
 **/
func (o *OSSObject) ReadImage(client *oss.Client) (*image.Image, error) {
	if (o == nil) || (client == nil) {
		return nil, errors.New("nil object or client")
	}
	if o.Key == "" {
		return nil, errors.New("empty key")
	}
	request := &oss.GetObjectRequest{
		Bucket: oss.Ptr(o.BuckName),
		Key:    oss.Ptr(o.Key),
	}
	result, err := client.GetObject(context.TODO(), request)
	if err != nil {
		global.Logger.Error("failed to get object2 %v", err)
		return nil, err
	}
	defer result.Body.Close()
	i, err := util.ReadImage(result.Body)
	if err != nil {
		global.Logger.Errorf("failed to decode image %v", err)
		return nil, err
	}
	return i, nil
}

/**
 * @Function: OSSListObject
 * @Description: List objects in the specified directory
 * @Param: client *oss.Client, buckName, directory string
 * @Return: *[]OSSObject, error
 **/
func OSSListObject(client *oss.Client, buckName, directory string) (*[]OSSObject, error) {
	// Create the Paginator for the ListObjectsV2 operation.
	p := client.NewListObjectsV2Paginator(&oss.ListObjectsV2Request{
		Bucket: oss.Ptr(buckName),
		Prefix: oss.Ptr(directory),
	})
	var objects []OSSObject
	// Iterate through the object pages
	for p.HasNext() {
		page, err := p.NextPage(context.TODO())
		if err != nil {
			global.Logger.Errorf("failed to get page  %v", err)
			return nil, err
		}
		// Print the objects found
		for _, obj := range page.Contents {
			key := oss.ToString(obj.Key)
			size := obj.Size
			lastModified := oss.ToTime(obj.LastModified)
			if key == directory || size == 0 {
				// 应该是文件夹，跳过
				continue
			}
			objs := OSSObject{
				BuckName:     buckName,
				Key:          oss.ToString(obj.Key),
				Filename:     ExtractOSSFilename(key),
				LastModified: lastModified,
			}
			objects = append(objects, objs)
		}
	}
	return &objects, nil
}

/**
 * @Function: ExtractOSSFilename
 * @Description: Extract the filename from the path
 * @Param: path string
 * @Return: string
 **/
func ExtractOSSFilename(path string) string {
	// 从path中提取文件名
	s := strings.Split(path, "/")
	if len(s) == 0 {
		// 无法提取文件名
		return ""
	}
	return s[len(s)-1]
}
