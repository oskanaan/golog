package logdisplay

import "github.com/oskanaan/golog/logreader"

func logreaderConfig(file string, sizes []int) logreader.LogReaderConfig{
	return logreader.LogReaderConfig{
		Files: []logreader.LogFile{{file, "Name"}},
		Seperator: "~",
		Headers: []logreader.Header{
			{"Date", sizes[0]},
			{"Thread", sizes[1]},
			{"Package", sizes[2]},
		},
	}
}

func logdisplayConfig() *LogDisplayConfig{
	return &LogDisplayConfig{
		Severities: []Severity {
			{`\bERROR\b`, []interface{}{1, 1}},
			{`\bWARN\b`, []interface{}{3, 1}},
			{`\bTRACE\b`, []interface{}{6, 5}},
			{`\bINFO\b`, []interface{}{2, 1}},
			{`\bDEBUG\b`, []interface{}{0, 1}},
		},
	}
}

