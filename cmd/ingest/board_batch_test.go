package main

import (
	"testing"

	"github.com/strelov1/freehire/internal/ingest/sources"
)

func TestSliceBoards(t *testing.T) {
	makeBoards := func(n int) []sources.CompanyEntry {
		boards := make([]sources.CompanyEntry, n)
		for i := range boards {
			boards[i] = sources.CompanyEntry{Board: "b" + string(rune('0'+i))}
		}
		return boards
	}

	t.Run("limit zero leaves it unchanged", func(t *testing.T) {
		boards := makeBoards(5)
		got := sliceBoards(boards, 1, 0)
		if len(got) != len(boards) {
			t.Fatalf("len(got) = %d, want %d", len(got), len(boards))
		}
	})

	t.Run("offset beyond length gives empty", func(t *testing.T) {
		boards := makeBoards(3)
		got := sliceBoards(boards, 99, 5)
		if len(got) != 0 {
			t.Fatalf("len(got) = %d, want 0", len(got))
		}
	})

	t.Run("offset plus limit beyond length is clamped", func(t *testing.T) {
		boards := makeBoards(5)
		got := sliceBoards(boards, 3, 10)
		if len(got) != 2 {
			t.Fatalf("len(got) = %d, want 2", len(got))
		}
		if got[0].Board != "b3" || got[1].Board != "b4" {
			t.Fatalf("got = %#v, want b3,b4", got)
		}
	})

	t.Run("normal slice", func(t *testing.T) {
		boards := makeBoards(7)
		got := sliceBoards(boards, 2, 3)
		want := []string{"b2", "b3", "b4"}
		if len(got) != len(want) {
			t.Fatalf("len(got) = %d, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i].Board != want[i] {
				t.Fatalf("got[%d] = %q, want %q", i, got[i].Board, want[i])
			}
		}
	})
}

func TestParseBoardBatchEnv(t *testing.T) {
	cases := []struct {
		name string
		env  string
		want int
	}{
		{"unset uses default", "", 0},
		{"valid integer", "15", 15},
		{"non-numeric falls back to default", "abc", 0},
		{"negative falls back to default", "-3", 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseBoardBatchEnv("INGEST_BOARD_LIMIT", tc.env, 0)
			if got != tc.want {
				t.Fatalf("parseBoardBatchEnv(%q) = %d, want %d", tc.env, got, tc.want)
			}
		})
	}
}

