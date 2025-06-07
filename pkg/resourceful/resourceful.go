package resourceful

type Resource[IDType, DataModel any] struct {
	CursorParameter CursorParameter
	result          *Result[IDType, DataModel]
	metadata        any
}

type CursorParameter struct {
	Cursor  any
	Limit   int
	OrderBy string
}

type MetadataParameter struct {
	TotalCount int
	NextCursor string
}

type Result[IDType, DataModel any] struct {
	PaginationResult []DataModel
}

func GetMetadata(parameter MetadataParameter) Metadata {
	return Metadata(parameter)
}

func NewResource[IDType, DataModel any](cursorParameter *CursorParameter) *Resource[IDType, DataModel] {
	if cursorParameter == nil {
		cursorParameter = &CursorParameter{
			Cursor: "",
			Limit:  0,
		}
	}

	return &Resource[IDType, DataModel]{
		CursorParameter: *cursorParameter,
		result:          &Result[IDType, DataModel]{},
	}
}

func (r *Resource[IDType, DataModel]) SetMetadata(metadata any) {
	r.metadata = metadata
}

func (r *Resource[IDType, DataModel]) SetResult(result []DataModel) {
	if r.result == nil {
		r.result = &Result[IDType, DataModel]{}
	}

	r.result.PaginationResult = result
}

func (r *Resource[IDType, DataModel]) GetResult() *Result[IDType, DataModel] {
	if r.result == nil {
		r.result = &Result[IDType, DataModel]{}
	}

	return r.result
}

func (r *Resource[IDType, DataModel]) Metadata() *Metadata {
	if r.metadata == nil {
		return &Metadata{}
	}

	metadata, ok := r.metadata.(*Metadata)
	if !ok {
		return &Metadata{}
	}

	return metadata
}

func (r *Resource[IDType, DataModel]) Response() *Response[IDType, DataModel] {
	return &Response[IDType, DataModel]{
		Metadata: r.metadata,
		Result:   r.GetResult(),
	}
}

func (r *Resource[IDType, DataModel]) SetDefaultCursor(cursor any) {
	if r.CursorParameter.Cursor == nil || r.CursorParameter.Cursor == "" {
		r.CursorParameter.Cursor = cursor
	}
	if r.CursorParameter.Limit <= 0 {
		r.CursorParameter.Limit = 10 // Default limit
	}
	if r.CursorParameter.OrderBy == "" {
		r.CursorParameter.OrderBy = "created_at" // Default order by
	}
}
