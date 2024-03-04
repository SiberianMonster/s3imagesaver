package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/SiberianMonster/s3imageserver/s3imageserver"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	_ = os.Mkdir("cache", os.ModePerm)
	s3imageserver.Run(nil)
}

func verifyToken(tokenString string) bool {
	publicKey, err := ioutil.ReadFile("verification.key")
	if err != nil {
		fmt.Print("Error:", err)
		return false
	}
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		errType := err.(*jwt.ValidationError)
		switch errType.Errors {
		case jwt.ValidationErrorMalformed:
			fmt.Println("malformed")
		case jwt.ValidationErrorUnverifiable:
			fmt.Println("unverifiable")
		case jwt.ValidationErrorSignatureInvalid:
			fmt.Println("signature invalid")
		case jwt.ValidationErrorExpired:
			fmt.Println("expired")
		case jwt.ValidationErrorNotValidYet:
			fmt.Println("not valid yet")
		}
	} else {
		return true
	}
	return false
}
