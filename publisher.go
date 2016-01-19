package ceausescu

import "github.com/garyburd/redigo/redis"

type Publisher struct {
	connectionPool *redis.Pool
}

func NewPublisher(config Config) Publisher {
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", config.RedisAddress)
		if err != nil {
			panic(err.Error())
		}
		return c, err
	}, config.MaxConnections)

	return Publisher{
		connectionPool :redisPool,
	}
}
func (publisher *Publisher) Close() {
	publisher.connectionPool.Close()
}
func (publisher *Publisher) Publish(queue string, value string) error {
	connection := publisher.connectionPool.Get()
	defer connection.Close()
	_, err := connection.Do("LPUSH", "ceausescu:" + queue, value)
	return err
}