package utils

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"unsafe"
	"os"
	"strings"
	"strconv"
	"errors"
)

var POW2 []int = []int{
	128,
	64,
	32,
	16,
	8,
	4,
	2,
	1,
}

func Uint16ToBytes(i uint16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, i)
	return buf
}

func RandomString(length int) (string) {
	buff := make([]byte, length)
	rand.Read(buff)
	return string(buff)
}

func SetProcessName(name string) error {
	argv0str := (*reflect.StringHeader)(unsafe.Pointer(&os.Args[0]))
	argv0 := (*[1 << 30]byte)(unsafe.Pointer(argv0str.Data))[:argv0str.Len]

	n := copy(argv0, name)
	if n < len(argv0) {
		argv0[n] = 0
	}

	return nil
}

func GenAddress(ip string, port int) string {
	return strings.Join([]string{ip, strconv.Itoa(port)}, ":")
}

func ByteNumToStr(num int) (string) {
	var a byte = byte(num)
	return string(a)
}

// num : from head to end,the bit locates.(first is 0)
func GetBitInBytes(buf []byte, num int) (int, error) {
	length := len(buf)
	if num >= length*8 {
		return -1, errors.New("num is too big!")
	}

	b := buf[num/8]
	r := int(b) & POW2[num]
	if r != 0 {
		r = 1
	}
	return r, nil
}
