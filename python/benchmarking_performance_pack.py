#!/usr/bin/env python3
"""
Benchmarking & Performance Analysis Pack
Comprehensive performance benchmarking, system analysis, and optimization for AWS Research Wizard
"""

import json
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
from enum import Enum

class BenchmarkingDomain(Enum):
    HPC_BENCHMARKING = "hpc_benchmarking"
    APPLICATION_PROFILING = "application_profiling"
    SYSTEM_PERFORMANCE = "system_performance"
    NETWORK_ANALYSIS = "network_analysis"
    STORAGE_BENCHMARKING = "storage_benchmarking"
    GPU_PERFORMANCE = "gpu_performance"
    CLOUD_OPTIMIZATION = "cloud_optimization"
    SCALABILITY_TESTING = "scalability_testing"
    COMPARATIVE_ANALYSIS = "comparative_analysis"

@dataclass
class BenchmarkingWorkload:
    """Benchmarking and performance analysis workload characteristics"""
    domain: BenchmarkingDomain
    benchmark_scope: str     # Single Node, Multi-Node, Cluster, Cloud, Hybrid
    target_systems: List[str] # CPU, GPU, Memory, Storage, Network, Application
    workload_types: List[str] # Compute, Memory, I/O, Network, Mixed
    measurement_focus: str   # Throughput, Latency, Scalability, Efficiency, Cost
    comparison_baseline: str # Previous Generation, Competitor, Theoretical Peak
    duration: str           # Minutes, Hours, Days, Continuous
    data_volume_tb: float   # Expected data volume from benchmarks
    analysis_complexity: str # Basic, Intermediate, Advanced, Expert

class BenchmarkingPerformancePack:
    """
    Comprehensive benchmarking and performance analysis environments optimized for AWS
    Supports HPC benchmarking, application profiling, and system optimization
    """
    
    def __init__(self):
        self.benchmarking_configurations = {
            "hpc_benchmark_suite": self._get_hpc_benchmarking_config(),
            "application_profiler": self._get_application_profiling_config(),
            "system_performance": self._get_system_performance_config(),
            "network_analyzer": self._get_network_analysis_config(),
            "storage_benchmark": self._get_storage_benchmarking_config(),
            "gpu_performance": self._get_gpu_performance_config(),
            "cloud_optimizer": self._get_cloud_optimization_config(),
            "scalability_tester": self._get_scalability_testing_config(),
            "comparative_analyzer": self._get_comparative_analysis_config()
        }
    
    def _get_hpc_benchmarking_config(self) -> Dict[str, Any]:
        """HPC benchmarking and performance evaluation platform"""
        return {
            "name": "HPC Benchmarking & Performance Evaluation Platform",
            "description": "Comprehensive HPC benchmarking suite for compute, memory, and parallel performance",
            "spack_packages": [
                # HPC benchmark suites
                "hpcc@1.5.0 %gcc@11.4.0 +openmp +mpi",      # HPC Challenge Benchmark
                "hpl@2.3 %gcc@11.4.0 +openmp +mpi",         # High Performance Linpack
                "stream@5.10 %gcc@11.4.0 +openmp",          # Memory bandwidth benchmark
                "osu-micro-benchmarks@6.2 %gcc@11.4.0 +mpi", # MPI benchmarks
                
                # Standard benchmarks
                "spec-cpu@2017.1.1 %gcc@11.4.0",            # SPEC CPU benchmarks
                "spec-mpi@2007 %gcc@11.4.0 +mpi",           # SPEC MPI benchmarks
                "nas-parallel-benchmarks@3.4.2 %gcc@11.4.0 +openmp +mpi", # NPB
                "graph500@3.0.0 %gcc@11.4.0 +mpi",          # Graph500 benchmark
                
                # Memory and cache benchmarks
                "lmbench@3.0-a9 %gcc@11.4.0",               # Memory latency benchmark
                "cachebench@0.2.0 %gcc@11.4.0",             # Cache performance
                "membench@1.0.2 %gcc@11.4.0",               # Memory subsystem benchmark
                "bandwidth@1.5.1 %gcc@11.4.0",              # Memory bandwidth testing
                
                # I/O benchmarks
                "ior@3.3.0 %gcc@11.4.0 +mpi +hdf5 +netcdf", # Parallel I/O benchmark
                "mdtest@3.4.0 %gcc@11.4.0 +mpi",            # Metadata benchmark
                "fio@3.35 %gcc@11.4.0",                     # Flexible I/O tester
                "bonnie++@2.00a %gcc@11.4.0",               # Filesystem benchmark
                
                # Network benchmarks
                "netperf@2.7.0 %gcc@11.4.0",                # Network performance
                "iperf3@3.13 %gcc@11.4.0",                  # Network throughput
                "perftest@24.01 %gcc@11.4.0",               # InfiniBand performance
                "mpiP@3.5 %gcc@11.4.0 +mpi",                # MPI profiling
                
                # AWS-optimized parallel computing libraries
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers +pmix +pmi +fabrics",
                "libfabric@1.18.1 %gcc@11.4.0 +verbs +mlx +efa",  # EFA support
                "aws-ofi-nccl@1.7.0 %gcc@11.4.0",  # AWS OFI plugin for NCCL
                "ucx@1.14.1 %gcc@11.4.0 +verbs +mlx +ib_hw_tm",  # Unified Communication X
                "mpich@4.1.2 %gcc@11.4.0 +pmi +slurm +libfabric",
                "intel-oneapi-mpi@2021.10.0 %gcc@11.4.0",
                
                # AWS ParallelCluster integration
                "aws-parallelcluster@3.7.0 %gcc@11.4.0",
                "slurm@23.02.5 %gcc@11.4.0 +pmix +numa",
                
                # Performance analysis tools
                "likwid@5.2.2 %gcc@11.4.0 +accessdaemon",   # Performance monitoring
                "papi@7.0.1 %gcc@11.4.0",                   # Performance API
                "tau@2.32.1 %gcc@11.4.0 +mpi +openmp +papi", # Performance tuning
                "scalasca@2.6.1 %gcc@11.4.0 +mpi",          # Scalability analysis
                
                # Python performance tools
                "python@3.11.5 %gcc@11.4.0",
                "py-perfplot@0.10.2 %gcc@11.4.0",           # Performance plotting
                "py-benchmark@1.0.8 %gcc@11.4.0",           # Benchmarking framework
                "py-performance@1.0.2 %gcc@11.4.0",         # Performance testing
                "py-memory-profiler@0.61.0 %gcc@11.4.0",    # Memory profiling
                
                # Data analysis and visualization
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                
                # Statistical analysis
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-benchmark@0.3-6 %gcc@11.4.0",            # R benchmarking
                "r-microbenchmark@1.4.10 %gcc@11.4.0",      # Micro-benchmarking
                
                # Database systems for results
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development and build tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0",
                "gfortran@11.4.0",
                "intel-oneapi-compilers@2023.2.1"
            ],
            "aws_instance_recommendations": {
                "single_node": {
                    "instance_type": "c6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage_gb": 500,
                    "cost_per_hour": 1.36,
                    "use_case": "Single-node CPU benchmarking and profiling"
                },
                "memory_intensive": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Memory-intensive benchmarks and large datasets"
                },
                "hpc_cluster": {
                    "instance_type": "hpc6a.48xlarge",
                    "vcpus": 96,
                    "memory_gb": 384,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.88,
                    "efa_enabled": True,
                    "placement_group": "cluster",
                    "enhanced_networking": "sr-iov",
                    "use_case": "HPC cluster benchmarking with EFA and enhanced networking"
                },
                "compute_optimized": {
                    "instance_type": "c6a.24xlarge",
                    "vcpus": 96,
                    "memory_gb": 192,
                    "storage_gb": 3800,
                    "cost_per_hour": 3.46,
                    "use_case": "Compute-intensive benchmarks and parallel workloads"
                }
            },
            "estimated_cost": {
                "compute": 800,
                "storage": 200,
                "network": 100,
                "data_transfer": 50,
                "total": 1150
            },
            "research_capabilities": [
                "SPEC CPU and MPI benchmark execution",
                "HPL and HPCC performance measurement",
                "Memory bandwidth and latency analysis",
                "Parallel I/O and filesystem benchmarking",
                "Network performance and scalability testing",
                "Cache hierarchy performance analysis",
                "Compiler and optimization comparison",
                "Performance regression detection"
            ],
            "aws_data_sources": [
                "CloudWatch performance metrics",
                "AWS X-Ray tracing data",
                "Enhanced networking statistics",
                "EBS and S3 performance metrics"
            ]
        }
    
    def _get_application_profiling_config(self) -> Dict[str, Any]:
        """Application profiling and optimization platform"""
        return {
            "name": "Application Profiling & Optimization Platform",
            "description": "Comprehensive application performance profiling, optimization, and tuning",
            "spack_packages": [
                # Profiling tools
                "intel-oneapi-vtune@2023.2.0 %gcc@11.4.0",  # Intel VTune Profiler
                "advisor@2023.2.0 %gcc@11.4.0",             # Intel Advisor
                "gprof@2.40 %gcc@11.4.0",                   # GNU profiler
                "valgrind@3.21.0 %gcc@11.4.0 +mpi +boost", # Memory debugging
                
                # Performance analysis frameworks
                "tau@2.32.1 %gcc@11.4.0 +mpi +openmp +papi +cuda",
                "scalasca@2.6.1 %gcc@11.4.0 +mpi +cube",
                "score-p@8.0 %gcc@11.4.0 +mpi +openmp +cuda +papi",
                "extrae@4.0.6 %gcc@11.4.0 +mpi +openmp +cuda",
                
                # Hardware counters and monitoring
                "papi@7.0.1 %gcc@11.4.0 +cuda +nvml",
                "likwid@5.2.2 %gcc@11.4.0 +accessdaemon +nvidia",
                "perf@6.3 %gcc@11.4.0",                     # Linux perf tool
                "oprofile@1.4.0 %gcc@11.4.0",               # System profiler
                
                # Python profiling tools
                "python@3.11.5 %gcc@11.4.0",
                "py-line-profiler@4.1.1 %gcc@11.4.0",       # Line-by-line profiling
                "py-memory-profiler@0.61.0 %gcc@11.4.0",
                "py-py-spy@0.3.14 %gcc@11.4.0",             # Sampling profiler
                "py-scalene@1.5.26 %gcc@11.4.0",            # High-performance profiler
                "py-austin@3.6.0 %gcc@11.4.0",              # Frame stack sampler
                
                # Code analysis and optimization
                "intel-oneapi-inspector@2023.2.0 %gcc@11.4.0", # Memory/threading errors
                "callgrind@3.21.0 %gcc@11.4.0",             # Call graph profiler
                "massif@3.21.0 %gcc@11.4.0",                # Heap profiler
                "helgrind@3.21.0 %gcc@11.4.0",              # Thread error detector
                
                # Parallel profiling
                "mpiP@3.5 %gcc@11.4.0 +mpi",
                "mpipi@3.4.1 %gcc@11.4.0 +mpi",
                "vampir@10.2.0 %gcc@11.4.0",                # MPI trace analysis
                "paraver@4.11.1 %gcc@11.4.0",               # Performance visualization
                
                # Application-specific profilers
                "hpctoolkit@2023.08.01 %gcc@11.4.0 +mpi +cuda +papi",
                "caliper@2.10.0 %gcc@11.4.0 +mpi +cuda +papi",
                "timemory@3.3.0 %gcc@11.4.0 +mpi +cuda +python",
                
                # Performance modeling
                "py-perfmodel@1.3.0 %gcc@11.4.0",           # Performance modeling
                "py-roofline@2.1.0 %gcc@11.4.0",            # Roofline model analysis
                "py-kerncraft@0.8.15 %gcc@11.4.0",          # Loop kernel analysis
                
                # Optimization tools
                "icc@2021.10.0 %gcc@11.4.0",                # Intel C++ compiler
                "ifort@2021.10.0 %gcc@11.4.0",              # Intel Fortran compiler
                "aocc@4.0.0 %gcc@11.4.0",                   # AMD optimizing compiler
                
                # Data analysis and visualization
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0",
                "gfortran@11.4.0"
            ],
            "aws_instance_recommendations": {
                "cpu_profiling": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage_gb": 300,
                    "cost_per_hour": 0.68,
                    "use_case": "CPU-intensive application profiling"
                },
                "memory_profiling": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Memory-intensive application analysis"
                },
                "parallel_profiling": {
                    "instance_type": "c6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage_gb": 1000,
                    "cost_per_hour": 1.36,
                    "use_case": "Parallel application profiling and optimization"
                },
                "comprehensive": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Comprehensive application performance analysis"
                }
            ],
            "estimated_cost": {
                "compute": 600,
                "storage": 150,
                "data_transfer": 40,
                "total": 790
            },
            "research_capabilities": [
                "CPU and memory usage profiling",
                "Cache miss and memory access pattern analysis",
                "Thread synchronization and parallel efficiency",
                "Hot spot identification and optimization",
                "Memory leak and error detection",
                "Compiler optimization effectiveness",
                "Roofline model performance analysis",
                "Application scalability assessment"
            ]
        }
    
    def _get_system_performance_config(self) -> Dict[str, Any]:
        """System performance monitoring and analysis platform"""
        return {
            "name": "System Performance Monitoring & Analysis Platform",
            "description": "Comprehensive system-level performance monitoring, analysis, and optimization",
            "spack_packages": [
                # System monitoring tools
                "htop@3.2.2 %gcc@11.4.0",                   # Interactive process viewer
                "iotop@0.6 %gcc@11.4.0",                    # I/O monitoring
                "nethogs@0.8.7 %gcc@11.4.0",               # Network bandwidth monitoring
                "sysstat@12.7.4 %gcc@11.4.0",              # System activity tools
                
                # Performance monitoring frameworks
                "collectd@5.12.0 %gcc@11.4.0 +python +postgresql",
                "telegraf@1.27.3 %gcc@11.4.0",             # Metrics collection
                "node-exporter@1.6.1 %gcc@11.4.0",         # Hardware metrics
                "prometheus@2.45.0 %gcc@11.4.0",           # Monitoring system
                
                # System profiling
                "systemtap@4.9 %gcc@11.4.0",               # Dynamic tracing
                "bcc@0.29.1 %gcc@11.4.0 +python",          # BPF tools
                "bpftrace@0.19.1 %gcc@11.4.0",             # High-level tracing
                "ebpf@1.2.0 %gcc@11.4.0",                  # Extended BPF
                
                # Hardware monitoring
                "lm-sensors@3.6.0 %gcc@11.4.0",            # Hardware monitoring
                "smartmontools@7.4 %gcc@11.4.0",           # S.M.A.R.T. monitoring
                "dmidecode@3.5 %gcc@11.4.0",               # Hardware information
                "hwloc@2.9.2 %gcc@11.4.0 +cuda +opencl",   # Hardware topology
                
                # Python system tools
                "python@3.11.5 %gcc@11.4.0",
                "py-psutil@5.9.5 %gcc@11.4.0",             # System information
                "py-py-cpuinfo@9.0.0 %gcc@11.4.0",         # CPU information
                "py-gpustat@1.1.1 %gcc@11.4.0",            # GPU monitoring
                "py-nvidia-ml-py@12.535.77 %gcc@11.4.0",   # NVIDIA GPU stats
                
                # Network monitoring
                "iftop@1.0pre4 %gcc@11.4.0",               # Network bandwidth
                "nload@0.7.4 %gcc@11.4.0",                 # Network load monitor
                "bandwhich@0.20.0 %gcc@11.4.0",            # Network utilization
                "speedtest-cli@2.1.3 %gcc@11.4.0",         # Network speed test
                
                # Storage monitoring
                "iostat@12.7.4 %gcc@11.4.0",               # I/O statistics
                "blktrace@1.3.0 %gcc@11.4.0",              # Block I/O tracing
                "fio@3.35 %gcc@11.4.0",                    # I/O testing
                "ioping@1.3 %gcc@11.4.0",                  # Disk latency monitor
                
                # Performance analysis
                "py-performance-analysis@2.1.0 %gcc@11.4.0",
                "py-system-metrics@1.5.0 %gcc@11.4.0",
                "py-resource-monitor@2.3.0 %gcc@11.4.0",
                
                # Time series databases
                "influxdb@2.7.1 %gcc@11.4.0",              # Time series database
                "graphite@1.1.10 %gcc@11.4.0 +python",    # Metrics database
                "victoria-metrics@1.93.1 %gcc@11.4.0",     # Fast metrics database
                
                # Visualization
                "grafana@10.0.3 %gcc@11.4.0",              # Monitoring dashboards
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-dash@2.13.0 %gcc@11.4.0",
                
                # Data analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "monitoring": {
                    "instance_type": "t3.medium",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "storage_gb": 100,
                    "cost_per_hour": 0.042,
                    "use_case": "Basic system monitoring and alerting"
                },
                "analysis": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 300,
                    "cost_per_hour": 0.34,
                    "use_case": "Performance analysis and optimization"
                },
                "comprehensive": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Comprehensive system performance analysis"
                },
                "enterprise": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Enterprise-scale performance monitoring"
                }
            },
            "estimated_cost": {
                "compute": 400,
                "storage": 100,
                "monitoring": 150,
                "data_transfer": 50,
                "total": 700
            },
            "research_capabilities": [
                "Real-time system performance monitoring",
                "Resource utilization analysis and optimization",
                "Performance bottleneck identification",
                "System capacity planning and forecasting",
                "Hardware health monitoring and alerting",
                "Application performance correlation",
                "Cost optimization through resource analysis",
                "Performance regression detection"
            ]
        }
    
    def _get_network_analysis_config(self) -> Dict[str, Any]:
        """Network performance analysis and optimization platform"""
        return {
            "name": "Network Performance Analysis & Optimization Platform",
            "description": "Comprehensive network performance testing, analysis, and optimization",
            "spack_packages": [
                # Network benchmarking tools
                "iperf3@3.13 %gcc@11.4.0",                  # Network throughput
                "netperf@2.7.0 %gcc@11.4.0",               # Network performance
                "nuttcp@8.2.2 %gcc@11.4.0",                # Network throughput
                "qperf@0.4.11 %gcc@11.4.0",                # RDMA and IP performance
                
                # Network monitoring
                "wireshark@4.0.7 %gcc@11.4.0 +qt",         # Network protocol analyzer
                "tcpdump@4.99.4 %gcc@11.4.0",              # Packet analyzer
                "nmap@7.94 %gcc@11.4.0",                   # Network discovery
                "mtr@0.95 %gcc@11.4.0",                    # Network diagnostics
                
                # AWS EFA and high-performance networking
                "perftest@24.01 %gcc@11.4.0",              # InfiniBand performance  
                "rdma-core@46.0 %gcc@11.4.0",              # RDMA libraries
                "libfabric@1.18.1 %gcc@11.4.0 +verbs +mlx +efa",  # EFA support
                "aws-efa-installer@1.26.1 %gcc@11.4.0",    # AWS EFA drivers
                "efa-profile@1.0.0 %gcc@11.4.0",           # EFA profiling tools
                
                # MPI network testing
                "osu-micro-benchmarks@6.2 %gcc@11.4.0 +mpi",
                "mpitests@3.2.20 %gcc@11.4.0 +mpi",        # MPI correctness tests
                "mpiping@4.0 %gcc@11.4.0 +mpi",            # MPI ping-pong test
                
                # Python network tools
                "python@3.11.5 %gcc@11.4.0",
                "py-speedtest-cli@2.1.3 %gcc@11.4.0",      # Network speed testing
                "py-network-benchmark@1.4.0 %gcc@11.4.0",  # Network benchmarking
                "py-scapy@2.5.0 %gcc@11.4.0",              # Packet manipulation
                "py-paramiko@3.2.0 %gcc@11.4.0",           # SSH connections
                
                # Network simulation
                "ns-3@3.38 %gcc@11.4.0 +python +mpi",      # Network simulator
                "omnet++@6.0.1 %gcc@11.4.0",               # Discrete event simulator
                "mininet@2.3.0 %gcc@11.4.0 +python",       # Network emulator
                
                # Traffic generation
                "iperf@2.1.9 %gcc@11.4.0",                 # Network testing tool
                "hping@3.0.0 %gcc@11.4.0",                 # TCP/IP packet assembler
                "ostinato@1.3.0 %gcc@11.4.0 +qt",          # Packet generator
                "trex@3.02 %gcc@11.4.0",                   # Traffic generator
                
                # Network analysis
                "py-networkx@3.1 %gcc@11.4.0",             # Network analysis
                "py-igraph@0.10.6 %gcc@11.4.0",            # Graph analysis
                "py-graph-tool@2.45 %gcc@11.4.0",          # Network analysis
                
                # Data analysis and visualization
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "network_testing": {
                    "instance_type": "c6in.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.48,
                    "use_case": "Network performance testing and analysis"
                },
                "high_bandwidth": {
                    "instance_type": "c6in.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage_gb": 500,
                    "cost_per_hour": 1.93,
                    "use_case": "High-bandwidth network testing (up to 50 Gbps)"
                },
                "ultra_high_network": {
                    "instance_type": "c6in.16xlarge",
                    "vcpus": 64,
                    "memory_gb": 128,
                    "storage_gb": 1000,
                    "cost_per_hour": 3.86,
                    "use_case": "Ultra-high network performance testing (up to 100 Gbps)"
                },
                "cluster_networking": {
                    "instance_type": "hpc6a.48xlarge",
                    "vcpus": 96,
                    "memory_gb": 384,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.88,
                    "use_case": "HPC cluster network performance analysis"
                }
            ],
            "estimated_cost": {
                "compute": 500,
                "network": 200,
                "storage": 100,
                "data_transfer": 150,
                "total": 950
            },
            "research_capabilities": [
                "Network throughput and latency measurement",
                "MPI communication pattern analysis",
                "High-bandwidth network testing (up to 100 Gbps)",
                "InfiniBand and Ethernet performance comparison",
                "Network topology optimization",
                "Application communication profiling",
                "Network congestion analysis",
                "Cloud network performance optimization"
            ]
        }
    
    def _get_storage_benchmarking_config(self) -> Dict[str, Any]:
        """Storage performance benchmarking platform"""
        return {
            "name": "Storage Performance Benchmarking Platform",
            "description": "Comprehensive storage system benchmarking and I/O performance analysis",
            "spack_packages": [
                # Storage benchmarking tools
                "fio@3.35 %gcc@11.4.0",                     # Flexible I/O tester
                "ior@3.3.0 %gcc@11.4.0 +mpi +hdf5 +netcdf", # Parallel I/O benchmark
                "mdtest@3.4.0 %gcc@11.4.0 +mpi",            # Metadata benchmark
                "bonnie++@2.00a %gcc@11.4.0",               # Filesystem benchmark
                
                # Parallel filesystem benchmarks
                "elbencho@3.0.0 %gcc@11.4.0 +mpi",          # Distributed storage benchmark
                "ddbench@1.0.0 %gcc@11.4.0 +mpi",           # Distributed database benchmark
                "macsio@1.1 %gcc@11.4.0 +mpi +hdf5",        # Multi-interface I/O
                
                # Database benchmarks
                "sysbench@1.0.20 %gcc@11.4.0 +mysql",       # Database benchmark
                "pgbench@15.4 %gcc@11.4.0",                 # PostgreSQL benchmark
                "tpcc@1.2.0 %gcc@11.4.0",                   # TPC-C benchmark
                "ycsb@0.17.0 %gcc@11.4.0",                  # Yahoo! Cloud benchmark
                
                # I/O tracing and analysis
                "blktrace@1.3.0 %gcc@11.4.0",              # Block I/O tracing
                "iotrace@1.0.0 %gcc@11.4.0",               # I/O trace analysis
                "iostat@12.7.4 %gcc@11.4.0",               # I/O statistics
                "iozone@3.506 %gcc@11.4.0",                # Filesystem benchmark
                
                # Cloud storage benchmarks
                "s3-benchmark@2.0.0 %gcc@11.4.0",          # S3 performance testing
                "azure-storage-bench@1.0.0 %gcc@11.4.0",   # Azure storage benchmark
                "gcs-benchmark@1.5.0 %gcc@11.4.0",         # Google Cloud Storage
                
                # Python storage tools
                "python@3.11.5 %gcc@11.4.0",
                "py-h5py@3.9.0 %gcc@11.4.0",               # HDF5 Python interface
                "py-netcdf4@1.6.4 %gcc@11.4.0",            # NetCDF Python interface
                "py-zarr@2.16.1 %gcc@11.4.0",              # Chunked arrays
                "py-boto3@1.28.25 %gcc@11.4.0",            # AWS SDK
                
                # Performance analysis
                "py-storage-benchmark@2.1.0 %gcc@11.4.0",  # Storage benchmarking
                "py-io-analysis@1.6.0 %gcc@11.4.0",        # I/O pattern analysis
                "py-filesystem-bench@1.3.0 %gcc@11.4.0",   # Filesystem benchmarking
                
                # Parallel I/O libraries
                "hdf5@1.14.2 %gcc@11.4.0 +mpi +fortran +cxx",
                "netcdf-c@4.9.2 %gcc@11.4.0 +mpi +parallel-netcdf",
                "parallel-netcdf@1.12.3 %gcc@11.4.0",
                "adios2@2.9.1 %gcc@11.4.0 +mpi +hdf5",
                
                # Data analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "mysql@8.0.34 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "basic_storage": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 500,
                    "cost_per_hour": 0.34,
                    "use_case": "Basic storage benchmarking and I/O testing"
                },
                "high_iops": {
                    "instance_type": "i4i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 1875,
                    "cost_per_hour": 0.60,
                    "use_case": "High IOPS storage testing with NVMe SSD"
                },
                "large_storage": {
                    "instance_type": "i4i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 7500,
                    "cost_per_hour": 2.40,
                    "use_case": "Large-scale storage benchmarking"
                },
                "parallel_io": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.05,
                    "use_case": "Parallel I/O and distributed storage testing"
                }
            },
            "estimated_cost": {
                "compute": 400,
                "storage": 300,
                "data_transfer": 100,
                "total": 800
            },
            "research_capabilities": [
                "Filesystem and block device benchmarking",
                "Parallel I/O performance analysis",
                "Cloud storage performance optimization",
                "Database I/O pattern analysis",
                "Storage scalability testing",
                "I/O bottleneck identification",
                "Data format performance comparison",
                "Storage cost-performance optimization"
            ]
        }
    
    def _get_gpu_performance_config(self) -> Dict[str, Any]:
        """GPU performance benchmarking and analysis platform"""
        return {
            "name": "GPU Performance Benchmarking & Analysis Platform",
            "description": "Comprehensive GPU performance testing, analysis, and optimization",
            "spack_packages": [
                # GPU benchmarking suites
                "gpu-burn@1.1 %gcc@11.4.0 +cuda",          # GPU stress testing
                "cuda-samples@12.2 %gcc@11.4.0 +cuda",     # CUDA sample programs
                "stream-gpu@1.0 %gcc@11.4.0 +cuda",        # GPU memory bandwidth
                "gpumembench@1.3.0 %gcc@11.4.0 +cuda",     # GPU memory benchmark
                
                # Deep learning benchmarks
                "pytorch@2.0.1 %gcc@11.4.0 +cuda",         # PyTorch framework
                "tensorflow@2.13.0 %gcc@11.4.0 +cuda",     # TensorFlow framework
                "mlperf@3.1.0 %gcc@11.4.0 +cuda +python",  # MLPerf benchmarks
                "nccl-tests@2.13.4 %gcc@11.4.0 +cuda",     # Multi-GPU communication
                
                # GPU computing libraries
                "cuda@12.2.2 %gcc@11.4.0",                 # CUDA toolkit
                "cupy@12.2.0 %gcc@11.4.0 +cuda",           # NumPy-like library for GPU
                "cudf@23.08.00 %gcc@11.4.0 +cuda",         # GPU DataFrame library
                "thrust@2.1.0 %gcc@11.4.0 +cuda",          # Parallel algorithms
                
                # GPU profiling tools
                "nvtop@3.0.1 %gcc@11.4.0",                 # GPU process monitor
                "nvprof@12.2 %gcc@11.4.0 +cuda",           # CUDA profiler
                "nsight-systems@2023.2.3 %gcc@11.4.0 +cuda", # System profiler
                "nsight-compute@2023.2.0 %gcc@11.4.0 +cuda", # Kernel profiler
                
                # AWS-optimized Multi-GPU benchmarks with EFA
                "nccl@2.18.3 %gcc@11.4.0 +cuda",           # NVIDIA Collective Communications
                "aws-ofi-nccl@1.7.0 %gcc@11.4.0 +cuda",    # AWS OFI plugin for NCCL over EFA
                "openmpi@4.1.5 %gcc@11.4.0 +cuda +fabrics", # MPI with CUDA and EFA support
                "ucx@1.14.1 %gcc@11.4.0 +cuda +verbs +mlx +ib_hw_tm", # Unified communication with EFA
                
                # Python GPU tools
                "python@3.11.5 %gcc@11.4.0",
                "py-nvidia-ml-py@12.535.77 %gcc@11.4.0",   # GPU monitoring
                "py-gpustat@1.1.1 %gcc@11.4.0",            # GPU utilization
                "py-gpu-benchmark@1.5.0 %gcc@11.4.0",      # GPU benchmarking
                "py-numba@0.57.1 %gcc@11.4.0 +cuda",       # JIT compilation for GPU
                
                # OpenCL benchmarks
                "opencl@3.0 %gcc@11.4.0",                  # OpenCL framework
                "clpeak@1.1.2 %gcc@11.4.0 +opencl",        # OpenCL peak performance
                "luxmark@3.1 %gcc@11.4.0 +opencl",         # OpenCL benchmark
                
                # Scientific computing benchmarks
                "hpgmg@1.0 %gcc@11.4.0 +cuda +mpi",        # Multigrid benchmark
                "minimd@2.0 %gcc@11.4.0 +cuda",            # Molecular dynamics
                "minife@2.1.0 %gcc@11.4.0 +cuda",          # Finite element
                "comd@1.1 %gcc@11.4.0 +cuda",              # Classical molecular dynamics
                
                # Performance analysis
                "py-performance-gpu@2.0.0 %gcc@11.4.0",    # GPU performance analysis
                "py-cuda-profiler@1.3.0 %gcc@11.4.0",      # CUDA profiling tools
                "py-roofline-gpu@1.1.0 %gcc@11.4.0",       # GPU roofline model
                
                # Data analysis and visualization
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0",
                "nvcc@12.2.2"
            ],
            "aws_instance_recommendations": {
                "single_gpu": {
                    "instance_type": "g5.xlarge",
                    "vcpus": 4,
                    "memory_gb": 16,
                    "gpu_count": 1,
                    "gpu_memory": "24 GB",
                    "cost_per_hour": 1.006,
                    "use_case": "Single GPU benchmarking and development"
                },
                "multi_gpu": {
                    "instance_type": "g5.12xlarge",
                    "vcpus": 48,
                    "memory_gb": 192,
                    "gpu_count": 4,
                    "gpu_memory": "96 GB total",
                    "cost_per_hour": 5.672,
                    "use_case": "Multi-GPU scaling and communication testing"
                },
                "high_memory_gpu": {
                    "instance_type": "p4d.24xlarge",
                    "vcpus": 96,
                    "memory_gb": 1152,
                    "gpu_count": 8,
                    "gpu_memory": "320 GB total",
                    "cost_per_hour": 32.77,
                    "use_case": "High-memory GPU workloads and large model training"
                },
                "gpu_cluster": {
                    "instance_type": "g5.48xlarge",
                    "vcpus": 192,
                    "memory_gb": 768,
                    "gpu_count": 8,
                    "gpu_memory": "192 GB total",
                    "cost_per_hour": 16.288,
                    "use_case": "GPU cluster performance and distributed computing"
                }
            },
            "estimated_cost": {
                "compute": 1200,
                "gpu": 800,
                "storage": 200,
                "data_transfer": 100,
                "total": 2300
            },
            "research_capabilities": [
                "Single and multi-GPU performance benchmarking",
                "Deep learning training and inference optimization",
                "GPU memory bandwidth and latency analysis",
                "CUDA kernel performance profiling",
                "Multi-GPU communication and scaling",
                "GPU cluster performance evaluation",
                "Machine learning framework comparison",
                "GPU cost-performance optimization"
            ]
        }
    
    def _get_cloud_optimization_config(self) -> Dict[str, Any]:
        """Cloud optimization and cost analysis platform"""
        return {
            "name": "Cloud Optimization & Cost Analysis Platform",
            "description": "AWS cost optimization, resource utilization analysis, and performance tuning",
            "spack_packages": [
                # Cloud monitoring and analysis
                "python@3.11.5 %gcc@11.4.0",
                "py-boto3@1.28.25 %gcc@11.4.0",            # AWS SDK
                "py-botocore@1.31.25 %gcc@11.4.0",         # AWS core library
                "py-aws-cli@2.13.17 %gcc@11.4.0",          # AWS command line
                
                # Cost optimization tools
                "py-cloud-custodian@0.9.24 %gcc@11.4.0",   # Cloud governance
                "py-aws-cost-optimizer@2.1.0 %gcc@11.4.0", # Cost optimization
                "py-rightsizing@1.5.0 %gcc@11.4.0",        # Instance rightsizing
                "py-spot-optimizer@1.3.0 %gcc@11.4.0",     # Spot instance optimization
                
                # Performance monitoring
                "py-cloudwatch@2.4.0 %gcc@11.4.0",         # CloudWatch integration
                "py-aws-performance@1.7.0 %gcc@11.4.0",    # Performance monitoring
                "py-resource-monitor@2.3.0 %gcc@11.4.0",   # Resource monitoring
                
                # Infrastructure as Code
                "terraform@1.5.5 %gcc@11.4.0",             # Infrastructure as Code
                "py-pulumi@3.78.1 %gcc@11.4.0",            # Modern IaC
                "py-troposphere@4.4.1 %gcc@11.4.0",        # CloudFormation generation
                
                # Container optimization
                "docker@24.0.5 %gcc@11.4.0",               # Container platform
                "py-kubernetes@27.2.0 %gcc@11.4.0",        # Kubernetes client
                "py-container-optimizer@1.6.0 %gcc@11.4.0", # Container optimization
                
                # Data analysis and reporting
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-dash@2.13.0 %gcc@11.4.0",
                
                # Machine learning for optimization
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-tensorflow@2.13.0 %gcc@11.4.0",
                "py-optuna@3.3.0 %gcc@11.4.0",             # Hyperparameter optimization
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "cost_analysis": {
                    "instance_type": "t3.medium",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "storage_gb": 100,
                    "cost_per_hour": 0.042,
                    "use_case": "Cost analysis and optimization recommendations"
                },
                "optimization": {
                    "instance_type": "c6i.large",
                    "vcpus": 2,
                    "memory_gb": 4,
                    "storage_gb": 200,
                    "cost_per_hour": 0.085,
                    "use_case": "Performance optimization and resource planning"
                },
                "enterprise": {
                    "instance_type": "r6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 64,
                    "storage_gb": 500,
                    "cost_per_hour": 0.51,
                    "use_case": "Enterprise-scale cloud optimization"
                }
            },
            "estimated_cost": {
                "compute": 150,
                "storage": 50,
                "monitoring": 100,
                "data_transfer": 25,
                "total": 325
            },
            "research_capabilities": [
                "AWS cost analysis and optimization",
                "Instance rightsizing recommendations",
                "Spot instance optimization strategies",
                "Resource utilization monitoring",
                "Performance-cost trade-off analysis",
                "Auto-scaling optimization",
                "Reserved instance planning",
                "Multi-cloud cost comparison"
            ]
        }
    
    def _get_scalability_testing_config(self) -> Dict[str, Any]:
        """Scalability testing and analysis platform"""
        return {
            "name": "Scalability Testing & Analysis Platform",
            "description": "Application and system scalability testing, weak and strong scaling analysis",
            "spack_packages": [
                # Scalability testing frameworks
                "python@3.11.5 %gcc@11.4.0",
                "py-scalability@2.5.0 %gcc@11.4.0",        # Scalability testing
                "py-parallel-scaling@1.8.0 %gcc@11.4.0",   # Parallel scaling analysis
                "py-load-testing@2.2.0 %gcc@11.4.0",       # Load testing tools
                
                # Load testing tools
                "jmeter@5.5 %gcc@11.4.0 +java",            # Load testing tool
                "locust@2.16.1 %gcc@11.4.0 +python",       # Load testing framework
                "wrk@4.2.0 %gcc@11.4.0",                   # HTTP benchmarking
                "siege@4.1.6 %gcc@11.4.0",                 # HTTP load testing
                
                # Parallel programming frameworks
                "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers +pmix",
                "mpich@4.1.2 %gcc@11.4.0 +pmi",
                "openmp@12.0.1 %gcc@11.4.0",               # OpenMP runtime
                
                # Container orchestration
                "kubernetes@1.28.0 %gcc@11.4.0",           # Container orchestration
                "docker@24.0.5 %gcc@11.4.0",               # Container platform
                "py-kubernetes@27.2.0 %gcc@11.4.0",        # Kubernetes client
                
                # Distributed computing
                "spark@3.4.1 %gcc@11.4.0 +hadoop +yarn",   # Distributed computing
                "hadoop@3.3.6 %gcc@11.4.0",                # Distributed storage
                "py-dask@2023.7.1 %gcc@11.4.0",            # Parallel computing
                "py-ray@2.6.1 %gcc@11.4.0",                # Distributed computing
                
                # Performance monitoring
                "prometheus@2.45.0 %gcc@11.4.0",           # Monitoring system
                "grafana@10.0.3 %gcc@11.4.0",              # Visualization
                "collectd@5.12.0 %gcc@11.4.0 +python",     # System metrics
                
                # Benchmark analysis
                "py-scaling-analysis@1.9.0 %gcc@11.4.0",   # Scaling analysis tools
                "py-performance-model@1.4.0 %gcc@11.4.0",  # Performance modeling
                "py-parallel-efficiency@2.1.0 %gcc@11.4.0", # Parallel efficiency
                
                # Statistical analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-dash@2.13.0 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0",
                "gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "small_scale": {
                    "instance_type": "c6i.2xlarge",
                    "vcpus": 8,
                    "memory_gb": 16,
                    "storage_gb": 200,
                    "cost_per_hour": 0.34,
                    "use_case": "Small-scale scalability testing (up to 8 cores)"
                },
                "medium_scale": {
                    "instance_type": "c6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 64,
                    "storage_gb": 500,
                    "cost_per_hour": 1.36,
                    "use_case": "Medium-scale testing (up to 32 cores)"
                },
                "large_scale": {
                    "instance_type": "c6a.24xlarge",
                    "vcpus": 96,
                    "memory_gb": 192,
                    "storage_gb": 1000,
                    "cost_per_hour": 3.46,
                    "use_case": "Large-scale testing (up to 96 cores)"
                },
                "cluster_scale": {
                    "instance_type": "hpc6a.48xlarge",
                    "vcpus": 96,
                    "memory_gb": 384,
                    "storage_gb": 2000,
                    "cost_per_hour": 2.88,
                    "use_case": "Cluster-scale testing with enhanced networking"
                }
            },
            "estimated_cost": {
                "compute": 600,
                "storage": 150,
                "network": 100,
                "monitoring": 50,
                "total": 900
            },
            "research_capabilities": [
                "Strong and weak scaling analysis",
                "Parallel efficiency measurement",
                "Load testing and stress testing",
                "Container and microservice scalability",
                "Distributed system performance",
                "Auto-scaling optimization",
                "Amdahl's and Gustafson's law validation",
                "Scalability bottleneck identification"
            ]
        }
    
    def _get_comparative_analysis_config(self) -> Dict[str, Any]:
        """Comparative analysis and benchmarking platform"""
        return {
            "name": "Comparative Analysis & Benchmarking Platform",
            "description": "Cross-platform performance comparison and competitive benchmarking",
            "spack_packages": [
                # Comparative analysis tools
                "python@3.11.5 %gcc@11.4.0",
                "py-benchmark-comparison@2.3.0 %gcc@11.4.0", # Benchmark comparison
                "py-performance-comparison@1.7.0 %gcc@11.4.0", # Performance comparison
                "py-cross-platform@1.5.0 %gcc@11.4.0",     # Cross-platform analysis
                
                # Statistical comparison
                "r@4.3.1 %gcc@11.4.0 +external-lapack",
                "r-benchmark@0.3-6 %gcc@11.4.0",           # R benchmarking
                "r-multcomp@1.4-25 %gcc@11.4.0",           # Multiple comparisons
                "r-pwr@1.3-0 %gcc@11.4.0",                 # Power analysis
                
                # Multi-architecture support
                "gcc@11.4.0",                               # GNU compiler
                "intel-oneapi-compilers@2023.2.1",         # Intel compilers
                "aocc@4.0.0 %gcc@11.4.0",                  # AMD compilers
                "nvhpc@23.7 %gcc@11.4.0",                  # NVIDIA HPC SDK
                
                # Cross-platform libraries
                "openmpi@4.1.5 %gcc@11.4.0",
                "mpich@4.1.2 %gcc@11.4.0",
                "intel-oneapi-mpi@2021.10.0",
                "mvapich2@2.3.7 %gcc@11.4.0",
                
                # Performance analysis
                "py-statistical-analysis@2.4.0 %gcc@11.4.0", # Statistical analysis
                "py-hypothesis-testing@1.6.0 %gcc@11.4.0",  # Hypothesis testing
                "py-effect-size@1.3.0 %gcc@11.4.0",        # Effect size calculation
                
                # Visualization and reporting
                "py-comparison-plots@1.8.0 %gcc@11.4.0",   # Comparison visualization
                "py-report-generator@2.1.0 %gcc@11.4.0",   # Automated reporting
                "py-dashboard@1.5.0 %gcc@11.4.0",          # Interactive dashboards
                
                # Data analysis
                "py-pandas@2.0.3 %gcc@11.4.0",
                "py-numpy@1.25.2 %gcc@11.4.0",
                "py-scipy@1.11.2 %gcc@11.4.0",
                "py-scikit-learn@1.3.0 %gcc@11.4.0",
                "py-statsmodels@0.14.0 %gcc@11.4.0",
                
                # Visualization
                "py-matplotlib@3.7.2 %gcc@11.4.0",
                "py-plotly@5.15.0 %gcc@11.4.0",
                "py-seaborn@0.12.2 %gcc@11.4.0",
                "py-bokeh@3.2.2 %gcc@11.4.0",
                
                # Database systems
                "postgresql@15.4 %gcc@11.4.0",
                "sqlite@3.42.0 %gcc@11.4.0",
                "py-sqlalchemy@2.0.19 %gcc@11.4.0",
                
                # Development tools
                "git@2.41.0 %gcc@11.4.0",
                "cmake@3.27.4 %gcc@11.4.0"
            ],
            "aws_instance_recommendations": {
                "comparison": {
                    "instance_type": "c6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 32,
                    "storage_gb": 500,
                    "cost_per_hour": 0.68,
                    "use_case": "Multi-platform performance comparison"
                },
                "statistical_analysis": {
                    "instance_type": "r6i.4xlarge",
                    "vcpus": 16,
                    "memory_gb": 128,
                    "storage_gb": 500,
                    "cost_per_hour": 1.02,
                    "use_case": "Statistical analysis and hypothesis testing"
                },
                "comprehensive": {
                    "instance_type": "r6i.8xlarge",
                    "vcpus": 32,
                    "memory_gb": 256,
                    "storage_gb": 1000,
                    "cost_per_hour": 2.05,
                    "use_case": "Comprehensive comparative analysis"
                }
            },
            "estimated_cost": {
                "compute": 500,
                "storage": 100,
                "analysis": 150,
                "reporting": 50,
                "total": 800
            },
            "research_capabilities": [
                "Cross-platform performance comparison",
                "Compiler and optimization effectiveness",
                "Hardware architecture comparison",
                "Statistical significance testing",
                "Performance regression analysis",
                "Competitive benchmarking",
                "Technology adoption recommendations",
                "ROI analysis for infrastructure upgrades"
            ]
        }
    
    def generate_benchmarking_recommendation(self, workload: BenchmarkingWorkload) -> Dict[str, Any]:
        """Generate optimized AWS infrastructure recommendation for benchmarking research"""
        
        # Select appropriate configuration based on domain
        domain_config_map = {
            BenchmarkingDomain.HPC_BENCHMARKING: "hpc_benchmark_suite",
            BenchmarkingDomain.APPLICATION_PROFILING: "application_profiler",
            BenchmarkingDomain.SYSTEM_PERFORMANCE: "system_performance",
            BenchmarkingDomain.NETWORK_ANALYSIS: "network_analyzer",
            BenchmarkingDomain.STORAGE_BENCHMARKING: "storage_benchmark",
            BenchmarkingDomain.GPU_PERFORMANCE: "gpu_performance",
            BenchmarkingDomain.CLOUD_OPTIMIZATION: "cloud_optimizer",
            BenchmarkingDomain.SCALABILITY_TESTING: "scalability_tester",
            BenchmarkingDomain.COMPARATIVE_ANALYSIS: "comparative_analyzer"
        }
        
        config_name = domain_config_map.get(workload.domain, "hpc_benchmark_suite")
        base_config = self.benchmarking_configurations[config_name].copy()
        
        # Adjust configuration based on workload characteristics
        self._optimize_for_scope(base_config, workload)
        self._optimize_for_targets(base_config, workload)
        self._optimize_for_analysis_complexity(base_config, workload)
        
        # Generate cost estimates
        base_config["estimated_cost"] = self._calculate_benchmarking_costs(workload, base_config)
        
        # Add optimization recommendations
        base_config["optimization_recommendations"] = self._generate_optimization_recommendations(workload)
        
        return {
            "configuration": base_config,
            "workload_analysis": {
                "domain": workload.domain.value,
                "benchmark_scope": workload.benchmark_scope,
                "target_systems": workload.target_systems,
                "measurement_focus": workload.measurement_focus,
                "analysis_complexity": workload.analysis_complexity,
                "data_volume": f"{workload.data_volume_tb} TB"
            },
            "deployment_recommendations": self._generate_deployment_recommendations(workload),
            "estimated_cost": base_config["estimated_cost"]
        }
    
    def _optimize_for_scope(self, config: Dict[str, Any], workload: BenchmarkingWorkload):
        """Optimize configuration based on benchmark scope"""
        scope_multipliers = {
            "Single Node": 1.0,
            "Multi-Node": 2.0,
            "Cluster": 4.0,
            "Cloud": 6.0,
            "Hybrid": 8.0
        }
        
        multiplier = scope_multipliers.get(workload.benchmark_scope, 1.0)
        
        # Adjust instance recommendations based on scope
        if "aws_instance_recommendations" in config:
            for instance_config in config["aws_instance_recommendations"].values():
                if multiplier > 2.0:
                    # Scale up for multi-node and cluster benchmarks
                    instance_config["storage_gb"] = int(instance_config["storage_gb"] * multiplier)
    
    def _optimize_for_targets(self, config: Dict[str, Any], workload: BenchmarkingWorkload):
        """Optimize configuration based on target systems"""
        if "GPU" in workload.target_systems:
            # Add GPU-specific optimizations
            if "spack_packages" in config:
                config["spack_packages"].extend([
                    "cuda@12.2.2 %gcc@11.4.0",
                    "nvtop@3.0.1 %gcc@11.4.0"
                ])
        
        if "Network" in workload.target_systems:
            # Add network benchmarking tools
            if "spack_packages" in config:
                config["spack_packages"].extend([
                    "iperf3@3.13 %gcc@11.4.0",
                    "netperf@2.7.0 %gcc@11.4.0"
                ])
    
    def _optimize_for_analysis_complexity(self, config: Dict[str, Any], workload: BenchmarkingWorkload):
        """Optimize configuration based on analysis complexity"""
        if workload.analysis_complexity in ["Advanced", "Expert"]:
            # Add advanced analysis tools
            if "spack_packages" in config:
                config["spack_packages"].extend([
                    "r@4.3.1 %gcc@11.4.0 +external-lapack",
                    "py-scikit-learn@1.3.0 %gcc@11.4.0"
                ])
    
    def _calculate_benchmarking_costs(self, workload: BenchmarkingWorkload, config: Dict[str, Any]) -> Dict[str, float]:
        """Calculate estimated costs for benchmarking infrastructure"""
        base_compute = 500
        base_storage = 150
        base_network = 75
        base_data_transfer = 50
        
        # Scale costs based on benchmark scope
        scope_multipliers = {"Single Node": 1.0, "Multi-Node": 2.0, "Cluster": 4.0, "Cloud": 6.0, "Hybrid": 8.0}
        multiplier = scope_multipliers.get(workload.benchmark_scope, 1.0)
        
        # Adjust for target systems
        target_multiplier = 1.0 + (len(workload.target_systems) * 0.2)
        
        # Adjust for analysis complexity
        complexity_multipliers = {"Basic": 0.5, "Intermediate": 1.0, "Advanced": 2.0, "Expert": 3.0}
        complexity_mult = complexity_multipliers.get(workload.analysis_complexity, 1.0)
        
        compute_cost = base_compute * multiplier * target_multiplier * complexity_mult
        storage_cost = base_storage * (1 + workload.data_volume_tb / 2.0)
        network_cost = base_network * multiplier if "Network" in workload.target_systems else base_network * 0.5
        data_transfer_cost = base_data_transfer * multiplier
        
        return {
            "compute": compute_cost,
            "storage": storage_cost,
            "network": network_cost,
            "data_transfer": data_transfer_cost,
            "total": compute_cost + storage_cost + network_cost + data_transfer_cost
        }
    
    def _generate_optimization_recommendations(self, workload: BenchmarkingWorkload) -> List[str]:
        """Generate optimization recommendations for benchmarking workloads"""
        recommendations = []
        
        if workload.benchmark_scope in ["Cluster", "Cloud", "Hybrid"]:
            recommendations.append("Consider using Spot Instances for cost-effective large-scale benchmarking")
            recommendations.append("Implement auto-scaling for variable benchmark workloads")
        
        if "GPU" in workload.target_systems:
            recommendations.append("Use GPU-optimized instances (G5, P4) for GPU performance testing")
            recommendations.append("Consider AWS Batch for parallel GPU benchmarking")
        
        if "Network" in workload.target_systems:
            recommendations.append("Use EFA-enabled instances (hpc6a, hpc6id, c6in) for ultra-low latency MPI")
            recommendations.append("Configure cluster placement groups for optimal EFA performance")
            recommendations.append("Enable libfabric with EFA provider for best network performance")
            recommendations.append("Use OSU micro-benchmarks to validate EFA network performance")
        
        if workload.analysis_complexity in ["Advanced", "Expert"]:
            recommendations.append("Use memory-optimized instances for complex statistical analysis")
            recommendations.append("Consider AWS ParallelCluster for distributed benchmarking")
        
        if workload.data_volume_tb > 5.0:
            recommendations.append("Use S3 for cost-effective storage of large benchmark datasets")
            recommendations.append("Consider EFS for shared storage in multi-node benchmarks")
        
        return recommendations
    
    def _generate_deployment_recommendations(self, workload: BenchmarkingWorkload) -> Dict[str, Any]:
        """Generate deployment recommendations for benchmarking research"""
        return {
            "deployment_strategy": "cluster" if workload.benchmark_scope in ["Cluster", "Cloud", "Hybrid"] else "single-node",
            "backup_strategy": "automated_snapshots_after_benchmarks",
            "monitoring": ["CloudWatch for infrastructure metrics", "Custom benchmarking dashboards", "Performance alerting"],
            "security": ["VPC with private subnets", "IAM roles for benchmark access", "Data encryption"],
            "networking": "efa_enabled" if "Network" in workload.target_systems else "enhanced_networking",
            "storage": "high_iops" if "Storage" in workload.target_systems else "standard"
        }

if __name__ == "__main__":
    # Example usage
    pack = BenchmarkingPerformancePack()
    
    # Example workload
    workload = BenchmarkingWorkload(
        domain=BenchmarkingDomain.HPC_BENCHMARKING,
        benchmark_scope="Cluster",
        target_systems=["CPU", "Memory", "Network"],
        workload_types=["Compute", "Memory"],
        measurement_focus="Scalability",
        comparison_baseline="Previous Generation",
        duration="Hours",
        data_volume_tb=1.0,
        analysis_complexity="Advanced"
    )
    
    recommendation = pack.generate_benchmarking_recommendation(workload)
    print(json.dumps(recommendation, indent=2))