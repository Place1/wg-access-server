package authutil

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func RandomString(size int) string {
	blk := make([]byte, size)
	_, err := rand.Read(blk)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to make a random string"))
	}
	return base64.StdEncoding.EncodeToString(blk)
}
