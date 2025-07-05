package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"novissima/internal/logging"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

type Content struct {
	ID          uuid.UUID  `json:"id"`
	TextEnglish string     `json:"text_english"`
	TextLatin   string     `json:"text_latin"`
	ImageURL    string     `json:"image_url"`
	LastSent    *time.Time `json:"last_sent"`
	Theme       string     `json:"theme"`
	ImageSource      *string    `json:"image_source"`
	TextSource      *string    `json:"text_source"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Service struct {
	dbClient        *supabase.Client
	loggingService  *logging.Service
	bucketName      string
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
		
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("content/%d%s", time.Now().UnixNano(), ext)
		
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return Content{}, fmt.Errorf("failed to read file: %w", err)
		}

		bucketName := s.bucketName
		_, err = s.dbClient.Storage.UploadFile(bucketName, filename, bytes.NewReader(fileBytes))
		if err != nil {
			return Content{}, fmt.Errorf("failed to upload to storage: %w", err)
		}

		imageURL = s.dbClient.Storage.GetPublicUrl(bucketName, filename).SignedURL
	}

	content := Content{
		TextEnglish: contentEnglish,
		TextLatin:   contentLatin,
		ImageURL:    imageURL,
		LastSent:    nil,
		Theme:       theme,
		ImageSource: &imageSource,
		TextSource: &textSource,
	}

	data, _, err := s.dbClient.From("content").Insert(content, true, "", "", "").Execute()
	if err != nil {
		log.Printf("Supabase error: %v", err)
		if data != nil {
			log.Printf("Supabase response: %s", string(data))
		}
		return Content{}, fmt.Errorf("failed to add content: %w", err)
	}

	s.loggingService.LogContentCreated(content.ID, content.TextEnglish, content.TextLatin, content.ImageURL, content.Theme, *content.ImageSource, *content.TextSource)
	log.Printf("Successfully added content: %s", data)

	return content, nil
}

func (s *Service) GetDailyContent() (*Content, error) {
	themes := []string{"heaven", "hell", "judgment", "death"}
	startDate := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	
	now := time.Now()
	
	cycleDay := int(now.Sub(startDate).Hours() / 24) % len(themes)

	currentTheme := themes[cycleDay]
	
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	
	content, _, err := s.dbClient.From("content").
		Select("*", "", false).
		Eq("theme", currentTheme).
		Or("last_sent.is.null", "").
		Or("last_sent.lt."+thirtyDaysAgo.Format("2006-01-02"), "").
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