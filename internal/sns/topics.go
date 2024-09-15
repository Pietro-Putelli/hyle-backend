package sns

type TopicNamesStruct struct {
	PushNotification string
}

var TopicNames = TopicNamesStruct{
	PushNotification: "push-notification-topic",
}
