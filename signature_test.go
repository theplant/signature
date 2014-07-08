package signature

import (
	"bytes"

	"reflect"
	// "encoding/base64"
	// "encoding/gob"

	"testing"
)

type Value struct {
	Name string
	Age  int
}

type NestedValue struct {
	GroupName string
	People    []Value
	Leader    Value
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

func TestEncodeAndDecodeForNestedStruct(t *testing.T) {
	var secret = "234234"

	var err error

	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf, secret)
	v1 := &NestedValue{
		GroupName: "丐帮",
		People: []Value{
			Value{"马大元", 46},
			Value{"白世镜", 54},
		},
		Leader: Value{"乔峰", 30},
	}
	err = enc.Encode(v1)

	if err != nil {
		t.Error(err)
	}

	var v2 *NestedValue
	dec := NewDecoder(buf, secret)
	err = dec.Decode(&v2)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(v1, v2) {
		t.Error(v1, v2)
	}

}

func TestEncodeMapToString(t *testing.T) {
	v := map[string]string{"name": "Felix", "age": "18", "color": "red"}
	secret := "q3244214"
	s1, _ := EncodeToString(&v, secret)
	s2, _ := EncodeToString(&v, secret)
	for i := 0; i < 100; i++ {
		if s1 != s2 {
			t.Error(s1, s2)
		}

		s1, _ = EncodeToString(&v, secret)
		s2, _ = EncodeToString(&v, secret)
	}
}

func TestDecodeStringToMap(t *testing.T) {
	v := map[string]string{"name": "Felix", "age": "18"}
	secret := "q3244214"

	s1, _ := EncodeToString(&v, secret)
	v1 := map[string]string{}
	err := DecodeString(s1, &v1, secret)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(v1, v) {
		t.Error(v1, v)
	}
}

func TestDecodeToMapWithCustumizeInterface(t *testing.T) {
	type Pet struct {
		Name   string
		Age    int
		Weight float32
	}

	v := map[string]interface{}{
		"name": "Felix",
		"age":  18,
		"pet":  Pet{"peter", 2, 23.4},
	}
	secret := "q3244214"

	s1, _ := EncodeToString(&v, secret)
	v1 := map[string]interface{}{}
	err := DecodeString(s1, &v1, secret)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(v1, v) {
		t.Error(v1, v)
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
