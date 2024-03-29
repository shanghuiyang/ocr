# word-recognizer
ocr implements Optical Character Recognition(OCR) using the APIs provided by web services like baidu and google.

## usage
```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/shanghuiyang/oauth"
	"github.com/shanghuiyang/ocr"
)

const (
	baiduOcrApiKey    = "your_baidu_ocr_api_key"
	baiduOcrSecretKey = "your_baidu_ocr_secret_key"
)

func main() {
	imgf := "images/test.jpg"
	img, err := ioutil.ReadFile(imgf)
	if err != nil {
		log.Printf("failed to read file %v, error: %v", imgf, err)
		os.Exit(1)
	}

	auth := oauth.NewBaiduOauth(baiduOcrApiKey, baiduOcrSecretKey, oauth.NewCacheImp())
	r := ocr.NewBaiduOCR(auth)
	text, err := r.Recognize(img)
	if err != nil {
		log.Printf("failed to recognize words from the image, error: %v", err)
		os.Exit(1)
	}
	fmt.Println(text)
}
```
