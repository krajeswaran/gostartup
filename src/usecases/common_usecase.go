package usecases

import (
	"github.com/franela/goreq"
	"gostartup/src/models"
	"math/big"
	"crypto/rand"
	"net/url"
	"crypto/sha256"
	"encoding/base64"
)

type CommonUsecase struct {}

func (c *CommonUsecase) CreateApiSignature(req *goreq.Request, acct *models.Account) {
	// TODO use jwt token
	sign, nonce := signature(req.Uri, acct.AuthToken)
	req.AddHeader(HEADER_SIGNATURE, sign)
	req.AddHeader(HEADER_SIGNATURE_NONCE, nonce)
}

func nonce() string {
	//Max random value, a 64-bits integer, i.e 2^64 - 1
	max := new(big.Int)
	max, _ = max.SetString("18446744073709551616", 10)

	//Generate cryptographically strong pseudo-random between 0 - max
	n, _ := rand.Int(rand.Reader, max)

	return n.String()
}

// baseurl + authtoken + nonce
func signature(callbackUrl string, authToken string) (string, string) {
	u, _ := url.Parse(callbackUrl)

	h := sha256.New()
	h.Write([]byte(u.Scheme + u.Host + u.Port() + u.EscapedPath()))
	h.Write([]byte(authToken))
	nonce := nonce()
	h.Write([]byte(nonce))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature, nonce
}
