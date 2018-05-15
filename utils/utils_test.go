package utils

import (
	"testing"
	"fmt"
)

func TestInt2Bytes(t *testing.T) {
	var a uint16 = 0
	buf := Uint16ToBytes(a)
	fmt.Println(buf)
}

func TestByteNumToStr(t *testing.T) {
	str := ByteNumToStr(1)
	fmt.Println([]byte(str))
	fmt.Println(str)
}

func TestGetBitInBytes(t *testing.T) {
	buf := []byte{0x01}
	res, err := GetBitInBytes(buf, 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("result:%d\r\n", res)

	buf = []byte{0x0F}
	res, err = GetBitInBytes(buf, 5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("result:%d\r\n", res)
}
