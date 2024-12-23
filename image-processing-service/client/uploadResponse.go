package client

type UploadResponse struct {
	Metadata map[string]string
	Url      string
}

type ImageListResponse struct {
	Images    []UploadResponse `json:"images"`
	NextToken string           `json:"next_token"`
	Page      int              `json:"page"`
}
