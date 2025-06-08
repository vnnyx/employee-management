package docshelper

import "github.com/vnnyx/employee-management/pkg/resourceful"

type Response[IDType, Model any, Metadata any] struct {
	RequestID string                          `json:"request_id" extensions:"x-order=0"`
	Metadata  Metadata                        `json:"metadata" extensions:"x-order=1"`
	Data      resourceful.Data[IDType, Model] `json:"data" extensions:"x-order=2"`
}
