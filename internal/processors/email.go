package processors

import (
	"errors"
	"fmt"
	"time"
)

type EmailJob struct {
	To      string
	Subject string
}

func (e EmailJob) Type() string {
	return "email"
}

func (e EmailJob) Process() (string, error) {
	time.Sleep(2 * time.Second)

	if e.To == "" {
		return "", errors.New("missing email recipient")
	}

	return fmt.Sprintf("email sent to %s with subject %q", e.To, e.Subject), nil
}
