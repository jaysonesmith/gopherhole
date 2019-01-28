package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew_Dimensions(t *testing.T) {
	testCases := []struct {
		name string
		x    int
		y    int
		xl   int
		yl   int
	}{
		{name: "0 defaults to 4x4", x: 0, y: 0, xl: 4, yl: 4},
		{name: "1 defaults to 4x4", x: 1, y: 1, xl: 4, yl: 4},
		{name: "2 defaults to 4x4", x: 2, y: 2, xl: 4, yl: 4},
		{name: "3 defaults to 4x4", x: 3, y: 3, xl: 4, yl: 4},
		{name: "4 returns 4x4", x: 4, y: 4, xl: 4, yl: 4},
		{name: "10 returns 10x10", x: 10, y: 10, xl: 10, yl: 10},
		{name: "1000 returns 10x10", x: 1000, y: 1000, xl: 1000, yl: 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual := New(tc.x, tc.y)

			assert.Equal(tt, tc.xl, len(actual))
			assert.Equal(tt, tc.yl, len(actual[0]))
		})
	}
}

func TestDifficulty(t *testing.T) {
	testCases := []struct {
		name     string
		x        int
		y        int
		d        int
		expected int
	}{
		{name: "10x10 Empty", x: 10, y: 10, d: 0, expected: 0},
		{name: "10x10 Easy", x: 10, y: 10, d: 1, expected: 15},
		{name: "10x10 Medium", x: 10, y: 10, d: 2, expected: 25},
		{name: "4x4 Easy", x: 4, y: 4, d: 1, expected: 3},
		{name: "4x4 Medium", x: 4, y: 4, d: 2, expected: 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {

			actual := difficulty(tc.x, tc.y, tc.d)

			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func TestFillBoard(t *testing.T) {
	testCases := []struct {
		name string
		d    int
		b    *Board
	}{
		{
			name: "Empty",
			d:    0,
			b: &Board{
				[]string{" ", " ", " ", " "},
				[]string{" ", " ", " ", " "},
				[]string{" ", " ", " ", " "},
				[]string{" ", " ", " ", " "},
			},
		},
		{
			name: "Easy",
			d:    1,
			b: &Board{
				[]string{" ", " ", "o", "g"},
				[]string{" ", " ", " ", "o"},
				[]string{" ", "o", "g", "g"},
				[]string{" ", " ", " ", " "},
			},
		},
		{
			name: "Medium",
			d:    2,
			b: &Board{
				[]string{" ", " ", "o", "g"},
				[]string{"g", " ", " ", "o"},
				[]string{"o", "o", "g", "g"},
				[]string{" ", " ", " ", " "},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			setRand(1)
			b := New(4, 4)

			fillBoard(tc.d, 1, &b)

			printBoard(&b)
			assert.Equal(tt, tc.b, &b)
		})
	}
}

func TestCanPlace_Boundaries(t *testing.T) {
	testCases := []struct {
		name string
		bd   int
		x    int
		y    int
		ok   bool
	}{
		{name: "Out of upper bounds", bd: 0, x: -1, y: 0, ok: false},
		{name: "Out of right bounds", bd: 4, x: 0, y: 4, ok: false},
		{name: "Out of lower bounds", bd: 4, x: 4, y: 0, ok: false},
		{name: "Out of left bounds", bd: 4, x: 0, y: -1, ok: false},
		{name: "Inside upper bounds", bd: 4, x: 1, y: 0, ok: true},
		{name: "Inside right bounds", bd: 4, x: 0, y: 0, ok: true},
		{name: "Inside lower bounds", bd: 4, x: 3, y: 0, ok: true},
		{name: "Inside left bounds", bd: 4, x: 0, y: 1, ok: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			b := New(tc.bd, tc.bd)

			ok := canPlace(tc.x, tc.y, b)

			assert.Equal(tt, tc.ok, ok)
		})
	}
}

func TestHoleArea(t *testing.T) {
	setRand(1)
	testCases := []struct {
		name string
		x    int
		y    int
		hx   int
		hy   int
		b    Board
	}{
		{
			name: "No available spots",
			x:    0, y: 0, hx: -1, hy: -1,
			b: Board{
				[]string{" ", "o", "o", "o"},
			},
		},
		{
			name: "Can only place to the right",
			x:    0, y: 0, hx: 0, hy: 1,
			b: Board{
				[]string{" ", " ", "o", "o"},
			},
		},
		{
			name: "Can only place downward",
			x:    0, y: 0, hx: 1, hy: 0,
			b: Board{
				[]string{" ", "o", "o", "o"},
				[]string{" ", "o", "o", "o"},
			},
		},
		{
			name: "Can only place upward",
			x:    1, y: 1, hx: 0, hy: 1,
			b: Board{
				[]string{"o", " ", "o", "o"},
				[]string{"o", " ", "o", "o"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			x, y := holeArea(tc.x, tc.y, tc.b)

			assert.Equal(tt, tc.hx, x)
			assert.Equal(tt, tc.hy, y)
		})
	}
}

func TestSurroundingGopher(t *testing.T) {
	testCases := []struct {
		name     string
		x        int
		y        int
		b        Board
		expected bool
	}{
		{
			name: "Middle of board, all empty",
			x:    1,
			y:    1,
			b: Board{
				[]string{" ", " ", " "},
				[]string{" ", " ", " "},
				[]string{" ", " ", " "},
			},
			expected: false,
		},
		{
			name: "Middle of board, one found above",
			x:    1,
			y:    1,
			b: Board{
				[]string{" ", "g", " "},
				[]string{" ", " ", " "},
				[]string{" ", " ", " "},
			},
			expected: true,
		},
		{
			name: "Left edge of board, empty",
			x:    0,
			y:    1,
			b: Board{
				[]string{" ", " ", " "},
				[]string{" ", " ", " "},
				[]string{" ", " ", " "},
			},
			expected: false,
		},
		{
			name: "Left edge of board, gopher below",
			x:    0,
			y:    1,
			b: Board{
				[]string{" ", " ", " "},
				[]string{" ", " ", " "},
				[]string{"g", " ", " "},
			},
			expected: false,
		},
		{
			name: "Bottom-right corner of board, empty",
			x:    2,
			y:    2,
			b: Board{
				[]string{" ", " ", " "},
				[]string{" ", " ", " "},
				[]string{" ", " ", " "},
			},
			expected: false,
		},
		{
			name: "Top-right corner of board, two found",
			x:    0,
			y:    2,
			b: Board{
				[]string{" ", "g", " "},
				[]string{" ", " ", "g"},
				[]string{" ", " ", " "},
			},
			expected: true,
		},
		{
			name: "Lower-left corner of board, one at diagonal",
			x:    2,
			y:    0,
			b: Board{
				[]string{" ", " ", " "},
				[]string{" ", "g", " "},
				[]string{" ", " ", " "},
			},
			expected: true,
		},
		{
			name: "Lower-left corner of board, one out of range",
			x:    2,
			y:    0,
			b: Board{
				[]string{" ", " ", "g"},
				[]string{" ", " ", " "},
				[]string{" ", " ", " "},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual := surroundingGopher(tc.x, tc.y, tc.b)

			assert.Equal(tt, tc.expected, actual)
		})
	}
}
