package imageops

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func FetchSource(ctx context.Context, sourceURL string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("IMAGE_FETCH_FAILED")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("IMAGE_FETCH_FAILED")
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, "", fmt.Errorf("IMAGE_FETCH_FAILED")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", fmt.Errorf("IMAGE_FETCH_FAILED")
	}
	return body, res.Header.Get("Content-Type"), nil
}
