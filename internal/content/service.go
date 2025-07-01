package content

import (
	"bytes"
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
	Source      *string    `json:"source"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Service struct {
	client *supabase.Client
	loggingService *logging.Service
}

func NewService(client *supabase.Client, loggingService *logging.Service) *Service {
	return &Service{
		client: client,
		loggingService: loggingService,
	}
}

func (s *Service) AddContent(contentEnglish string, contentLatin string, file multipart.File, header *multipart.FileHeader, theme string, source string) (Content, error) {
	var imageURL string
	
	if file != nil && header != nil {
		defer file.Close()
		
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("content/%d%s", time.Now().UnixNano(), ext)
		
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return Content{}, fmt.Errorf("failed to read file: %w", err)
		}

		bucketName := "content-images"
		_, err = s.client.Storage.UploadFile(bucketName, filename, bytes.NewReader(fileBytes))
		if err != nil {
			return Content{}, fmt.Errorf("failed to upload to storage: %w", err)
		}

		imageURL = s.client.Storage.GetPublicUrl(bucketName, filename).SignedURL
	}

	content := Content{
		TextEnglish: contentEnglish,
		TextLatin:   contentLatin,
		ImageURL:    imageURL,
		LastSent:    nil,
		Theme:       theme,
		Source:      &source,
	}

	data, _, err := s.client.From("content").Insert(content, true, "", "", "").Execute()
	if err != nil {
		log.Printf("Supabase error: %v", err)
		if data != nil {
			log.Printf("Supabase response: %s", string(data))
		}
		return Content{}, fmt.Errorf("failed to add content: %w", err)
	}

	s.loggingService.LogContentCreated(content.ID, content.TextEnglish, content.TextLatin, content.ImageURL, content.Theme, *content.Source)
	log.Printf("Successfully added content: %s", data)

	return content, nil
}

func (s *Service) GetDailyContent() error {	
	
	return nil
}

func (s *Service) SendDailyContent() error {
	
	return nil
}