package ceausescu_test
import "testing"
import (
	"./.."
	"sync"
	"strings"
	"log"
	"os"
)

func TestShouldCheckTheSubscriberConcurently(t *testing.T) {
	config := ceausescu.Config{
		MaxConnections:100,
		RedisAddress:os.Getenv("redis"),
	}
	log.Println("Building work!")
	var wg sync.WaitGroup
	publisher := ceausescu.NewPublisher(config)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		publisher.Publish("test", "laptecuorez")
	}
	publisher.Close()
	log.Println("Done! building work!")

	subscriber := ceausescu.NewSubscriber(config)
	log.Println("Building workers!")
	subscriber.Work("test", 100, func(data string, err error) {
		if err != nil {
			t.Fatal(err.Error())
		}
		if strings.Compare(data, "laptecuorez") != 0 {
			t.Fatalf("strings don't match expected:laptecuorez got:%s", data)
		}
		defer wg.Done()
	})
	log.Println("Done!Building workers!")
	log.Println("Waiting for completion!")
	wg.Wait()
	subscriber.Close()
	log.Println("Test finished!")
}