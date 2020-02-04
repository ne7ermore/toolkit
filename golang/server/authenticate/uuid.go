package authenticate

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/gofrs/uuid"
)

func UUID() (string, error) {
	tempUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return Md5([]byte(tempUUID.String())), nil
}

func Md5(bs []byte) string {
	m := md5.New()
	m.Write(bs)
	return hex.EncodeToString(m.Sum(nil))
}
