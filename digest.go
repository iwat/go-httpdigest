// go-httpdigest - A simple Go(lang) drop-in replacement for "net/http/Client"
// with Digest Auth capability.

// Copyright (c) 2015 Chaiwat Shuetrakoonpaiboon. All rights reserved.
//
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

// Package httpdigest provides a simple Go(lang) drop-in replacement for
// "net/http/Client" with Digest Auth capability.
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
	Nc        int16
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

func (d DigestChallenge) ApplyAuth(user, pass string, req *http.Request) {
	d.Nc += 0x1
	d.Cnonce = randomNonce()
	d.Method = req.Method
	path := req.URL.RequestURI()

	AuthHeader := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc=%08x, qop=%s, response="%s", algorithm=%s`,
		user, d.Realm, d.Nonce, path, d.Cnonce, d.Nc, d.Qop, d.calculateResponse(user, pass, path), d.Algorithm)

	if d.Opaque != "" {
		AuthHeader = fmt.Sprintf(`%s, opaque="%s"`, AuthHeader, d.Opaque)
	}

	req.Header.Set("Authorization", AuthHeader)
}

func (d DigestChallenge) calculateResponse(user, pass, path string) string {
	ha1, ha2 := d.calculateChecksum(user, pass, path)

	message := strings.Join([]string{ha1, d.Nonce, fmt.Sprintf("%08x", d.Nc), d.Cnonce, d.Qop, ha2}, ":")
	return fmt.Sprintf("%x", md5.Sum([]byte(message)))
}

func (d DigestChallenge) calculateChecksum(user, pass, path string) (string, string) {
	switch d.Algorithm {
	case "MD5":
		a1 := fmt.Sprintf("%s:%s:%s", user, d.Realm, pass)
		ha1 := fmt.Sprintf("%x", md5.Sum([]byte(a1)))

		a2 := fmt.Sprintf("%s:%s", d.Method, path)
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
