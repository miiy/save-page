package redisdoc

import (
	"log"
	"testing"
)

func TestRedisDoc_Get(t *testing.T) {
	redisDoc := NewRedisDoc()
	if err := redisDoc.Get(); err != nil {
		log.Print(err)
	}
}