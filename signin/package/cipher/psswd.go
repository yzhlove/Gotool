package cipher

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	errHashPasswdEmpty   = errors.New("password is empty")
	errHashPassedInvalid = errors.New("password is incorrect")
	errArgonVersion      = errors.New("argon2 version not supported")
)

var argonFields = []string{"argon2id", "b64hash", "b64salt"}

const (
	memorySize  = 16 * 1024
	iterations  = 3
	parallelism = 4
	saltLength  = 16
	keyLength   = 32
)

func GenerateHashPasswd(passwd string) (string, error) {

	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(passwd), salt,
		iterations,
		memorySize,
		parallelism,
		keyLength)

	var encodeHash strings.Builder
	encodeHash.WriteString(fmt.Sprintf("$argon2id=%d", argon2.Version))
	encodeHash.WriteString(fmt.Sprintf("$b64salt=%s", base64.StdEncoding.EncodeToString(salt)))
	encodeHash.WriteString(fmt.Sprintf("$b64hash=%s", base64.StdEncoding.EncodeToString(hash)))
	return encodeHash.String(), nil
}

func CompareHashPasswd(passwd string, encodedHash string) (bool, error) {

	if len(encodedHash) == 0 {
		return false, errHashPasswdEmpty
	}

	parts := strings.Split(encodedHash, "$")
	if len(parts) != 4 {
		return false, errHashPassedInvalid
	}

	fields := make(map[string]string, len(argonFields))
	for _, value := range parts {
		if parse := strings.SplitN(value, "=", 2); len(parse) == 2 {
			fields[parse[0]] = parse[1]
		}
	}

	tags := slices.Collect(maps.Keys(fields))
	slices.Sort(tags)
	if !slices.Equal(argonFields, tags) {
		return false, errHashPassedInvalid
	}

	if strconv.FormatInt(argon2.Version, 10) != fields[argonFields[0]] {
		return false, errArgonVersion
	}

	decodeHash, err := base64.StdEncoding.DecodeString(fields[argonFields[1]])
	if err != nil {
		return false, err
	}

	decodeSalt, err := base64.StdEncoding.DecodeString(fields[argonFields[2]])
	if err != nil {
		return false, err
	}

	compareHash := argon2.IDKey([]byte(passwd), decodeSalt,
		iterations,
		memorySize,
		parallelism,
		uint32(len(decodeHash)))
	return hmac.Equal(decodeHash, compareHash), nil
}
