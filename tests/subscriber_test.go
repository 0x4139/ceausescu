package ceausescu_test
import "testing"
import (
	"./.."
	"strings"
	"log"
	"os"
)

func TestShouldCheckTheSubscriberConcurently(t *testing.T) {
	config := ceausescu.Config{
		RedisAddress:os.Getenv("redis"),
	}
	log.Println("Building work!")
	publisher := ceausescu.NewPublisher(config)
	for i := 0; i < 1024; i++ {
		err := publisher.Publish("test", "laptecuorez")
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	publisher.Close()
	log.Println("Done! building work!")
	subscriber := ceausescu.NewSubscriber(config)
	log.Println("Building workers!")
	subscriber.Work("test", 1024, func(data string) {
		log.Println(data)
		
		if strings.Compare(data, "laptecuorez") != 0 {
			t.Fatalf("strings don't match expected:laptecuorez got:%s", data)
		}
	})
	log.Println("Waiting for completion!")
	log.Println("Test finished!")
	subscriber.Wait()
}