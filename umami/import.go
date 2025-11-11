package umami

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rtfmkiesel/umami-importer/config"
)

func (c *Client) Import(config *config.ImportConfig, ctx context.Context) error {
	var filesToProcess []string

	log.Info("Starting import for '%s'", config.Website.BaseURL)

	// Gather files to process
	for _, logpath := range config.Logs.Paths {
		log.Debug("Trying logpath '%s'", logpath)
		info, err := os.Stat(logpath)
		if err != nil {
			return log.NewError(err)
		}

		if info.IsDir() {
			log.Debug("Identified '%s' as dir, checking for files (recursive=%t)", logpath, config.Logs.Recursive)

			files, err := getFilesFromDirectory(logpath, config.Logs.IncludeExtension, config.Logs.Recursive)
			if err != nil {
				return log.NewError(err)
			}
			log.Debug("Found %d files to import from '%s'", len(files), logpath)
			filesToProcess = append(filesToProcess, files...)
		} else {
			log.Debug("Identified '%s' as single file, importing a single file", logpath)
			filesToProcess = append(filesToProcess, logpath)
		}
	}

	for _, filePath := range filesToProcess {
		select {
		case <-ctx.Done():
			return nil
		default:
			log.Info("Processing file '%s'", filePath)
			if err := c.processFile(filePath, config, ctx); err != nil {
				return log.NewError(err) // Hard return, do not process further
			}
		}
	}

	return nil
}

func (c *Client) importEventWithRetries(e *Event) error {
	for attempt := 1; attempt <= c.config.Retries; attempt++ {
		if attempt > 1 {
			time.Sleep(time.Duration(attempt) * time.Second)
			log.Warning("Retrying import for '%s' (attempt %d/%d)", e.Url, attempt, c.config.Retries)
		}

		err := c.importEventOnce(e)
		if err == nil {
			return nil
		}
		log.Error(err)
	}

	return fmt.Errorf("request for '%s' not imported: retries exhausted", e.Url)
}

func (c *Client) importEventOnce(e *Event) error {
	jsonBody, err := e.ToUmamiRequestBodyJson()
	if err != nil {
		return err
	}
	log.Debug("Importing event: payload='%s'", jsonBody)

	req, err := http.NewRequest(http.MethodPost, c.config.CollectionURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	if c.config.CustomHTTPHeaders != nil {
		for key, value := range c.config.CustomHTTPHeaders {
			req.Header.Set(key, value)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", e.UserAgent) // Set the User-Agent header of the request to the one from the event

	log.Debug("Sending request to Umami: payload='%s'", jsonBody)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to import event: request failed with status %d: %s", resp.StatusCode, body)
	}

	return nil
}
