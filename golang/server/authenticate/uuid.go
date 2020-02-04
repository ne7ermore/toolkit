package authenticate

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/gofrs/uuid"
)

func UUID() (string, error) {
	tempUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	md5 := Md5([]byte(tempUUID.String()))
	upper := strings.ToUpper(md5)
	return upper, nil
}

func Md5(bs []byte) string {
	m := md5.New()
	m.Write(bs)
	return hex.EncodeToString(m.Sum(nil))
}
