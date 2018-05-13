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
