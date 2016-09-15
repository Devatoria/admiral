package token

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/base32"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

// ClaimsAccess represents an access to a repository
type ClaimsAccess struct {
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Actions []string `json:"actions"`
}

// Claims represents a custom JWT claims with a list of accesses
type Claims struct {
	Access []ClaimsAccess `json:"access"`
	jwt.StandardClaims
}

// NewToken returns a new token based on given claims
func NewToken(service, subject string, accesses []ClaimsAccess) *jwt.Token {
	claims := Claims{
		accesses,
		jwt.StandardClaims{
			Issuer:    viper.GetString("auth.issuer"),
			Subject:   subject,
			Audience:  service,
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.token-expiration")) * time.Minute).Unix(),
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        uuid.NewV4().String(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
}

// SignToken signs the given token with the registry certificate and returns the signed token as string
func SignToken(token *jwt.Token) (string, error) {
	// Read certificate
	certBytes, err := ioutil.ReadFile(viper.GetString("auth.certificate"))
	if err != nil {
		return "", err
	}

	certPem, _ := pem.Decode(certBytes)
	if certPem == nil {
		return "", errors.New("Failed to parse certificate")
	}

	cert, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		return "", err
	}

	// Compute libtrust fingerprint
	derBytes, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return "", err
	}

	hasher := crypto.SHA256.New()
	_, err = hasher.Write(derBytes)
	if err != nil {
		return "", err
	}

	token.Header["kid"], err = keyIDEncode(hasher.Sum(nil)[:30])
	if err != nil {
		return "", err
	}

	// Read certificate private key
	keyBytes, err := ioutil.ReadFile(viper.GetString("auth.private-key"))
	if err != nil {
		return "", err
	}

	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil {
		return "", errors.New("Failed to parse private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func keyIDEncode(b []byte) (string, error) {
	s := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")
	var buf bytes.Buffer
	var i int
	var err error
	for i = 0; i < len(s)/4-1; i++ {
		start := i * 4
		end := start + 4
		_, err = buf.WriteString(s[start:end] + ":")
		if err != nil {
			return "", err
		}
	}

	_, err = buf.WriteString(s[i*4:])
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
