package functional

import (
	"sort"
	"strings"
)

func Sorted[T any](slice []T, less func(a, b T) bool) []T {
	sorted := make([]T, len(slice))
	copy(sorted, slice)

	sort.Slice(sorted, func(i, j int) bool {
		return less(sorted[i], sorted[j])
	})

	return sorted
}

func GetLines(value string) []string {
	return strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
}

func CountIf[T any](slice []T, predicate func(val T) bool) int {
	count := 0
	for _, el := range slice {
		if predicate(el) {
			count += 1
		}
	}
	return count
}

func Transform[T any, TFrom any](splice []TFrom, transform func(value TFrom) (T, error)) ([]T, error) {
	transformed := []T{}
	var err error

	for _, el := range splice {
		transformed, err = TransformAppend(transformed, el, transform)
		if err != nil {
			return nil, err
		}
	}

	return transformed, nil
}

func TransformAppend[T any, TFrom any](splice []T, value TFrom, transform func(value TFrom) (T, error)) ([]T, error) {
	transformed, err := transform(value)
	if err != nil {
		return nil, err
	}
	return append(splice, transformed), nil
}

func Reduce[T any](slice []T, initial T, predicate func(index int, lhs, rhs T) T) T {
	for index, el := range slice {
		initial = predicate(index, initial, el)
	}

	return initial
}
