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

func (subscriber *Subscriber) doWork(fn Worker, queue string, con redis.Conn) {
	for {
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
		con := subscriber.connectionPool.Get()
		go subscriber.doWork(fn, queueName, con)
	}
}

func (subscriber *Subscriber) Start()  {
	subscriber.wg.Add(1)
}
func (subscriber *Subscriber) Wait() {
	subscriber.wg.Wait()
}
func (subscriber *Subscriber) Stop() {
	subscriber.wg.Done()
}