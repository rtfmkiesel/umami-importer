package umami

import (
	"bufio"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rtfmkiesel/umami-importer/config"
	"github.com/rtfmkiesel/umami-importer/db"
	"github.com/rtfmkiesel/umami-importer/hash"
)

const (
	// How many goroutines process log lines into events.
	// Should be smaller as MaxRequests because it is way faster
	// to process logs then to send requests
	logLineWorkers = 15
)

func (c *Client) processFile(filePath string, importConfig *config.ImportConfig, ctx context.Context) error {
	fileHash, err := hash.File(filePath)
	if err != nil {
		return err
	}

	// Check if the file itself was already processed fully
	known, err := db.IsKnown(fileHash)
	if err != nil {
		return err
	}

	if known {
		log.Debug("Skipping file '%s', already processed", filePath)
		return nil
	}

	chanLogLines := make(chan string)
	chanEvents := make(chan *Event)
	wgLogEntryWorkers := new(sync.WaitGroup)
	wgEventWorkers := new(sync.WaitGroup)

	parser := newLogParser(importConfig)

	// Spawns goroutines for reading log lines
	for range logLineWorkers {
		wgLogEntryWorkers.Go(func() {
			for line := range chanLogLines {
				logEntry, err := parser.ParseLine(line)
				if err != nil {
					log.Warning("Line '%s': %v", line, err)
					return
				}

				event, err := entrytoUmamiEvent(logEntry, importConfig.Website.BaseURL, importConfig.Website.ID)
				if err != nil {
					log.Warning("Line '%s': %v", line, err)
					return
				}
				event.Hash = hash.String(line)

				chanEvents <- event
			}
		})
	}

	// Spawn goroutines for sending events to Umami
	for range c.config.MaxRequests {
		wgEventWorkers.Go(func() {
			for event := range chanEvents {
				if err := c.importEventWithRetries(event); err != nil {
					log.Error(err)
					return
				}

				if err := db.Add(event.Hash); err != nil {
					log.Error(err)
					return
				}
			}
		})
	}

	fh, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fh.Close() //nolint:errcheck

	scanner := bufio.NewScanner(fh)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	// Read the log file
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
			line := scanner.Text()

			lineHash := hash.String(line)
			known, err := db.IsKnown(lineHash)
			if err != nil {
				return err
			}

			// Skipped already processed lines
			if known {
				continue
			}

			chanLogLines <- line
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	close(chanLogLines)
	wgLogEntryWorkers.Wait()
	close(chanEvents)
	wgEventWorkers.Wait()

	// Mark the document hash as completed
	if err := db.Add(fileHash); err != nil {
		return err
	}

	log.Debug("Completed processing file '%s'", filePath)
	return nil
}

// Simple file walker with extension filter and recursive option
func getFilesFromDirectory(dirPath, filterExt string, recursive bool) ([]string, error) {
	var files []string

	if filterExt != "" && !strings.HasPrefix(filterExt, ".") {
		filterExt = "." + filterExt
	}

	if !recursive {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, err
		}

		files = make([]string, 0, len(entries))

		if filterExt == "" {
			for _, e := range entries {
				if !e.IsDir() {
					files = append(files, filepath.Join(dirPath, e.Name()))
				}
			}
			return files, nil
		}

		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if filepath.Ext(name) == filterExt {
				files = append(files, filepath.Join(dirPath, name))
			}
		}
		return files, nil
	}

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filterExt == "" {
			files = append(files, path)
			return nil
		}

		if filepath.Ext(d.Name()) == filterExt {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
