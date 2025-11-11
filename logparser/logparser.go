package logparser

import (
	"time"

	"github.com/go-playground/validator/v10"

	logger "github.com/rtfmkiesel/kisslog"
)

var (
	validate = validator.New()
	log      = logger.New("logparser")
)

type Entry struct {
	RemoteAddr string    `validate:"required"`
	Timestamp  time.Time `validate:"required"`
	Url        string    `validate:"required"`
	Referrer   string    // Optional
	UserAgent  string    `validate:"required"`
}

func (entry *Entry) validate() error {
	return validate.Struct(entry)
}

// Log parsers must target this interface
type Parser interface {
	ParseLine(line string) (*Entry, error)
}
