package redis_queue

import "github.com/gomodule/redigo/redis"

// Queue is a data structure with a Key, which is the name of the Redis endpoint / database name. The structure also contains a connection object that we use to interact with Redis
// The items in the queue are all stored in Redis so they are not represented in the struct.
// Currently the queue only accepts strings.
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

// Length provides the amount of items currently in the queue
func (q *Queue) Length() (int, error) {
	length, err := redis.Int(q.Conn.Do("LLEN", q.Key))
	if err != nil {
		return 0, err
	}

	return length, nil
}

// Peek checks for items in the queue, returning false if the queue is empty or true otherwise
func (q *Queue) Peek() (bool, error) {
	length, err := redis.Int(q.Conn.Do("LLEN", q.Key))
	if err != nil {
		return false, err
	}

	return length > 0, nil
}

// Close closes the Redis connection. Using defer q.Close() after you create a new queue is a good idea
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
