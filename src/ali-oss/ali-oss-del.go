package alioss

import (
	"alioss/src/api"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func OssDelFile() {

	cfg := api.LoadConfig()

	OSSEndpoint := fmt.Sprint(cfg.Section("oss").Key("OSSEndpoint"))
	AccessKeyId := fmt.Sprint(cfg.Section("oss").Key("AccessKeyId"))
	AccessKeySecret := fmt.Sprint(cfg.Section("oss").Key("AccessKeySecret"))

	// 创建OSSClient实例。
	client, err := oss.New(OSSEndpoint, AccessKeyId, AccessKeySecret)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写存储空间名称。
	OSSBucket := fmt.Sprint(cfg.Section("oss").Key("OSSBucket"))
	bucket, err := client.Bucket(OSSBucket)

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 列举包含指定前缀的文件。
	OSSDir := fmt.Sprint(cfg.Section("oss").Key("OSSDir"))
	lsRes, err := bucket.ListObjects(oss.Prefix(OSSDir), oss.MaxKeys(1000))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	StartTimeStr := fmt.Sprint(cfg.Section("oss").Key("StartTime"))
	MonthsAgoStr := fmt.Sprint(cfg.Section("oss").Key("MonthsAgo"))
	MonthsAgoInt, err := strconv.Atoi(MonthsAgoStr)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// t1time, _ := time.ParseInLocation("2006-01-02 15:04:05", t1str, time.Local)
	MonthsAgo := time.Now().AddDate(0, MonthsAgoInt, 0).Format("2006-01-02")
	StartTimeTime, _ := time.ParseInLocation("2006-01-02", StartTimeStr, time.Local)
	MonthsAgoTime, _ := time.ParseInLocation("2006-01-02", MonthsAgo, time.Local)

	for _, object := range lsRes.Objects {
		//LastModifiedTime, _ := time.ParseInLocation("2006-01-02", object.LastModified.String(), time.Local)
		// Before 在之前
		if StartTimeTime.Before(object.LastModified) {
			if object.LastModified.Before(MonthsAgoTime) {
				fmt.Println("删除", StartTimeTime, "到", MonthsAgoTime, "之间的文件", object.Key)
				// After 在之后
				//if MonthsAgoTime.After(object.LastModified) {
				err := bucket.DeleteObject(object.Key)
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(-1)
				}
				fmt.Println("删除文件成功:", object.Key)
			}
		}
	}
}
