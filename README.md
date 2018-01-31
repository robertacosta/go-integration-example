# go-integration-example

Examples of integrating with dependent web services using REST wrapped with Hystrix and sending events using Kafka.

Meetup Walkthrough

Call example on the meetup service
`curl 127.0.0.1:50000/example`

this will make a request to the meetuptest service at the endpoint
`curl 127.0.0.1:50010/message`

when a message is returned, it is immediately pushed onto a Kafka topic.

The meetupworker listens for messages on the topic and will simple log the message.
The meetupworker then marks the offset, acknowledging it received the message and thus not consuming this message later

Testing Hystrix:
1 - With meetup and meetuptest started, make a request to /example
2 - Should get the message "Hello Go Meetup"
3 - Shutdown the meetuptest service with `docker-compose stop meetuptest`
4 - Make more requests to /example
5 - Should get the message "Keep Calm and Eat Pizza"
6 - Restart the meetuptest service with `docker-compose restart meetuptest`
7 - Within 10 seconds, make more requests to /example
8 - Should expect to see "Keep Calm and Eat Pizza" until the circuit closes at which point "Hello Go Meetup"