package bowling

import (
	"reflect"
	"strings"
	"testing"
)

func TestPointsForGame(t *testing.T) {
	tests := []struct {
		name string
		game string
		want int
	}{
		{"Normal numbers", "1234", 10},
		{"A spare", "1/", 10},
		{"A strike", "X", 10},
		{"A spare and a normal round", "1/5-", 20},
		{"Strike and normal", "X5-", 20},
		{"Double and normal", "XX5-", 45},
		{"Full game all normal", "12345123451234512345", 60},
		{"Normal game", "9-9-9-9-9-9-9-9-9-9-", 90},
		{"All spares", "5/5/5/5/5/5/5/5/5/5/5", 150},
		{"Double", "XX", 30},
		{"Turkey", "XXX", 60},
		{"10 strikes and 1 gutter", "XXXXXXXXXX--", 270},
		{"9 strikes, 1 spare, 1 strike", "XXXXXXXXX9/X", 279},
		{"Perfect game", "XXXXXXXXXXXX", 300},
	}
	for _, tt := range tests {
		got := pointsForGame(tt.game)
		if got != tt.want {
			t.Errorf("%v, got : %v, want : %v", tt.name, got, tt.want)
		}
	}
}

func TestGameToSlice(t *testing.T) {
	tests := []struct {
		game string
		want []string
	}{
		{"1234", []string{"1", "2", "3", "4"}},
		{"XXXX", []string{"X", "X", "X", "X"}},
		{"9-9-", []string{"9", "-", "9", "-"}},
		{"5/5/", []string{"5", "/", "5", "/"}},
	}

	for _, tt := range tests {
		got := gameToSlice(tt.game)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got : %v, want : %v", got, tt.want)
		}
	}
}

func TestCheckStrike(t *testing.T) {
	tests := []struct {
		throw string
		want  bool
	}{
		{"X", true},
		{"-", false},
		{"1", false},
		{"/", false},
	}

	for _, tt := range tests {
		got := checkStrike(tt.throw)
		if got != tt.want {
			t.Errorf("got : %v, want : %v", got, tt.want)
		}
	}
}

func TestCheckSpare(t *testing.T) {
	tests := []struct {
		throw string
		want  bool
	}{
		{"54", false},
		{"1/", true},
	}

	for _, tt := range tests {
		arr := strings.Split(tt.throw, "")
		got := checkSpare(arr[0], arr[1])
		if got != tt.want {
			t.Errorf("got : %v, want : %v", got, tt.want)
		}
	}
}

func TestCountPointsForRound(t *testing.T) {
	tests := []struct {
		throw1 string
		throw2 string
		want   int
	}{
		{"1", "2", 3},
		{"9", "-", 9},
		{"5", "/", 10},
	}
	for _, tt := range tests {
		got := pointsForRound(tt.throw1, tt.throw2)
		if got != tt.want {
			t.Errorf("got : %v, want : %v", got, tt.want)
		}
	}
}

func TestPointsForThrow(t *testing.T) {
	tests := []struct {
		throw string
		want  int
	}{
		{"1", 1},
		{"-", 0},
	}
	for _, tt := range tests {
		got := pointsForThrow(tt.throw)
		if got != tt.want {
			t.Errorf("got : %v, want : %v", got, tt.want)
		}
	}
}

func TestHandleSpare(t *testing.T) {
	tests := []struct {
		game string
		want int
	}{
		{"5/5-", 15},
		{"5/X", 20},
	}

	for _, tt := range tests {
		arr := strings.Split(tt.game, "")
		got := handleSpare(arr)
		if got != tt.want {
			t.Errorf("got : %v, want : %v", got, tt.want)
		}
	}
}

func TestHandleStrike(t *testing.T) {
	tests := []struct {
		game string
		want int
	}{
		{"X5-", 15},
		{"XX5-", 25},
		{"XXX5-", 30},
		{"XXXX", 30},
		{"X5/", 20},
		{"X5/5-", 20},
		{"XX5/5-", 25},
		{"XXX5/5-", 30},
		{"X5/X--", 20},
	}

	for _, tt := range tests {
		rounds := strings.Split(tt.game, "")
		got := handleStrike(rounds)
		if got != tt.want {
			t.Errorf("got : %v, want : %v, game : %v", got, tt.want, tt.game)
		}
	}
}
