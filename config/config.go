package config

import "strings"

// Config type
type Config struct {
	Verbosity       int
	ReportInterval  int
	PushEndpoint    string
	MetricsEndpoint string
	Metrics         string
}

// StringToSlice return slice from "x,y,z"
func StringToSlice(in string) []string {
	rez := []string{}
	for _, val := range strings.Split(in, ",") {
		rez = append(rez, val)
	}

	return rez
}
