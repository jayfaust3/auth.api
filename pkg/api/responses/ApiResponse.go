package responses

type ApiResponse[TData any] struct {
	Data TData
}
