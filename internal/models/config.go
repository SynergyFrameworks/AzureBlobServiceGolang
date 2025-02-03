package models

type Config struct {
	Server struct {
		Port int
		Host string
	}
	Azure struct {
		AccountName   string
		AccountKey    string
		ContainerName string
	}
	Kafka struct {
		Brokers       []string
		ConsumerGroup string
		Topics struct {
			StorageEvents string
		}
		Producer struct {
			RequiredAcks int
			Compression  int
			Retries     int
		}
	}
}