package signature

import (
	"bytes"
	// "encoding/base64"
	// "encoding/gob"
	// "fmt"
	"testing"
)

type Value struct {
	Name string
	Age  int
}

func TestEncodeAndDecode(t *testing.T) {
	var secret = "234234"

	var err error

	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, secret)
	v1 := &Value{"Azuma", 25}
	err = enc.Encode(v1)

	if err != nil {
		t.Error(err)
	}

	var v2 Value
	dec := NewDecoder(buf, secret)
	err = dec.Decode(&v2)

	if err != nil {
		t.Error(err)
	}

	if v1.Name != v2.Name {
		t.Error(v1, v2)
	}

}

func TestDecodeStringAndEncodeToString(t *testing.T) {
	v1 := &Value{"Felix", 18}
	var v2 Value
	secret := "q3244214"
	s, _ := EncodeToString(v1, secret)
	DecodeString(s, &v2, secret)
	if v1.Name != v2.Name {
		t.Error(v1, v2)
	}
}

// func TestGobBase64(t *testing.T) {
// 	v1 := &Value{"Azuma", 25}

// 	buf := bytes.NewBuffer(nil)
// 	enc := gob.NewEncoder(base64.NewEncoder(base64.StdEncoding, buf))
// 	enc.Encode(v1)
// 	fmt.Println(buf.String())

// 	buf2 := bytes.NewBuffer(nil)
// 	enc2 := gob.NewEncoder(buf2)
// 	enc2.Encode(v1)

// 	fmt.Println("==========")
// 	fmt.Println(buf2.Bytes())
// 	fmt.Println(base64.StdEncoding.EncodeToString(buf2.Bytes()))

// }

// func TestBase64(t *testing.T) {
// 	fmt.Println("")
// 	var b []byte
// 	for i := 1; i < 60; i++ {
// 		fmt.Println(i)
// 		b = append(b, 'A')
// 		fmt.Println(b)
// 		fmt.Println(base64.URLEncoding.EncodeToString(b))
// 	}
// }
