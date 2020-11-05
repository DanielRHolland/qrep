package main

import (
	"bufio"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
)

type authenticator struct {
	userPasswords map[string]string
	tokenUsers    map[string]string //tokens should timeout

}

func newAuthFromFile(filename string) *authenticator {
	auth := authenticator{
		userPasswords: make(map[string]string),
		tokenUsers:    make(map[string]string),
	}
	if filename == "" {
		auth.userPasswords["admin"] = "admin"
	} else {
		readFile, err := os.Open(filename)
		if err != nil {
			log.Fatal("Could not read users file: ", filename)
		}
		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {
			temp := strings.Split(fileScanner.Text(), ",")
			auth.userPasswords[temp[0]] = temp[1]
		}
		readFile.Close()
	}
	return &auth
}

func (a *authenticator) checkUserPasswordValid(username string, password string) bool {
	//TODO implement properly
	pw, exists := a.userPasswords[username]
	return exists && pw == password
}

func (a *authenticator) validateUser(username string, password string) (token string, valid bool) {
	if a.checkUserPasswordValid(username, password) {
		token := uuid.New().String()
		a.tokenUsers[token] = username
		return token, true
	} else {
		return "", false
	}
}

func (a *authenticator) tokenValid(token string) (username string, valid bool) {
	username, valid = a.tokenUsers[token]
	return username, valid
}
