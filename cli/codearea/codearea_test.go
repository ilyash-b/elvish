package codearea

import (
	"errors"
	"reflect"
	"testing"

	"github.com/elves/elvish/cli/term"
	"github.com/elves/elvish/edit/ui"
	"github.com/elves/elvish/styled"
)

var bb = ui.NewBufferBuilder

var renderTests = []struct {
	name    string
	widget  *Widget
	width   int
	height  int
	wantBuf *ui.BufferBuilder
}{
	{
		"prompt only",
		&Widget{State: State{
			Prompt: styled.MakeText("~>", "bold")}},
		10, 24,
		bb(10).WriteStringSGR("~>", "1").SetDotToCursor(),
	},
	{
		"rprompt only",
		&Widget{State: State{
			RPrompt: styled.MakeText("RP", "inverse")}},
		10, 24,
		bb(10).SetDotToCursor().WriteSpacesSGR(8, "").WriteStringSGR("RP", "7"),
	},
	{
		"code only with dot at beginning",
		&Widget{State: State{
			CodeBuffer: CodeBuffer{Content: "code", Dot: 0}}},
		10, 24,
		bb(10).SetDotToCursor().WritePlain("code"),
	},
	{
		"code only with dot at middle",
		&Widget{State: State{
			CodeBuffer: CodeBuffer{Content: "code", Dot: 2}}},
		10, 24,
		bb(10).WritePlain("co").SetDotToCursor().WritePlain("de"),
	},
	{
		"code only with dot at end",
		&Widget{State: State{
			CodeBuffer: CodeBuffer{Content: "code", Dot: 4}}},
		10, 24,
		bb(10).WritePlain("code").SetDotToCursor(),
	},
	{
		"prompt, code and rprompt",
		&Widget{State: State{
			Prompt:     styled.Plain("~>"),
			CodeBuffer: CodeBuffer{Content: "code", Dot: 4},
			RPrompt:    styled.Plain("RP"),
		}},
		10, 24,
		bb(10).WritePlain("~>code").SetDotToCursor().WritePlain("  RP"),
	},
	{
		"rprompt too long",
		&Widget{State: State{
			Prompt:     styled.Plain("~>"),
			CodeBuffer: CodeBuffer{Content: "code", Dot: 4},
			RPrompt:    styled.Plain("1234"),
		}},
		10, 24,
		bb(10).WritePlain("~>code").SetDotToCursor(),
	},
	{
		"highlighted code",
		&Widget{
			State: State{CodeBuffer: CodeBuffer{Content: "code", Dot: 4}},
			Highlighter: func(code string) (styled.Text, []error) {
				return styled.MakeText(code, "bold"), nil
			},
		},
		10, 24,
		bb(10).WriteStringSGR("code", "1").SetDotToCursor(),
	},
	{
		"static errors in code",
		&Widget{
			State: State{CodeBuffer: CodeBuffer{Content: "code", Dot: 4}},
			Highlighter: func(code string) (styled.Text, []error) {
				err := errors.New("static error")
				return styled.Plain(code), []error{err}
			},
		},
		10, 24,
		bb(10).WritePlain("code").SetDotToCursor().
			Newline().WritePlain("static error"),
	},
	{
		"pending code inserting at the dot",
		&Widget{State: State{
			CodeBuffer:  CodeBuffer{Content: "code", Dot: 4},
			PendingCode: PendingCode{From: 4, To: 4, Content: "x"},
		}},
		10, 24,
		bb(10).WritePlain("code").WriteStringSGR("x", "4").SetDotToCursor(),
	},
	{
		"pending code replacing at the dot",
		&Widget{State: State{
			CodeBuffer:  CodeBuffer{Content: "code", Dot: 2},
			PendingCode: PendingCode{From: 2, To: 4, Content: "x"},
		}},
		10, 24,
		bb(10).WritePlain("co").WriteStringSGR("x", "4").SetDotToCursor(),
	},
	{
		"pending code to the left of the dot",
		&Widget{State: State{
			CodeBuffer:  CodeBuffer{Content: "code", Dot: 4},
			PendingCode: PendingCode{From: 1, To: 3, Content: "x"},
		}},
		10, 24,
		bb(10).WritePlain("c").WriteStringSGR("x", "4").WritePlain("e").SetDotToCursor(),
	},
	{
		"pending code to the right of the cursor",
		&Widget{State: State{
			CodeBuffer:  CodeBuffer{Content: "code", Dot: 1},
			PendingCode: PendingCode{From: 2, To: 3, Content: "x"},
		}},
		10, 24,
		bb(10).WritePlain("c").SetDotToCursor().WritePlain("o").
			WriteStringSGR("x", "4").WritePlain("e"),
	},
	{
		"ignore invalid pending code",
		&Widget{State: State{
			CodeBuffer:  CodeBuffer{Content: "code", Dot: 4},
			PendingCode: PendingCode{From: 2, To: 1, Content: "x"},
		}},
		10, 24,
		bb(10).WritePlain("code").SetDotToCursor(),
	},
	{
		"prioritize lines before the cursor with small height",
		&Widget{State: State{
			CodeBuffer: CodeBuffer{Content: "a\nb\nc\nd", Dot: 3},
		}},
		10, 2,
		bb(10).WritePlain("a").Newline().WritePlain("b").SetDotToCursor(),
	},
	{
		"show only the cursor line when height is 1",
		&Widget{State: State{
			CodeBuffer: CodeBuffer{Content: "a\nb\nc\nd", Dot: 3},
		}},
		10, 1,
		bb(10).WritePlain("b").SetDotToCursor(),
	},
	{
		"show lines after the cursor when all lines before the cursor are shown",
		&Widget{State: State{
			CodeBuffer: CodeBuffer{Content: "a\nb\nc\nd", Dot: 3},
		}},
		10, 3,
		bb(10).WritePlain("a").Newline().WritePlain("b").SetDotToCursor().
			Newline().WritePlain("c"),
	},
}

func TestRender(t *testing.T) {
	for _, test := range renderTests {
		t.Run(test.name, func(t *testing.T) {
			buf := test.widget.Render(test.width, test.height)
			wantBuf := test.wantBuf.Buffer()
			if !reflect.DeepEqual(buf, wantBuf) {
				t.Errorf("got buf %v, want %v", buf, wantBuf)
			}
		})
	}
}

var handleTests = []struct {
	name      string
	widget    *Widget
	events    []term.Event
	wantState State
}{
	{
		"simple inserts",
		&Widget{},
		[]term.Event{term.K('c'), term.K('o'), term.K('d'), term.K('e')},
		State{CodeBuffer: CodeBuffer{Content: "code", Dot: 4}},
	},
	{
		"unicode inserts",
		&Widget{},
		[]term.Event{term.K('你'), term.K('好')},
		State{CodeBuffer: CodeBuffer{Content: "你好", Dot: 6}},
	},
	{
		"unterminated paste",
		&Widget{},
		[]term.Event{term.PasteSetting(true), term.K('"'), term.K('x')},
		State{},
	},
	{
		"literal paste",
		&Widget{},
		[]term.Event{
			term.PasteSetting(true),
			term.K('"'), term.K('x'),
			term.PasteSetting(false)},
		State{CodeBuffer: CodeBuffer{Content: "\"x", Dot: 2}},
	},
	{
		"literal paste swallowing functional keys",
		&Widget{},
		[]term.Event{
			term.PasteSetting(true),
			term.K('a'), term.K(ui.F1), term.K('b'),
			term.PasteSetting(false)},
		State{CodeBuffer: CodeBuffer{Content: "ab", Dot: 2}},
	},
	{
		"quoted paste",
		&Widget{QuotePaste: func() bool { return true }},
		[]term.Event{
			term.PasteSetting(true),
			term.K('"'), term.K('x'),
			term.PasteSetting(false)},
		State{CodeBuffer: CodeBuffer{Content: "'\"x'", Dot: 4}},
	},
	{
		"backspace at end of code",
		&Widget{},
		[]term.Event{
			term.K('c'), term.K('o'), term.K('d'), term.K('e'),
			term.K(ui.Backspace)},
		State{CodeBuffer: CodeBuffer{Content: "cod", Dot: 3}},
	},
	{
		"backspace at middle of buffer",
		&Widget{State: State{
			CodeBuffer: CodeBuffer{Content: "code", Dot: 2}}},
		[]term.Event{term.K(ui.Backspace)},
		State{CodeBuffer: CodeBuffer{Content: "cde", Dot: 1}},
	},
	{
		"backspace at beginning of buffer",
		&Widget{State: State{
			CodeBuffer: CodeBuffer{Content: "code", Dot: 0}}},
		[]term.Event{term.K(ui.Backspace)},
		State{CodeBuffer: CodeBuffer{Content: "code", Dot: 0}},
	},
	{
		"backspace deleting unicode character",
		&Widget{},
		[]term.Event{
			term.K('你'), term.K('好'), term.K(ui.Backspace)},
		State{CodeBuffer: CodeBuffer{Content: "你", Dot: 3}},
	},
	{
		"abbreviation expansion",
		&Widget{
			Abbreviations: func(f func(abbr, full string)) {
				f("dn", "/dev/null")
			},
		},
		[]term.Event{term.K('d'), term.K('n')},
		State{CodeBuffer: CodeBuffer{Content: "/dev/null", Dot: 9}},
	},
	{
		"abbreviation expansion preferring longest",
		&Widget{
			Abbreviations: func(f func(abbr, full string)) {
				f("n", "none")
				f("dn", "/dev/null")
			},
		},
		[]term.Event{term.K('d'), term.K('n')},
		State{CodeBuffer: CodeBuffer{Content: "/dev/null", Dot: 9}},
	},
	{
		"abbreviation expansion interrupted by function key",
		&Widget{
			Abbreviations: func(f func(abbr, full string)) {
				f("dn", "/dev/null")
			},
		},
		[]term.Event{term.K('d'), term.K(ui.F1), term.K('n')},
		State{CodeBuffer: CodeBuffer{Content: "dn", Dot: 2}},
	},
}

func TestHandle(t *testing.T) {
	for _, test := range handleTests {
		t.Run(test.name, func(t *testing.T) {
			w := test.widget
			for _, event := range test.events {
				w.Handle(event)
			}
			if !reflect.DeepEqual(w.State, test.wantState) {
				t.Errorf("got state %v, want %v", w.State, test.wantState)
			}
		})
	}
}

var unhandledEvents = []term.Event{
	// Mouse events are unhandled
	term.MouseEvent{},
	// Function keys are unhandled (except Backspace)
	term.K(ui.F1),
	term.K('X', ui.Ctrl),
}

func TestHandle_UnhandledEvents(t *testing.T) {
	w := &Widget{}
	for _, event := range unhandledEvents {
		handled := w.Handle(event)
		if handled {
			t.Errorf("event %v got handled", event)
		}
	}
}

func TestHandle_AbbreviationExpansionInterruptedByExternalMutation(t *testing.T) {
	w := &Widget{
		Abbreviations: func(f func(abbr, full string)) {
			f("dn", "/dev/null")
		},
	}
	w.Handle(term.K('d'))
	w.MutateCodeAreaState(func(s *State) { s.CodeBuffer.InsertAtDot("d") })
	w.Handle(term.K('n'))
	wantState := State{CodeBuffer: CodeBuffer{Content: "ddn", Dot: 3}}
	if !reflect.DeepEqual(w.State, wantState) {
		t.Errorf("got state %v, want %v", w.State, wantState)
	}
}

func TestHandle_EnterEmitsSubmit(t *testing.T) {
	var submitted string
	w := &Widget{
		State:    State{CodeBuffer: CodeBuffer{Content: "code", Dot: 4}},
		OnSubmit: func(code string) { submitted = code },
	}
	w.Handle(term.K('\n'))
	if submitted != "code" {
		t.Errorf("code not submitted")
	}
}

func TestHandle_DefaultNoopSubmit(t *testing.T) {
	w := &Widget{State: State{
		CodeBuffer: CodeBuffer{Content: "code", Dot: 4}}}
	w.Handle(term.K('\n'))
	// No panic, we are good
}

func TestState(t *testing.T) {
	w := &Widget{}
	w.MutateCodeAreaState(func(s *State) { s.CodeBuffer.Content = "code" })
	if w.CopyState().CodeBuffer.Content != "code" {
		t.Errorf("state not mutated")
	}
}
