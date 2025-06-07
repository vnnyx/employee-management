package resourceful

type Response[IDType, DataModel any] struct {
	Metadata any                        `json:"metadata"`
	Result   *Result[IDType, DataModel] `json:"result"`
}

type Metadata struct {
	TotalCount int    `json:"total_count"`
	NextCursor string `json:"next_cursor,omitempty"`
}
