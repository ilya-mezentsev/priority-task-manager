package queue

type TaskReceiver interface {
	Add(task Task[any]) error
}

type Task[T any] struct {
	Priority    int
	UUID        string `json:"uuid"`
	AccountHash string `json:"account_hash"`
	Type        string `json:"type"`
	Data        T      `json:"data"`
}
