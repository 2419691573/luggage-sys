package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"luggage-sys2/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOService struct {
	client     *minio.Client
	bucketName string
}

var minioServiceInstance *MinIOService

// GetMinIOService 获取 MinIO 服务单例
func GetMinIOService() (*MinIOService, error) {
	if minioServiceInstance != nil {
		return minioServiceInstance, nil
	}

	// 连接 MinIO 服务器
	client, err := minio.New(config.MinIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinIOAccessKeyID, config.MinIOSecretAccessKey, ""),
		Secure: config.MinIOUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	service := &MinIOService{
		client:     client,
		bucketName: config.MinIOBucketName,
	}

	// 尝试检查存储桶，如果失败则直接使用（存储桶已存在，只是可能没有检查权限）
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.MinIOBucketName)
	if err != nil {
		// 检查失败，可能是权限问题，但存储桶已存在
		// 直接使用存储桶，如果不存在会在上传时失败
		log.Printf("Warning: Cannot check bucket existence (may be permission issue): %v, will use bucket directly", err)
	} else if !exists {
		// 存储桶不存在，尝试创建
		log.Printf("Bucket %s does not exist, attempting to create...", config.MinIOBucketName)
		err = client.MakeBucket(ctx, config.MinIOBucketName, minio.MakeBucketOptions{})
		if err != nil {
			// 创建失败，可能是权限不足或存储桶已存在
			// 继续使用，如果存储桶实际不存在会在上传时失败
			log.Printf("Warning: Failed to create bucket (may already exist or permission issue): %v, will try to use it anyway", err)
		} else {
			log.Printf("Successfully created MinIO bucket: %s", config.MinIOBucketName)

			// 尝试设置存储桶为公开读取（可选）
			policy := `{
				"Version": "2012-10-17",
				"Statement": [{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + config.MinIOBucketName + `/*"]
				}]
			}`
			err = client.SetBucketPolicy(ctx, config.MinIOBucketName, policy)
			if err != nil {
				log.Printf("Warning: Failed to set bucket policy (this is optional): %v", err)
			}
		}
	} else {
		log.Printf("Bucket %s exists and is accessible", config.MinIOBucketName)
	}
	
	// 无论检查或创建是否成功，都继续初始化服务
	// 如果存储桶实际不存在，会在上传时失败并给出明确错误

	minioServiceInstance = service
	log.Printf("MinIO service initialized: endpoint=%s, bucket=%s", config.MinIOEndpoint, config.MinIOBucketName)
	return service, nil
}

// SaveImageToMinIO 上传图片到 MinIO
func (s *MinIOService) SaveImageToMinIO(fileHeader *multipart.FileHeader, maxBytes int64) (*UploadResult, error) {
	if fileHeader == nil {
		return nil, ErrMissingFileField
	}

	if maxBytes > 0 && fileHeader.Size > maxBytes {
		return nil, ErrFileTooLarge
	}

	// 打开文件
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 检测文件类型
	head := make([]byte, 512)
	n, _ := io.ReadFull(src, head)
	detected := http.DetectContentType(head[:n])

	// 重新打开文件
	_ = src.Close()
	src, err = fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 检查文件类型
	ext, ok := imageExtFromContentType(detected)
	if !ok {
		return nil, ErrInvalidFileType
	}

	// 生成文件路径（按年月组织，保持和本地存储一样的格式）
	now := time.Now()
	objectName := fmt.Sprintf("%04d/%02d/%s%s",
		now.Year(),
		int(now.Month()),
		randomHex(16),
		ext)

	// 上传到 MinIO
	ctx := context.Background()
	info, err := s.client.PutObject(ctx, s.bucketName, objectName, src, fileHeader.Size, minio.PutObjectOptions{
		ContentType: detected,
	})
	if err != nil {
		// 提供更详细的错误信息
		if err.Error() == "The specified bucket does not exist." {
			return nil, fmt.Errorf("bucket '%s' does not exist or AccessKey has no permission to access it. Please check: 1) bucket name is correct, 2) AccessKey has permission to access this bucket", s.bucketName)
		}
		return nil, fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	// 返回相对URL（格式和本地存储一样：/uploads/2026/01/xxx.jpg）
	// 这样前端代码完全不需要改
	relativeURL := "/uploads/" + objectName

	return &UploadResult{
		RelativeURL: relativeURL,
		FileName:    info.Key,
		Size:        info.Size,
		ContentType: detected,
	}, nil
}

// GetObject 从 MinIO 获取文件（用于代理访问）
func (s *MinIOService) GetObject(objectPath string) (*minio.Object, error) {
	// 去掉 /uploads/ 前缀，得到 MinIO 中的对象路径
	// 例如：/uploads/2026/01/abc.jpg -> 2026/01/abc.jpg
	objectName := objectPath
	if len(objectPath) > 9 && objectPath[:9] == "/uploads/" {
		objectName = objectPath[9:]
	}

	ctx := context.Background()
	obj, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from MinIO: %w", err)
	}

	return obj, nil
}
