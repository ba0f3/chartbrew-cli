package body

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type Source struct {
	Data     string
	DataFile string
	Stdin    io.Reader
}

func ReadJSON(src Source) ([]byte, error) {
	count := 0
	if src.Data != "" {
		count++
	}
	if src.DataFile != "" {
		count++
	}
	if count == 0 {
		return nil, errors.New("one of --data or --data-file is required")
	}
	if count > 1 {
		return nil, errors.New("use only one of --data or --data-file")
	}

	var data []byte
	var err error
	switch {
	case src.Data != "":
		data = []byte(src.Data)
	case src.DataFile == "-":
		if src.Stdin == nil {
			return nil, errors.New("stdin is required when --data-file is -")
		}
		data, err = io.ReadAll(src.Stdin)
	default:
		data, err = os.ReadFile(src.DataFile)
	}
	if err != nil {
		return nil, fmt.Errorf("read request body: %w", err)
	}
	if !json.Valid(data) {
		return nil, errors.New("request body must be valid JSON")
	}
	return data, nil
}
