package logparser

import (
	"regexp"
	"time"
)

var (
	// Regex for the default nginx log profile
	reNginxLogLine = regexp.MustCompile(`^(\S+) - (\S+) \[([^\]]+)\] "(\S+) ([^"]*) ([^"]*)" (\d+) (\S+) "([^"]*)" "([^"]*)"`)
)

// Parses the default Nginx logs
type ParserNginx struct{}

func (p *ParserNginx) ParseLine(line string) (*Entry, error) {
	matches := reNginxLogLine.FindStringSubmatch(line)

	if len(matches) != 11 {
		return nil, log.NewError("failed to parse log line (matches=%d,wanted=11)", len(matches))
	}

	entry := &Entry{
		RemoteAddr: matches[1],
		Url:        matches[5],
		Referrer:   matches[9],
		UserAgent:  matches[10],
	}

	timeStr := matches[3]
	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", timeStr)
	if err != nil {
		return nil, log.NewError(err)
	}
	entry.Timestamp = timestamp

	if entry.Referrer == "-" {
		entry.Referrer = ""
	}

	if err := entry.validate(); err != nil {
		return nil, err
	}

	return entry, nil
}
