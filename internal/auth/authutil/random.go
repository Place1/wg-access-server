package authutil

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/sirupsen/logrus"
)

func RandomString(size int) string {
	blk := make([]byte, size)
	_, err := rand.Read(blk)
	if err != nil {
		logrus.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(blk)
}
