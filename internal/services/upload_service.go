package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadService struct {
}

func NewUploadService() *UploadService {
	return &UploadService{}
}

type UploadResult struct {
	RelativeURL string `json:"relative_url"`
	FileName    string `json:"file_name"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

var (
	ErrFileTooLarge     = errors.New("file too large")
	ErrInvalidFileType  = errors.New("invalid file type")
	ErrMissingFileField = errors.New("missing file")
)

// SaveImageFile saves an uploaded image to local disk under:
//   uploads/YYYY/MM/<random>.<ext>
//
// It returns a relative URL like: /uploads/2026/01/xxxx.jpg
func (s *UploadService) SaveImageFile(fileHeader *multipart.FileHeader, maxBytes int64) (*UploadResult, error) {
	if fileHeader == nil {
		return nil, ErrMissingFileField
	}
	if maxBytes > 0 && fileHeader.Size > maxBytes {
		return nil, ErrFileTooLarge
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// sniff content-type
	head := make([]byte, 512)
	n, err := io.ReadFull(src, head)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, fmt.Errorf("failed to read file header: %w", err)
	}
	if n == 0 {
		return nil, errors.New("file is empty")
	}
	detected := http.DetectContentType(head[:n])

	// reset reader by reopening (multipart.File may not support Seek)
	_ = src.Close()
	src, err = fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	ext, ok := imageExtFromContentType(detected)
	if !ok {
		return nil, ErrInvalidFileType
	}

	now := time.Now()
	dir := filepath.Join("uploads", fmt.Sprintf("%04d", now.Year()), fmt.Sprintf("%02d", int(now.Month())))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	name := randomHex(16) + ext
	dstPath := filepath.Join(dir, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	written, err := io.Copy(dst, io.LimitReader(src, maxBytes))
	if err != nil {
		return nil, err
	}

	relativeURL := "/" + filepath.ToSlash(filepath.Join(dir, name))
	return &UploadResult{
		RelativeURL: relativeURL,
		FileName:    name,
		Size:        written,
		ContentType: detected,
	}, nil
}

func imageExtFromContentType(ct string) (string, bool) {
	ct = strings.ToLower(ct)
	switch ct {
	case "image/jpeg":
		return ".jpg", true
	case "image/png":
		return ".png", true
	case "image/webp":
		return ".webp", true
	default:
		return "", false
	}
}

func randomHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

