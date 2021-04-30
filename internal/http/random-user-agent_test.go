package http_test

import (
	"regexp"
	"testing"

	"github.com/slysterous/scrapmon/internal/http"
)

func TestGenerateRandomUserAgent(t *testing.T) {
	for i := 0; i < 10; i++ {
		ua := http.GenerateRandomUserAgent()
		r, _ := regexp.Compile(".+?[/\\s][\\d.]+")
		isUserAgent := r.MatchString(ua)
		if !isUserAgent {
			t.Errorf("expected valid user agent, got: %s", ua)
		}
	}
}
