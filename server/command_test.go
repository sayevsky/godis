package server

import (
	"testing"
	"strings"
	"bufio"
	"reflect"
)
func TestParseGet(t *testing.T) {
	in := "GET\r\n6\r\ngolang\r\n"

	reader := strings.NewReader(in)
	want := &Get{"golang", BaseCommand{false,  make(chan Response)}}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *Get:
		if got.Key != want.Key || got.GetBaseCommand().IsAsync != got.GetBaseCommand().IsAsync {
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
	in := "SET\r\n6\r\ngolang\r\n8\r\n@awesome\r\n10ns\r\n0\r\n"
	reader := strings.NewReader(in)
	want := &SetUpd{"golang", "awesome", 10, false, BaseCommand{false, make(chan Response)}}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure", err)
		t.FailNow()
	}

	switch got := got.(type) {
	case *SetUpd:
		if got.Key != want.Key || got.Value != want.Value ||
			got.duration != want.duration || got.update != want.update ||
			got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}
	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseSetArray(t *testing.T) {
	in := "SET\r\n6\r\ngolang\r\n0\r\n*2\r\n7\r\nawesome\r\n3\r\nyep\r\n10ns\r\n0\r\n"
	reader := strings.NewReader(in)
	array := []string{"awesome", "yep"}
	want := &SetUpd{"golang", array, 10, false, BaseCommand{false, make(chan Response)}}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure", err)
		t.FailNow()
	}

	switch got := got.(type) {
	case *SetUpd:

		if got.Key != want.Key || !reflect.DeepEqual(got.Value, want.Value) ||
			got.duration != want.duration || got.update != want.update ||
			got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}
	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseSetDict(t *testing.T) {
	in := "SET\r\n6\r\ngolang\r\n0\r\n>2\r\n7\r\nawesome\r\n3\r\nyep\r\n3\r\nabc\r\n1\r\na\r\n10ns\r\n0\r\n"
	reader := strings.NewReader(in)
	dict := make(map[string]string)
	dict["awesome"] = "yep"
	dict["abc"] = "a"
	want := &SetUpd{"golang", dict, 10, false, BaseCommand{false, make(chan Response)}}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure", err)
		t.FailNow()
	}

	switch got := got.(type) {
	case *SetUpd:
		equals := reflect.DeepEqual(got.Value, want.Value)

		if got.Key != want.Key || !equals ||
			got.duration != want.duration || got.update != want.update ||
			got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}
	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseUpdate(t *testing.T) {
	in := "UPD\r\n6\r\ngolang\r\n8\r\n@awesome\r\n10ns\r\n1\r\n"
	reader := strings.NewReader(in)
	want := &SetUpd{"golang", "awesome", 10, true, BaseCommand{true, make(chan Response)}}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure", err)
		t.FailNow()
	}

	switch got := got.(type) {
	case *SetUpd:
		if got.Key != want.Key || got.Value != want.Value ||
			got.duration != want.duration || got.update != want.update ||
			got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync{
			t.Errorf("Commands aren't the same!", in, got, want)
		}
	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}



func TestParseDelete(t *testing.T) {
	in := "DEL\r\n6\r\ngolang\r\n0\r\n"
	reader := strings.NewReader(in)
	want := &Del{"golang", BaseCommand{false, make(chan Response)}}
	got, err := ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *Del:
		if got.Key != want.Key || got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}
	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseKeys(t *testing.T) {
	in := "KEYS\r\n2\r\n.*\r\n"
	reader := strings.NewReader(in)
	want := &Keys{".*",BaseCommand{false, make(chan Response)}}
	got, err := ParseCommand(bufio.NewReader(reader))

	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *Keys:
		if got.Pattern != want.Pattern || got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}

	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseCount(t *testing.T) {
	in := "COUNT\r\n"
	reader := strings.NewReader(in)
	want := &Count{BaseCommand{false, make(chan Response)}}
	got, err := ParseCommand(bufio.NewReader(reader))

	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *Count:
		if got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}

	default: t.Errorf("Commands aren't the same!", in, got, want)
	}
}