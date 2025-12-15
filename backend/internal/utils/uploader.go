package utils

import (
	"bytes"
	"e-commerce/backend/internal/config"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	storage_go "github.com/supabase-community/storage-go"
)

// UploadToSupabase mengupload file ke bucket Supabase Storage
func UploadToSupabase(file *multipart.FileHeader, folderName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %w", err)
	}
	defer src.Close()

	buffer := bytes.NewBuffer(nil)
	if _, err := buffer.ReadFrom(src); err != nil {
		return "", fmt.Errorf("gagal membaca file: %w", err)
	}

	if buffer.Len() > 5*1024*1024 {
		return "", fmt.Errorf("ukuran file terlalu besar (max 5MB)")
	}

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	folderName = strings.ReplaceAll(folderName, " ", "_")
	folderName = strings.Trim(folderName, "/")

	fullPath := fmt.Sprintf("%s/%s", folderName, fileName)

	contentType := getContentType(ext)

	fileOptions := storage_go.FileOptions{
		ContentType: &contentType,
	}

	fileReader := bytes.NewReader(buffer.Bytes())

	_, err = config.SupabaseStorage.UploadFile("product", fullPath, fileReader, fileOptions)
	if err != nil {
		return "", fmt.Errorf("gagal upload ke supabase: %w", err)
	}

	publicURL := config.SupabaseStorage.GetPublicUrl("product", fullPath)

	return publicURL.SignedURL, nil
}

// DeleteFromSupabase menghapus file dari Supabase Storage
// Parameter fileURL: full URL dari file yang akan dihapus
func DeleteFromSupabase(fileURL string) error {
	if fileURL == "" {
		return nil 
	}

	filePath := extractFilePathFromURL(fileURL)
	if filePath == "" {
		return fmt.Errorf("invalid file URL: %s", fileURL)
	}

	_, err := config.SupabaseStorage.RemoveFile("product", []string{filePath})
	if err != nil {
		return fmt.Errorf("gagal menghapus file dari supabase: %w", err)
	}

	return nil
}

func DeleteMultipleFromSupabase(fileURLs []string) error {
	if len(fileURLs) == 0 {
		return nil
	}

	filePaths := make([]string, 0, len(fileURLs))
	for _, url := range fileURLs {
		if url != "" {
			path := extractFilePathFromURL(url)
			if path != "" {
				filePaths = append(filePaths, path)
			}
		}
	}

	if len(filePaths) == 0 {
		return nil
	}

	_, err := config.SupabaseStorage.RemoveFile("product", filePaths)
	if err != nil {
		return fmt.Errorf("gagal menghapus files dari supabase: %w", err)
	}

	return nil
}


func extractFilePathFromURL(fileURL string) string {
	if fileURL == "" {
		return ""
	}

	parts := strings.Split(fileURL, "/product/")
	if len(parts) < 2 {
		return ""
	}

	return parts[1]
}

func ReplaceFile(oldFileURL string, newFile *multipart.FileHeader, folderName string) (string, error) {

	newFileURL, err := UploadToSupabase(newFile, folderName)
	if err != nil {
		return "", fmt.Errorf("gagal upload file baru: %w", err)
	}

	if oldFileURL != "" {
		go func() {
			if err := DeleteFromSupabase(oldFileURL); err != nil {
				fmt.Printf("Warning: gagal menghapus file lama: %v\n", err)
			}
		}()
	}

	return newFileURL, nil
}

// GetFileInfo mendapatkan informasi file dari URL
func GetFileInfo(fileURL string) map[string]string {
	filePath := extractFilePathFromURL(fileURL)
	if filePath == "" {
		return nil
	}

	parts := strings.Split(filePath, "/")
	fileName := parts[len(parts)-1]
	ext := filepath.Ext(fileName)

	return map[string]string{
		"path":      filePath,
		"filename":  fileName,
		"extension": ext,
		"folder":    strings.TrimSuffix(filePath, "/"+fileName),
	}
}

// IsValidImageURL memeriksa apakah URL adalah image yang valid
func IsValidImageURL(fileURL string) bool {
	if fileURL == "" {
		return false
	}

	ext := strings.ToLower(filepath.Ext(fileURL))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}

	return false
}

// getContentType menentukan MIME type berdasarkan extension
func getContentType(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
