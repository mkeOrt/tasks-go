package database

import "testing"

func TestNewSqliteDB(t *testing.T) {
	t.Parallel()
	sql, err := NewSqliteDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if sql == nil {
		t.Fatal("expected sqlite db")
	}
}
