// Package policies defines the core policy model used by the policy engine.
//
// A Policy contains metadata (resource, effect, period) and a root
// PolicyCondition that is evaluated against an attribute context. The
// package provides types and helpers to validate, evaluate and manipulate
// policies and conditions.
package policies
