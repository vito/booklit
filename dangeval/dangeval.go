// Package dangeval embeds the Dang interpreter so Booklit documents can
// evaluate {expr} interpolations inside JSX. An Evaluator holds a long-
// lived type and value environment so repeated snippets share a scope.
//
// Discovery follows Dang's own conventions:
//   - dang.toml (walked upward from the project directory) configures
//     extra GraphQL imports.
//   - dagger.json (in the project directory or above) is auto-imported
//     as the default Dagger module if present.
//
// Per-section / file-local scope is not supported yet; every snippet
// runs against the same global env. See jsx-dang.md (Phase 3) for the
// follow-ups still on the list.
package dangeval

import (
	"context"
	"fmt"

	"github.com/vito/dang/pkg/dang"
	"github.com/vito/dang/pkg/hm"
)

// Evaluator parses and evaluates Dang expressions against a held-open
// environment. Construct one per build with New, defer Close, then call
// Eval as many times as needed.
type Evaluator struct {
	ctx      context.Context
	typeEnv  hm.Env
	evalEnv  dang.EvalEnv
	services *dang.ServiceRegistry
}

// New constructs an Evaluator rooted at projectDir. dang.toml and
// dagger.json discovery walk up from projectDir following Dang's
// conventions.
func New(ctx context.Context, projectDir string) (*Evaluator, error) {
	services := &dang.ServiceRegistry{}
	ctx = dang.ContextWithServices(ctx, services)

	var configs []dang.ImportConfig

	configPath, projCfg, err := dang.FindProjectConfig(projectDir)
	if err != nil {
		return nil, fmt.Errorf("finding dang.toml: %w", err)
	}
	if projCfg != nil {
		ctx = dang.ContextWithProjectConfig(ctx, configPath, projCfg)
		resolved, err := dang.ResolveImportConfigs(ctx, projCfg, projectDirOf(configPath))
		if err != nil {
			return nil, fmt.Errorf("resolving dang.toml imports: %w", err)
		}
		configs = resolved
	}

	configs = dang.ResolveDaggerImport(ctx, configs, projectDir)

	if len(configs) > 0 {
		ctx = dang.ContextWithImportConfigs(ctx, configs...)
	}

	typeEnv, evalEnv := dang.BuildEnvFromImports("booklit", configs)

	return &Evaluator{
		ctx:      ctx,
		typeEnv:  typeEnv,
		evalEnv:  evalEnv,
		services: services,
	}, nil
}

// Close shuts down any subprocesses started by dang.toml imports (e.g.
// a `dagger session`). Safe to call multiple times.
func (e *Evaluator) Close() {
	if e == nil || e.services == nil {
		return
	}
	e.services.StopAll()
	e.services = nil
}

// Eval parses raw as a sequence of Dang forms, type-checks against the
// held type env, and evaluates against the held value env. The returned
// Value is the result of the last form; for a typical single-expression
// snippet that's the value of the expression.
func (e *Evaluator) Eval(raw string) (dang.Value, error) {
	parsed, err := dang.ParseWithRecovery("<jsx>", []byte(raw))
	if err != nil {
		return nil, err
	}

	block, ok := parsed.(*dang.ModuleBlock)
	if !ok {
		return nil, fmt.Errorf("dang parser returned %T, expected *ModuleBlock", parsed)
	}
	if len(block.Forms) == 0 {
		return nil, fmt.Errorf("empty Dang expression")
	}

	fresh := hm.NewSimpleFresher()
	if _, err := dang.InferFormsWithPhases(e.ctx, block.Forms, e.typeEnv, fresh); err != nil {
		return nil, dang.ConvertInferError(err)
	}

	var result dang.Value
	for _, form := range block.Forms {
		v, err := dang.EvalNode(e.ctx, e.evalEnv, form)
		if err != nil {
			return nil, err
		}
		result = v
	}
	return result, nil
}

// projectDirOf returns the directory containing dang.toml. Split out so
// the import path resolution mirrors Dang's own.
func projectDirOf(configPath string) string {
	for i := len(configPath) - 1; i >= 0; i-- {
		if configPath[i] == '/' || configPath[i] == '\\' {
			return configPath[:i]
		}
	}
	return "."
}
