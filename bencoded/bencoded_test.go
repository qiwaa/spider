package bencoded

import (
	"testing"
	"fmt"
)

func TestDecode(t *testing.T) {
	_, err := Decode(nil)
	if err != nil {
		fmt.Printf("err:%v", err)
	}
}
