package container

import (
	"time"
	"github.com/chenminjian/spider/utils"
)

// blockedItem represents a blocked node.
type blockedItem struct {
	ip         string
	port       int
	createTime time.Time
}

// blackList manages the blocked nodes including which sends bad information
// and can't ping out.
type BlackList struct {
	list         *SyncedMap
	maxSize      int
	expiredAfter time.Duration
}

// newBlackList returns a blackList pointer.
func NewBlackList(size int) *BlackList {
	return &BlackList{
		list:         NewSyncedMap(),
		maxSize:      size,
		expiredAfter: time.Hour * 1,
	}
}

// genKey returns a key. If port is less than 0, the key wil be ip. Ohterwise
// it will be `ip:port` format.
func (bl *BlackList) genKey(ip string, port int) string {
	key := ip
	if port >= 0 {
		key = utils.GenAddress(ip, port)
	}
	return key
}

// insert adds a blocked item to the blacklist.
func (bl *BlackList) Insert(ip string, port int) {
	if bl.list.Len() >= bl.maxSize {
		return
	}

	bl.list.Set(bl.genKey(ip, port), &blockedItem{
		ip:         ip,
		port:       port,
		createTime: time.Now(),
	})
}

// delete removes blocked item form the blackList.
func (bl *BlackList) delete(ip string, port int) {
	bl.list.Delete(bl.genKey(ip, port))
}

// validate checks whether ip-port pair is in the block nodes list.
func (bl *BlackList) Has(ip string, port int) bool {
	if _, ok := bl.list.Get(ip); ok {
		return true
	}

	key := bl.genKey(ip, port)

	v, ok := bl.list.Get(key)
	if ok {
		if time.Now().Sub(v.(*blockedItem).createTime) < bl.expiredAfter {
			return true
		}
		bl.list.Delete(key)
	}
	return false
}

// clear cleans the expired items every 10 minutes.
func (bl *BlackList) clear() {
	for _ = range time.Tick(time.Minute*10) {
		keys := make([]interface{}, 0, 100)

		for item := range bl.list.Iter() {
			if time.Now().Sub(
				item.Val.(*blockedItem).createTime) > bl.expiredAfter {

				keys = append(keys, item.key)
			}
		}

		bl.list.DeleteMulti(keys)
	}
}
