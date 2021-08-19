package ocr

type OCR interface {
	Recognize(image []byte) (string, error)
}
