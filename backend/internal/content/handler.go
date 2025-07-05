package content

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strings"
)

func (s *Service) HandleCreateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	textEnglish := strings.TrimSpace(r.FormValue("textEnglish"))
	textLatin := strings.TrimSpace(r.FormValue("textLatin"))
	theme := strings.TrimSpace(r.FormValue("theme"))
	imageSource := strings.TrimSpace(r.FormValue("imageSource"))
	textSource := strings.TrimSpace(r.FormValue("textSource"))
	if textEnglish == "" {
		http.Error(w, "textEnglish is required", http.StatusBadRequest)
		return
	}

	if textLatin == "" {
		http.Error(w, "textLatin is required", http.StatusBadRequest)
		return
	}

	if theme == "" {
		http.Error(w, "theme is required", http.StatusBadRequest)
		return
	}

	var file multipart.File
	var header *multipart.FileHeader
	var err error
	
	if file, header, err = r.FormFile("image"); err != nil {
		file = nil
		header = nil
	} else {

		defer file.Close()
		
		if header.Size > 5<<20 {
			http.Error(w, "Image file too large (max 5MB)", http.StatusBadRequest)
			return
		}

		contentType := header.Header.Get("Content-Type")
		if !isValidImageType(contentType) {
			http.Error(w, "Invalid image type. Only JPEG, PNG, and GIF are allowed", http.StatusBadRequest)
			return
		}
	}

	content, err := s.AddContent(textEnglish, textLatin, file, header, theme, imageSource, textSource)
	if err != nil {
		http.Error(w, "Failed to add content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Successfully added content: " + content.TextLatin,
		"imageURL": content.ImageURL,
	})
}
