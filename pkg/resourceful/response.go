package resourceful

type ResourceResponse[IDType, Model any] struct {
	RequestID string              `json:"request_id"`
	Metadata  any                 `json:"metadata"`
	Data      Data[IDType, Model] `json:"data"`
}

type Data[IDType, Model any] struct {
	PaginatedResults []Model  `json:"paginated_results"`
	IDs              []IDType `json:"ids"`
}
