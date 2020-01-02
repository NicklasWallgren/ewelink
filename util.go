package ewelink

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"math/rand"
)

func generateNonce() string {
	return "1"
}

func getRandomNumber(n1 int, n2 int) int {
	return rand.Intn(n2-n1) + n1
}

func calculateHash(subject []byte) string {
	mac := hmac.New(sha256.New, []byte("6Nz4n0xA8s8qdxQf2GqurZj2Fs55FUvM"))
	mac.Write(subject)

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func readerToString(reader io.Reader) string {
	buf := new(bytes.Buffer)

	buf.ReadFrom(reader)

	return buf.String()
}
