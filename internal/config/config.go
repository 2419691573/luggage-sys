package config

import (
	"fmt"
	"os"
	"strings"
)

var (
	DBDSN     string
	JWTSecret string
	Port      string

	// MinIO 配置
	MinIOEndpoint        string
	MinIOAccessKeyID     string
	MinIOSecretAccessKey string
	MinIOUseSSL          bool
	MinIOBucketName      string
)

func Init() {
	DBDSN = os.Getenv("DB_DSN")
	if DBDSN == "" {
		DBDSN = "root:123456@tcp(127.0.0.1:3306)/hotel_luggage?charset=utf8mb4&parseTime=True&loc=Local"
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		JWTSecret = "your-secret-key"
	}

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "10.154.39.253:8080"
	}

	// MinIO 配置
	MinIOEndpoint = os.Getenv("MINIO_ENDPOINT")
	if MinIOEndpoint == "" {
		MinIOEndpoint = "minio.2huo.tech:443" // 默认地址（使用 HTTPS 端口）
	}

	MinIOAccessKeyID = os.Getenv("MINIO_ACCESS_KEY_ID")
	if MinIOAccessKeyID == "" {
		MinIOAccessKeyID = "i8IuD8lJYxE5kAL1HOwS" // 默认 AccessKey
	}

	MinIOSecretAccessKey = os.Getenv("MINIO_SECRET_ACCESS_KEY")
	if MinIOSecretAccessKey == "" {
		MinIOSecretAccessKey = "lAfdJNMqAQDmNrK8peuIwu5un6PFI0EtgWlae7jv" // 默认 SecretKey
	}

	// 读取 SSL 配置，如果环境变量未设置，根据端口判断
	sslEnv := os.Getenv("MINIO_USE_SSL")
	if sslEnv == "true" {
		MinIOUseSSL = true
	} else if sslEnv == "false" {
		MinIOUseSSL = false
	} else {
		// 如果环境变量未设置，根据端口自动判断：443 用 HTTPS，其他用 HTTP
		if MinIOEndpoint != "" {
			MinIOUseSSL = strings.Contains(MinIOEndpoint, ":443") || strings.HasPrefix(MinIOEndpoint, "https://")
		}
	}

	MinIOBucketName = os.Getenv("MINIO_BUCKET_NAME")
	if MinIOBucketName == "" {
		MinIOBucketName = "training-hotel" // 默认桶名
	}
	
	// 打印 MinIO 配置（用于调试）
	fmt.Printf("MinIO Config: endpoint=%s, bucket=%s, accessKey=%s, useSSL=%v\n", 
		MinIOEndpoint, MinIOBucketName, MinIOAccessKeyID, MinIOUseSSL)
}
