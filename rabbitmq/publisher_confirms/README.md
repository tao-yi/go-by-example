## Enabling Publisher Confirms on a Channel

Publisher confirms are a RabbitMQ extension to the AMQP 0.9.1 protocol, so they are not enabled by default. Publisher confirms are enabled at the channel level with the `confirmSelect` method:

```go
ch, err := conn.Channel()
failOnError(err, "Failed to open a channel")
defer ch.Close()
// Confirm puts this channel into confirm mode so that the client can ensure all publishings have successfully been received by the server.
ch.Confirm(false)
```
