package utility

import (
	"math"
)

/*
	The Levenshtein algorithm, also known as Levenshtein distance or edit distance, is a metric for
	measuring the difference between two sequences (usually strings). It calculates the minimum number
	of single-character edits required to change one string into the other. The allowable edits are insertion,
	deletion, or substitution of a single character.
*/

func levenshteinDistance(a, b string) int {
	aLen := len(a)
	bLen := len(b)
	if aLen == 0 {
		return bLen
	}
	if bLen == 0 {
		return aLen
	}

	// Initialize the matrix
	matrix := make([][]int, aLen+1)
	for i := range matrix {
		matrix[i] = make([]int, bLen+1)
	}

	// Fill the base case values
	for i := 0; i <= aLen; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= bLen; j++ {
		matrix[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= aLen; i++ {
		for j := 1; j <= bLen; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			matrix[i][j] = int(math.Min(math.Min(
				float64(matrix[i-1][j]+1),  // Deletion
				float64(matrix[i][j-1]+1)), // Insertion
				float64(matrix[i-1][j-1]+cost), // Substitution
			))
		}
	}

	return matrix[aLen][bLen]
}

// Function to calculate the percentage of changed characters based on Levenshtein distance
func percentageChanged(prev, next string) float64 {
	editDistance := levenshteinDistance(prev, next)
	longerLen := math.Max(float64(len(prev)), float64(len(next)))
	percentage := (float64(editDistance) / longerLen) * 100
	return percentage
}

// Function to check if at least 20% of the characters have changed
func AtLeast20PercentChanged(prev, next string) bool {
	percentChanged := percentageChanged(prev, next)
	return percentChanged >= 20
}
