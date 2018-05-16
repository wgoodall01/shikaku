package shikaku

import (
	"testing"

	"reflect"
)

func TestBoardParseUnevenLengths(t *testing.T) {
	_, err := NewBoardFromString(`
		-- 02 --
		-- -- -- 04
		-- 03 --
	`)

	if err == nil {
		t.Fatal("Didn't blow up parsing a board with uneven row lengths")
	}
}

func TestBoardParseNotIntegers(t *testing.T) {
	_, err := NewBoardFromString(`
		-- nope what
	`)

	if err == nil {
		t.Fatal("Didn't blow up parsing a board with ")
	}
}

func TestBoardParseValid(t *testing.T) {
	b, err := NewBoardFromString(`
		-- -- 05 -- --
		-- 04 -- -- --
		03 02 -- -- --
		-- -- -- 06 --
		-- 05 -- -- --
	`)

	if err != nil {
		t.Fatal("Failed to parse valid board")
	}

	t1 := b.Grid[0][2].(Given).Area == 5
	t2 := b.Grid[2][1].(Given).Area == 2
	t3 := reflect.TypeOf(b.Grid[4][4]) == reflect.TypeOf(Blank{})
	t4 := reflect.TypeOf(b.Grid[0][0]) == reflect.TypeOf(Blank{})

	if !(t1 && t2 && t3 && t4) {
		t.Fatal("Didn't parse valid board correctly.")
	}

	t.Log(b.String())
}
