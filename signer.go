package ekasa

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	emptyBodyHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	algorithm     = "NWS4-HMAC-SHA256"
)

type signer struct {
	publicKey    string
	privateKey   string
	timeProvider func() time.Time
}

func newSigner(publicKey, privateKey string) *signer {
	return &signer{
		publicKey:    publicKey,
		privateKey:   privateKey,
		timeProvider: func() time.Time { return time.Now().UTC() },
	}
}

func (s *signer) sign(method, rawURL string, headers map[string]string, body []byte) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	now := s.timeProvider()
	longDate := now.Format("2006-01-02T15:04:05.000") + "Z"
	shortDate := now.Format("2006-01-02")

	var bodyHash string
	if len(body) == 0 {
		bodyHash = emptyBodyHash
	} else {
		h := sha256.New()
		h.Write(body)
		bodyHash = hex.EncodeToString(h.Sum(nil))
	}

	// Set required headers for signing
	headers["x-nd-date"] = longDate
	headers["x-nd-content-sha256"] = bodyHash
	headers["host"] = u.Host

	// 1. Canonical Resource Path
	resourcePath := u.Path
	if resourcePath == "" {
		resourcePath = "/"
	}
	// Path should be encoded but / preserved
	parts := strings.Split(resourcePath, "/")
	for i, p := range parts {
		parts[i] = customRawUrlEncode(p)
	}
	canonicalResourcePath := strings.Join(parts, "/")

	// 2. Canonical Query Parameters
	query := u.Query()
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var queryParts []string
	for _, k := range keys {
		vals := query[k]
		sort.Strings(vals)
		for _, v := range vals {
			queryParts = append(queryParts, customRawUrlEncode(k)+"="+customRawUrlEncode(v))
		}
	}
	canonicalQueryString := strings.Join(queryParts, "&")

	// 3. Canonical Headers and Signed Header Names
	// PHP sorts header names for "SignedHeaders" using strtoupper comparison
	signedHeaderNames := make([]string, 0, len(headers))
	for k := range headers {
		signedHeaderNames = append(signedHeaderNames, k)
	}
	sort.Slice(signedHeaderNames, func(i, j int) bool {
		return strings.ToUpper(signedHeaderNames[i]) < strings.ToUpper(signedHeaderNames[j])
	})

	signedHeadersParts := make([]string, 0, len(signedHeaderNames))
	for _, k := range signedHeaderNames {
		signedHeadersParts = append(signedHeadersParts, strings.ToLower(k))
	}
	signedHeaders := strings.Join(signedHeadersParts, ";")

	// PHP sorts the "canonicalHeaders" (key:value pairs) using lowercase lexical order
	canonicalHeaderKeys := make([]string, 0, len(headers))
	for k := range headers {
		canonicalHeaderKeys = append(canonicalHeaderKeys, strings.ToLower(k))
	}
	sort.Strings(canonicalHeaderKeys)

	var canonicalHeadersParts []string
	for _, lk := range canonicalHeaderKeys {
		// Find original value
		var originalValue string
		for ok, ov := range headers {
			if strings.ToLower(ok) == lk {
				originalValue = ov
				break
			}
		}
		canonicalHeadersParts = append(canonicalHeadersParts, lk+":"+strings.TrimSpace(originalValue))
	}
	canonicalHeaders := strings.Join(canonicalHeadersParts, "\n")

	// 4. Canonical Request
	canonicalRequestParts := []string{
		method,
		canonicalResourcePath,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		bodyHash,
	}
	canonicalRequest := strings.Join(canonicalRequestParts, "\n")

	hReq := sha256.New()
	hReq.Write([]byte(canonicalRequest))
	canonicalRequestHash := hex.EncodeToString(hReq.Sum(nil))

	// 5. String to Sign
	stringToSignParts := []string{
		algorithm,
		longDate,
		s.publicKey,
		canonicalRequestHash,
	}
	stringToSign := strings.Join(stringToSignParts, "\n")

	// 6. Compute Signature
	signingKey := hmacSHA256([]byte("NWS4"+s.privateKey), []byte(shortDate))
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	// 7. Authorization Header
	authHeaderParts := []string{
		"Credential=" + s.publicKey,
		"SignedHeaders=" + url.QueryEscape(signedHeaders),
		"Timestamp=" + url.QueryEscape(longDate),
		"Signature=" + signature,
	}

	// PHP uses rawurlencode on the whole comma-separated string then prepends algorithm
	authValue := strings.Join(authHeaderParts, ",")
	// rawurlencode in PHP for , is %2C, for = is %3D
	// url.QueryEscape might be different, let's be careful.
	// PHP rawurlencode is RFC 3986. url.PathEscape or custom.

	// Replicating PHP rawurlencode behavior
	authHeader := algorithm + " " + customRawUrlEncode(authValue)
	headers["Authorization"] = authHeader

	return nil
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func customRawUrlEncode(str string) string {
	// PHP rawurlencode encodes everything except - . _ ~
	var result strings.Builder
	for i := 0; i < len(str); i++ {
		c := str[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '.' || c == '_' || c == '~' {
			result.WriteByte(c)
		} else {
			fmt.Fprintf(&result, "%%%02X", c)
		}
	}
	return result.String()
}
