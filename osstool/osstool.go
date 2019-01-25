package osstool

import (
	"app/datastruct"
	"app/log"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func CreateOSSBucket() *oss.Bucket {
	// 创建OSSClient实例。
	client, err := oss.New(datastruct.OSSEndpoint, datastruct.OSSAccessKeyId, datastruct.OSSAccessKeySecret)
	if err != nil {
		log.Debug("Error:%v", err)
		os.Exit(-1)
	}

	bucketName := datastruct.OSSBucketName

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Debug("Error:%v", err)
		os.Exit(-1)
	}

	return bucket
}

func DeleteFile(bucket *oss.Bucket, objectName string) {
	// 删除单个文件。
	err := bucket.DeleteObject(objectName)
	if err != nil {
		log.Debug("Error:%v", err)
	}
}

// func SignedURL(bucket *oss.Bucket, objectName string) string {
// 	// signedURL, err := bucket.SignURL(objectName, oss.HTTPGet, 6000)
// 	// //oss.Process("image/resize,h_100"))
// 	// if err != nil {
// 	// 	log.Debug("SignedURL err:%v", err.Error())
// 	// 	return ""
// 	// }
// 	signedURL := fmt.Sprintf("https://rouge999.oss-cn-shenzhen.aliyuncs.com/%s", objectName)
// 	return signedURL
// }

func CreateOSSURL(objectName string) string {
	url := fmt.Sprintf("https://rouge999.oss-cn-shenzhen.aliyuncs.com/%s", objectName)
	return url
}
