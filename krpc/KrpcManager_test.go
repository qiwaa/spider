package krpc

import (
	"testing"
	"fmt"
)

func TestGetNewTransactionId(t *testing.T) {
	mgr := NewKrpcManager()
	str := mgr.getNewTransactionId()
	fmt.Println([]byte(str))
	fmt.Println(str)
}
