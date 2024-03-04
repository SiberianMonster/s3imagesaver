# S3imageserver package for Go

The S3imageserver package for Go provides with a scalable API server that fetches images from a S3 bucket, resizes them, displays then and caches them on the instance. In case of errors it optionally displays a fallback image, but still returns a 404 response, thus preventing caching.

### Usage

Run the server, pass the optional configuration parameter:

	./s3imageserver -c=config.json


You can use it as a package like this:

	package main

	import (
		"fmt"

		"github.com/SiberianMonster/s3imageserver/s3imageserver"
		"github.com/dgrijalva/jwt-go"
	)

	func main() {
		s3imageserver.Run(nil)
	}


There is also an option to pass a handler for validation, so it's easy to implement JWT client verification:

	package main

	import (
		"fmt"
		"io/ioutil"

		"github.com/SiberianMonster/s3imageserver/s3imageserver"
		"github.com/dgrijalva/jwt-go"
	)

	func main() {
		s3imageserver.Run(verifyToken)
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


### Install

Install the package:

	go get github.com/SiberianMonster/s3imageserver