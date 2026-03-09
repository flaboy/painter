package api

type ImageSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Resize struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Fit    string `json:"fit"`
}

type ImageResult struct {
	MimeType    string `json:"mimeType"`
	Format      string `json:"format"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	BytesBase64 string `json:"bytesBase64"`
}

type GenerateImageRequest struct {
	RequestID      string    `json:"requestId"`
	Prompt         string    `json:"prompt"`
	NegativePrompt string    `json:"negativePrompt,omitempty"`
	Provider       string    `json:"provider,omitempty"`
	Style          string    `json:"style,omitempty"`
	Size           ImageSize `json:"size"`
	Format         string    `json:"format,omitempty"`
	Background     string    `json:"background,omitempty"`
	Quality        string    `json:"quality,omitempty"`
}

type EditImageRequest struct {
	RequestID string    `json:"requestId"`
	Mode      string    `json:"mode"`
	SourceUrl string    `json:"sourceUrl"`
	MaskUrl   string    `json:"maskUrl,omitempty"`
	Prompt    string    `json:"prompt,omitempty"`
	Provider  string    `json:"provider,omitempty"`
	Size      ImageSize `json:"size"`
	Format    string    `json:"format,omitempty"`
}

type ConvertImageRequest struct {
	RequestID  string `json:"requestId"`
	SourceUrl  string `json:"sourceUrl"`
	Format     string `json:"format"`
	Resize     Resize `json:"resize"`
	Quality    int    `json:"quality,omitempty"`
	Background string `json:"background,omitempty"`
}

func (r GenerateImageRequest) Validate() error {
	if r.Prompt == "" {
		return invalidRequestError("prompt is required")
	}
	if r.Size.Width <= 0 || r.Size.Height <= 0 {
		return invalidRequestError("size.width and size.height are required")
	}
	if !isSupportedFormat(r.Format, true) {
		return invalidRequestError("format must be png or jpeg")
	}
	return nil
}

func (r EditImageRequest) Validate() error {
	if r.Mode == "" {
		return invalidRequestError("mode is required")
	}
	if r.SourceUrl == "" {
		return invalidRequestError("sourceUrl is required")
	}
	if !isSupportedFormat(r.Format, true) {
		return invalidRequestError("format must be png or jpeg")
	}
	return nil
}

func (r ConvertImageRequest) Validate() error {
	if r.SourceUrl == "" {
		return invalidRequestError("sourceUrl is required")
	}
	if r.Format == "" {
		return invalidRequestError("format is required")
	}
	if !isSupportedFormat(r.Format, false) {
		return invalidRequestError("format must be png or jpeg")
	}
	return nil
}

type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func invalidRequestError(message string) *ValidationError {
	return &ValidationError{Code: "INVALID_REQUEST", Message: message}
}

func isSupportedFormat(format string, allowEmpty bool) bool {
	switch format {
	case "":
		return allowEmpty
	case "png", "jpeg", "jpg":
		return true
	default:
		return false
	}
}
