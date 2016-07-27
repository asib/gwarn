package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	filepathlib "path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	out = buf
	src := `package main

	//:warning test

	func main() {
	//:warning test2
	}`

	f, err := ioutil.TempFile("", "testParseFile.go")
	if err != nil {
		t.Fatal("could not create temp file:", err)
	}
	defer os.Remove(f.Name())
	if _, err = f.Write([]byte(src)); err != nil {
		t.Fatal("could not write src to temp file:", err)
	}
	if err = parseFile(f.Name()); err != nil {
		t.Fatal("could not parse file:", err)
	}
	if err = f.Close(); err != nil {
		t.Fatal("could not close temp file:", err)
	}

	warn1 := fmt.Sprintf("%s:3: test", f.Name())
	warn2 := fmt.Sprintf("%s:6: test2", f.Name())

	// Check `buf`
	if l, err := buf.ReadBytes('\n'); err != nil {
		t.Error("could not read from buffer:", err)
	} else if trimmed := bytes.TrimSpace(l); !bytes.Equal(trimmed, []byte(warn1)) {
		t.Errorf("expected %q, got %q\n", warn1, string(trimmed))
	}
	if l, err := buf.ReadBytes('\n'); err != nil {
		t.Error("could not read from buffer:", err)
	} else if trimmed := bytes.TrimSpace(l); !bytes.Equal(trimmed, []byte(warn2)) {
		t.Errorf("expected %q, got %q\n", warn2, string(trimmed))
	}
}

func TestParseFileWontParseCommentsInStrings(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	out = buf
	src := `package main

	func main() {
		test := ` + "`" + `
			this is a nested multiline string with a
			//:warning inside
		` + "`" + `
	}`

	f, err := ioutil.TempFile("", "testParseFile")
	if err != nil {
		t.Fatal("could not create temp file:", err)
	}
	defer os.Remove(f.Name())
	if _, err = f.Write([]byte(src)); err != nil {
		t.Fatal("could not write src to temp file:", err)
	}
	if err = parseFile(f.Name()); err != nil {
		t.Fatal("could not parse file:", err)
	}
	if err = f.Close(); err != nil {
		t.Fatal("could not close temp file:", err)
	}

	if _, err = buf.ReadBytes('\n'); err == nil {
		t.Error("should return error")
	} else if err != io.EOF {
		t.Error("should return io.EOF error")
	}
}

func TestParseDir(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	out = buf
	src1 := `package main

	func main() {
		//:warning test
	}`
	src2 := `package main

	func helper() {
		//:warning test2
	}`

	dir, err := ioutil.TempDir("", "testParseDir")
	if err != nil {
		t.Fatal("could not create temp dir:", err)
	}
	defer os.RemoveAll(dir)

	tmpfn1 := filepathlib.Join(dir, "1.go")
	tmpfn2 := filepathlib.Join(dir, "2.go")
	if err := ioutil.WriteFile(tmpfn1, []byte(src1), 0666); err != nil {
		t.Fatal("could not write temp file 1:", err)
	}
	if err := ioutil.WriteFile(tmpfn2, []byte(src2), 0666); err != nil {
		t.Fatal("could not write temp file 2:", err)
	}

	if err = parseDir(dir); err != nil {
		t.Fatal("could not parse directory:", err)
	}

	warn1 := fmt.Sprintf("%s:4: test", tmpfn1)
	warn2 := fmt.Sprintf("%s:4: test2", tmpfn2)

	if l, err := buf.ReadBytes('\n'); err != nil {
		t.Error("should not return an error:", err)
	} else if trimmed := bytes.TrimSpace(l); !bytes.Equal(trimmed, []byte(warn1)) {
		t.Errorf("expected %s, got %s\n", warn1, string(trimmed))
	}
	if l, err := buf.ReadBytes('\n'); err != nil {
		t.Error("should not return an error:", err)
	} else if trimmed := bytes.TrimSpace(l); !bytes.Equal(trimmed, []byte(warn2)) {
		t.Errorf("expected %s, got %s\n", warn2, string(trimmed))
	}
}
