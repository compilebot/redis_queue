package redis_queue

import (
	"testing"
)

func TestNewQueue(t *testing.T) {

}

func TestRedisConnection(t *testing.T) {
	conn, err := redisConnection("localhost:6379")

	if err != nil {
		t.Errorf(err.Error())
	}

	_ = conn
}
