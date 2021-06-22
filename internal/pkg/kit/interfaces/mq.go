package interfaces

type MQ interface {
	NewProducer() error
	NewConsumer(topic, channel string) error

	Send(topic, value string) error
	Consume() (value string, err error)
}
