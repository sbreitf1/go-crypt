package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/sbreitf1/go-console"
)

var (
	app       = kingpin.New("crypt", "Util application to crypt and check passwords")
	argMethod = kingpin.Flag("method", "Bcrypt version (2a, sha256, sha512)").Short('m').String()
	argInput  = kingpin.Flag("input", "Use input parameter instead of StdIn").Short('i').String()
	argPrompt = kingpin.Flag("prompt", "Prompt user to enter password").Short('p').Bool()
	argSalt   = kingpin.Flag("salt", "Use specific salt instead of random").Short('s').String()
	argVerify = kingpin.Flag("verify", "Return wether the input password matches the given hash").Short('v').String()
)

func main() {
	kingpin.Parse()

	if err := appMain(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func appMain() error {
	password, err := getPassword()
	if err != nil {
		return fmt.Errorf("failed to read password: %s", err.Error())
	}

	if len(*argVerify) > 0 {
		return doVerify(password)
	}
	return doCrypt(password)
}

func getPassword() (string, error) {
	if len(*argInput) > 0 {
		if *argPrompt {
			return "", fmt.Errorf("conflicting arguments input and prompt are set")
		}

		return *argInput, nil
	}

	if *argPrompt {
		fmt.Print("Password: ")
		pass, err := console.ReadPassword()
		if err != nil {
			return "", fmt.Errorf("failed to prompt for password: %s", err.Error())
		}
		return pass, nil
	}

	sb := &strings.Builder{}
	_, err := io.Copy(sb, os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read password from StdIn: %s", err.Error())
	}
	return sb.String(), nil
}
