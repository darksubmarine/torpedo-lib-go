package sql_utils_test

import (
	"github.com/darksubmarine/torpedo-lib-go/storage/sql_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertStatementFromDMO(t *testing.T) {

	dmo := fullDMO
	statement := sql_utils.InsertStatementFromDMO("comments", &dmo)

	expectedStatement := "INSERT INTO comments (id,created,updated,string,number,boolean,slice,custom) VALUES (:id,:created,:updated,:string,:number,:boolean,:slice,:custom)"
	assert.EqualValues(t, expectedStatement, statement)
}

func TestUpdateStatementFromDMO(t *testing.T) {
	dmo := fullDMO
	statement := sql_utils.UpdateStatementFromDMO("comments", &dmo)

	expectedStatement := "UPDATE comments SET created = :created, updated = :updated, string = :string, number = :number, boolean = :boolean, slice = :slice, custom = :custom WHERE id = :id"
	assert.EqualValues(t, expectedStatement, statement)
}
