package trn

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	sep = "::"
)

var (
	// ErrInvalidStringFormat invalid TRN string format
	ErrInvalidStringFormat = errors.New("invalid TRN string format")
)

// TRN is a Torpedo Resource Name
type TRN struct {
	id   string `json:"id"`
	kind string `json:"kind"`
	name string `json:"name"`
}

// New creates a new TRN instance
func New(kind, name, id string) *TRN {
	return &TRN{kind: kind, name: name, id: id}
}

// NewFromString returns a TRN instance from the given string representation
func NewFromString(trn string) (*TRN, error) {
	partitions := strings.Split(trn, sep)
	if len(partitions) != 4 {
		return nil, ErrInvalidStringFormat
	}

	return &TRN{id: partitions[3], kind: partitions[1], name: partitions[2]}, nil
}

// Kind returns the resource type
func (a *TRN) Kind() string { return a.kind }

// Name returns the resource name
func (a *TRN) Name() string { return a.name }

// ID returns the resource id
func (a *TRN) ID() string { return a.id }

// String returns the TRN string representation
func (a *TRN) String() string {
	return fmt.Sprintf("trn::%s::%s::%s", a.kind, a.name, a.id)
	return fmt.Sprintf("URN:%s:%s@%s", a.kind, a.name, a.id)
}

// Equals returns true if the given TRN is equals
func (a *TRN) Equals(comp TRN) bool {
	if a.ID() == comp.ID() && a.Name() == comp.Name() && a.Kind() == comp.Kind() {
		return true
	}
	return false
}

// JSON returns the TRN json representation
func (a *TRN) JSON() string {

	type ExportedJSONMetadata struct {
		ObjectType   string `json:"objectType"`
		ExportedDate int64  `json:"exportedDateMillis"`
	}

	type ExportedJson struct {
		Metadata ExportedJSONMetadata `json:"metadata"`
		Id       string               `json:"id"`
		Kind     string               `json:"kind"`
		Name     string               `json:"name"`
	}

	toJson := ExportedJson{
		Metadata: ExportedJSONMetadata{
			ObjectType:   "TRN",
			ExportedDate: time.Now().UTC().UnixMilli(),
		},
		Id:   a.id,
		Kind: a.kind,
		Name: a.name,
	}

	data, _ := json.Marshal(toJson)
	return string(data)
}
