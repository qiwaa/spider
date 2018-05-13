package utils

import (
	"encoding/binary"
	"math/rand"
	"reflect"
	"unsafe"
	"os"
	"strings"
	"strconv"
)

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
