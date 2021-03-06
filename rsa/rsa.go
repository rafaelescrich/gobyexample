package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
)

const maxEncodeLength = 117

var pubkey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCsYmysgY/RHSUMkfXk2Tt/g9sv
JssYzBGD9YjCddSCZbVSTEZX9zcC9eRrhLWx1zO/wvnkGIzipe3qakasmv3wECPw
bJf0bHiY429Z2tH65s+LZWjSGoxL7S4uNO+hAD//aiKYPJnhfjnbtxKnfJkcEdxG
B4/44oI4vC4xn00/zwIDAQAB
-----END PUBLIC KEY-----`)

func getPubKey() (pub *rsa.PublicKey, err error) {
	var (
		block        *pem.Block
		pubInterface interface{}
	)

	block, _ = pem.Decode(pubkey)
	if block == nil {
		err = errors.New("public key invalid")
		return
	}

	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return
	}

	pub = pubInterface.(*rsa.PublicKey)

	return
}

func Encrypt(data *map[string]interface{}) (encrypted string, err error) {
	var (
		jsonByte []byte
		encrypt  []byte
		pub      *rsa.PublicKey
		sliceLen = maxEncodeLength
	)

	pub, err = getPubKey()
	if err != nil {
		return
	}

	jsonByte, err = json.Marshal(data)
	if err != nil {
		return
	}

	jsonByteLen := len(jsonByte)

	encrypts := make([][]byte, jsonByteLen/maxEncodeLength+1)

	for i := 0; i < jsonByteLen; i = i + sliceLen {
		length := jsonByteLen - i
		if length < maxEncodeLength {
			sliceLen = length
		}
		encrypt, err = rsa.EncryptPKCS1v15(rand.Reader, pub, jsonByte[i:i+sliceLen])
		if err != nil {
			return
		}

		encrypts = append(encrypts, encrypt)
	}

	encrypted = base64.StdEncoding.EncodeToString(bytes.Join(encrypts, []byte("")))

	return
}
