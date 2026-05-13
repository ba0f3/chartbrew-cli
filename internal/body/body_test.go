package body

import (
	"os"
	"strings"
	"testing"
)

func TestReadInlineJSON(t *testing.T) {
	got, err := ReadJSON(Source{Data: `{"name":"demo"}`})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"name":"demo"}` {
		t.Fatalf("body = %s", got)
	}
}

func TestReadJSONFile(t *testing.T) {
	path := t.TempDir() + "/body.json"
	t.Cleanup(func() {})
	if err := osWriteFile(path, `{"name":"file"}`); err != nil {
		t.Fatal(err)
	}
	got, err := ReadJSON(Source{DataFile: path})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"name":"file"}` {
		t.Fatalf("body = %s", got)
	}
}

func TestReadJSONFromStdin(t *testing.T) {
	got, err := ReadJSON(Source{DataFile: "-", Stdin: strings.NewReader(`{"name":"stdin"}`)})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"name":"stdin"}` {
		t.Fatalf("body = %s", got)
	}
}

func TestRejectInvalidJSON(t *testing.T) {
	if _, err := ReadJSON(Source{Data: `{`}); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func TestRequireExactlyOneBodySource(t *testing.T) {
	if _, err := ReadJSON(Source{}); err == nil {
		t.Fatal("expected missing body source error")
	}
	if _, err := ReadJSON(Source{Data: `{}`, DataFile: "-"}); err == nil {
		t.Fatal("expected multiple body source error")
	}
}

func osWriteFile(path, data string) error {
	return os.WriteFile(path, []byte(data), 0o600)
}
