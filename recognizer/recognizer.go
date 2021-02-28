package recognizer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/shanghuiyang/go-speech/oauth"
)

const (
	baiduURL = "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic"
)

// Recognizer ...
type Recognizer struct {
	auth *oauth.Oauth
}

type response struct {
	LogID     uint64   `json:"log_id"`
	Results   []result `json:"words_result"`
	ResultNum uint32   `json:"words_result_num"`
	ErrorCode int      `json:"error_code"`
	ErrorMsg  string   `json:"error_msg"`
}

type result struct {
	Words string `json:"words"`
}

// New ...
func New(auth *oauth.Oauth) *Recognizer {
	return &Recognizer{
		auth: auth,
	}
}

// Recognize ...
func (r *Recognizer) Recognize(imageFile string) (string, error) {
	token, err := r.auth.GetToken()
	if err != nil {
		return "", err
	}

	b64img, err := r.b64Image(imageFile)
	if err != nil {
		return "", err
	}

	formData := url.Values{
		"access_token": {token},
		"image":        {b64img},
	}
	resp, err := http.PostForm(baiduURL, formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	if res.ErrorCode > 0 {
		return "", fmt.Errorf("error_code: %v, error_msg: %v\n", res.ErrorCode, res.ErrorMsg)
	}

	// fmt.Printf("num: %v\n", res.ResultNum)
	// fmt.Printf("results: %v\n", res.Results)

	if res.ResultNum == 0 {
		return "", fmt.Errorf("no results")
	}

	var words string
	for _, r := range res.Results {
		if words != "" {
			words += "\n"
		}
		words += r.Words
	}
	return words, nil
}

func (r *Recognizer) b64Image(imageFile string) (string, error) {
	file, err := os.Open(imageFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	image, err := ioutil.ReadAll(file)
	if err != nil {
		return "nil", err
	}
	b64img := base64.StdEncoding.EncodeToString(image)
	return b64img, nil
}
