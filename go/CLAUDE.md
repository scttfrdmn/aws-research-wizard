# Claude Development Guidelines for AWS Research Wizard

## Core Development Rules

### **CRITICAL RULE: Never Work Around Problems - Fix Them**

**This is the fundamental rule for all development work on this project.**

When encountering compilation errors, runtime issues, or any technical problems:

1. **NEVER use workarounds** - temporary fixes, commented-out code, or "hacks"
2. **ALWAYS investigate the root cause** and implement proper solutions
3. **PROPERLY fix compilation errors** by understanding type systems, imports, and function signatures
4. **Resolve dependency issues** rather than bypassing them
5. **Fix failing tests** rather than ignoring them

### Examples of What NOT to Do:
- Commenting out failing code instead of fixing it
- Using underscore assignments to ignore unused variables without understanding why they're unused
- Working around type errors instead of understanding the proper type system
- Bypassing import issues instead of resolving package dependencies
- Making temporary changes and marking things as "completed" when they're not

### Examples of Proper Problem Solving:
- **Type Issues**: Understand package structure and use correct type references (e.g., `data.ProjectConfig` vs `ProjectConfig`)
- **Import Problems**: Resolve package dependencies and understand Go module structure
- **Function Signatures**: Check actual function definitions and use correct parameter counts
- **Variable Shadowing**: Rename variables to avoid shadowing package names
- **Build Errors**: Fix all compilation issues before proceeding

## Project Context

### Intelligent Data Movement System
This project implements a sophisticated, configuration-driven data movement system for research environments. Key components:

- **Transfer Engine Abstraction**: Unified interface for s5cmd, rclone, aws CLI, etc.
- **Pattern Analysis**: Smart detection of small file problems and data characteristics
- **Cost Optimization**: S3 pricing analysis and bundling recommendations
- **Workflow Orchestration**: Declarative YAML configuration for complex transfers
- **Domain Optimization**: Research-specific optimizations (genomics, climate, ML)

### Package Structure
```
/Users/scttfrdmn/src/aws-research-wizard/go/
├── internal/data/              # Core data movement engines and types
├── internal/commands/data/     # CLI commands for data operations
├── internal/aws/              # AWS service integration
├── examples/                  # Configuration examples
└── cmd/                       # Application entry points
```

### Key Files
- `internal/data/project_config.go` - Declarative configuration types
- `internal/data/workflow_engine.go` - Workflow orchestration engine
- `internal/data/transfer_engine.go` - Transfer engine interface
- `internal/commands/data/workflow.go` - CLI workflow commands
- `examples/genomics-project.yaml` - Real-world configuration example

## Development History

The system was built progressively:
1. **Phase 1**: Transfer engine abstraction and AWS integration
2. **Phase 2**: Pattern analysis and cost optimization
3. **Phase 3**: Intelligent recommendations and warnings
4. **Phase 4**: Workflow orchestration and CLI integration
5. **Phase 5**: Research domain optimizations and bundling

## Testing and Validation

- Always run `go build ./cmd/main.go` to verify compilation
- Run `go test ./internal/data/... -v` to validate functionality
- Test CLI commands with real configurations from `examples/`
- Ensure all code follows Go best practices and is production-ready

## Current Status

The intelligent data movement system is fully implemented and functional:
- ✅ All core engines implemented (s5cmd, rclone, suitcase)
- ✅ Pattern analysis and cost optimization working
- ✅ Workflow orchestration complete
- ✅ CLI integration functional
- ✅ All compilation errors resolved properly

## Future Development

When adding new features:
1. Follow the established patterns and architecture
2. Add comprehensive tests for new functionality
3. Update configuration schemas as needed
4. Maintain the declarative configuration approach
5. **Always fix problems properly - never work around them**

---

**Remember: The quality of this codebase depends on solving problems correctly rather than working around them. This is not negotiable.**
