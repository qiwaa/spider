package route

import (
	"github.com/chenminjian/spider/utils"
	"fmt"
)

const NUM_OF_SON = 2

type sentinel struct {
	root  *trieNode
	depth int
}

func newSentinel() (*sentinel) {
	s := &sentinel{
		root:  NewTrieNode(0),
		depth: 20,
	}
	return s
}

func (s *sentinel) insert(key string, val interface{}) (error) {
	err := s.insertRecur(key, 0, s.root, val)
	if err != nil {
		return err
	}
	return nil
}

func (s *sentinel) insertRecur(key string, num int, parent *trieNode, val interface{}) (error) {
	lr, err := utils.GetBitInBytes([]byte(key), num)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var son *trieNode
	if lr == 0 {
		if parent.left == nil {
			parent.left = NewTrieNode(0)
		}
		son = parent.left
	} else if lr == 1 {
		if parent.right == nil {
			parent.right = NewTrieNode(1)
		}
		son = parent.right
	}

	err = s.insertRecur(key, num+1, son, val)
	if err != nil {
		return err
	}

	return nil
}
