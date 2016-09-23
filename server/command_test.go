package server

import "testing"
import "strings"
import "bufio"
func TestParseGet(t *testing.T) {
	in := "GET\r\n6\r\ngolang\r\n"

	reader := strings.NewReader(in)
	want := &Get{"golang", BaseCommand{false,  make(chan Resulter)}}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *Get:
		if got.Key != want.Key {
		t.Errorf("Commands to the same", in, got, want)
	}
	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseBadCommand(t *testing.T) {
	in := "\n"

	reader := strings.NewReader(in)
	_, err := ParseCommand(bufio.NewReader(reader))
	if err == nil {
		t.Errorf("error detected as expected")
	}

}

func TestParseSet(t *testing.T) {
	in := "SET\r\n6\r\ngolang\r\n7\r\nawesome\r\n10\r\n"
	reader := strings.NewReader(in)
	want := &SetUpd{"golang", "awesome", 10, false, BaseCommand{false, make(chan Resulter)}}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *SetUpd:
		if got.Key != want.Key && got.Value == want.Value && got.TTL == want.TTL && got.update == want.update {
			t.Errorf("Commands to the same", in, got, want)
		}
	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseUpdate(t *testing.T) {
	in := "UPD\r\n6\r\ngolang\r\n7\r\nawesome\r\n10\r\n"
	reader := strings.NewReader(in)
	want := &SetUpd{"golang", "awesome", 10, true, BaseCommand{false, make(chan Resulter)}}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *SetUpd:
		if got.Key != want.Key && got.Value == want.Value && got.TTL == want.TTL && got.update == want.update {
			t.Errorf("Commands to the same", in, got, want)
		}
	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseDelete(t *testing.T) {
	in := "DEL\r\n6\r\ngolang\r\n"
	reader := strings.NewReader(in)
	want := &Del{"golang", BaseCommand{false, make(chan Resulter)}}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *Del:
		if got.Key != want.Key {
			t.Errorf("Commands to the same", in, got, want)
		}
	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseKeys(t *testing.T) {
	in := "KEYS\r\n"
	reader := strings.NewReader(in)
	want := &Keys{BaseCommand{false, make(chan Resulter)}}
	got, err := ParseCommand(bufio.NewReader(reader))

	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *Keys:

	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}
