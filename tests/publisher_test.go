package ceausescu_test
import "testing"
import (
	"./.."
	"github.com/garyburd/redigo/redis"
	"strings"
	"os"
)

func TestPublisherShouldInsertValueInQueue(t *testing.T) {
	config := ceausescu.Config{
		MaxConnections:100,
		RedisAddress:os.Getenv("redis"),
	}
	publisher := ceausescu.NewPublisher(config)
	publisher.Publish("publisherTest", "laptecuorez")

	redisConnection, err := redis.Dial("tcp", config.RedisAddress)

	if err != nil {
		panic(err)
	}
	data, err := redisConnection.Do("RPOP", "ceausescu:publisherTest")
	if err != nil {
		panic(err)
	}
	returnValue, err := redis.String(data, nil)
	if err != nil {
		panic(err)
	}
	if strings.Compare(returnValue, "laptecuorez") != 0 {
		t.Fatalf("String wore no match expected:%s got:%s", "laptecuorez", returnValue)
	}

}