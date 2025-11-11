package umami

import (
	"encoding/json"
	"time"
)

type Event struct {
	// Default event params

	WebsiteId string `json:"website"`
	Hostname  string `json:"hostname"`
	Referrer  string `json:"referrer"`
	Url       string `json:"url"`

	// The following values can be overriden by specifying them in the payload
	// as seen in in umami/src/app/api/send/route.ts

	UserAgent  string    `json:"useragent"`
	RemoteAddr string    `json:"ip"`
	TimeStamp  time.Time `json:"timestamp"`

	// Hash required to save event as processed inside the DB
	Hash string `json:"-"`
}

// Internal request for sending umami events as detailed in the Umami API
// https://umami.is/docs/api/sending-stats
type clientUmamiRequest struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

func (e *Event) ToUmamiRequestBodyJson() ([]byte, error) {
	payload := map[string]any{
		"website":   e.WebsiteId,
		"hostname":  e.Hostname,
		"referrer":  e.Referrer,
		"url":       e.Url,
		"userAgent": e.UserAgent,
		"ip":        e.RemoteAddr,
		"timestamp": e.TimeStamp.Unix(),
	}

	req := clientUmamiRequest{
		Type:    "event",
		Payload: payload,
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return jsonBody, nil
}
