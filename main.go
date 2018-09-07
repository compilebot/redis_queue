package redis_queue

import "github.com/gomodule/redigo/redis"

type Queue struct {
	uri  string
	key  string
	conn redis.Conn
}

// New Queue generates a new Redis list with Queue methods.
func NewQueue(uri, key string) (*Queue, error) {
	conn, err := redisConnection(uri)
	if err != nil {
		return nil, err
	}

	q := &Queue{uri: uri, key: key, conn: conn}

	return q, nil
}

func (q Queue) Enqueue() {

}

func (q Queue) Dequeue() {

}

func redisConnection(uri string) (redis.Conn, error) {
	conn, err := redis.Dial("tcp", uri)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
