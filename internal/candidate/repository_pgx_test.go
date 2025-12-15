package candidate

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func Test_scanJSON_ToleratesTypeMismatchObjectVsArray(t *testing.T) {
	var m Media
	if err := scanJSON([]byte("[]"), &m); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func Test_scanJSON_ToleratesTypeMismatchArrayVsObject(t *testing.T) {
	var programs []MainProgram
	if err := scanJSON([]byte("{}"), &programs); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func Test_scanJSON_WrapsSingleStringIntoSlice(t *testing.T) {
	var missions []string
	if err := scanJSON([]byte("\"Mission 1\""), &missions); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(missions) != 1 || missions[0] != "Mission 1" {
		t.Fatalf("unexpected missions: %#v", missions)
	}
}

func Test_isUndefinedColumn(t *testing.T) {
	err := &pgconn.PgError{
		Code:       "42703",
		ColumnName: "photo_media_id",
		Message:    "column \"photo_media_id\" does not exist",
	}
	if !isUndefinedColumn(err, "photo_media_id") {
		t.Fatalf("expected true")
	}
	if isUndefinedColumn(err, "other") {
		t.Fatalf("expected false")
	}
}
