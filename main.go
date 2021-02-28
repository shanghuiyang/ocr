package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shanghuiyang/go-speech/oauth"
	"github.com/shanghuiyang/word-recognizer/recognizer"
)

const (
	// replace your_app_key and your_secret_key with yours
	appKey    = "hc1tZ5KDCGk3tGHMhlvQRclv"
	secretKey = "0RNCrXGRfWSsC5BDktH5ujBXxK6ixpnx"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("error: invalid args")
		fmt.Println("usage: image-recognizer test.jpg")
		os.Exit(1)
	}
	imgf := os.Args[1]

	auth := oauth.New(appKey, secretKey, oauth.NewCacheMan())
	r := recognizer.New(auth)
	text, err := r.Recognize(imgf)
	if err != nil {
		log.Printf("failed to recognize the image, error: %v", err)
		os.Exit(1)
	}
	fmt.Println(text)
}
