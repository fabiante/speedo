package app

import (
	"encoding/json"
	"io"
	"time"

	"github.com/simonfrey/jsonl"
)

type LogLine struct {
	Time         time.Time `json:"time"`
	DownloadMbps float64   `json:"downloadMbps"`
	UploadMbps   float64   `json:"uploadMbps"`
}

func DecodeLog(r io.Reader) ([]*LogLine, error) {
	reader := jsonl.NewReader(r)

	var lines []*LogLine

	err := reader.ReadLines(func(data []byte) error {
		var line LogLine
		err := json.Unmarshal(data, &line)
		if err != nil {
			return err
		} else {
			lines = append(lines, &line)
			return nil
		}
	})
	if err != nil {
		return nil, err
	}

	return lines, err
}
