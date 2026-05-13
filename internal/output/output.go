package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type Format string

const (
	JSON     Format = "json"
	Markdown Format = "markdown"
	Raw      Format = "raw"
)

type ErrorEnvelope struct {
	Error   bool   `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Write(w io.Writer, format Format, value any) error {
	switch format {
	case "", JSON:
		data, err := marshalJSON(value, true)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, string(data))
		return err
	case Markdown:
		data, err := marshalJSON(value, true)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "```json\n%s\n```\n", data)
		return err
	case Raw:
		data, err := marshalJSON(value, false)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, string(data))
		return err
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
}

func WriteError(w io.Writer, code int, message string) error {
	data, err := json.Marshal(ErrorEnvelope{Error: true, Code: code, Message: message})
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

func marshalJSON(value any, indent bool) ([]byte, error) {
	if raw, ok := value.(json.RawMessage); ok {
		if !json.Valid(raw) {
			return nil, fmt.Errorf("response is not valid JSON")
		}
		if !indent {
			return raw, nil
		}
		var v any
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, err
		}
		return json.MarshalIndent(v, "", "  ")
	}
	if indent {
		return json.MarshalIndent(value, "", "  ")
	}
	return json.Marshal(value)
}
