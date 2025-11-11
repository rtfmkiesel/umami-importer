package logparser

import (
	"regexp"
	"time"
)

var (
	// Regex for the default apache log profile
	reApacheLogLine = regexp.MustCompile(`^(\S+) (\S+) (\S+) \[([^\]]+)\] "(\S+) ([^"]*) ([^"]*)" (\d+) (\S+) "([^"]*)" "([^"]*)"`)
)

// Parses the default Apache2 logs
type ParserApache2 struct{}

func (p *ParserApache2) ParseLine(line string) (*Entry, error) {
	matches := reApacheLogLine.FindStringSubmatch(line)

	if len(matches) != 12 {
		return nil, log.NewError("failed to parse log line (matches=%d,wanted=12)", len(matches))
	}

	entry := &Entry{
		RemoteAddr: matches[1],
		Url:        matches[6],
		Referrer:   matches[10],
		UserAgent:  matches[11],
	}

	timeStr := matches[4]
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
