package ceausescu
import "github.com/garyburd/redigo/redis"

type Subscriber struct {
	connectionPool *redis.Pool
	config         Config
}

type Worker func(string, error)

func (subscriber *Subscriber) doWork(fn Worker,queue string) {
	for {
		data, err := subscriber.connectionPool.Get().Do("RPOP", "ceausescu/"+queue)
		if err != nil {
			fn("", err)
			continue
		}
		returnValue, err := redis.String(data, nil)
		fn(returnValue, err)
	}
}
func (subscriber *Subscriber) Close() {
	subscriber.connectionPool.Close()
}
func NewSubscriber(config Config) Subscriber {
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", config.RedisAddress)
		if err != nil {
			panic(err)
		}
		return c, nil
	}, config.MaxConnections)

	return Subscriber{
		connectionPool :redisPool,
	}
}

func (subscriber *Subscriber) Work(queueName string, concurency int, fn Worker) {
	for i := 0; i < concurency; i++ {
		go subscriber.doWork(fn,queueName)
	}
}