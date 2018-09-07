package redis_queue

import (
	"testing"

	"github.com/gomodule/redigo/redis"
)

func TestNewQueue(t *testing.T) {
	queue, err := NewQueue("localhost:6379", "_test_queue")
	defer queue.Close()

	if err != nil {
		t.Error(err.Error())
	}

	_, err = queue.Conn.Do("LPUSH", queue.Key, "test_item_1")
	if err != nil {
		t.Error(err.Error())
	}

	item, err := redis.String(queue.Conn.Do("RPOP", queue.Key))

	if err != nil {
		t.Error(err.Error())
	}

	if item != "test_item_1" {
		t.Error("incorrect value retrieved")
	}
}

func TestEnqueue(t *testing.T) {
	q, err := NewQueue("localhost:6379", "_push_test_queue")
	defer q.Close()
	q.Conn.Do("DEL", q.Key)

	if err != nil {
		t.Error(err.Error())
	}

	err = q.Enqueue("push_item_1")
	if err != nil {
		t.Error(err.Error())
	}

	err = q.Enqueue("push_item_2")
	if err != nil {
		t.Error(err.Error())
	}

	count, err := redis.Int(q.Conn.Do("LLEN", q.Key))
	if err != nil {
		t.Error(err.Error())
	}

	if count != 2 {
		t.Error("incorrect length in enqueue method")
	}

}

func TestDequeue(t *testing.T) {
	q, err := NewQueue("localhost:6379", "_pop_test_queue")
	defer q.Close()
	q.Conn.Do("DEL", q.Key)

	if err != nil {
		t.Error(err.Error())
	}

	_, err = q.Conn.Do("LPUSH", q.Key, "test_item_1")
	_, err = q.Conn.Do("LPUSH", q.Key, "test_item_2")
	_, err = q.Conn.Do("LPUSH", q.Key, "test_item_3")

	item, err := q.Dequeue()

	if item != "test_item_1" {
		t.Error("incorrect item returned in dequeue method")
	}

	count, err := redis.Int(q.Conn.Do("LLEN", q.Key))
	if err != nil {
		t.Error(err.Error())
	}

	if count != 2 {
		t.Error("incorrect length in dequeue method")
	}

}

func TestPollQueue(t *testing.T) {
	q, err := NewQueue("localhost:6379", "_poll_test_queue")
	defer q.Close()

	if err != nil {
		t.Error(err.Error())
	}

	q.Conn.Do("DEL", q.Key)

	hasItem, err := q.PollQueue()
	if err != nil {
		t.Error(err.Error())
	}

	if hasItem {
		t.Error("Queue should be empty")
	}

	q.Enqueue("test_item_1")

	hasItem, err = q.PollQueue()
	if err != nil {
		t.Error(err.Error())
	}

	if !hasItem {
		t.Error("Queue should not be empty")
	}

}

func TestRedisConnection(t *testing.T) {
	conn, err := redisConnection("localhost:6379")
	defer conn.Close()
	defer conn.Do("FLUSHALL")

	if err != nil {
		t.Errorf(err.Error())
	}
}
