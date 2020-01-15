package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/simia-tech/crypt"
)

type verifyFunc func(password, expectedHash string) (bool, error)

func doVerify(password string) error {
	hash := *argVerify

	if len(hash) < 3 || !strings.HasPrefix(hash, "$") || strings.Count(hash, "$") < 2 {
		return fmt.Errorf("invalid crypt value")
	}

	// find second $ to extract prefix
	pos := strings.Index(hash[1:], "$")
	prefix := hash[:pos+2] // $ is contained as count is checked to be at least 2

	f, err := getVerifyHandler(prefix)
	if err != nil {
		return err
	}

	ok, err := f(password, hash)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("the supplied password does not match the given crypt value")
	}
	fmt.Println("password is ok")
	return nil
}

func getVerifyHandler(prefix string) (verifyFunc, error) {
	switch prefix {
	case "$1$":
		return verifyMD5, nil

	case "$2a$":
		return verifyBcrypt2a, nil

	case "$5$":
		return verifySHA256, nil

	case "$6$":
		return verifySHA512, nil

	default:
		return nil, fmt.Errorf("unregognized crypt prefix %q", prefix)
	}
}

func verifyMD5(password, expectedHash string) (bool, error) {
	parts := strings.Split(expectedHash, "$")
	if len(parts) != 4 || len(parts[2]) != 8 {
		return false, fmt.Errorf("invalid md5 crypt value")
	}

	salt := parts[2]
	hash, err := crypt.Crypt(password, fmt.Sprintf("$1$%s", salt))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(hash) == expectedHash, nil
}

func verifyBcrypt2a(password, expectedHash string) (bool, error) {
	parts := strings.Split(expectedHash, "$")
	if len(parts) != 4 || len(parts[3]) != 53 {
		return false, fmt.Errorf("invalid sha256 crypt value")
	}

	rounds := parts[2]
	salt := parts[3][:22]
	hash, err := crypt.Crypt(password, fmt.Sprintf("$2a$%s$%s", rounds, salt))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(hash) == expectedHash, nil
}

func verifySHA256(password, expectedHash string) (bool, error) {
	parts := strings.Split(expectedHash, "$")
	if len(parts) < 4 || len(parts) > 5 {
		return false, fmt.Errorf("invalid sha256 crypt value")
	}

	var settings strings.Builder
	if len(parts) == 5 {
		if strings.HasPrefix(parts[2], "rounds=") {
			settings.WriteString(parts[2])
			settings.WriteString("$")
		} else {
			return false, fmt.Errorf("unknown crypt setting %q", parts[2])
		}
	}

	salt := parts[len(parts)-2]
	hash, err := crypt.Crypt(password, fmt.Sprintf("$5$%s%s", settings.String(), salt))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(hash) == expectedHash, nil
}

func verifySHA512(password, expectedHash string) (bool, error) {
	parts := strings.Split(expectedHash, "$")
	if len(parts) < 4 || len(parts) > 5 {
		return false, fmt.Errorf("invalid sha512 crypt value")
	}

	var settings strings.Builder
	if len(parts) == 5 {
		if strings.HasPrefix(parts[2], "rounds=") {
			settings.WriteString(parts[2])
			settings.WriteString("$")
		} else {
			return false, fmt.Errorf("unknown crypt setting %q", parts[2])
		}
	}

	salt := parts[len(parts)-2]
	hash, err := crypt.Crypt(password, fmt.Sprintf("$6$%s%s", settings.String(), salt))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(hash) == expectedHash, nil
}
