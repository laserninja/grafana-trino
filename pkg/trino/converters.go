package trino

import (
	"regexp"

	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

// converters returns the type converters for mapping Trino SQL types to Grafana data types.
func converters() []sqlutil.Converter {
	nullStringConverter := sqlutil.NullStringConverter
	nullStringConverter.InputTypeRegex = regexp.MustCompile("char|varchar|varbinary|json|interval year to month|interval day to second|decimal|ipaddress|unknown")

	nullDecimalConverter := sqlutil.NullDecimalConverter
	nullDecimalConverter.InputTypeRegex = regexp.MustCompile("real|double")

	nullInt64Converter := sqlutil.NullInt64Converter
	nullInt64Converter.InputTypeRegex = regexp.MustCompile("tinyint|smallint|integer|bigint")

	nullTimeConverter := sqlutil.NullTimeConverter
	nullTimeConverter.InputTypeRegex = regexp.MustCompile("date|time|time with time zone|timestamp|timestamp with time zone")

	nullBoolConverter := sqlutil.NullBoolConverter
	nullBoolConverter.InputTypeName = "boolean"

	return []sqlutil.Converter{
		nullStringConverter,
		nullDecimalConverter,
		nullInt64Converter,
		nullTimeConverter,
		nullBoolConverter,
	}
}
