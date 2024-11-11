package entity_test

type RelationshipDTO struct {
	// required
	Id_      *string `json:"id"`
	Created_ *int64  `json:"created"`
	Updated_ *int64  `json:"updated"`
}

type DTOEntityBase struct {
	Id_      *string `json:"id"`
	Created_ *int64  `json:"created"`
	Updated_ *int64  `json:"updated"`
	DTOEntityBasePartial
}

type DTOEntityBasePartial struct {
	String_  *string `json:"string,omitempty" read_method:"String"`
	Int_     *int    `json:"number,omitempty"`
	Boolean_ *bool   `json:"boolean,omitempty"`
	Slice_   []int   `json:"slice"`
	Inlist_  *string `json:"inlist,omitempty"`

	Hasone_  *RelationshipDTO   `json:"hasone,omitempty"`
	Hasmany_ []*RelationshipDTO `json:"hasmany,omitempty"`
}

type DTOEntity struct {
	DTOEntityBase
	Name_ *string `json:"name,omitempty"`
}
