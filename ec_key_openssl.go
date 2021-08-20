package libtrust

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io"
)

// Sign signs the data read from the io.Reader using a signature algorithm supported
// by the elliptic curve private key. If the specified hashing algorithm is
// supported by this key, that hash function is used to generate the signature
// otherwise the the default hashing algorithm for this key is used. Returns
// the signature and the name of the JWK signature algorithm used, e.g.,
// "ES256", "ES384", "ES512".
func (k *ecPrivateKey) Sign(data io.Reader, hashID crypto.Hash) (signature []byte, alg string, err error) {
	// Generate a signature of the data using the internal alg.
	// The given hashId is only a suggestion, and since EC keys only support
	// on signature/hash algorithm given the curve name, we disregard it for
	// the elliptic curve JWK signature implementation.
	hId := k.signatureAlgorithm.HashID()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(data)
	if err != nil {
		return nil, "", fmt.Errorf("error reading data: %s", err)
	}

	r, s, err := ecdsa.HashSign(rand.Reader, k.PrivateKey, buf.Bytes(), hId)
	if err != nil {
		return nil, "", fmt.Errorf("error producing signature: %s", err)
	}
	rBytes, sBytes := r.Bytes(), s.Bytes()
	octetLength := (k.ecPublicKey.Params().BitSize + 7) >> 3
	// MUST include leading zeros in the output
	rBuf := make([]byte, octetLength-len(rBytes), octetLength)
	sBuf := make([]byte, octetLength-len(sBytes), octetLength)

	rBuf = append(rBuf, rBytes...)
	sBuf = append(sBuf, sBytes...)

	signature = append(rBuf, sBuf...)
	alg = k.signatureAlgorithm.HeaderParam()

	return
}
