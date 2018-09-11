package redis_queue

import "github.com/gomodule/redigo/redis"

type Queue struct {
	Key  string
	Conn redis.Conn
}

// New Queue generates a new Redis list with Queue methods.
func NewQueue(uri, key string) (*Queue, error) {
	conn, err := redisConnection(uri)
	if err != nil {
		return nil, err
	}

	q := &Queue{Key: key, Conn: conn}

	return q, nil
}

// Enqueue adds a string to the start of the queue
func (q *Queue) Enqueue(item string) error {

	_, err := q.Conn.Do("LPUSH", q.Key, item)

	if err != nil {
		return err
	}
	return nil
}

// Dequeue removes the earliest item added to the queue
func (q *Queue) Dequeue() (string, error) {

	item, err := redis.String(q.Conn.Do("RPOP", q.Key))

	if err != nil {
		return item, err
	}

	return item, nil
}

func (q *Queue) Peek() (bool, error) {
	length, err := redis.Int(q.Conn.Do("LLEN", q.Key))
	if err != nil {
		return false, err
	}

	return length > 0, nil
}

func (q *Queue) Close() {
	q.Conn.Close()
}

func redisConnection(uri string) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", uri)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
