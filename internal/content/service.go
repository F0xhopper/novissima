package content

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/supabase-community/supabase-go"
)

type Content struct {
	TextEnglish string
	TextLatin string
	ImageURL string
	LastSent *time.Time
	Theme string
	Source *string
	UpdatedAt time.Time
	CreatedAt time.Time
}

type Service struct {
	client *supabase.Client
	
}

func NewService(client *supabase.Client) *Service {
	return &Service{
		client: client,
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
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	data, _, err := s.client.From("content").Insert(content, true, "", "", "").Execute()
	if err != nil {
		log.Printf("Supabase error: %v", err)
		if data != nil {
			log.Printf("Supabase response: %s", string(data))
		}
		return Content{}, fmt.Errorf("failed to add content: %w", err)
	}

	log.Printf("Successfully added content: %s", data)

	return content, nil
}

func (s *Service) GetDailyContent() error {	
	
	return nil
}
