package ceausescu

import "github.com/garyburd/redigo/redis"

type Config struct {
	RedisAddress   string
}

type Publisher struct {
	connection redis.Conn
	config     Config
}


func NewPublisher(config Config) Publisher {
	c, err := redis.Dial("tcp", config.RedisAddress)
	if err != nil {
		panic(err.Error())
	}
	return Publisher{
		connection :c,
	}
}
func (publisher *Publisher) Close() {
	publisher.connection.Close()
}
func (publisher *Publisher) Publish(queue string, value string) error {
	_, err := publisher.connection.Do("LPUSH", "ceausescu:" + queue, value)
	return err
}