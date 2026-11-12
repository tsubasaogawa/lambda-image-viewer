package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/tsubasaogawa/lambda-image-viewer/src/viewer/internal/model"
)

const privateURLLifetime = 15 * time.Minute

type privateURLSigner interface {
	Sign(key string) (string, error)
}

type cloudFrontURLSigner struct {
	signer *sign.URLSigner
	domain string
	now    func() time.Time
}

var (
	signerCacheLock sync.RWMutex
	cachedSigner    privateURLSigner
	loadSignerFunc  = loadPrivateURLSigner

	newPrivateURLSigner = getCachedPrivateURLSigner
)

func getCachedPrivateURLSigner() (privateURLSigner, error) {
	signerCacheLock.RLock()
	if cachedSigner != nil {
		s := cachedSigner
		signerCacheLock.RUnlock()
		return s, nil
	}
	signerCacheLock.RUnlock()

	signerCacheLock.Lock()
	defer signerCacheLock.Unlock()

	if cachedSigner != nil {
		return cachedSigner, nil
	}

	signer, err := loadSignerFunc()
	if err != nil {
		return nil, err
	}
	cachedSigner = signer
	return cachedSigner, nil
}

func resetSignerCache() {
	signerCacheLock.Lock()
	cachedSigner = nil
	signerCacheLock.Unlock()
}

func loadPrivateURLSigner() (privateURLSigner, error) {
	secretARN := os.Getenv("PRIVATE_SIGNING_KEY_SECRET_ARN")
	keyID := os.Getenv("CLOUDFRONT_SIGNING_KEY_ID")
	domain := os.Getenv("ORIGIN_DOMAIN")
	if secretARN == "" || keyID == "" || domain == "" {
		return nil, errors.New("private URL signing is not configured")
	}
	sess, err := session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("create AWS session: %w", err)
	}
	secret, err := secretsmanager.New(sess).GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretARN),
	})
	if err != nil {
		return nil, fmt.Errorf("retrieve signing key: %w", err)
	}
	block, _ := pem.Decode([]byte(aws.StringValue(secret.SecretString)))
	if block == nil {
		return nil, errors.New("signing key is not PEM encoded")
	}
	var key *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		var parsed interface{}
		parsed, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err == nil {
			var ok bool
			key, ok = parsed.(*rsa.PrivateKey)
			if !ok {
				err = errors.New("signing key is not RSA")
			}
		}
	default:
		err = errors.New("unsupported signing key format")
	}
	if err != nil {
		return nil, fmt.Errorf("parse signing key: %w", err)
	}
	if key.N.BitLen() != 2048 {
		return nil, errors.New("signing key must be RSA 2048")
	}
	return &cloudFrontURLSigner{signer: sign.NewURLSigner(keyID, key), domain: domain, now: time.Now}, nil
}

func validPrivateKey(key string) bool {
	if !model.IsPrivateKey(key) || strings.ContainsAny(key, "\\%") {
		return false
	}
	for _, r := range key {
		if unicode.IsControl(r) {
			return false
		}
	}
	for _, part := range strings.Split(key, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}

func (s *cloudFrontURLSigner) Sign(key string) (string, error) {
	if !validPrivateKey(key) {
		return "", errors.New("invalid private image key")
	}
	u := url.URL{Scheme: "https", Host: s.domain, Path: "/" + key}
	return s.signer.Sign(u.String(), s.now().Add(privateURLLifetime))
}

type privateURLs struct {
	ImageURL     string `json:"image_url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

func signPrivateURLs(signer privateURLSigner, key string) (privateURLs, error) {
	if !validPrivateKey(key) {
		return privateURLs{}, errors.New("invalid private image key")
	}
	image, err := signer.Sign(key)
	if err != nil {
		return privateURLs{}, err
	}
	thumbnail, err := signer.Sign("thumbnail/" + key)
	if err != nil {
		return privateURLs{}, err
	}
	return privateURLs{ImageURL: image, ThumbnailURL: thumbnail}, nil
}
