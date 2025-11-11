package umami

import (
	"net/url"
	"regexp"

	"github.com/rtfmkiesel/umami-importer/config"
	"github.com/rtfmkiesel/umami-importer/logparser"
)

// Select a logparser based on *config.ImportConfig
func newLogParser(importConfig *config.ImportConfig) logparser.Parser {
	switch importConfig.Logs.Type {
	case "apache":
		return &logparser.ParserApache2{}
	case "nginx":
		return &logparser.ParserNginx{}
	case "custom":
		return &logparser.ParserCustom{
			ReLogLine:  regexp.MustCompile(importConfig.Logs.TypeCustomRegex), // Covered by go-playground/validator
			TimeFmtStr: importConfig.Logs.TypeCustomTimestamp,
		}
	default:
		// Covered by go-playground/validator
		panic("invalid logs.type")
	}
}

// Convert a *logparser.Entry to a *Event
func entrytoUmamiEvent(entry *logparser.Entry, websiteUrl, websiteId string) (*Event, error) {
	u, err := url.Parse(websiteUrl)
	if err != nil {
		return nil, log.NewError(err)
	}

	return &Event{
		WebsiteId:  websiteId,
		Hostname:   u.Host,
		Referrer:   entry.Referrer,
		Url:        entry.Url,
		UserAgent:  entry.UserAgent,
		RemoteAddr: entry.RemoteAddr,
		TimeStamp:  entry.Timestamp,
	}, nil
}
