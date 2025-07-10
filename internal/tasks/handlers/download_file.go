package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/logging"
)

type DownloadFilePayload struct {
	Url      string `json:"url"`
	Filename string `json:"filename"`
}

func DownloadFile(ctx context.Context, url string, filename string) error {
	if url == "" || filename == "" {
		return fmt.Errorf("url and filename must not be empty")
	}

	outputDir := path.Join(config.ProjectRootPath, "storage")
	outputPath := path.Join(outputDir, filename)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer outFile.Close()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Remove(outFile.Name())
		return fmt.Errorf("non-200 response: %s", resp.Status)
	}

	n, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	logging.DebugLog(fmt.Sprintf("Downloaded %d bytes to %s\n", n, outputPath))

	return nil
}
