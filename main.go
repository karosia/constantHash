package main

import (
	"fmt"
	"math/rand"
	"time"

	"consistentHash/ConsistentHash"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func randomString(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func main() {
	hashRing := ConsistentHash.NewConsistentHash(10, nil)
	hashRing.Add("NodeA", "NodeB", "NodeC", "NodeD", "NodeE", "NodeF")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	keys := make([]string, 100)
	for i := 0; i < 100; i++ {
		keys[i] = randomString(r, 6)
	}

	fmt.Println(keys)

	for _, key := range keys {
		node := hashRing.Get(key)
		fmt.Printf("Key %q is assigned to node %q\n", key, node)
	}

	//hashRing.GetRanges()

}
