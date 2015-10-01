package php

import (
	"os"
	"path"
	"strconv"
	"testing"
)

var testDir string

var execTests = []struct {
	file     string // Filename to run
	expected string // Expected output
}{
	{"echo.php", "Hello World"},
}

var bindTests = []struct {
	value    interface{} // Value to bind
	expected string      // Serialized form of value
}{
	{42, "i:42;"},
	{3.14159, "d:3.1415899999999999;"},
	{true, "b:1;"},
	{"Such bind", `s:9:"Such bind";`},
	{[]string{"this", "that"}, `a:2:{i:0;s:4:"this";i:1;s:4:"that";}`},
	{[][]int{[]int{1, 2}, []int{3, 4}}, `a:2:{i:0;a:2:{i:0;i:1;i:1;i:2;}i:1;a:2:{i:0;i:3;i:1;i:4;}}`},
}

type MockWriter struct {
	buffer []byte
}

func (m *MockWriter) Write(p []byte) (int, error) {
	if m.buffer == nil {
		m.buffer = p
	} else {
		m.buffer = append(m.buffer, p...)
	}

	return len(p), nil
}

func (m *MockWriter) String() string {
	if m.buffer == nil {
		return ""
	}

	return string(m.buffer)
}

func (m *MockWriter) Reset() {
	if m.buffer != nil {
		m.buffer = m.buffer[:0]
	}
}

func TestContextExec(t *testing.T) {
	var w MockWriter

	e, _ := New()
	ctx, _ := e.NewContext(&w)

	defer e.Destroy()
	defer ctx.Destroy()

	for _, tt := range execTests {
		file := path.Join(testDir, tt.file)
		if err := ctx.Exec(file); err != nil {
			t.Errorf("Context.Exec(%s): %s", tt.file, err)
		}

		actual := w.String()
		w.Reset()

		if actual != tt.expected {
			t.Errorf("Context.Exec(%s): expected '%s', actual '%s'", tt.file, tt.expected, actual)
		}
	}
}

func TestContextBind(t *testing.T) {
	var w MockWriter

	e, _ := New()
	ctx, _ := e.NewContext(&w)

	defer e.Destroy()
	defer ctx.Destroy()

	for i, tt := range bindTests {
		if err := ctx.Bind(strconv.FormatInt(int64(i), 10), tt.value); err != nil {
			t.Errorf("Context.Bind(%v): %s", tt.value, err)
		}

		ctx.Exec(path.Join(testDir, "bind.php"))

		actual := w.String()
		w.Reset()

		if actual != tt.expected {
			t.Errorf("Context.Bind(%v): expected '%s', actual '%s'", tt.value, tt.expected, actual)
		}
	}
}

func init() {
	wd, _ := os.Getwd()
	testDir = path.Join(wd, ".tests")
}
