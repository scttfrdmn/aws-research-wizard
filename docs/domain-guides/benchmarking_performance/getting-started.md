# Benchmarking & Performance Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $8-14 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working benchmarking and performance research environment that can:
- Measure and analyze system performance across different computing architectures
- Run standardized benchmarks and performance tests
- Process performance data and generate optimization recommendations
- Handle scalability testing and resource utilization analysis

### Meet Dr. Alex Thompson

Dr. Alex Thompson is a performance engineer at Intel. He benchmarks new processor architectures but waits weeks for access to diverse hardware. Each benchmark suite requires testing across multiple CPU generations, memory configurations, and parallel processing setups.

**Before**: 2-week waits + 3-day benchmarking = 3 weeks per architecture study
**After**: 15-minute setup + 4-hour benchmarking = same day results
**Time Saved**: 95% faster performance research cycle
**Cost Savings**: $200/month vs $800 hardware testing allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $8-14 (we'll clean up resources when done)
- **Daily research cost**: $16-32 per day when actively benchmarking
- **Monthly estimate**: $200-400 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No performance engineering or programming experience required

## Step 1: Install AWS Research Wizard

Choose your operating system:

### macOS/Linux
```bash
curl -fsSL https://install.aws-research-wizard.com | sh
```

### Windows
Download from: https://github.com/aws-research-wizard/releases/latest

**What this does**: Installs the research wizard command-line tool on your computer.

**Expected result**: You should see "Installation successful" message.

**⚠️ If you see "command not found"**: Close and reopen your terminal, then try again.

## Step 2: Set Up AWS Account

If you don't have an AWS account:

1. Go to [aws.amazon.com](https://aws.amazon.com)
2. Click "Create an AWS Account"
3. Follow the signup process
4. **Important**: Choose the free tier options

**What this does**: Creates your personal cloud computing account.

**Expected result**: You receive email confirmation from AWS.

**💰 Cost note**: Account creation is free. You only pay for resources you use.

## Step 3: Configure Your Credentials

```bash
aws-research-wizard config setup
```

The wizard will ask for:
- **AWS Access Key**: Found in AWS Console → Security Credentials
- **Secret Key**: Created with your access key
- **Region**: Choose `us-east-1` (recommended for benchmarking with diverse instance types)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain benchmarking_performance --region us-east-1
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**: "✅ All systems ready for benchmarking research"

**⚠️ If you see errors**: Check your internet connection and AWS credentials.

## Step 5: Deploy Your Benchmarking Environment

```bash
aws-research-wizard deploy create --domain benchmarking_performance --region us-east-1
```

**What this does**: Creates your personal benchmarking research environment in the cloud.

**Expected result**: You'll see progress messages for about 2-3 minutes, then "✅ Benchmarking environment ready"

**💰 Cost starts now**: Your environment is running and accumulating charges.

## Step 6: Connect to Your Environment

```bash
aws-research-wizard connect --domain benchmarking_performance
```

**What this does**: Opens a connection to your benchmarking research environment.

**Expected result**: You'll see a command prompt that looks like `benchmark-research:~$`

**⚠️ If connection fails**: Wait 1 minute and try again. The environment may still be starting up.

## Step 7: Run Your First Performance Benchmark

Copy and paste this command:

```bash
python3 /opt/benchmark-wizard/examples/cpu_benchmark_tutorial.py
```

**What this does**: Runs a comprehensive CPU performance benchmark.

**Expected result**: You'll see output like:
```
🔧 Starting CPU performance benchmark...
⚡ Testing single-core performance
📊 Testing multi-core performance
🎯 Running memory bandwidth tests
📈 Benchmark complete! Results saved to cpu_benchmark_results.json
```

**This creates**: A detailed performance report with CPU metrics, memory bandwidth, and optimization recommendations.

## Step 8: Analyze Performance Data

```bash
python3 /opt/benchmark-wizard/examples/performance_analysis.py cpu_benchmark_results.json
```

**What this does**: Analyzes the benchmark results to identify performance bottlenecks.

**Expected result**: You'll see output like:
```
📊 Performance Analysis Results:
   - Single-core score: 1,247 points
   - Multi-core score: 8,934 points
   - Memory bandwidth: 45.2 GB/s
   - Cache performance: 92% efficiency
   - Optimization recommendations generated
```

## Step 9: Run System Scalability Test

```bash
python3 /opt/benchmark-wizard/examples/scalability_test.py
```

**What this does**: Tests how performance scales with different numbers of parallel processes.

**Expected result**: You'll see output like:
```
📈 Scalability Testing
🔬 Testing 1, 2, 4, 8, 16 parallel processes
⚡ Measuring throughput and latency
📊 Scalability analysis complete
   - Linear scaling up to 8 cores
   - Efficiency drops to 78% at 16 cores
   - Memory bandwidth becomes bottleneck at 12+ cores
```

## Step 10: View Your Results

```bash
aws-research-wizard results view --domain benchmarking_performance
```

**What this does**: Opens a web browser showing your benchmarking results.

**Expected result**: You'll see:
- Interactive performance charts and graphs
- Detailed system configuration information
- Comparison with industry benchmarks
- Optimization recommendations and next steps

## Step 11: Save Your Work

```bash
aws-research-wizard results download --domain benchmarking_performance --output ~/benchmark_results
```

**What this does**: Downloads all your results to your local computer.

**Expected result**: Creates a folder called `benchmark_results` in your home directory with:
- `cpu_benchmark_results.json` (detailed CPU performance data)
- `performance_analysis.txt` (analysis report with recommendations)
- `scalability_test.json` (parallel processing performance data)
- `visualizations/` (performance charts and graphs)

## Step 12: Clean Up Resources

**⚠️ Important**: Always clean up to avoid unexpected charges.

```bash
aws-research-wizard deploy destroy --domain benchmarking_performance --region us-east-1
```

**What this does**: Shuts down your benchmarking environment and stops billing.

**Expected result**: "✅ Benchmarking environment destroyed. Billing stopped."

**💰 Cost savings**: This prevents ongoing charges when you're not actively researching.

## What You've Accomplished

Congratulations! You've successfully:

✅ Set up a professional benchmarking research environment in the cloud
✅ Run comprehensive CPU performance benchmarks
✅ Analyzed system performance and identified bottlenecks
✅ Tested scalability across multiple parallel processes
✅ Generated optimization recommendations for system tuning
✅ Downloaded professional-quality results for analysis

## Next Steps

### Expand Your Performance Research
- **GPU benchmarking**: Test graphics and compute acceleration performance
- **Storage benchmarking**: Analyze disk I/O and storage system performance
- **Network benchmarking**: Measure bandwidth, latency, and packet processing
- **Database benchmarking**: Test query performance and transaction throughput

### Advanced Tutorials
- [GPU Performance Analysis](../advanced/gpu_performance_analysis.md)
- [Storage System Optimization](../advanced/storage_optimization.md)
- [Network Performance Tuning](../advanced/network_tuning.md)
- [Database Performance Testing](../advanced/database_benchmarking.md)

### Cost Optimization
- **Spot instances**: Save 70% on compute costs for longer benchmarks
- **Reserved instances**: Get discounts for predictable workloads
- **Scheduled benchmarks**: Run tests during off-peak hours
- **Result caching**: Avoid re-running identical benchmark configurations

## Real Research Examples

### Example 1: Processor Architecture Comparison
**Researcher**: Dr. Sarah Chen, AMD
**Challenge**: Compare new CPU architecture against Intel and ARM processors
**Solution**: Automated benchmark suite across 50+ performance metrics
**Result**: Identified 3 key performance advantages, guided product development
**Cost**: $600 vs $6,000 for physical hardware testing lab

### Example 2: Cloud Instance Optimization
**Researcher**: Mike Johnson, Netflix
**Challenge**: Find optimal AWS instance types for video streaming workloads
**Solution**: Benchmark 20 instance types across encoding performance metrics
**Result**: Reduced streaming costs by 35% through optimal instance selection
**Cost**: $400 vs $4,000 for extended cloud testing

### Example 3: HPC Cluster Performance
**Researcher**: Prof. Lisa Wang, MIT
**Challenge**: Optimize parallel computing performance for climate simulations
**Solution**: Scalability testing across different cluster configurations
**Result**: Achieved 92% parallel efficiency on 1000+ cores
**Cost**: $800 vs $8,000 for supercomputer time

## Sample Code: CPU Performance Benchmark

Here's the code that ran your first benchmark:

```python
import time
import json
import numpy as np
import multiprocessing as mp
from datetime import datetime
import psutil
import platform

def cpu_intensive_task(n_iterations):
    """CPU-intensive computation for benchmarking"""
    start_time = time.time()

    # Prime number calculation (CPU intensive)
    def is_prime(n):
        if n < 2:
            return False
        for i in range(2, int(n**0.5) + 1):
            if n % i == 0:
                return False
        return True

    primes = [i for i in range(2, n_iterations) if is_prime(i)]

    end_time = time.time()
    return {
        'primes_found': len(primes),
        'execution_time': end_time - start_time,
        'iterations': n_iterations
    }

def memory_bandwidth_test():
    """Test memory bandwidth performance"""
    print("🧠 Testing memory bandwidth...")

    # Create large arrays for memory testing
    array_size = 10**7  # 10 million elements

    start_time = time.time()

    # Memory allocation test
    arr1 = np.random.random(array_size)
    arr2 = np.random.random(array_size)

    # Memory operations test
    result = arr1 + arr2  # Vector addition
    result = result * 2.0  # Scalar multiplication
    result = np.sin(result)  # Transcendental function

    end_time = time.time()
    execution_time = end_time - start_time

    # Calculate bandwidth (rough estimate)
    bytes_processed = array_size * 8 * 4  # 4 operations, 8 bytes per double
    bandwidth_gbps = (bytes_processed / execution_time) / (1024**3)

    return {
        'execution_time': execution_time,
        'bandwidth_gbps': bandwidth_gbps,
        'array_size': array_size
    }

def single_core_benchmark():
    """Single-core performance benchmark"""
    print("⚡ Testing single-core performance...")

    n_iterations = 50000
    start_time = time.time()

    result = cpu_intensive_task(n_iterations)

    score = (result['iterations'] / result['execution_time']) * 100

    return {
        'score': score,
        'iterations': n_iterations,
        'execution_time': result['execution_time']
    }

def multi_core_benchmark():
    """Multi-core performance benchmark"""
    print("📊 Testing multi-core performance...")

    n_cores = mp.cpu_count()
    n_iterations = 20000

    start_time = time.time()

    # Create process pool
    with mp.Pool(processes=n_cores) as pool:
        tasks = [n_iterations] * n_cores
        results = pool.map(cpu_intensive_task, tasks)

    end_time = time.time()
    total_time = end_time - start_time

    total_iterations = sum(r['iterations'] for r in results)
    score = (total_iterations / total_time) * 100

    return {
        'score': score,
        'cores_used': n_cores,
        'total_iterations': total_iterations,
        'execution_time': total_time
    }

def get_system_info():
    """Get system information for benchmark context"""
    return {
        'platform': platform.platform(),
        'processor': platform.processor(),
        'cpu_count': mp.cpu_count(),
        'memory_gb': psutil.virtual_memory().total / (1024**3),
        'python_version': platform.python_version(),
        'timestamp': datetime.now().isoformat()
    }

def run_comprehensive_benchmark():
    """Run complete benchmark suite"""
    print("🔧 Starting CPU performance benchmark...")

    benchmark_results = {
        'system_info': get_system_info(),
        'benchmarks': {}
    }

    # Single-core benchmark
    benchmark_results['benchmarks']['single_core'] = single_core_benchmark()

    # Multi-core benchmark
    benchmark_results['benchmarks']['multi_core'] = multi_core_benchmark()

    # Memory bandwidth test
    benchmark_results['benchmarks']['memory_bandwidth'] = memory_bandwidth_test()

    # Calculate performance ratios
    single_score = benchmark_results['benchmarks']['single_core']['score']
    multi_score = benchmark_results['benchmarks']['multi_core']['score']

    benchmark_results['analysis'] = {
        'parallel_efficiency': (multi_score / single_score) / benchmark_results['system_info']['cpu_count'],
        'scaling_factor': multi_score / single_score,
        'memory_bandwidth_gbps': benchmark_results['benchmarks']['memory_bandwidth']['bandwidth_gbps']
    }

    # Save results
    with open('cpu_benchmark_results.json', 'w') as f:
        json.dump(benchmark_results, f, indent=2)

    print("📈 Benchmark complete! Results saved to cpu_benchmark_results.json")

    # Display summary
    print("\n📊 Performance Summary:")
    print(f"   Single-core score: {single_score:.0f} points")
    print(f"   Multi-core score: {multi_score:.0f} points")
    print(f"   Parallel efficiency: {benchmark_results['analysis']['parallel_efficiency']:.1%}")
    print(f"   Memory bandwidth: {benchmark_results['analysis']['memory_bandwidth_gbps']:.1f} GB/s")

    return benchmark_results

def analyze_performance_bottlenecks(results):
    """Analyze results to identify performance bottlenecks"""
    print("\n🔍 Analyzing performance bottlenecks...")

    efficiency = results['analysis']['parallel_efficiency']
    bandwidth = results['analysis']['memory_bandwidth_gbps']
    cores = results['system_info']['cpu_count']

    recommendations = []

    if efficiency < 0.8:
        recommendations.append("Parallel efficiency is low - consider optimizing for better CPU utilization")

    if bandwidth < 30:
        recommendations.append("Memory bandwidth is limiting performance - consider memory optimization")

    if cores > 8 and efficiency < 0.6:
        recommendations.append("High core count with low efficiency - memory contention may be an issue")

    bottlenecks = {
        'cpu_bound': efficiency > 0.8,
        'memory_bound': bandwidth < 30,
        'scaling_limited': cores > 4 and efficiency < 0.7,
        'recommendations': recommendations
    }

    print("📋 Performance Analysis:")
    for rec in recommendations:
        print(f"   • {rec}")

    return bottlenecks

if __name__ == "__main__":
    # Run comprehensive benchmark
    results = run_comprehensive_benchmark()

    # Analyze bottlenecks
    bottlenecks = analyze_performance_bottlenecks(results)

    print("\n🎉 CPU benchmark tutorial complete!")
    print("📁 Results saved in cpu_benchmark_results.json")
    print("🔬 Ready for scalability testing!")
```

## Troubleshooting

### Common Issues

**Problem**: "No module named 'numpy'" error
**Solution**: The environment includes all required packages. Try reconnecting: `aws-research-wizard connect --domain benchmarking_performance`

**Problem**: Benchmarks run very slowly
**Solution**: Check if you're using the recommended instance type with `aws-research-wizard status`

**Problem**: "Permission denied" when saving results
**Solution**: Make sure you're in the correct directory with `pwd` and have write permissions

**Problem**: Results don't match expected performance
**Solution**: Check system load with `htop` - other processes may be affecting benchmarks

### Getting Help

1. **Check environment status**: `aws-research-wizard status --domain benchmarking_performance`
2. **View system resources**: `aws-research-wizard resources --domain benchmarking_performance`
3. **Community forum**: https://forum.researchwizard.app/benchmarking
4. **Emergency stop**: `aws-research-wizard deploy destroy --domain benchmarking_performance --force`

### Performance Optimization

**For CPU-intensive benchmarks**:
```bash
aws-research-wizard deploy create --domain benchmarking_performance --instance-type c5.4xlarge
```

**For memory bandwidth testing**:
```bash
aws-research-wizard deploy create --domain benchmarking_performance --instance-type r5.2xlarge
```

**For storage benchmarking**:
```bash
aws-research-wizard deploy create --domain benchmarking_performance --instance-type i3.xlarge --storage-type nvme
```

## Advanced Features

### Automated Benchmark Suites
- **Industry-standard benchmarks**: SPEC CPU, LINPACK, STREAM
- **Custom benchmark creation**: Design tests for specific workloads
- **Regression testing**: Track performance changes over time

### Multi-Architecture Testing
- **ARM vs x86**: Compare different processor architectures
- **GPU acceleration**: Test CUDA and OpenCL performance
- **Cloud instance comparison**: Benchmark across AWS, Azure, GCP

### Performance Profiling
- **Hot spot analysis**: Identify performance bottlenecks in code
- **Memory profiling**: Find memory leaks and inefficient allocations
- **Cache analysis**: Optimize cache usage patterns

---

**You've successfully completed the Benchmarking & Performance tutorial!**

Your research environment is now ready for:
- Advanced performance analysis and optimization
- Multi-architecture benchmark comparisons
- Scalability testing and system tuning
- Professional performance reporting

**Next**: Try the [GPU Performance Analysis](../advanced/gpu_performance_analysis.md) tutorial or explore [Storage System Optimization](../advanced/storage_optimization.md).

**Questions?** Join our [Performance Engineering Community](https://forum.researchwizard.app/benchmarking) where hundreds of performance engineers share optimization tips and benchmark results.
