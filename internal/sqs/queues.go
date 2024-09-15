package sqs

type QueueNamesStruct struct {
	PickKeywords string
}

var QueueNames = QueueNamesStruct{
	PickKeywords: "pick-keywords",
}
