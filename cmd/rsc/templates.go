package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	ts "github.com/na4ma4/go-timestring"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// basicFunctions are the set of initial functions provided to every template.
func basicFunctions(extra ...template.FuncMap) template.FuncMap {
	o := template.FuncMap{
		"json": func(v interface{}) string {
			buf := &bytes.Buffer{}
			enc := json.NewEncoder(buf)
			enc.SetEscapeHTML(false)
			_ = enc.Encode(v) //nolint:errchkjson // template function
			// Remove the trailing new line added by the encoder
			return strings.TrimSpace(buf.String())
		},
		"split":    strings.Split,
		"join":     strings.Join,
		"title":    cases.Title(language.English).String,
		"lower":    cases.Lower(language.English).String,
		"upper":    cases.Upper(language.English).String,
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

	if len(extra) > 0 {
		for _, add := range extra {
			for k, v := range add {
				o[k] = v
			}
		}
	}

	return o
}

// padToLength adds whitespace to pad to the supplied length.
func padToMaxLength(source interface{}) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", 0), source)
}

// padToLength adds whitespace to pad to the supplied length.
func padToLength(source interface{}, prefix int) string {
	return fmt.Sprintf(fmt.Sprintf("%%-%ds", prefix), source)
}

// padWithSpace adds whitespace to the input if the input is non-empty.
func padWithSpace(source interface{}, prefix, suffix int) string {
	src := fmt.Sprintf("%s", source)

	if src == "" {
		return src
	}

	return strings.Repeat(" ", prefix) + src + strings.Repeat(" ", suffix)
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

func stringBool(source interface{}, yes, no string) string {
	switch val := source.(type) {
	case *bool:
		if val != nil {
			return stringBool(*val, yes, no)
		}

		return "nil"
	case bool:
		return stringBool(val, yes, no)
	default:
		return fmt.Sprintf("%s", val)
	}
}

// stringTrueFalse returns "true" or "false" for boolean input.
func stringTrueFalse(source interface{}) string {
	return stringBool(source, "true", "false")
}

// stringYesNo returns "yes" or "no" for boolean input.
func stringYesNo(source bool) string {
	return stringBool(source, "yes", "no")
}

// stringTab returns a tab character.
func stringTab() string {
	return "\t"
}

// truncateWithLength truncates the source string up to the length provided by the input.
func truncateWithLength(source interface{}, length int) string {
	src := fmt.Sprintf("%s", source)

	if len(src) < length {
		return src
	}

	return src[:length]
}
