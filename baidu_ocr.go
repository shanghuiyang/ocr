package ocr

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/shanghuiyang/oauth"
)

const (
	baiduOcrAPI = "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic"
)

// BaiduOCR ...
type BaiduOCR struct {
	auth oauth.Oauth
}

type baiduResponse struct {
	LogID     uint64        `json:"log_id"`
	Results   []baiduResult `json:"words_result"`
	ResultNum uint32        `json:"words_result_num"`
	ErrorCode int           `json:"error_code"`
	ErrorMsg  string        `json:"error_msg"`
}

type baiduResult struct {
	Words string `json:"words"`
}

// NewBaiduOCR ...
func NewBaiduOCR(auth oauth.Oauth) *BaiduOCR {
	return &BaiduOCR{
		auth: auth,
	}
}

// Recognize ...
func (r *BaiduOCR) Recognize(image []byte) (string, error) {
	token, err := r.auth.Token()
	if err != nil {
		return "", err
	}

	b64img := base64.StdEncoding.EncodeToString(image)
	formData := url.Values{
		"access_token": {token},
		"image":        {b64img},
	}
	resp, err := http.PostForm(baiduOcrAPI, formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res baiduResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	if res.ErrorCode > 0 {
		return "", fmt.Errorf("error_code: %v, error_msg: %v", res.ErrorCode, res.ErrorMsg)
	}

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
