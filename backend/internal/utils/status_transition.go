package utils

import "fmt"

type StatusTransition struct {
	allowed map[string][]string
}

func NewStatusTransition(allowed map[string][]string) *StatusTransition {
	return &StatusTransition{allowed: allowed}
}

// IsValid mengecek apakah transisi from -> to diizinkan.
func (s *StatusTransition) IsValid(from, to string) bool {
	allowed, ok := s.allowed[from]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == to {
			return true
		}
	}
	return false
}

func (s *StatusTransition) Validate(from, to string) error {
	if _, ok := s.allowed[from]; !ok {
		return fmt.Errorf("invalid current status: %s", from)
	}
	if !s.IsValid(from, to) {
		return fmt.Errorf("invalid status transition from %s to %s", from, to)
	}
	return nil
}
