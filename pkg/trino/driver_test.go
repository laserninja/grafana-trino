package trino

import (
	"testing"
)

func TestParseRoles_Empty(t *testing.T) {
	roles, err := parseRoles("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Errorf("expected empty map, got %v", roles)
	}
}

func TestParseRoles_Whitespace(t *testing.T) {
	roles, err := parseRoles("  \t  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Errorf("expected empty map, got %v", roles)
	}
}

func TestParseRoles_SingleRole(t *testing.T) {
	roles, err := parseRoles("system:admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if roles["system"] != "admin" {
		t.Errorf("got %v, want system:admin", roles)
	}
}

func TestParseRoles_MultipleRoles(t *testing.T) {
	roles, err := parseRoles("system:admin;catalog1:reader;catalog2:writer")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 3 {
		t.Fatalf("expected 3 roles, got %d", len(roles))
	}
	if roles["system"] != "admin" {
		t.Errorf("system role: got %q, want %q", roles["system"], "admin")
	}
	if roles["catalog1"] != "reader" {
		t.Errorf("catalog1 role: got %q, want %q", roles["catalog1"], "reader")
	}
	if roles["catalog2"] != "writer" {
		t.Errorf("catalog2 role: got %q, want %q", roles["catalog2"], "writer")
	}
}

func TestParseRoles_WithSpaces(t *testing.T) {
	roles, err := parseRoles(" system : admin ; catalog1 : reader ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if roles["system"] != "admin" {
		t.Errorf("system role: got %q, want %q", roles["system"], "admin")
	}
	if roles["catalog1"] != "reader" {
		t.Errorf("catalog1 role: got %q, want %q", roles["catalog1"], "reader")
	}
}

func TestParseRoles_InvalidFormat(t *testing.T) {
	_, err := parseRoles("invalid-no-colon")
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestParseRoles_ColonInRole(t *testing.T) {
	roles, err := parseRoles("catalog:role:with:colons")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// SplitN with 2 means everything after first colon is the role
	if roles["catalog"] != "role:with:colons" {
		t.Errorf("got %q, want %q", roles["catalog"], "role:with:colons")
	}
}

func TestConverters(t *testing.T) {
	c := converters()
	if len(c) != 5 {
		t.Fatalf("expected 5 converters, got %d", len(c))
	}

	// Verify the string converter matches expected types
	stringConv := c[0]
	if stringConv.InputTypeRegex == nil {
		t.Fatal("string converter has nil regex")
	}
	for _, typeName := range []string{"char", "varchar", "varbinary", "json", "decimal", "ipaddress"} {
		if !stringConv.InputTypeRegex.MatchString(typeName) {
			t.Errorf("string converter should match %q", typeName)
		}
	}

	// Verify decimal converter
	decimalConv := c[1]
	if decimalConv.InputTypeRegex == nil {
		t.Fatal("decimal converter has nil regex")
	}
	for _, typeName := range []string{"real", "double"} {
		if !decimalConv.InputTypeRegex.MatchString(typeName) {
			t.Errorf("decimal converter should match %q", typeName)
		}
	}

	// Verify int64 converter
	int64Conv := c[2]
	if int64Conv.InputTypeRegex == nil {
		t.Fatal("int64 converter has nil regex")
	}
	for _, typeName := range []string{"tinyint", "smallint", "integer", "bigint"} {
		if !int64Conv.InputTypeRegex.MatchString(typeName) {
			t.Errorf("int64 converter should match %q", typeName)
		}
	}

	// Verify time converter
	timeConv := c[3]
	if timeConv.InputTypeRegex == nil {
		t.Fatal("time converter has nil regex")
	}
	for _, typeName := range []string{"date", "time", "timestamp", "time with time zone", "timestamp with time zone"} {
		if !timeConv.InputTypeRegex.MatchString(typeName) {
			t.Errorf("time converter should match %q", typeName)
		}
	}

	// Verify bool converter
	boolConv := c[4]
	if boolConv.InputTypeName != "boolean" {
		t.Errorf("bool converter: got type name %q, want %q", boolConv.InputTypeName, "boolean")
	}
}
