package httpdigest

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type DigestChallenge struct {
	Realm     string
	Qop       string
	Method    string
	Nonce     string
	Opaque    string
	Algorithm string
	Cnonce    string
	Path      string
	Nc        int16
	username  string
	password  string
}

func ChallengeFromResponse(resp *http.Response) DigestChallenge {
	params := parseWWWAuth(resp)

	d := DigestChallenge{}
	d.Realm = params["realm"]
	d.Qop = params["qop"]
	d.Nonce = params["nonce"]
	d.Opaque = params["opaque"]
	d.Algorithm = "MD5"

	if params["algorithm"] != "" {
		d.Algorithm = params["algorithm"]
	}

	d.Nc = 0x0

	return d
}

func (d DigestChallenge) ApplyAuth(req *http.Request) {
	d.Nc += 0x1
	d.Cnonce = randomNonce()
	d.Method = req.Method
	d.Path = req.URL.RequestURI()

	AuthHeader := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc=%08x, qop=%s, response="%s", algorithm=%s`,
		d.username, d.Realm, d.Nonce, d.Path, d.Cnonce, d.Nc, d.Qop, d.calculateResponse(), d.Algorithm)

	if d.Opaque != "" {
		AuthHeader = fmt.Sprintf(`%s, opaque="%s"`, AuthHeader, d.Opaque)
	}

	req.Header.Set("Authorization", AuthHeader)
}

func (d DigestChallenge) calculateResponse() string {
	ha1, ha2 := d.calculateChecksum()

	message := strings.Join([]string{ha1, d.Nonce, fmt.Sprintf("%08x", d.Nc), d.Cnonce, d.Qop, ha2}, ":")
	return fmt.Sprintf("%x", md5.Sum([]byte(message)))
}

func (d DigestChallenge) calculateChecksum() (string, string) {
	switch d.Algorithm {
	case "MD5":
		a1 := fmt.Sprintf("%s:%s:%s", d.username, d.Realm, d.password)
		ha1 := fmt.Sprintf("%x", md5.Sum([]byte(a1)))

		a2 := fmt.Sprintf("%s:%s", d.Method, d.Path)
		ha2 := fmt.Sprintf("%x", md5.Sum([]byte(a2)))
		return ha1, ha2

	case "MD5-sess":
	default:
		//token
	}

	return "", ""
}

func parseWWWAuth(r *http.Response) map[string]string {
	s := strings.SplitN(r.Header.Get("Www-Authenticate"), " ", 2)

	if len(s) != 2 || s[0] != "Digest" {
		return nil
	}

	result := map[string]string{}

	for _, kv := range strings.Split(s[1], ",") {
		parts := strings.SplitN(kv, "=", 2)

		if len(parts) != 2 {
			continue
		}

		result[strings.Trim(parts[0], "\" ")] = strings.Trim(parts[1], "\" ")
	}

	return result
}

func randomNonce() string {
	k := make([]byte, 12)

	for bytes := 0; bytes < len(k); {
		n, err := rand.Read(k[bytes:])

		if err != nil {
			panic("rand.Read() failed")
		}

		bytes += n
	}

	return base64.StdEncoding.EncodeToString(k)
}
