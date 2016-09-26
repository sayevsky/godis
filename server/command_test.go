package server

import (
	"bufio"
	"github.com/sayevsky/godis/internal"
	"reflect"
	"strings"
	"testing"
)

func TestParseGet(t *testing.T) {
	in := "GET\r\n6\r\ngolang\r\n"

	reader := strings.NewReader(in)
	want := &internal.Get{"golang", internal.BaseCommand{false, make(chan internal.Response)}}
	got, err := internal.ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *internal.Get:
		if got.Key != want.Key || got.GetBaseCommand().IsAsync != got.GetBaseCommand().IsAsync {
			t.Errorf("Commands to the same", in, got, want)
		}
	default:
		t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseBadCommand(t *testing.T) {
	in := "\n"

	reader := strings.NewReader(in)
	_, err := internal.ParseCommand(bufio.NewReader(reader))
	if err == nil {
		t.Errorf("error detected as expected")
	}

}

func TestParseSet(t *testing.T) {
	in := "SET\r\n6\r\ngolang\r\n8\r\n@awesome\r\n10ns\r\n0\r\n"
	reader := strings.NewReader(in)
	want := &internal.SetUpd{"golang", "awesome", 10, false, internal.BaseCommand{false, make(chan internal.Response)}}
	got, err := internal.ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure", err)
		t.FailNow()
	}

	switch got := got.(type) {
	case *internal.SetUpd:
		if got.Key != want.Key || got.Value != want.Value ||
			got.Duration != want.Duration || got.Update != want.Update ||
			got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}
	default:
		t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseSetArray(t *testing.T) {
	in := "SET\r\n6\r\ngolang\r\n0\r\n*2\r\n7\r\nawesome\r\n3\r\nyep\r\n10ns\r\n0\r\n"
	reader := strings.NewReader(in)
	array := []string{"awesome", "yep"}
	want := &internal.SetUpd{"golang", array, 10, false, internal.BaseCommand{false, make(chan internal.Response)}}
	got, err := internal.ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure", err)
		t.FailNow()
	}

	switch got := got.(type) {
	case *internal.SetUpd:

		if got.Key != want.Key || !reflect.DeepEqual(got.Value, want.Value) ||
			got.Duration != want.Duration || got.Update != want.Update ||
			got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}
	default:
		t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseSetDict(t *testing.T) {
	in := "SET\r\n6\r\ngolang\r\n0\r\n>2\r\n7\r\nawesome\r\n3\r\nyep\r\n3\r\nabc\r\n1\r\na\r\n10ns\r\n0\r\n"
	reader := strings.NewReader(in)
	dict := make(map[string]string)
	dict["awesome"] = "yep"
	dict["abc"] = "a"
	want := &internal.SetUpd{"golang", dict, 10, false, internal.BaseCommand{false, make(chan internal.Response)}}
	got, err := internal.ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure", err)
		t.FailNow()
	}

	switch got := got.(type) {
	case *internal.SetUpd:
		equals := reflect.DeepEqual(got.Value, want.Value)

		if got.Key != want.Key || !equals ||
			got.Duration != want.Duration || got.Update != want.Update ||
			got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}
	default:
		t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseUpdate(t *testing.T) {
	in := "UPD\r\n6\r\ngolang\r\n8\r\n@awesome\r\n10ns\r\n1\r\n"
	reader := strings.NewReader(in)
	want := &internal.SetUpd{"golang", "awesome", 10, true, internal.BaseCommand{true, make(chan internal.Response)}}
	got, err := internal.ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure", err)
		t.FailNow()
	}

	switch got := got.(type) {
	case *internal.SetUpd:
		if got.Key != want.Key || got.Value != want.Value ||
			got.Duration != want.Duration || got.Update != want.Update ||
			got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}
	default:
		t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseDelete(t *testing.T) {
	in := "DEL\r\n6\r\ngolang\r\n0\r\n"
	reader := strings.NewReader(in)
	want := &internal.Del{"golang", internal.BaseCommand{false, make(chan internal.Response)}}
	got, err := internal.ParseCommand(bufio.NewReader(reader))
	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *internal.Del:
		if got.Key != want.Key || got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}
	default:
		t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseKeys(t *testing.T) {
	in := "KEYS\r\n2\r\n.*\r\n"
	reader := strings.NewReader(in)
	want := &internal.Keys{".*", internal.BaseCommand{false, make(chan internal.Response)}}
	got, err := internal.ParseCommand(bufio.NewReader(reader))

	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *internal.Keys:
		if got.Pattern != want.Pattern || got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}

	default:
		t.Errorf("Commands aren't the same!", in, got, want)
	}
}

func TestParseCount(t *testing.T) {
	in := "COUNT\r\n"
	reader := strings.NewReader(in)
	want := &internal.Count{internal.BaseCommand{false, make(chan internal.Response)}}
	got, err := internal.ParseCommand(bufio.NewReader(reader))

	if err != nil {
		t.Errorf("parsed with failure")
	}

	switch got := got.(type) {
	case *internal.Count:
		if got.GetBaseCommand().IsAsync != want.GetBaseCommand().IsAsync {
			t.Errorf("Commands aren't the same!", in, got, want)
		}

	default:
		t.Errorf("Commands aren't the same!", in, got, want)
	}
}
