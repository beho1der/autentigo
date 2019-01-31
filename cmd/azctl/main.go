package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mcluseau/autorizo/client"
)

var (
	az *client.Client

	termIn  = bufio.NewReader(os.Stdin)
	termOut = os.Stderr
)

func main() {
	defaultServer := os.Getenv("AZCTL_SERVER")
	if len(defaultServer) == 0 {
		defaultServer = "http://localhost:8080"
	}

	serverURL := flag.String("server", defaultServer, "Autorizo server URL")
	flag.Parse()

	az = client.New(*serverURL)

	args := flag.Args()
	if len(args) < 1 {
		fail(errors.New("need a command"))
	}

	// handle termination signals
	sig := make(chan os.Signal, 1)
	go func() {
		<-sig
		resetTerm()
		os.Exit(1)
	}()

	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	// execute command
	switch v := args[0]; v {
	case "login":
		login()

	case "validate":
		if len(args) < 2 {
			fail(errors.New("validate command needs a token"))
		}

		validate(args[1])
	default:
		fail(fmt.Errorf("unknown command: %s", v))
	}
}

func login() {
	termOut.WriteString("username: ")
	username, err := termIn.ReadString('\n')
	fail(err)

	termOut.WriteString("password: \x1b[8m")
	password, err := termIn.ReadString('\n')
	fail(err)

	// remove trailing \n
	username = username[0 : len(username)-1]
	password = password[0 : len(password)-1]

	resetTerm()

	res, err := az.Login(username, password)
	fail(err)

	fmt.Println(res.Token)
}

func fail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		resetTerm()
		os.Exit(255)
	}
}

func resetTerm() {
	termOut.WriteString("\x1b[0m")
}

func validate(token string) {
	claims := jwt.MapClaims{}
	ok, err := az.Validate(token, claims)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if ok {
		fmt.Println("Token is valid")
		os.Exit(0)
	} else {
		fmt.Println("Token is NOT valid")
		os.Exit(1)
	}
}