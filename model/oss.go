package model

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type OSSObject struct {
	Key          string    // Object key
	Filename     string    // Object filename
	LastModified time.Time // Object last modified time
}

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
			log.Fatalf("failed to get page  %v", err)
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
				Key:          oss.ToString(obj.Key),
				Filename:     ExtractOSSFilename(key),
				LastModified: lastModified,
			}
			objects = append(objects, objs)
		}
	}
	return &objects, nil
}

func ExtractOSSFilename(path string) string {
	// 从path中提取文件名
	s := strings.Split(path, "/")
	if len(s) == 0 {
		// 无法提取文件名
		return ""
	}
	return s[len(s)-1]
}
