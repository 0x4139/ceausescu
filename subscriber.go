package ceausescu
import (
	"github.com/garyburd/redigo/redis"
	"sync"
//	"log"
)

type Subscriber struct {
	config    Config
	wg        sync.WaitGroup
	isRunning bool
}

type Worker func(string, error)

func (subscriber *Subscriber) doWork(fn Worker, queue string) {
	con, err := redis.Dial("tcp", subscriber.config.RedisAddress)
	defer con.Close()
	defer subscriber.wg.Done()
	if err != nil {
		fn("", err)
		return
	}
	for subscriber.isRunning {
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
func NewSubscriber(config Config) Subscriber {
	return Subscriber{
		config:config,
	}
}
func (subscriber *Subscriber) Work(queueName string, concurency int, fn Worker) {
	subscriber.isRunning = true
	for i := 0; i < concurency; i++ {
		subscriber.wg.Add(1)
		go subscriber.doWork(fn, queueName)
	}
}

func (subscriber *Subscriber) Wait() {
	subscriber.wg.Wait()
}
func (subscriber *Subscriber) Stop() {
	subscriber.isRunning = false
	subscriber.wg.Wait()
}