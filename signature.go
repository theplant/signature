package signature

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"hash"
	"io"
)

var (
	InvalidSignature = errors.New("invalid signature")
)

type Decoder struct {
	br     io.Reader
	sh1    hash.Hash
	secret string
}

func NewDecoder(r io.Reader, secret string) *Decoder {
	return &Decoder{
		br:     base64.NewDecoder(base64.StdEncoding, r),
		sh1:    sha1.New(),
		secret: secret,
	}
}

func (dec *Decoder) Decode(e interface{}) (err error) {
	// read head 24 byte, which is signature.
	var sig = make([]byte, 24)
	_, err = dec.br.Read(sig)
	if err != nil {
		return
	}

	_, err = io.WriteString(dec.sh1, dec.secret)
	if err != nil {
		return
	}

	teeReader := io.TeeReader(dec.br, dec.sh1)
	gd := gob.NewDecoder(teeReader)
	err = gd.Decode(e)

	if err != nil {
		return
	}

	if bytes.Compare(dec.sh1.Sum(nil), sig[:20]) != 0 {
		return InvalidSignature
	}

	return
}

type Encoder struct {
	bw     io.WriteCloser
	sh1    hash.Hash
	secret string
}

func NewEncoder(w io.Writer, secret string) *Encoder {
	return &Encoder{
		bw:     base64.NewEncoder(base64.StdEncoding, w),
		sh1:    sha1.New(),
		secret: secret,
	}
}

func (enc *Encoder) Encode(e interface{}) (err error) {
	defer enc.bw.Close()

	io.WriteString(enc.sh1, enc.secret)
	ge := gob.NewEncoder(enc.sh1)
	ge.Encode(e)

	var sig = enc.sh1.Sum(nil)
	sig = append(sig, []byte{0, 0, 0, 0}...)

	_, err = enc.bw.Write(sig)
	if err != nil {
		return
	}

	ge1 := gob.NewEncoder(enc.bw)
	err = ge1.Encode(e)

	return
}
