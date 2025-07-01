# Testing Tutorial Guard AI Integration

## Prerequisites

1. Set your Anthropic API key as an environment variable:
   ```bash
   export ANTHROPIC_API_KEY="your-api-key-here"
   ```

2. Ensure Go 1.21+ is installed

## Basic AI Integration Test

Test the core AI capabilities with simple examples:

```bash
# Build the test
go build -o test-ai cmd/test-ai/main.go

# Run the test
./test-ai
```

This test validates:
- Health check connectivity
- Instruction parsing
- Expectation validation
- Error interpretation
- Context compression

## Real-World Tutorial Test

Test on actual tutorial content from Spack-Manager-Go:

```bash
# Build the real tutorial test
go build -o test-real-tutorial cmd/test-real-tutorial/main.go

# Run the test
./test-real-tutorial
```

This test demonstrates:
- Complete tutorial interpretation (9 instructions across 3 sections)
- Complex instruction parsing
- Error recovery scenarios
- Performance metrics and cost analysis

## Expected Results

### Basic Test Output
```
🧠 Testing Tutorial Guard AI Integration with Claude...
🔍 Testing health check...
✅ Health check passed!

📖 Testing instruction parsing...
🎯 AI Intent: Create and enter a new directory for project initialization
📊 Confidence: 1.00
...
```

### Real Tutorial Test Output
```
🔬 Testing Tutorial Guard on Real Spack-Manager-Go Tutorial
📖 Tutorial: Installing Spack Manager Go
🎯 Success Rate: 100% (9/9 instructions understood)
📊 Confidence: 0.95-1.0 range (average 0.98)
💰 Total Cost: ~$0.06 for comprehensive analysis
```

## Performance Benchmarks

Typical costs per operation:
- Instruction Parsing: $0.0068 per instruction
- Validation: $0.0041 per check
- Error Interpretation: $0.0089 per error
- Context Compression: $0.0032 per compression

## Troubleshooting

**API Key Issues:**
- Ensure `ANTHROPIC_API_KEY` is set correctly
- Verify the API key has sufficient credits
- Check for network connectivity to api.anthropic.com

**Build Issues:**
- Ensure Go 1.21+ is installed
- Run `go mod tidy` to resolve dependencies
- Check that you're in the tutorial-guard directory

**Model Issues:**
- Default model is `claude-3-5-sonnet-20241022`
- If unavailable, update the model name in `pkg/ai/claude.go`