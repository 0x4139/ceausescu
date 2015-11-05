package ceausescu
import (
	"github.com/garyburd/redigo/redis"
	"sync"
//	"log"
	"log"
)

type Subscriber struct {
	config      Config
	wg          sync.WaitGroup
	isRunning   bool
	workChannel chan string
}

type Worker func(string)

func (subscriber *Subscriber) doWork(fn Worker) {
	defer subscriber.wg.Done()
	for subscriber.isRunning {
		fn(<-subscriber.workChannel)
	}
}
func NewSubscriber(config Config) Subscriber {
	return Subscriber{
		config:config,
	}
}
func (subscriber *Subscriber) Work(queueName string, concurency int, fn Worker) {
	subscriber.isRunning = true
	subscriber.workChannel = make(chan string)
	for i := 0; i < concurency; i++ {
		subscriber.wg.Add(1)
		go subscriber.doWork(fn)
	}
	go func() {
		con, err := redis.Dial("tcp", subscriber.config.RedisAddress)
		defer con.Close()
		if err != nil {
			panic(err)
		}
		for subscriber.isRunning {
			data, err := con.Do("BRPOP", "ceausescu:" + queueName, 0)
			if err != nil {
				log.Println(err)
				continue
			}
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