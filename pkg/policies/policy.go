package policies

import (
	"fmt"
	"time"

	"github.com/tavaresphil/go-policy-engine/pkg/timerange"
	"github.com/tavaresphil/go-policy-engine/pkg/utils"
)

type Effect string

const (
	// EffectDeny denies access when the condition matches.
	EffectDeny Effect = "deny"
	// EffectAllow allows access when the condition matches.
	EffectAllow Effect = "allow"
)

// Policy represents an access control policy for a resource (and optionally a
// resource ID). Policies have an effect (allow/deny), an active period and a
// root condition that determines applicability.
type Policy struct {
	ID         string               `json:"id,omitempty"`
	Resource   string               `json:"resource,omitempty"`
	ResourceID string               `json:"resource_id,omitempty"`
	Effect     Effect               `json:"effect,omitempty"`
	Condition  PolicyCondition      `json:"condition,omitempty"`
	Version    string               `json:"version,omitempty"`
	DryRun     bool                 `json:"dry_run,omitempty"`
	Period     *timerange.TimeRange `json:"period,omitempty"`
}

func (p Policy) IsActiveAt(t time.Time) bool {
	if p.Period == nil {
		return true
	}
	return p.Period.Contains(t)
}

func (p Policy) IsActive() bool {
	return p.IsActiveAt(time.Now())
}

func (p *Policy) Deactivate() error {
	// check if already deactivated
	if p.Period.End() != nil && p.Period.End().Before(time.Now()) {
		return nil
	}

	newPeriod, err := p.Period.WithEnd(utils.Ptr(time.Now()))
	if err != nil {
		return err
	}
	p.Period = newPeriod
	return nil
}

// Validate validates the policy structure and its condition
func (p Policy) Validate() error {
	if p.Resource == "" {
		return fmt.Errorf("policy resource is required")
	}
	if p.Effect != EffectAllow && p.Effect != EffectDeny {
		return fmt.Errorf("invalid effect: %s", p.Effect)
	}
	if p.Period == nil {
		return fmt.Errorf("policy period is required")
	}
	return p.Condition.Validate()
}

// IsExpired checks if the policy has expired
func (p Policy) IsExpired() bool {
	if p.Period == nil || p.Period.End() == nil {
		return false
	}
	return p.Period.End().Before(time.Now())
}

// IsExpiredAt checks if the policy was expired at a specific time
func (p Policy) IsExpiredAt(t time.Time) bool {
	if p.Period == nil || p.Period.End() == nil {
		return false
	}
	return p.Period.End().Before(t)
}

// WillExpireIn checks if the policy will expire within a duration
func (p Policy) WillExpireIn(d time.Duration) bool {
	if p.Period == nil || p.Period.End() == nil {
		return false
	}
	return p.Period.End().Before(time.Now().Add(d))
}

// IsScheduled checks if the policy is scheduled to start in the future
func (p Policy) IsScheduled() bool {
	if p.Period == nil {
		return false
	}
	return p.Period.Start().After(time.Now())
}

// Matches checks if the policy applies to a given resource and resourceID
func (p Policy) Matches(resource, resourceID string) bool {
	return p.Resource == resource && p.ResourceID == resourceID
}

// MatchesResource checks if the policy applies to a given resource (ignoring resourceID)
func (p Policy) MatchesResource(resource string) bool {
	return p.Resource == resource
}

// AppliesTo checks if policy is active and matches the resource
func (p Policy) AppliesTo(resource, resourceID string) bool {
	return p.IsActive() && p.Matches(resource, resourceID)
}

// IsDeny returns true if the policy effect is deny
func (p Policy) IsDeny() bool {
	return p.Effect == EffectDeny
}

// IsAllow returns true if the policy effect is allow
func (p Policy) IsAllow() bool {
	return p.Effect == EffectAllow
}

// ShouldBlock returns true if the condition matches and effect is deny
// or if condition doesn't match and effect is allow
func (p Policy) ShouldBlock(conditionMatches bool) bool {
	if p.Effect == EffectDeny {
		return conditionMatches
	}
	return !conditionMatches
}

// Activate sets the policy period to start now if it's scheduled
func (p *Policy) Activate() error {
	if !p.IsScheduled() {
		return nil // Already active or expired
	}

	newPeriod, err := p.Period.WithStart(time.Now())
	if err != nil {
		return err
	}
	p.Period = newPeriod
	return nil
}

// ExtendBy extends the policy period by a duration
func (p *Policy) ExtendBy(d time.Duration) error {
	if p.Period == nil {
		return fmt.Errorf("policy has no period")
	}

	var newEnd *time.Time
	if p.Period.End() != nil {
		t := p.Period.End().Add(d)
		newEnd = &t
	} else {
		t := time.Now().Add(d)
		newEnd = &t
	}

	newPeriod, err := p.Period.WithEnd(newEnd)
	if err != nil {
		return err
	}
	p.Period = newPeriod
	return nil
}

// SetEndDate sets a specific end date for the policy
func (p *Policy) SetEndDate(endTime time.Time) error {
	if p.Period == nil {
		return fmt.Errorf("policy has no period")
	}

	newPeriod, err := p.Period.WithEnd(&endTime)
	if err != nil {
		return err
	}
	p.Period = newPeriod
	return nil
}

// RemainingDuration returns the duration until the policy expires
// Returns 0 if policy has no end or is already expired
func (p Policy) RemainingDuration() time.Duration {
	if p.Period == nil || p.Period.End() == nil {
		return 0
	}

	remaining := time.Until(*p.Period.End())
	if remaining < 0 {
		return 0
	}
	return remaining
}

// String returns a human-readable representation of the policy
func (p Policy) String() string {
	status := "active"
	if p.IsExpired() {
		status = "expired"
	} else if p.IsScheduled() {
		status = "scheduled"
	}

	dryRunStr := ""
	if p.DryRun {
		dryRunStr = " [DRY-RUN]"
	}

	return fmt.Sprintf("Policy[%s] %s/%s: %s (%s)%s",
		p.ID, p.Resource, p.ResourceID, p.Effect, status, dryRunStr)
}

// Clone creates a deep copy of the policy
func (p Policy) Clone() Policy {
	clone := p
	if p.Period != nil {
		periodCopy := *p.Period
		clone.Period = &periodCopy
	}
	return clone
}

// WithDryRun returns a copy of the policy with DryRun set to the specified value
func (p Policy) WithDryRun(dryRun bool) Policy {
	clone := p.Clone()
	clone.DryRun = dryRun
	return clone
}

// IsSameResource checks if two policies target the same resource
func (p Policy) IsSameResource(other Policy) bool {
	return p.Resource == other.Resource && p.ResourceID == other.ResourceID
}

// HasConflict checks if two policies have conflicting effects on the same resource
func (p Policy) HasConflict(other Policy) bool {
	return p.IsSameResource(other) && p.Effect != other.Effect
}

// Priority returns a priority score for conflict resolution
// Higher priority = more specific/restrictive
func (p Policy) Priority() int {
	priority := 0

	// Deny policies have higher priority
	if p.IsDeny() {
		priority += 100
	}

	// Non-dry-run policies have higher priority
	if !p.DryRun {
		priority += 50
	}

	// Policies with specific resource IDs have higher priority
	if p.ResourceID != "" && p.ResourceID != "*" {
		priority += 25
	}

	return priority
}
