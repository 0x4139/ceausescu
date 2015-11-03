package ceausescu

import "github.com/garyburd/redigo/redis"

type Config struct {
	RedisAddress   string
	MaxConnections int
}

type Publisher struct {
	connectionPool *redis.Pool
	config         Config
}


func NewPublisher(config Config) Publisher {
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", config.RedisAddress)
		if err != nil {
			panic(err)
		}
		return c, nil
	}, config.MaxConnections)

	return Publisher{
		connectionPool :redisPool,
	}
}
func (publisher *Publisher) Close() {
	publisher.connectionPool.Close()
}
func (publisher *Publisher) Publish(queue string, value string) error {
	_, err := publisher.connectionPool.Get().Do("LPUSH", "ceausescu:" + queue, value)
	return err
}