package api

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/base32"
	"encoding/pem"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Devatoria/admiral/auth"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

// getToken returns a JWT bearer token to the registry containing the user accesses
func getToken(c *gin.Context) {
	service := c.Query("service")
	user, err := auth.GetCurrentUser(c)
	if err != nil {
		panic(err)
	}

	// Scope is empty only for authentication
	var claimsAccesses []ClaimsAccess
	scope := c.Query("scope")
	if scope != "" {
		//TODO: retrieve rights
	}

	// Create bearer token
	claims := Claims{
		claimsAccesses,
		jwt.StandardClaims{
			Issuer:    viper.GetString("auth.issuer"),
			Subject:   user.Username,
			Audience:  service,
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.token-expiration")) * time.Minute).Unix(),
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			Id:        uuid.NewV4().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Read certificate
	certBytes, err := ioutil.ReadFile(viper.GetString("auth.certificate"))
	if err != nil {
		panic(err)
	}

	certPem, _ := pem.Decode(certBytes)
	if certPem == nil {
		panic("Failed to parse certificate")
	}

	cert, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		panic(err)
	}

	// Compute libtrust fingerprint
	derBytes, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		panic(err)
	}

	hasher := crypto.SHA256.New()
	_, err = hasher.Write(derBytes)
	if err != nil {
		panic(err)
	}

	token.Header["kid"] = keyIDEncode(hasher.Sum(nil)[:30])

	// Read certificate private key
	keyBytes, err := ioutil.ReadFile(viper.GetString("auth.private-key"))
	if err != nil {
		panic(err)
	}

	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil {
		panic("Failed to parse private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	tokenString, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func keyIDEncode(b []byte) string {
	s := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")
	var buf bytes.Buffer
	var i int
	var err error
	for i = 0; i < len(s)/4-1; i++ {
		start := i * 4
		end := start + 4
		_, err = buf.WriteString(s[start:end] + ":")
		if err != nil {
			panic(err)
		}
	}

	_, err = buf.WriteString(s[i*4:])
	if err != nil {
		panic(err)
	}

	return buf.String()
}
