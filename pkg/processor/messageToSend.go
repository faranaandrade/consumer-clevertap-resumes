package processor

type MessageToSend[T any] struct {
	Topic   string `json:"topic" example:"aais-outdated-account-data"`
	Message T      `json:"message" example:"{\"userid\": \"user-test\", \"creationDate\": \"2009-11-10T23:00:00Z\"}"`
}
