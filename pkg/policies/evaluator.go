package policies

import (
	"context"
	"fmt"
)

type Engine interface {
	Eval(cond PolicyCondition, ctx Resolver) (bool, error)
}

type PolicyRepository interface {
	FindByResourceAndResourceID(ctx context.Context, resource, resourceID string) ([]Policy, error)
}

type EvaluatorRequest struct {
	Resource   string
	ResourceID string
	Context    MapAttributes
}

// Evaluator is a higher level component that retrieves policies from a
// repository and executes them using an Engine against a request context.
type Evaluator interface {
	Eval(ctx context.Context, req EvaluatorRequest) error
}

type evaluator struct {
	eng  Engine
	repo PolicyRepository
}

// NewEvaluator constructs a new Evaluator from an Engine and a PolicyRepository.
func NewEvaluator(eng Engine, repo PolicyRepository) Evaluator {
	return &evaluator{
		eng:  eng,
		repo: repo,
	}
}

func (e *evaluator) Eval(ctx context.Context, req EvaluatorRequest) error {
	pols, err := e.repo.FindByResourceAndResourceID(ctx, req.Resource, req.ResourceID)
	if err != nil {
		return err
	}

	for _, pol := range pols {
		ok, err := e.eng.Eval(pol.Condition, req.Context)
		if err != nil {
			return err
		}

		allowed := ok
		if pol.Effect == EffectDeny {
			allowed = !allowed
		}

		if pol.DryRun {
			return nil
		}

		if !allowed {
			return fmt.Errorf("execution is dained")
		}
	}
	return nil
}
