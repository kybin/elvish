package newedit

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/elves/elvish/eval"
	"github.com/elves/elvish/styled"
)

func TestCallPrompt_ConvertsValueOutput(t *testing.T) {
	testCallPrompt(t, "put PROMPT", styled.Unstyled("PROMPT"), false)
	testCallPrompt(t, "styled PROMPT red",
		styled.Transform(styled.Unstyled("PROMPT"), "red"), false)
}

func TestCallPrompt_ErrorsOnInvalidValueOutput(t *testing.T) {
	testCallPrompt(t, "put good; put [bad]", styled.Unstyled("good"), true)
}

func TestCallPrompt_ErrorsOnException(t *testing.T) {
	testCallPrompt(t, "fail error", nil, true)
}

func TestCallPrompt_ConvertsBytesOutput(t *testing.T) {
	testCallPrompt(t, "print PROMPT", styled.Unstyled("PROMPT"), false)
}

func testCallPrompt(t *testing.T, fsrc string, want styled.Text, wantErr bool) {
	ev := eval.NewEvaler()
	ev.EvalSource(eval.NewScriptSource(
		"[t]", "[t]", fmt.Sprintf("f = { %s }", fsrc)))
	f := ev.Global["f"].Get().(eval.Callable)
	ed := NewEditor(os.Stdin, os.Stdout, ev)

	content := callPrompt(ed.core, ev, f)
	if !reflect.DeepEqual(content, want) {
		t.Errorf("get prompt result %v, want %v", content, want)
	}

	notes := ed.core.State.Raw.Notes
	if wantErr {
		if len(notes) == 0 {
			t.Errorf("got no error, want errors")
		}
	} else {
		if len(notes) > 0 {
			t.Errorf("got errors %v, want none", notes)
		}
	}
}
