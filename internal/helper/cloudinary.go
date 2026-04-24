package helper

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

const maxFileSize = 5 << 20 // 5MB

// ValidatePhoto memvalidasi ukuran dan ekstensi file foto
func ValidatePhoto(file *multipart.FileHeader) error {
	if file.Size > maxFileSize {
		return fmt.Errorf("file_too_large")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("invalid_file_format")
	}

	return nil
}

// UploadReportPhoto mengupload foto laporan ke Cloudinary dan mengembalikan URL-nya.
// Instance Cloudinary di-inject dari luar (dari config), bukan dibuat ulang tiap call.
func UploadReportPhoto(ctx context.Context, cld *cloudinary.Cloudinary, file *multipart.FileHeader, reportID string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	publicID := fmt.Sprintf("reports/%s", reportID)

	result, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID: publicID,
		Folder:   "temukan",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to cloudinary: %w", err)
	}

	return result.SecureURL, nil
}
