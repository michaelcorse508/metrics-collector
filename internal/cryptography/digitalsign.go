package cryptography

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
)

var ErrSignCalculatingError = errors.New("cannot calculate sign on given message")
var ErrInvalidSign = errors.New("invalid sign")

func SignMessageHMAC(message []byte, key []byte) ([]byte, error) {
	signer := hmac.New(sha256.New, key)
	_, err := signer.Write(message)
	if err != nil {
		return nil, ErrSignCalculatingError
	}

	return signer.Sum(nil), nil
}

func CheckMessageHMAC(message []byte, clientDigitalSign []byte, key []byte) error {
	serverDigitalSign, err := SignMessageHMAC(message, key)
	if err != nil {
		return err
	}

	result := hmac.Equal(serverDigitalSign, clientDigitalSign)
	if result {
		return nil
	}

	return ErrInvalidSign
}
