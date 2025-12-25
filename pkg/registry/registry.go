package registry

import (
	"github.com/tavaresphil/go-policy-engine/pkg/evaluators/expr"
	"github.com/tavaresphil/go-policy-engine/pkg/evaluators/native"
	"github.com/tavaresphil/go-policy-engine/pkg/policies"
)

type EngineKind string

const (
	EngineExpr   EngineKind = "expr"
	EngineNative EngineKind = "native"
)

var handlers = make(map[EngineKind]policies.Engine)

func Register(kind EngineKind, eng policies.Engine) {
	handlers[kind] = eng
}

func MustGet(kind EngineKind) policies.Engine {
	f, ok := handlers[kind]
	if !ok {
		panic("unknown policy engine: " + string(kind))
	}
	return f
}

func init() {
	Register(EngineExpr, expr.NewEngine())
	Register(EngineNative, native.NewEngine())
}
