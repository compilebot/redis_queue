package redis_queue

import (
	"fmt"
	"sync"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func TestManyQueues(t *testing.T) {
	queues := 30
	inserts := 100
	deletes := 20

	q, _ := NewQueue("localhost:6379", "_break_test_queue")
	q.Conn.Do("DEL", q.Key)

	var wg sync.WaitGroup

	for i := 0; i < queues; i++ {
		wg.Add(1)
		go func(s int) {
			q, _ := NewQueue("localhost:6379", "_break_test_queue")

			for j := 0; j < inserts; j++ {
				err := q.Enqueue(string(s))
				if err != nil {
					t.Error("Error in insert: ", err.Error())
				}
			}

			for k := 0; k < deletes; k++ {
				_, err := q.Dequeue()
				if err != nil {
					t.Error("Error in delete: ", err.Error())
				}
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	count, err := redis.Int(q.Conn.Do("LLEN", q.Key))
	if err != nil {
		t.Error("Error checking length", err.Error())
	}

	expected := queues*inserts - queues*deletes

	if count != expected {
		t.Error(fmt.Sprintf("Count is invalid in ManyQueues testing function. Expected %d, got %d", expected, count))
	}
}
