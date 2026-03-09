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
