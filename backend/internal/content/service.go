package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"novissima/internal/logging"

	"github.com/google/uuid"
	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
)

type Content struct {
	ID          uuid.UUID  `json:"id"`
	TextEnglish string     `json:"text_english"`
	TextLatin   *string    `json:"text_latin"`
	ImageURL    *string    `json:"image_url"`
	LastSent    *time.Time `json:"last_sent"`
	Theme       string     `json:"theme"`
	ImageSource *string    `json:"image_source"`
	TextSource  *string    `json:"text_source"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type ContentCreate struct {
	TextEnglish string  `json:"text_english"`
	TextLatin   *string `json:"text_latin"`
	ImageURL    *string `json:"image_url"`
	Theme       string  `json:"theme"`
	ImageSource *string `json:"image_source"`
	TextSource  *string `json:"text_source"`
}

type Service struct {
	dbClient       *supabase.Client
	loggingService *logging.Service
	bucketName     string
}

func NewService(dbClient *supabase.Client, loggingService *logging.Service, bucketName string) *Service {
	return &Service{
		dbClient:       dbClient,
		loggingService: loggingService,
		bucketName:     bucketName,
	}
}

func (s *Service) AddContent(contentEnglish string, contentLatin string, file multipart.File, header *multipart.FileHeader, theme string, imageSource string, textSource string) (Content, error) {

	var imageURL string

	if file != nil && header != nil {
		defer file.Close()

		fileExt := ""
		if idx := strings.LastIndex(header.Filename, "."); idx != -1 {
			fileExt = header.Filename[idx:]
		}
		filename := fmt.Sprintf("%s/%s%s", theme, uuid.New().String(), fileExt)

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return Content{}, fmt.Errorf("failed to read file: %w", err)
		}

		fileReader := bytes.NewReader(fileBytes)

		contentType := header.Header.Get("Content-Type")
		if contentType == "" {
			switch fileExt {
			case ".jpg", ".jpeg":
				contentType = "image/jpeg"
			case ".png":
				contentType = "image/png"
			case ".gif":
				contentType = "image/gif"
			default:
				contentType = "application/octet-stream"
			}
		}

		fileOptions := storage_go.FileOptions{
			ContentType: &contentType,
		}
		_, err = s.dbClient.Storage.UploadFile(s.bucketName, filename, fileReader, fileOptions)
		if err != nil {
			return Content{}, fmt.Errorf("failed to upload to storage: %w", err)
		}
		imageURL = s.dbClient.Storage.GetPublicUrl(s.bucketName, filename).SignedURL
	}

	content := ContentCreate{
		TextEnglish: contentEnglish,
		TextLatin:   &contentLatin,
		ImageURL:    &imageURL,
		Theme:       theme,
		ImageSource: &imageSource,
		TextSource:  &textSource,
	}

	data, _, err := s.dbClient.From("content").Insert(content, true, "", "", "").Execute()
	if err != nil {
		log.Printf("Supabase error: %v", err)
		if data != nil {
			log.Printf("Supabase response: %s", string(data))
		}
		return Content{}, fmt.Errorf("failed to add content: %w", err)
	}

	log.Printf("Successfully added content: %s", data)

	return Content{}, nil
}

func (s *Service) GetDailyContent() (*Content, error) {
	themes := []string{"death"}
	startDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)

	now := time.Now()

	cycleDay := int(now.Sub(startDate).Hours()/24) % len(themes)

	currentTheme := themes[cycleDay]

	// thirtyDaysAgo := now.AddDate(0, 0, -30)

	content, _, err := s.dbClient.From("content").
		Select("*", "", false).
		Eq("theme", currentTheme).
		Limit(1, "").
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get daily content: %w", err)
	}

	var contents []Content
	if err := json.Unmarshal(content, &contents); err != nil {
		return nil, fmt.Errorf("failed to parse content: %w", err)
	}

	if len(contents) == 0 {
		return nil, fmt.Errorf("no content available for today's theme: %s", currentTheme)
	}

	return &contents[0], nil
}

func (s *Service) UpdateLastSent(id uuid.UUID) error {
	_, _, err := s.dbClient.From("content").Update(map[string]interface{}{
		"last_sent": time.Now(),
	}, "", "").Eq("id", id.String()).Execute()
	return err
}

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

