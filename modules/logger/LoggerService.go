package logger

type LoggerService struct {
	logChannel chan string
}

func NewLoggerService() *LoggerService {
	logChannel := make(chan string)

	return &LoggerService{logChannel}
}

func (service *LoggerService) EnableLogging() {
	defer close(service.logChannel)

	for value := range service.logChannel {
		logger.Println(value)
	}
}

func (service *LoggerService) GetChannel() chan string {
	return service.logChannel
}
