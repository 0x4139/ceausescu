![logo](https://fbcdn-profile-a.akamaihd.net/hprofile-ak-xaf1/v/t1.0-1/c27.27.336.336/s160x160/164778_139702199420541_6459978_n.jpg?oh=ef653cb94ffeb08a02ac6d53e78b288c&oe=56A2BFBD&__gda__=1452718710_23d979a71b91cf35cc80eaac17d21e8a)
# Ceausescu
A golang queue backed by redis with concurrency support
## Installation
```go
go get github.com/0x4139/ceausescu
```
## Usage
### Publish a job
```go
config := ceausescu.Config{
		MaxConnections:100,
		RedisAddress:os.Getenv("redis"),
	}
publisher := ceausescu.NewPublisher(config)
publisher.Publish("foo", "bar")
```
### Process a job
```go
config := ceausescu.Config{
    MaxConnections:100,
    RedisAddress:os.Getenv("redis"),
}
subscriber := ceausescu.NewSubscriber(config)
subscriber.Work("test", 100, func(data string, err error) {
    if err != nil {
       panic(err)
    }
})
```

## Contributing
1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## History
TODO: 

## Credits
TODO: Write credits

## License
TODO: MIT
