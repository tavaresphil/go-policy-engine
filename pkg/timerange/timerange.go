package timerange

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrZeroStart      = errors.New("start time cannot be zero")
	ErrEndBeforeStart = errors.New("end time cannot be before start")
	ErrSplitOutside   = errors.New("split time outside range")
)

// TimeRange represents a half-open time interval [start, end).
// If end is nil, the range is open-ended.
type TimeRange struct {
	start time.Time
	end   *time.Time
}

// New creates a valid TimeRange enforcing all invariants.
func New(start time.Time, end *time.Time) (*TimeRange, error) {
	if start.IsZero() {
		return nil, ErrZeroStart
	}

	if end != nil && end.Before(start) {
		return nil, ErrEndBeforeStart
	}

	// defensive copy
	var endCopy *time.Time
	if end != nil {
		t := *end
		endCopy = &t
	}

	return &TimeRange{
		start: start,
		end:   endCopy,
	}, nil
}

// MustNew panics if the TimeRange is invalid.
func MustNew(start time.Time, end *time.Time) *TimeRange {
	tr, err := New(start, end)
	if err != nil {
		panic(err)
	}
	return tr
}

// Start returns the inclusive start of the range.
func (tr *TimeRange) Start() time.Time {
	return tr.start
}

// End returns the exclusive end of the range (nil if open-ended).
func (tr *TimeRange) End() *time.Time {
	if tr.end == nil {
		return nil
	}
	t := *tr.end
	return &t
}

// IsOpenEnded reports whether the range has no end.
func (tr *TimeRange) IsOpenEnded() bool {
	return tr.end == nil
}

// IsEmpty reports whether start == end.
func (tr *TimeRange) IsEmpty() bool {
	return tr.end != nil && tr.start.Equal(*tr.end)
}

// Duration returns the duration of the range.
// Open-ended ranges return 0.
func (tr *TimeRange) Duration() time.Duration {
	if tr.end == nil {
		return 0
	}
	return tr.end.Sub(tr.start)
}

// Contains reports whether the given time is within the range.
func (tr *TimeRange) Contains(t time.Time) bool {
	if t.Before(tr.start) {
		return false
	}
	if tr.end == nil {
		return true
	}
	return t.Before(*tr.end)
}

// Overlaps reports whether two ranges intersect.
func (tr *TimeRange) Overlaps(other TimeRange) bool {
	return tr.Intersect(other) != nil
}

// Intersect returns the intersection of two ranges, or nil if none exists.
func (tr *TimeRange) Intersect(other TimeRange) *TimeRange {
	start := maxTime(tr.start, other.start)

	var end *time.Time
	switch {
	case tr.end == nil && other.end == nil:
		end = nil
	case tr.end == nil:
		end = other.end
	case other.end == nil:
		end = tr.end
	default:
		e := minTime(*tr.end, *other.end)
		end = &e
	}

	if end != nil && !start.Before(*end) {
		return nil
	}

	result, _ := New(start, end)
	return result
}

// Clamp restricts the range to the provided bounds.
func (tr *TimeRange) Clamp(bounds TimeRange) *TimeRange {
	return tr.Intersect(bounds)
}

// Split divides the range at the given time.
// Returns left and right ranges.
func (tr *TimeRange) Split(at time.Time) (*TimeRange, *TimeRange, error) {
	if !tr.Contains(at) || at.Equal(tr.start) {
		return nil, nil, ErrSplitOutside
	}

	left, _ := New(tr.start, &at)

	var rightEnd *time.Time
	if tr.end != nil {
		rightEnd = tr.end
	}
	right, _ := New(at, rightEnd)

	return left, right, nil
}

// WithEnd returns a new TimeRange with a different end.
func (tr *TimeRange) WithEnd(end *time.Time) (*TimeRange, error) {
	return New(tr.start, end)
}

// WithStart returns a new TimeRange with a different start.
func (tr *TimeRange) WithStart(start time.Time) (*TimeRange, error) {
	return New(start, tr.end)
}

// Equals reports whether two ranges are equal.
func (tr *TimeRange) Equals(other TimeRange) bool {
	if !tr.start.Equal(other.start) {
		return false
	}

	if tr.end == nil && other.end == nil {
		return true
	}

	if tr.end == nil || other.end == nil {
		return false
	}

	return tr.end.Equal(*other.end)
}

type timeRangeDTO struct {
	Start time.Time  `json:"start"`
	End   *time.Time `json:"end,omitempty"`
}

func (tr TimeRange) MarshalJSON() ([]byte, error) {
	dto := timeRangeDTO{
		Start: tr.start,
		End:   tr.end,
	}
	return json.Marshal(dto)
}

func (tr *TimeRange) UnmarshalJSON(data []byte) error {
	var dto timeRangeDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return err
	}

	newTR, err := New(dto.Start, dto.End)
	if err != nil {
		return err
	}

	*tr = *newTR
	return nil
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
