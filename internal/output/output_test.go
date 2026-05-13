package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, JSON, map[string]string{"version": "dev"}); err != nil {
		t.Fatal(err)
	}
	var got map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["version"] != "dev" {
		t.Fatalf("version = %q", got["version"])
	}
}

func TestWriteRaw(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, Raw, json.RawMessage(`{"ok":true}`)); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) != `{"ok":true}` {
		t.Fatalf("raw output = %q", buf.String())
	}
}

func TestWriteErrorEnvelope(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteError(&buf, 1, "example"); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) != `{"error":true,"code":1,"message":"example"}` {
		t.Fatalf("error output = %q", buf.String())
	}
}
