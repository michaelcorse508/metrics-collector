package cryptography

import (
	"encoding/hex"
	"net/http"
)

const HeaderHashSHA256 = "HashSHA256"

type DsResponseWriter struct {
	http.ResponseWriter
	SecretKey []byte
}

func (rw *DsResponseWriter) Write(data []byte) (int, error) {
	digitalSign, err := SignMessageHMAC(data, rw.SecretKey)
	if err != nil {
		return rw.ResponseWriter.Write(data)
	}

	digitalSignString := hex.EncodeToString(digitalSign)

	rw.Header().Set(HeaderHashSHA256, digitalSignString)
	return rw.ResponseWriter.Write(data)
}
