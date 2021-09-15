package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	ts "github.com/na4ma4/go-timestring"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// basicFunctions are the set of initial functions provided to every template.
//nolint:gochecknoglobals // for reusability and mostly because it was copy/pasted from docker/cli
var basicFunctions = template.FuncMap{
	"json": func(v interface{}) string {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		_ = enc.Encode(v)
		// Remove the trailing new line added by the encoder
		return strings.TrimSpace(buf.String())
	},
	"split":    strings.Split,
	"join":     strings.Join,
	"title":    strings.Title,
	"lower":    strings.ToLower,
	"upper":    strings.ToUpper,
	"pad":      padWithSpace,
	"padlen":   padToLength,
	"padmax":   padToMaxLength,
	"truncate": truncateWithLength,
	"tf":       stringTrueFalse,
	"yn":       stringYesNo,
	"t":        stringTab,
	"age":      humanAgeFormat,
	"time":     timeFormat,
	"date":     dateFormat,
}

// padToLength adds whitespace to pad to the supplied length.
func padToMaxLength(source string) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", 0), source)
}

// padToLength adds whitespace to pad to the supplied length.
func padToLength(source string, prefix int) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", prefix), source)
}

// padWithSpace adds whitespace to the input if the input is non-empty.
func padWithSpace(source string, prefix, suffix int) string {
	if source == "" {
		return source
	}

	return strings.Repeat(" ", prefix) + source + strings.Repeat(" ", suffix)
}

// humanAgeFormat returns a duration in a human readable format.
func humanAgeFormat(source interface{}) string {
	switch s := source.(type) {
	case time.Time:
		return ts.LongProcess.Option(ts.Abbreviated, ts.ShowMSOnSeconds).String(time.Since(s))
	case timestamppb.Timestamp:
		return ts.LongProcess.Option(ts.Abbreviated, ts.ShowMSOnSeconds).String(time.Since(s.AsTime()))
	case *timestamppb.Timestamp:
		return ts.LongProcess.Option(ts.Abbreviated, ts.ShowMSOnSeconds).String(time.Since(s.AsTime()))
	default:
		return fmt.Sprintf("%s", source)
	}
}

// // ageFormat returns time in seconds ago.
// func ageFormat(source interface{}) string {
// 	switch s := source.(type) {
// 	case time.Time:
// 		return fmt.Sprintf("%0.2f", time.Since(s).Seconds())
// 		// return s.Format(time.RFC3339)
// 	case timestamppb.Timestamp:
// 		return fmt.Sprintf("%0.2f", time.Since(s.AsTime()).Seconds())
// 	case *timestamppb.Timestamp:
// 		return fmt.Sprintf("%0.2f", time.Since(s.AsTime()).Seconds())
// 	default:
// 		return fmt.Sprintf("%s", source)
// 	}
// }

// timeFormat returns time in RFC3339 format.
func timeFormat(source interface{}) string {
	switch s := source.(type) {
	case time.Time:
		return s.Format(time.RFC3339)
	case timestamppb.Timestamp:
		return s.AsTime().Format(time.RFC3339)
	case *timestamppb.Timestamp:
		return s.AsTime().Format(time.RFC3339)
	default:
		return fmt.Sprintf("%s", source)
	}
}

// dateFormat returns date in YYYY-MM-DD format.
func dateFormat(source interface{}) string {
	switch s := source.(type) {
	case time.Time:
		return s.Format("2006-01-02")
	case timestamppb.Timestamp:
		return s.AsTime().Format("2006-01-02")
	case *timestamppb.Timestamp:
		return s.AsTime().Format("2006-01-02")
	default:
		return fmt.Sprintf("%q", source)
	}
}

// stringTrueFalse returns "true" or "false" for boolean input.
func stringTrueFalse(source bool) string {
	if source {
		return "true"
	}

	return "false"
}

// stringYesNo returns "yes" or "no" for boolean input.
func stringYesNo(source bool) string {
	if source {
		return "yes"
	}

	return "no"
}

// stringTab returns a tab character.
func stringTab() string {
	return "\t"
}

// truncateWithLength truncates the source string up to the length provided by the input.
func truncateWithLength(source string, length int) string {
	if len(source) < length {
		return source
	}

	return source[:length]
}
