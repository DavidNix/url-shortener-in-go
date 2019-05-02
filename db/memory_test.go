package db

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDB_Track(t *testing.T) {
	db := NewMemory()
	db.Track("test1")
	require.EqualValues(t, 1, db.Visits("test1"))
	db.Track("test1")
	require.EqualValues(t, 2, db.Visits("test1"))
}
