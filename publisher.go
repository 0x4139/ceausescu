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
	redisPool := &redis.Pool{
		MaxIdle: config.MaxConnections,
		MaxActive: config.MaxConnections, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.RedisAddress)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
	return Publisher{
		connectionPool :redisPool,
	}
}
func (publisher *Publisher) Close() {
	publisher.connectionPool.Close()
}
func (publisher *Publisher) Publish(queue string, value string) error {
	con, err := publisher.connectionPool.Get()
	if err!nil {
		continue
	}
	defer con.Close()
	_, err = con.Do("LPUSH", "ceausescu:" + queue, value)
	return err
}