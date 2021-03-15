package pkg

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func setup() (*DB, sqlmock.Sqlmock) {
	conn, mock, _ := sqlmock.New()
	db, _ := NewDBWithMockForTest(false, conn)
	return db, mock
}

type user struct {
	ID   int64
	Name string
}

func TestDB_QueryOne(t *testing.T) {
	db, mock := setup()
	sql := "SELECT id, name FROM `test` WHERE name = (?)"
	rows := mock.NewRows([]string{"id", "name"}).
		AddRow(1, "TEST")
	mock.
		ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(1).
		WillReturnRows(rows)
	var u user
	err := db.QueryOne(&u, sql, 1)
	assert.Nil(t, err)
	if assert.NotEmpty(t, u) {
		assert.Equal(t, int64(1), u.ID)
		assert.Equal(t, "TEST", u.Name)
	}
}

func TestDB_QueryMore(t *testing.T) {
	db, mock := setup()
	sql := "SELECT id, name FROM `test` WHERE id BETWEEN ? AND ?;"
	id := []int64{1, 2, 3}
	name := []string{"FIRST", "Second", "Third"}
	mockRows := mock.NewRows([]string{"id", "name"}).
		AddRow(id[0], name[0]).
		AddRow(id[1], name[1]).
		AddRow(id[2], name[2])
	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(1, 3).
		WillReturnRows(mockRows)
	rows, err := db.QueryMore(sql, 1, 3)
	assert.Nil(t, err)
	var index = 0
	for rows.Next() {
		var u user
		err := db.ScanRows(rows, &u)
		assert.Nil(t, err)
		if assert.NotEmpty(t, u) {
			assert.Equal(t, id[index], u.ID)
			assert.Equal(t, name[index], u.Name)
			index++
		}
	}
}
