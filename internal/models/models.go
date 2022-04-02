package models

type LogLine struct {
	Source string `json:"source"`
	Line   string `json:"line" binding:"required"`
}

type LogLines struct {
	Logs []LogLine `json:"logs" binding:"required"`
}

func NewLogLine(source, line string) *LogLine {
	return &LogLine{source, line}
}
