package server

import "testing"
import "strings"
import "bufio"

func TestParseGet(t *testing.T) {
	in := "GET\r\n6\r\ngolang\r\n"
	reader := strings.NewReader(in)
	want := &Get{"golang"}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *Get:
		if *got != *want {
		t.Errorf("Commands to the same", in, got, want)
	}
	default: t.Errorf("Commands to the same!!!", in, got, want)
	}
}

func TestParseSet(t *testing.T) {
	// SET\r\n<numberOfBytes>\r\n<key>\r\n<numberOfBytes>\r\n<value>\r\n<TTL>\r\n
	in := "SET\r\n6\r\ngolang\r\n7\r\nawesome\r\n10"
	reader := strings.NewReader(in)
	want := &SetUpd{"golang", "awesome", 10, false}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *SetUpd:
		if *got != *want {
			t.Errorf("Commands to the same", in, got, want)
		}
	default: t.Errorf("Commands to the same!!!", in, got, want)
	}
}

func TestParseUpdate(t *testing.T) {
	// SET\r\n<numberOfBytes>\r\n<key>\r\n<numberOfBytes>\r\n<value>\r\n<TTL>\r\n
	in := "UPD\r\n6\r\ngolang\r\n7\r\nawesome\r\n10"
	reader := strings.NewReader(in)
	want := &SetUpd{"golang", "awesome", 10, true}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *SetUpd:
		if *got != *want {
			t.Errorf("Commands to the same", in, got, want)
		}
	default: t.Errorf("Commands to the same!!!", in, got, want)
	}
}

func TestParseDelete(t *testing.T) {
	// SET\r\n<numberOfBytes>\r\n<key>\r\n<numberOfBytes>\r\n<value>\r\n<TTL>\r\n
	in := "DEL\r\n6\r\ngolang\r\n"
	reader := strings.NewReader(in)
	want := &Del{"golang"}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *SetUpd:
		if *got != *want {
			t.Errorf("Commands to the same", in, got, want)
		}
	default: t.Errorf("Commands to the same!!!", in, got, want)
	}
}
