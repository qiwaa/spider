package route

type trieNode struct {
	key    int
	bucket *bucket
	left   *trieNode
	right  *trieNode
}

func NewTrieNode(k int) *trieNode {
	return &trieNode{
		key:    k,
		bucket: newBucket(),
		left:   nil,
		right:  nil,
	}
}
