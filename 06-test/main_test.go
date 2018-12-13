package main

import "testing"

const EXPECTED_CORRECT = 15
const EXPECTED_WRONG = 16

func Test_sum(t *testing.T) {
	result := sum([]int{1, 2, 3, 4, 5})
	if result != EXPECTED_CORRECT {
		t.Errorf("Expected %d, got %d", EXPECTED_CORRECT, result)
	}
}

func Test_sum_failure(t *testing.T) {
	result := sum([]int{1, 2, 3, 4, 5})
	if result != EXPECTED_WRONG {
		t.Errorf("Expected %d, got %d", EXPECTED_WRONG, result)
	}
}
