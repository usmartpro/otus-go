package app

type NotificationSource interface {
	GetNotificationChannel() (<-chan Notification, error)
}

type NotificationTransport interface {
	String() string
	Send(Notification) error
}

type NotificationSender struct {
	source NotificationSource
	logger Logger
}

func NewSender(
	source NotificationSource,
	logger Logger,
) *NotificationSender {
	return &NotificationSender{source, logger}
}

func (s *NotificationSender) Run() {
	s.logger.Info("! notification. run")
	channel, err := s.source.GetNotificationChannel()
	if err != nil {
		s.logger.Error("! notification. error get from channel: %s", err)
		return
	}

	for notification := range channel {
		s.logger.Info("! notification. %s", notification)
	}
}
