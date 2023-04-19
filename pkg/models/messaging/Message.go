package messaging

type Message[TData any] struct {
	Data TData `json:"data"`
}
