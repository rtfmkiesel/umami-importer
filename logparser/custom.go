package logparser

import (
	"regexp"
	"time"
)

// Parses custom log lines using named regex groups
type ParserCustom struct {
	ReLogLine  *regexp.Regexp
	TimeFmtStr string
}

func (p *ParserCustom) ParseLine(line string) (*Entry, error) {
	matches := p.ReLogLine.FindStringSubmatch(line)
	if matches == nil {
		return nil, log.NewError("failed to parse log line: regex did not match")
	}

	names := p.ReLogLine.SubexpNames()
	result := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" && i < len(matches) {
			result[name] = matches[i]
		}
	}

	requiredFields := []string{"remote_addr", "timestamp", "url", "user_agent"}
	for _, field := range requiredFields {
		if _, exists := result[field]; !exists {
			return nil, log.NewError("required field '%s' not found in regex", field)
		}
	}

	timeStr := result["timestamp"]
	var timestamp time.Time
	timestamp, err := time.Parse(p.TimeFmtStr, timeStr)
	if err != nil {
		return nil, log.NewError(err)
	}

	entry := &Entry{
		Url:        result["url"],
		UserAgent:  result["user_agent"],
		RemoteAddr: result["remote_addr"],
		Timestamp:  timestamp,
	}

	if referrer, exists := result["referrer"]; exists {
		entry.Referrer = referrer
	}
	if entry.Referrer == "-" {
		entry.Referrer = ""
	}

	if err := entry.validate(); err != nil {
		return nil, err
	}

	return entry, nil
}
