# Rabbitmq example
This is an example application that runs a published made in Go and a subscriber made in Java.

The Go publisher is a cli application that receives messages from CLI and sends it to rabbit.

The Java subscriber is a Quarkus application that listen to Rabbit messages and simulates real work for each message.

## TODO
- [x] Dockerize applications
- [ ] Document applications
