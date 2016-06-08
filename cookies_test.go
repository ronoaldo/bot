package bot

import (
	"testing"
)

func TestEncodeDecodeCookies(t *testing.T) {
	b := New()
	_, err := b.GET("http://www.microsoft.com/")
	if err != nil {
		t.Errorf("Failed to make HTTP GET: %v", err)
	}

	cookies, err := b.EncodeCookies()
	if err != nil {
		t.Errorf("Unable to encode cookies: %v", err)
	}
	t.Logf("Encoded values: %s", string(cookies))

	if len(cookies) > 0 {
		b2 := New()
		if err = b2.DecodeCookies(cookies); err != nil {
			t.Errorf("Unable to decode cookies: %v", err)
		}

		c2, err := b2.EncodeCookies()
		if err != nil {
			t.Errorf("Error reencoding cookies: %v", err)
		}
		t.Logf("Reencoded values: %s", string(c2))

		if len(c2) == 0 {
			t.Errorf("Reencoded cookies are empty!")
		}
	}
}
