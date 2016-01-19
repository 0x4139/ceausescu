package ceausescu

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"log"
)

type Subscriber struct {
	wg             sync.WaitGroup
	isRunning      bool
	workChannel    chan string
	connectionPool *redis.Pool
}

type Worker func(string)

func (subscriber *Subscriber) doWork(fn Worker) {
	defer subscriber.wg.Done()
	for subscriber.isRunning {
		fn(<-subscriber.workChannel)
	}
}
func NewSubscriber(config Config) Subscriber {
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", config.RedisAddress)
		if err != nil {
			panic(err.Error())
		}
		return c, err
	}, config.MaxConnections)
	return Subscriber{
		config:config,
		connectionPool:redisPool,
	}
}
func (subscriber *Subscriber) Work(queueName string, concurrency int, fn Worker) {
	subscriber.isRunning = true
	subscriber.workChannel = make(chan string)
	for i := 0; i < concurrency; i++ {
		subscriber.wg.Add(1)
		go subscriber.doWork(fn)
	}
	go func() {
		for subscriber.isRunning {
			connection := subscriber.connectionPool.Get()
			data, err := connection.Do("BRPOP", "ceausescu:" + queueName, 0)
			if err != nil {
				log.Println(err)
				continue
			}
			connection.Close()
			result, err := redis.StringMap(data, nil)
			if err != nil {
				log.Println(err)
				continue
			}
			subscriber.workChannel <- result["ceausescu:" + queueName]
		}
	}()
}

func (subscriber *Subscriber) Wait() {
	subscriber.wg.Wait()
}
func (subscriber *Subscriber) Stop() {
	subscriber.isRunning = false
	subscriber.wg.Wait()
}