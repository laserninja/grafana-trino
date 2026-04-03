package trino

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

func makeQuery(from, to time.Time) *sqlutil.Query {
	return &sqlutil.Query{
		TimeRange: backend.TimeRange{
			From: from,
			To:   to,
		},
	}
}

func TestMacroTimeFrom(t *testing.T) {
	from := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)
	to := time.Date(2023, 1, 15, 11, 30, 0, 0, time.UTC)
	q := makeQuery(from, to)

	got, err := macroTimeFrom(q, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "TIMESTAMP '2023-01-15 10:30:00'"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMacroTimeTo(t *testing.T) {
	from := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)
	to := time.Date(2023, 1, 15, 11, 30, 0, 0, time.UTC)
	q := makeQuery(from, to)

	got, err := macroTimeTo(q, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "TIMESTAMP '2023-01-15 11:30:00'"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMacroTimeFilter(t *testing.T) {
	from := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 6, 2, 0, 0, 0, 0, time.UTC)
	q := makeQuery(from, to)

	got, err := macroTimeFilter(q, []string{"created_at"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "created_at BETWEEN TIMESTAMP '2023-06-01 00:00:00' AND TIMESTAMP '2023-06-02 00:00:00'"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMacroTimeFilter_WithFormat(t *testing.T) {
	from := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 6, 2, 0, 0, 0, 0, time.UTC)
	q := makeQuery(from, to)

	got, err := macroTimeFilter(q, []string{"created_at", "'yyyy-MM-dd HH:mm:ss'"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "TIMESTAMP created_at BETWEEN TIMESTAMP '2023-06-01 00:00:00' AND TIMESTAMP '2023-06-02 00:00:00'"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMacroTimeFilter_NoArgs(t *testing.T) {
	q := makeQuery(time.Now(), time.Now())
	_, err := macroTimeFilter(q, []string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestMacroDateFilter(t *testing.T) {
	from := time.Date(2023, 3, 15, 10, 30, 0, 0, time.UTC)
	to := time.Date(2023, 3, 20, 14, 0, 0, 0, time.UTC)
	q := makeQuery(from, to)

	got, err := macroDateFilter(q, []string{"order_date"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "order_date BETWEEN date '2023-03-15' AND date '2023-03-20'"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMacroDateFilter_WrongArgCount(t *testing.T) {
	q := makeQuery(time.Now(), time.Now())
	_, err := macroDateFilter(q, []string{})
	if err == nil {
		t.Fatal("expected error for 0 args")
	}
	_, err = macroDateFilter(q, []string{"a", "b"})
	if err == nil {
		t.Fatal("expected error for 2 args")
	}
}

func TestMacroUnixEpochFilter(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	q := makeQuery(from, to)

	got, err := macroUnixEpochFilter(q, []string{"epoch_col"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "epoch_col BETWEEN 1672531200 AND 1672617600"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMacroUnixEpochFilter_WrongArgCount(t *testing.T) {
	q := makeQuery(time.Now(), time.Now())
	_, err := macroUnixEpochFilter(q, []string{})
	if err == nil {
		t.Fatal("expected error for 0 args")
	}
	_, err = macroUnixEpochFilter(q, []string{"a", "b"})
	if err == nil {
		t.Fatal("expected error for 2 args")
	}
}

func TestMacroTimeGroup(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)
	q := makeQuery(from, to)

	got, err := macroTimeGroup(q, []string{"created_at", "'1h'"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "FROM_UNIXTIME(FLOOR(TO_UNIXTIME(created_at)/3600)*3600)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMacroTimeGroup_WithFormat(t *testing.T) {
	q := makeQuery(time.Now(), time.Now())

	got, err := macroTimeGroup(q, []string{"ts", "'1d'", "'yyyy-MM-dd HH:mm:ss'"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "FROM_UNIXTIME(FLOOR(TO_UNIXTIME(TIMESTAMP ts)/86400)*86400)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMacroTimeGroup_TooFewArgs(t *testing.T) {
	q := makeQuery(time.Now(), time.Now())
	_, err := macroTimeGroup(q, []string{"col"})
	if err == nil {
		t.Fatal("expected error for 1 arg")
	}
}

func TestMacroTimeGroup_InvalidInterval(t *testing.T) {
	q := makeQuery(time.Now(), time.Now())
	_, err := macroTimeGroup(q, []string{"col", "'invalid'"})
	if err == nil {
		t.Fatal("expected error for invalid interval")
	}
}

func TestMacroUnixEpochGroup(t *testing.T) {
	q := makeQuery(time.Now(), time.Now())

	got, err := macroUnixEpochGroup(q, []string{"epoch_col", "'5m'"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "FROM_UNIXTIME(FLOOR(epoch_col/300)*300)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMacroParseTime(t *testing.T) {
	q := makeQuery(time.Now(), time.Now())

	// Default format
	got, err := macroParseTime(q, []string{"col"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "TIMESTAMP col"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	// Custom format
	got, err = macroParseTime(q, []string{"col", "'yyyy-MM-dd'"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want = "parse_datetime(col,'yyyy-MM-dd')"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMacroParseTime_NoArgs(t *testing.T) {
	q := makeQuery(time.Now(), time.Now())
	_, err := macroParseTime(q, []string{})
	if err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		target, format, want string
	}{
		{"col", "", "col"},
		{"col", "'yyyy-MM-dd HH:mm:ss'", "TIMESTAMP col"},
		{"col", "'yyyy-MM-dd'", "parse_datetime(col,'yyyy-MM-dd')"},
	}

	for _, tc := range tests {
		got := parseTime(tc.target, tc.format)
		if got != tc.want {
			t.Errorf("parseTime(%q, %q) = %q, want %q", tc.target, tc.format, got, tc.want)
		}
	}
}

func TestMacrosMap(t *testing.T) {
	expected := []string{
		"dateFilter", "parseTime", "unixEpochFilter", "timeFilter",
		"timeFrom", "timeGroup", "unixEpochGroup", "timeTo",
	}

	for _, name := range expected {
		if _, ok := macros[name]; !ok {
			t.Errorf("macro %q not found in macros map", name)
		}
	}

	if len(macros) != len(expected) {
		t.Errorf("macros map has %d entries, expected %d", len(macros), len(expected))
	}
}
