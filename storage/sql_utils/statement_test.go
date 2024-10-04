package sql_utils_test

import (
	"github.com/darksubmarine/torpedo-lib-go/storage/sql_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertStatementFromDMO(t *testing.T) {

	dmo := fullDMO
	statement := sql_utils.InsertStatementFromDMO("sqlite", "comments", &dmo, nil)

	expectedStatement := "INSERT INTO comments (id,created,updated,string,number,boolean,slice,custom) VALUES (:id,:created,:updated,:string,:number,:boolean,:slice,:custom)"
	assert.EqualValues(t, expectedStatement, statement)
}

func TestUpdateStatementFromDMO(t *testing.T) {
	dmo := fullDMO
	statement := sql_utils.UpdateStatementFromDMO("sqlite", "comments", &dmo, nil)

	expectedStatement := "UPDATE comments SET created = :created, updated = :updated, string = :string, number = :number, boolean = :boolean, slice = :slice, custom = :custom WHERE id = :id"
	assert.EqualValues(t, expectedStatement, statement)
}

func TestInsertStatementFromDMO_DriverName_MySQL(t *testing.T) {

	dmo := fullDMO
	statement := sql_utils.InsertStatementFromDMO("mysql", "comments", &dmo, nil)

	expectedStatement := "INSERT INTO comments (`id`,`created`,`updated`,`string`,`number`,`boolean`,`slice`,`custom`) VALUES (:id,:created,:updated,:string,:number,:boolean,:slice,:custom)"
	assert.EqualValues(t, expectedStatement, statement)
}

func TestUpdateStatementFromDMO_DriverName_MySQL(t *testing.T) {
	dmo := fullDMO
	statement := sql_utils.UpdateStatementFromDMO("mysql", "comments", &dmo, nil)

	expectedStatement := "UPDATE comments SET `created` = :created, `updated` = :updated, `string` = :string, `number` = :number, `boolean` = :boolean, `slice` = :slice, `custom` = :custom WHERE id = :id"
	assert.EqualValues(t, expectedStatement, statement)
}
