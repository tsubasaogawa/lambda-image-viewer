package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

func TestCloudFrontSigning(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Unix(1700000000, 0)
	signer := &cloudFrontURLSigner{signer: sign.NewURLSigner("key-id", key), domain: "cdn.example", now: func() time.Time { return now }}
	signed, err := signer.Sign("private/画像 space.jpg")
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(signed)
	if err != nil {
		t.Fatal(err)
	}
	if u.Path != "/private/画像 space.jpg" || u.Query().Get("Key-Pair-Id") != "key-id" ||
		u.Query().Get("Expires") != strconv.FormatInt(now.Add(15*time.Minute).Unix(), 10) ||
		u.Query().Get("Signature") == "" {
		t.Fatalf("Unexpected signed URL: %s", signed)
	}
	signature, err := base64.StdEncoding.DecodeString(strings.NewReplacer("-", "+", "_", "=", "~", "/").Replace(u.Query().Get("Signature")))
	if err != nil {
		t.Fatal(err)
	}
	u.RawQuery = ""
	policy, err := json.Marshal(sign.NewCannedPolicy(u.String(), now.Add(privateURLLifetime)))
	if err != nil {
		t.Fatal(err)
	}
	hash := sha1.Sum(policy)
	if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA1, hash[:], signature); err != nil {
		t.Fatalf("Signature verification failed: %v", err)
	}
	tampered, err := json.Marshal(sign.NewCannedPolicy(u.String()+"x", now.Add(privateURLLifetime)))
	if err != nil {
		t.Fatal(err)
	}
	tamperedHash := sha1.Sum(tampered)
	if rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA1, tamperedHash[:], signature) == nil {
		t.Fatal("Tampered URL verified successfully")
	}
	for _, invalid := range []string{"public/image.jpg", "../private/image.jpg", "private/../public.jpg", "private/%2Fimage.jpg", "/private/image.jpg", "private\\image.jpg", "private/image\x00.jpg"} {
		if _, err := signer.Sign(invalid); err == nil {
			t.Errorf("Expected rejection for %q", invalid)
		}
	}
}

func TestPrivateURLsHandler(t *testing.T) {
	mockPrivateSigner(t)
	for _, key := range []string{"public/image.jpg", "../private/image.jpg"} {
		resp, err := privateURLsHandler(key)
		if err != nil || resp.StatusCode != 400 {
			t.Fatalf("Expected rejection for %q, got %+v, %v", key, resp, err)
		}
	}
	resp, err := privateURLsHandler("private/image.jpg")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Failed to sign: %+v, %v", resp, err)
	}
	var urls privateURLs
	if err := json.Unmarshal([]byte(resp.Body), &urls); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(urls.ImageURL, "/private/image.jpg?") ||
		!strings.Contains(urls.ThumbnailURL, "/thumbnail/private/image.jpg?") ||
		resp.Headers["Cache-Control"] != "private, no-store" {
		t.Fatalf("Unexpected response: %+v", resp)
	}
}

func TestPrivateURLsHandlerSignerFailure(t *testing.T) {
	original := newPrivateURLSigner
	newPrivateURLSigner = func() (privateURLSigner, error) { return nil, errors.New("key unavailable") }
	t.Cleanup(func() { newPrivateURLSigner = original })
	resp, err := privateURLsHandler("private/image.jpg")
	if err == nil || resp.StatusCode != 500 || strings.Contains(resp.Body, "key unavailable") {
		t.Fatalf("Expected a safe failure, got %+v, %v", resp, err)
	}
}

func TestPrivateURLsRoute(t *testing.T) {
	mockPrivateSigner(t)
	resp, err := Index(events.LambdaFunctionURLRequest{
		RawPath:        "/cameraroll/sign",
		RawQueryString: "key=private%2Fimage.jpg",
	})
	if err != nil || resp.StatusCode != 200 || !strings.Contains(resp.Body, "thumbnail/private/image.jpg") {
		t.Fatalf("Failed to refresh private URLs: %+v, %v", resp, err)
	}
}

func TestCameraRollHTMLPrivateURLs(t *testing.T) {
	mockPrivateSigner(t)
	t.Setenv("ORIGIN_DOMAIN", "cdn.example")
	t.Setenv("VIEWER_DOMAIN", "viewer.example")
	thumbs := []model.Thumbnail{{Id: "private/a.jpg", Private: true}}
	resp, err := generateCameraRollHtml(&mockDB{thumbsToReturn: &thumbs}, true)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("Failed to render: %+v, %v", resp, err)
	}

	if strings.Contains(resp.Body, "SALT_FOR_PRIVATE_IMAGE") ||
		strings.Contains(resp.Body, "const salt") ||
		!strings.Contains(resp.Body, "thumbnail") ||
		!strings.Contains(resp.Body, "cdn.example") ||
		!strings.Contains(resp.Body, "Expires=9999999999") ||
		resp.Headers["Cache-Control"] != "private, no-store" {
		t.Fatal("Private page has missing signed URLs or insecure cache headers")
	}
	assertScriptValid(t, resp.Body)
	testRefreshScript(t, resp.Body)
}

func testRefreshScript(t *testing.T, body string) {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required to test the rendered JavaScript")
	}
	script := regexp.MustCompile(`(?s)<script>\s*(.*?)</script>`).FindStringSubmatch(body)
	if len(script) != 2 {
		t.Fatal("Inline script not found")
	}
	prefix := `
const assert = require('node:assert/strict');
const element = () => ({ style: {}, classList: {}, addEventListener() {} });
global.document = { getElementById: element, querySelector: element, addEventListener() {} };
global.window = { location: { origin: 'https://viewer.example', search: '' }, addEventListener() {} };
let fetchCount = 0;
let fail = false;
global.fetch = async (url, options) => {
  fetchCount++;
  assert.equal(url.pathname, '/cameraroll/sign');
  assert.equal(url.searchParams.get('key'), 'private/a.jpg');
  assert.equal(options.cache, 'no-store');
  return { ok: !fail, status: 403, json: async () => ({
    image_url: 'https://cdn.example/private/a.jpg?Expires=9999999999',
    thumbnail_url: 'https://cdn.example/thumbnail/private/a.jpg?Expires=9999999999'
  }) };
};
`
	suffix := `
(async () => {
  assert.equal(initialThumbnailsData[0].image_url, 'https://cdn.example/private/a.jpg?Expires=9999999999');
  assert.equal(initialThumbnailsData[0].thumbnail_url, 'https://cdn.example/thumbnail/private/a.jpg?Expires=9999999999');
  const thumb = { id: 'private/a.jpg', image_url: 'https://cdn.example/private/a.jpg?Expires=1' };
  await refreshPrivateURLs(thumb);
  assert.equal(fetchCount, 1);
  assert.match(thumb.thumbnail_url, /thumbnail\/private\/a.jpg/);
  await refreshPrivateURLs(thumb);
  assert.equal(fetchCount, 1);
  fail = true;
  await assert.rejects(refreshPrivateURLs(thumb, true));
})().catch(error => { console.error(error); process.exitCode = 1; });
`
	cmd := exec.Command(node, "-")
	cmd.Stdin = strings.NewReader(prefix + script[1] + suffix)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Private URL refresh failed: %v\n%s", err, output)
	}
}

func TestPrivateImageAndMetadataResponses(t *testing.T) {
	mockPrivateSigner(t)
	t.Setenv("ORIGIN_DOMAIN", "cdn.example")
	original := getMetadata
	getMetadata = func(string) (*model.Metadata, error) { return &model.Metadata{Title: "test"}, nil }
	t.Cleanup(func() { getMetadata = original })

	image, err := generateImageHtml("private/image.jpg")
	if err != nil || image.StatusCode != 200 || image.Headers["Cache-Control"] != "private, no-store" ||
		!strings.Contains(image.Body, "Expires=9999999999") || !strings.Contains(image.Body, "/cameraroll/sign") {
		t.Fatalf("Unexpected private image response: %+v, %v", image, err)
	}
	assertScriptValid(t, image.Body)
	metadata, err := generateMetadataJson("private/image.jpg")
	if err != nil || metadata.StatusCode != 200 || metadata.Headers["Cache-Control"] != "private, no-store" {
		t.Fatalf("Unexpected private metadata response: %+v, %v", metadata, err)
	}
	public, err := generateImageHtml("public/image.jpg")
	if err != nil || public.Headers["Cache-Control"] != "public, max-age=31536000" ||
		strings.Contains(public.Body, "/cameraroll/sign") {
		t.Fatalf("Unexpected public image response: %+v, %v", public, err)
	}
}

func assertScriptValid(t *testing.T, body string) {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required to validate the rendered JavaScript")
	}
	script := regexp.MustCompile(`(?s)<script>\s*(.*?)</script>`).FindStringSubmatch(body)
	if len(script) != 2 {
		t.Fatal("Inline script not found")
	}
	cmd := exec.Command(node, "--check")
	cmd.Stdin = bytes.NewBufferString(script[1])
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Rendered JavaScript is invalid: %v\n%s", err, output)
	}
}

func TestSignerCaching(t *testing.T) {
	resetSignerCache()
	defer resetSignerCache()

	origLoader := loadSignerFunc
	defer func() { loadSignerFunc = origLoader }()

	callCount := 0
	loadSignerFunc = func() (privateURLSigner, error) {
		callCount++
		return mockSigner{}, nil
	}

	s1, err := getCachedPrivateURLSigner()
	if err != nil || s1 == nil {
		t.Fatalf("Expected signer, got %v, error: %v", s1, err)
	}
	if callCount != 1 {
		t.Fatalf("Expected callCount=1, got %d", callCount)
	}

	for i := 0; i < 5; i++ {
		s, err := getCachedPrivateURLSigner()
		if err != nil || s != s1 {
			t.Fatalf("Expected cached signer %v, got %v, error: %v", s1, s, err)
		}
	}
	if callCount != 1 {
		t.Fatalf("Expected callCount to remain 1, got %d", callCount)
	}

	resetSignerCache()

	s2, err := getCachedPrivateURLSigner()
	if err != nil || s2 == nil {
		t.Fatalf("Expected signer after reset, got %v, error: %v", s2, err)
	}
	if callCount != 2 {
		t.Fatalf("Expected callCount=2 after cache reset, got %d", callCount)
	}
}

func TestSignerCachingErrorNotCached(t *testing.T) {
	resetSignerCache()
	defer resetSignerCache()

	origLoader := loadSignerFunc
	defer func() { loadSignerFunc = origLoader }()

	fail := true
	callCount := 0
	loadSignerFunc = func() (privateURLSigner, error) {
		callCount++
		if fail {
			return nil, errors.New("temporary error")
		}
		return mockSigner{}, nil
	}

	s, err := getCachedPrivateURLSigner()
	if err == nil || s != nil {
		t.Fatalf("Expected error, got %v", s)
	}
	if callCount != 1 {
		t.Fatalf("Expected callCount=1, got %d", callCount)
	}

	fail = false
	s, err = getCachedPrivateURLSigner()
	if err != nil || s == nil {
		t.Fatalf("Expected success on retry, got error: %v", err)
	}
	if callCount != 2 {
		t.Fatalf("Expected callCount=2 on retry, got %d", callCount)
	}
}
