package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"

	"github.com/simia-tech/crypt/bcrypt"

	"github.com/simia-tech/crypt"
)

/* relevant methods:
   $1$: MD5-based crypt ('md5crypt')
   $2$: Blowfish-based crypt ('bcrypt')
   $2a$: enforce UTF-8 and string null-terminator (default for go)
   $2x$ / $2y$: NOT relevant
   $2b$: fixed unsigned char bug
   $sha1$: SHA-1-based crypt ('sha1crypt')
   $5$: SHA-256-based crypt ('sha256crypt')
*/

type cryptFunc func(password string) (string, error)

func doCrypt(password string) error {
	f, err := getCryptHandler(*argMethod)
	if err != nil {
		return err
	}

	hash, err := f(password)
	if err != nil {
		return err
	}

	fmt.Println(hash)
	return nil
}

func getOrGenSalt(minLen, maxLen int) string {
	if len(*argSalt) > 0 {
		if len(*argSalt) < minLen {
			return fmt.Sprintf("%s%s", *argSalt, genSalt(minLen-len(*argSalt)))
		} else if len(*argSalt) > maxLen {
			return (*argSalt)[:maxLen]
		}
		return *argSalt
	}
	return genSalt(maxLen)
}

func genSalt(len int) string {
	buffer := make([]byte, len*3/4)
	rand.Read(buffer)
	return string(bcrypt.Base64Encode(buffer))[:len]
}

func getCryptHandler(method string) (cryptFunc, error) {
	switch strings.ToLower(method) {
	case "": // default method
		fallthrough
	case "default":
		fallthrough
	case "$2a$":
		fallthrough
	case "2a":
		return cryptBcrypt2a, nil

	case "256":
		fallthrough
	case "$5$":
		fallthrough
	case "5":
		fallthrough
	case "sha256":
		return cryptSHA256, nil

	case "512":
		fallthrough
	case "$6$":
		fallthrough
	case "6":
		fallthrough
	case "sha512":
		return cryptSHA512, nil

	default:
		return nil, fmt.Errorf("Unknown hash function %q", method)
	}
}

func cryptBcrypt2a(password string) (string, error) {
	hash, err := crypt.Crypt(password, fmt.Sprintf("$2a$10$%s", getOrGenSalt(22, 22)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(hash), nil
}

func cryptSHA256(password string) (string, error) {
	hash, err := crypt.Crypt(password, fmt.Sprintf("$5$%s", getOrGenSalt(0, 12)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(hash), nil
}

func cryptSHA512(password string) (string, error) {
	hash, err := crypt.Crypt(password, fmt.Sprintf("$6$%s", getOrGenSalt(0, 12)))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(hash), nil
}
