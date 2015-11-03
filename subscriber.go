package ceausescu
import (
	"github.com/garyburd/redigo/redis"
	"sync"
)

type Subscriber struct {
	connectionPool *redis.Pool
	config         Config
	wg             sync.WaitGroup
}

type Worker func(string, error)

func (subscriber *Subscriber) doWork(fn Worker, queue string) {
	for {
		con := subscriber.connectionPool.Get()

		defer con.Close()
		data, err := con.Do("BRPOP", "ceausescu:" + queue, 0)
		if err != nil {
			fn("", err)
			continue
		}
		result, err := redis.StringMap(data, nil)
		if err != nil {
			fn("", err)
			continue
		}
		fn(result["ceausescu:" + queue], err)
	}
}
func (subscriber *Subscriber) Close() {
	subscriber.connectionPool.Close()
}
func NewSubscriber(config Config) Subscriber {
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

	return Subscriber{
		connectionPool :redisPool,
	}
}

func (subscriber *Subscriber) Work(queueName string, concurency int, fn Worker) {
	for i := 0; i < concurency; i++ {
		subscriber.wg.Add(1)
		go subscriber.doWork(fn, queueName)
	}
}

func (subscriber *Subscriber) Wait() {
	subscriber.wg.Wait()
}

func (subscriber *Subscriber) Done() {
	subscriber.wg.Done()
}