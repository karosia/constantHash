package ConsistentHash

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type ConsistentHash struct {
	hash     Hash           // Hash function
	replicas int            // The number of nodes
	keys     []int          // hash keys
	hashMap  map[int]string // hash -> Node Name
}

func NewConsistentHash(replicas int, hash Hash) *ConsistentHash {
	h := &ConsistentHash{
		replicas: replicas,
		hash:     hash,
		hashMap:  make(map[int]string),
	}
	if h.hash == nil {
		h.hash = crc32.ChecksumIEEE
	}
	return h
}

func (h *ConsistentHash) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < h.replicas; i++ {
			hash := int(h.hash([]byte(strconv.Itoa(i) + node)))
			h.keys = append(h.keys, hash)
			h.hashMap[hash] = node
		}

	}
	sort.Ints(h.keys)
}

func (h *ConsistentHash) Get(key string) string {
	if len(h.keys) == 0 {
		return ""
	}
	hash := int(h.hash([]byte(key)))
	fmt.Printf("key: %-10s → hash: %10d\n", key, hash)
	idx := sort.Search(len(h.keys), func(i int) bool { return h.keys[i] >= hash })
	if idx == len(h.keys) {
		idx = 0
	}
	return h.hashMap[h.keys[idx]]
}

func (h *ConsistentHash) Remove(node string) {
	for i := 0; i < h.replicas; i++ {
		hash := int(h.hash([]byte(strconv.Itoa(i) + node)))
		delete(h.hashMap, hash)

		idx := sort.SearchInts(h.keys, hash)
		if idx < len(h.keys) && h.keys[idx] == hash {
			h.keys = append(h.keys[:idx], h.keys[idx+1:]...)
		}
	}
}

func (h *ConsistentHash) GetRanges() {
	for i, hash := range h.keys {
		currentNode := h.hashMap[hash]

		var nextHash int
		if i+1 < len(h.keys) {
			nextHash = h.keys[i+1]
		} else {
			nextHash = h.keys[0]
		}

		fmt.Printf("Range %10d ~ %10d → Node %s\n", hash, nextHash, currentNode)
	}
}
