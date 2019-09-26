package checksum

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
)

// GoModHash is the result of hash
type GoModHash struct {
	Hash                  string
	HashBase64            string
	HashSynthesized       string
	HashSynthesizedBase64 string
	GoSum                 string
}

// HashGoMod hash the go.mod file
func HashGoMod(goModFile string) (GoModHash, error) {
	h := sha256.New()
	data, err := ioutil.ReadFile(goModFile)
	if err != nil {
		return GoModHash{}, err
	}

	r := ioutil.NopCloser((bytes.NewReader(data)))
	if err != nil {
		return GoModHash{}, err
	}

	_, err = io.Copy(h, r)
	r.Close()
	if err != nil {
		return GoModHash{}, err
	}

	bArr := h.Sum(nil)
	hash := fmt.Sprintf("%x", bArr)

	hFinal := sha256.New()
	fmt.Fprintf(hFinal, "%x  %s\n", bArr, "go.mod")
	bArrFinal := hFinal.Sum(nil)
	hashSynthesized := fmt.Sprintf("%x", bArrFinal)

	allHash := GoModHash{
		Hash:                  hash,
		HashBase64:            Base64(bArr),
		HashSynthesized:       hashSynthesized,
		HashSynthesizedBase64: Base64(bArrFinal),
		GoSum:                 "h1:" + Base64(bArrFinal),
	}

	return allHash, nil

}
