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

func MustTransform[T any, TFrom any](splice []TFrom, transform func(value TFrom) T) []T {
	transformed, _ := Transform(splice, func(value TFrom) (T, error) { return transform(value), nil })
	return transformed
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

func Reduce[T any, TFrom any](slice []TFrom, initial T, predicate func(index int, lhs T, rhs TFrom) T) T {
	for index, el := range slice {
		initial = predicate(index, initial, el)
	}

	return initial
}

func All[T any](slice []T, predicate func(index int, el T) bool) bool {
	for index, el := range slice {
		if !predicate(index, el) {
			return false
		}
	}

	return true
}

func Any[T any](slice []T, predicate func(index int, el T) bool) bool {
	for index, el := range slice {
		if predicate(index, el) {
			return true
		}
	}

	return false
}

func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}
