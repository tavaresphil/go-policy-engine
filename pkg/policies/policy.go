package policies

import (
	"time"

	"github.com/tavaresphil/go-policy-engine/pkg/timerange"
	"github.com/tavaresphil/go-policy-engine/pkg/utils"
)

type Effect string

const (
	EffectDeny  Effect = "deny"
	EffectAllow Effect = "allow"
)

type Policy struct {
	ID         string               `json:"id"`
	Resource   string               `json:"resource"`
	ResourceID string               `json:"resource_id"`
	Effect     Effect               `json:"effect"`
	Condition  PolicyCondition      `json:"condition"`
	Version    string               `json:"version"`
	DryRun     bool                 `json:"dry_run"`
	Period     *timerange.TimeRange `json:"period"`
}

func (p Policy) IsActiveAt(t time.Time) bool {
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
