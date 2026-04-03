package trino

import (
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/gtime"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
)

const (
	timestampFormat   = "'yyyy-MM-dd HH:mm:ss'"
	goTimestampFormat = "2006-01-02 15:04:05"
)

// parseTime wraps a column reference in the appropriate Trino time parsing function.
// If format is empty, returns the column as-is.
// If format matches the default timestamp format, wraps with TIMESTAMP keyword.
// Otherwise, uses Trino's parse_datetime function.
func parseTime(target, format string) string {
	if format == "" {
		return target
	}
	if format == timestampFormat {
		return fmt.Sprintf("TIMESTAMP %s", target)
	}
	return fmt.Sprintf("parse_datetime(%s,%s)", target, format)
}

// parseTimeGroup extracts the interval duration and time variable from macro arguments.
func parseTimeGroup(query *sqlutil.Query, args []string) (time.Duration, string, error) {
	if len(args) < 2 {
		return 0, "", fmt.Errorf("macro $__timeGroup needs time column and interval, got %d args", len(args))
	}

	interval, err := gtime.ParseInterval(strings.Trim(args[1], `'`))
	if err != nil {
		return 0, "", fmt.Errorf("error parsing interval %v: %w", args[1], err)
	}

	timeVar := args[0]
	if len(args) == 3 {
		timeVar = parseTime(args[0], args[2])
	}

	return interval, timeVar, nil
}

func macroTimeGroup(query *sqlutil.Query, args []string) (string, error) {
	interval, timeVar, err := parseTimeGroup(query, args)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("FROM_UNIXTIME(FLOOR(TO_UNIXTIME(%s)/%v)*%v)", timeVar, interval.Seconds(), interval.Seconds()), nil
}

func macroUnixEpochGroup(query *sqlutil.Query, args []string) (string, error) {
	interval, timeVar, err := parseTimeGroup(query, args)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("FROM_UNIXTIME(FLOOR(%s/%v)*%v)", timeVar, interval.Seconds(), interval.Seconds()), nil
}

func macroParseTime(query *sqlutil.Query, args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("macro $__parseTime expected at least one argument")
	}

	column := args[0]
	timeFormat := timestampFormat

	if len(args) == 2 {
		timeFormat = args[1]
	}

	return parseTime(column, timeFormat), nil
}

func macroTimeFilter(query *sqlutil.Query, args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("macro $__timeFilter expected at least one argument")
	}

	column := args[0]
	timeFormat := ""
	from := query.TimeRange.From.UTC().Format(goTimestampFormat)
	to := query.TimeRange.To.UTC().Format(goTimestampFormat)

	if len(args) > 1 {
		timeFormat = args[1]
	}
	timeVar := parseTime(column, timeFormat)

	return fmt.Sprintf("%s BETWEEN TIMESTAMP '%s' AND TIMESTAMP '%s'", timeVar, from, to), nil
}

func macroUnixEpochFilter(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("macro $__unixEpochFilter expected one argument, got %d", len(args))
	}

	column := args[0]
	from := query.TimeRange.From.UTC().Unix()
	to := query.TimeRange.To.UTC().Unix()

	return fmt.Sprintf("%s BETWEEN %d AND %d", column, from, to), nil
}

func macroTimeFrom(query *sqlutil.Query, args []string) (string, error) {
	return fmt.Sprintf("TIMESTAMP '%s'", query.TimeRange.From.UTC().Format(goTimestampFormat)), nil
}

func macroTimeTo(query *sqlutil.Query, args []string) (string, error) {
	return fmt.Sprintf("TIMESTAMP '%s'", query.TimeRange.To.UTC().Format(goTimestampFormat)), nil
}

func macroDateFilter(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("macro $__dateFilter expected 1 argument, got %d", len(args))
	}

	column := args[0]
	from := query.TimeRange.From.UTC().Format("2006-01-02")
	to := query.TimeRange.To.UTC().Format("2006-01-02")

	return fmt.Sprintf("%s BETWEEN date '%s' AND date '%s'", column, from, to), nil
}

var macros = sqlutil.Macros{
	"dateFilter":      macroDateFilter,
	"parseTime":       macroParseTime,
	"unixEpochFilter": macroUnixEpochFilter,
	"timeFilter":      macroTimeFilter,
	"timeFrom":        macroTimeFrom,
	"timeGroup":       macroTimeGroup,
	"unixEpochGroup":  macroUnixEpochGroup,
	"timeTo":          macroTimeTo,
}
