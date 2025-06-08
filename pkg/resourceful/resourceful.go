package resourceful

import (
	"encoding/base64"

	"github.com/goccy/go-json"
	"github.com/vnnyx/employee-management/pkg/optional"
)

type Mode string

const (
	ModeOffset Mode = "offset"
	ModeCursor Mode = "cursor"
)

type Resource[IDType, Model any] struct {
	Parameter *Parameter
	result    *Result[IDType, Model]
	metadata  any
}

type Parameter struct {
	Limit          optional.Int64
	Page           optional.Int64
	Mode           Mode
	Cursor         *Cursor
	additionalData any // This field is used to pass additional data that might be needed for the resource.
}

type Cursor struct {
	Key   string
	Value string
}

type Result[IDType, Model any] struct {
	// CustomInformation *Model // Disabled for now, as it is not used in the current implementation.
	PaginationResult []Model
	IDs              []IDType
}

func NewResource[IDType, Model any](parameter *Parameter) *Resource[IDType, Model] {
	if parameter != nil {
		// Set default values if not provided
		if !parameter.Limit.IsPresent() || parameter.Limit.MustGet() == 0 {
			parameter.Limit = optional.NewInt64(1)
		}
		if !parameter.Page.IsPresent() || parameter.Page.MustGet() == 0 {
			parameter.Page = optional.NewInt64(1)
		}
	}

	return &Resource[IDType, Model]{
		Parameter: parameter,
		result:    &Result[IDType, Model]{},
	}
}

func (r *Resource[IDType, Model]) SetResult(result Result[IDType, Model]) {
	r.result = &result
}

func (r *Resource[IDType, Model]) SetMetadata(metadata any) {
	r.metadata = metadata
}

func (r *Resource[IDType, Model]) Response() *ResourceResponse[IDType, Model] {
	resp := &ResourceResponse[IDType, Model]{
		Metadata: r.metadata,
		Data: Data[IDType, Model]{
			PaginatedResults: r.result.PaginationResult,
			IDs:              r.result.IDs,
		},
	}

	return resp
}

// Disabled for now, as it is not used in the current implementation.
// func (r *Resource[IDType, Model, Item]) CustomResponse() ResourceResponse {
// 	resp := ResourceResponse{
// 		Metadata: r.metadata,
// 		Data: any(map[string]any{
// 			"ids":               r.result.IDs,
// 			"pagination_result": r.result.PaginationResult,
// 		}),
// 	}

// 	if r.result.CustomInformation != nil {
// 		switch customInfo := any(r.result.CustomInformation).(type) {
// 		case map[string]any:
// 			for k, v := range customInfo {
// 				respData := resp.Data.(map[string]any)
// 				respData[k] = v
// 			}
// 		default:
// 			infoBytes, err := json.Marshal(customInfo)
// 			if err != nil {
// 				panic(err)
// 			}

// 			var infoMap map[string]any
// 			if err := json.Unmarshal(infoBytes, &infoMap); err != nil {
// 				panic(err)
// 			}

// 			for k, v := range infoMap {
// 				respData := resp.Data.(map[string]any)
// 				respData[k] = v
// 			}
// 		}
// 	}

// 	return resp
// }

// Disabled for now, as it is not used in the current implementation.
// func (r *Resource[IDType, Model, Item]) DecodeCustomResponse(resp ResourceResponse) (*Result[IDType, Model, Item], error) {
// 	data, ok := resp.Data.(map[string]any)
// 	if !ok {
// 		return nil, nil
// 	}

// 	result := &Result[IDType, Model, Item]{
// 		IDs:              make([]IDType, 0),
// 		PaginationResult: make([]Item, 0),
// 	}

// 	if ids, exists := data["ids"]; exists {
// 		if idSlice, ok := ids.([]IDType); ok {
// 			result.IDs = idSlice
// 		}
// 	}

// 	if paginationResult, exists := data["pagination_result"]; exists {
// 		if itemSlice, ok := paginationResult.([]Item); ok {
// 			result.PaginationResult = itemSlice
// 		}
// 	}

// 	if customInfo, exists := data["custom_information"]; exists {
// 		if customInfoMap, ok := customInfo.(map[string]any); ok {
// 			var customInfoModel Model
// 			customInfoBytes, err := json.Marshal(customInfoMap)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if err := json.Unmarshal(customInfoBytes, &customInfoModel); err != nil {
// 				return nil, err
// 			}
// 			result.CustomInformation = &customInfoModel
// 		}
// 	}

// 	return result, nil
// }

func DecodeCursor(cursor optional.String) (*Cursor, error) {
	if !cursor.IsPresent() || cursor.MustGet() == "" {
		return nil, nil
	}

	// Decode from base64
	decoded, err := base64.StdEncoding.DecodeString(cursor.MustGet())
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a Cursor struct
	var c Cursor
	if err := json.Unmarshal(decoded, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func EncodeCursor(c *Cursor) (optional.String, error) {
	if c == nil {
		return optional.String{}, nil
	}

	// Marshal the Cursor struct to JSON
	data, err := json.Marshal(c)
	if err != nil {
		return optional.String{}, err
	}

	// Encode the JSON data to base64
	encoded := base64.StdEncoding.EncodeToString(data)

	return optional.NewString(encoded), nil
}

func (p *Parameter) SetAdditionalData(data any) {
	p.additionalData = data
}

func (p *Parameter) GetAdditionalData() any {
	if p.additionalData == nil {
		return nil
	}
	return p.additionalData
}
