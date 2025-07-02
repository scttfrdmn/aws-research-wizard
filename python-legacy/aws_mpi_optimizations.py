#!/usr/bin/env python3
"""
AWS MPI Optimizations Module
Comprehensive EFA, placement group, and instance optimization for AWS research workloads
"""

from typing import Dict, List, Any
from dataclasses import dataclass


@dataclass
class AWSMPIConfig:
    """AWS MPI configuration for research workloads"""
    efa_enabled: bool = True
    placement_group_strategy: str = "cluster"  # cluster, partition, spread
    enhanced_networking: str = "sr-iov"
    multi_gpu_communication: bool = False
    preferred_instances: List[str] = None
    slurm_integration: bool = True

    def __post_init__(self):
        if self.preferred_instances is None:
            self.preferred_instances = [
                "hpc6a.48xlarge",   # AMD EPYC, 100 Gbps networking, EFA
                "hpc6id.32xlarge",  # Intel, 200 Gbps networking, EFA
                "c6in.32xlarge",    # Intel, 200 Gbps networking, EFA
                "c6a.48xlarge",     # AMD EPYC, 50 Gbps networking, EFA
                "m6i.32xlarge",     # Intel, 50 Gbps networking, EFA
                "r6i.32xlarge"      # Intel, 50 Gbps networking, EFA
            ]


class AWSMPIOptimizer:
    """Comprehensive AWS MPI optimization for research computing"""

    @staticmethod
    def get_optimized_mpi_packages() -> List[str]:
        """Get AWS-optimized MPI package configurations"""
        return [
            # Core MPI with EFA support
            "openmpi@4.1.5 %gcc@11.4.0 +legacylaunchers +pmix +pmi +fabrics",
            "libfabric@1.18.1 %gcc@11.4.0 +verbs +mlx +efa",  # EFA support
            "aws-ofi-nccl@1.7.0 %gcc@11.4.0",  # AWS OFI plugin for NCCL
            "ucx@1.14.1 %gcc@11.4.0 +verbs +mlx +ib_hw_tm",  # Unified Communication X

            # Alternative MPI implementations
            "mpich@4.1.2 %gcc@11.4.0 +pmi +slurm +libfabric",
            "intel-oneapi-mpi@2021.10.0 %gcc@11.4.0 +libfabric",

            # Multi-GPU communication
            "nccl@2.18.3 %gcc@11.4.0 +cuda",
            "aws-ofi-nccl@1.7.0 %gcc@11.4.0 +cuda",

            # MPI-enabled parallel libraries
            "fftw@3.3.10 %gcc@11.4.0 +mpi +openmp +pfft_patches",
            "hdf5@1.14.2 %gcc@11.4.0 +mpi +threadsafe +fortran +cxx",
            "netcdf-c@4.9.2 %gcc@11.4.0 +mpi +parallel-netcdf",
            "parallel-netcdf@1.12.3 %gcc@11.4.0",

            # High-performance linear algebra
            "scalapack@2.2.0 %gcc@11.4.0 ^openmpi+fabrics",
            "mumps@5.6.1 %gcc@11.4.0 +mpi +parmetis",
            "petsc@3.19.4 %gcc@11.4.0 +mpi +hypre +metis +mumps",
            "hypre@2.29.0 %gcc@11.4.0 +mpi +openmp",

            # MPI benchmarking and profiling
            "osu-micro-benchmarks@6.2 %gcc@11.4.0 +mpi",
            "ior@3.3.0 %gcc@11.4.0 +mpi +hdf5 +netcdf",
            "mdtest@3.4.0 %gcc@11.4.0 +mpi",
            "mpiP@3.5 %gcc@11.4.0 +mpi",

            # Job scheduling and resource management
            "slurm@23.02.5 %gcc@11.4.0 +pmix +numa +nvml",
            "aws-parallelcluster@3.7.0 %gcc@11.4.0"
        ]

    @staticmethod
    def get_efa_optimized_instance_config(base_config: Dict[str, Any]) -> Dict[str, Any]:
        """Enhance instance configuration with EFA optimizations"""
        optimized_config = base_config.copy()

        # Add EFA-specific configuration
        optimized_config.update({
            "efa_enabled": True,
            "placement_group": "cluster",
            "enhanced_networking": "sr-iov",
            "network_performance": "Up to 100 Gbps" if "hpc6a" in base_config.get("instance_type", "") else "Up to 50 Gbps",
            "mpi_optimizations": [
                "EFA (Elastic Fabric Adapter) for low-latency networking",
                "Cluster placement group for minimal latency",
                "SR-IOV enhanced networking",
                "DPDK userspace networking support",
                "Hardware-accelerated MPI collective operations"
            ],
            "recommended_settings": {
                "FI_PROVIDER": "efa",
                "RDMAV_FORK_SAFE": "1",
                "FI_EFA_ENABLE_SHM_TRANSFER": "1",
                "NCCL_PROTO": "simple",
                "NCCL_ALGO": "ring"
            }
        })

        return optimized_config

    @staticmethod
    def get_mpi_performance_tuning() -> Dict[str, Any]:
        """Get MPI performance tuning recommendations"""
        return {
            "environment_variables": {
                # EFA-specific settings
                "FI_PROVIDER": "efa",
                "FI_EFA_ENABLE_SHM_TRANSFER": "1",
                "FI_EFA_USE_DEVICE_RDMA": "1",

                # OpenMPI tuning
                "OMPI_MCA_btl": "^vader,tcp,openib,uct",
                "OMPI_MCA_pml": "ucx",
                "OMPI_MCA_osc": "ucx",

                # UCX tuning for EFA
                "UCX_TLS": "rc,ud,sm,self",
                "UCX_NET_DEVICES": "efa0:1",

                # NCCL tuning for multi-GPU
                "NCCL_PROTO": "simple",
                "NCCL_ALGO": "ring",
                "NCCL_DEBUG": "WARN",

                # Memory and performance
                "RDMAV_FORK_SAFE": "1",
                "MALLOC_ARENA_MAX": "4"
            },

            "mpirun_flags": [
                "--mca pml ucx",
                "--mca btl ^vader,tcp,openib,uct",
                "--mca osc ucx",
                "--bind-to core",
                "--map-by socket:PE=1",
                "--report-bindings"
            ],

            "srun_flags": [
                "--mpi=pmix",
                "--ntasks-per-node=96",  # Adjust based on instance type
                "--cpus-per-task=1",
                "--cpu-bind=cores"
            ],

            "placement_optimization": {
                "cluster_placement_group": True,
                "single_az_deployment": True,
                "dedicated_tenancy": False,  # Usually not needed for HPC
                "enhanced_networking": True
            }
        }

    @staticmethod
    def get_efa_setup_commands() -> List[str]:
        """Get EFA setup and verification commands"""
        return [
            # EFA installation and setup
            "curl -O https://s3-us-west-2.amazonaws.com/aws-efa-installer/aws-efa-installer-1.26.1.tar.gz",
            "tar -xf aws-efa-installer-1.26.1.tar.gz",
            "cd aws-efa-installer && sudo ./efa_installer.sh -y -g",

            # Verify EFA installation
            "fi_info -p efa",
            "ibv_devinfo",
            "/opt/amazon/efa/bin/fi_pingpong",

            # Load EFA module
            "sudo modprobe efa",
            "lsmod | grep efa",

            # Test MPI over EFA
            "mpirun -n 2 -N 1 --mca pml ucx /opt/amazon/efa/bin/fi_pingpong",

            # Performance testing
            "mpirun -n 4 -N 2 --hostfile hosts --mca pml ucx osu_latency",
            "mpirun -n 4 -N 2 --hostfile hosts --mca pml ucx osu_bw"
        ]

    @staticmethod
    def generate_slurm_efa_config() -> str:
        """Generate Slurm configuration for EFA"""
        return """
# Slurm configuration for EFA-enabled clusters
NodeName=compute-[1-32] CPUs=96 Sockets=2 CoresPerSocket=24 ThreadsPerCore=2 State=UNKNOWN
PartitionName=hpc Nodes=compute-[1-32] Default=YES MaxTime=INFINITE State=UP

# MPI plugin configuration
MpiDefault=pmix
PrologFlags=Alloc,NoHold
Prolog=/opt/slurm/etc/prolog.d/50-setup-efa

# EFA-specific settings
LaunchParameters=use_interactive_step
TaskPlugin=task/affinity,task/cgroup
ProctrackType=proctrack/cgroup

# Performance optimizations
SelectType=select/cons_tres
SelectTypeParameters=CR_Core_Memory
"""

    @staticmethod
    def get_aws_parallelcluster_config_snippet() -> Dict[str, Any]:
        """Get AWS ParallelCluster configuration for EFA"""
        return {
            "Image": {
                "Os": "alinux2"
            },
            "HeadNode": {
                "InstanceType": "c6i.large",
                "Networking": {
                    "SubnetId": "subnet-12345678"
                },
                "Ssh": {
                    "KeyName": "my-key"
                }
            },
            "Scheduling": {
                "Scheduler": "slurm",
                "SlurmQueues": [
                    {
                        "Name": "hpc",
                        "ComputeResources": [
                            {
                                "Name": "hpc-queue",
                                "InstanceType": "hpc6a.48xlarge",
                                "MinCount": 0,
                                "MaxCount": 64,
                                "Efa": {
                                    "Enabled": True,
                                    "GdrSupport": False
                                },
                                "PlacementGroup": {
                                    "Enabled": True
                                }
                            }
                        ],
                        "Networking": {
                            "SubnetIds": ["subnet-12345678"],
                            "PlacementGroup": {
                                "Enabled": True
                            }
                        }
                    }
                ]
            },
            "SharedStorage": [
                {
                    "MountDir": "/shared",
                    "Name": "ebs-shared",
                    "StorageType": "Ebs",
                    "EbsSettings": {
                        "VolumeType": "gp3",
                        "Size": 1000,
                        "Throughput": 1000,
                        "Iops": 10000
                    }
                },
                {
                    "MountDir": "/scratch",
                    "Name": "fsx-scratch",
                    "StorageType": "FsxLustre",
                    "FsxLustreSettings": {
                        "StorageCapacity": 2400,
                        "DeploymentType": "SCRATCH_2",
                        "PerUnitStorageThroughput": 1000
                    }
                }
            ]
        }


def apply_aws_mpi_optimizations(spack_packages: List[str], instance_config: Dict[str, Any]) -> Dict[str, Any]:
    """Apply comprehensive AWS MPI optimizations to a research pack"""

    # Get optimized MPI packages
    optimized_packages = AWSMPIOptimizer.get_optimized_mpi_packages()

    # Merge with existing packages, avoiding duplicates
    all_packages = list(set(spack_packages + optimized_packages))

    # Get EFA-optimized instance configuration
    optimized_instance_config = AWSMPIOptimizer.get_efa_optimized_instance_config(instance_config)

    # Get performance tuning configuration
    performance_config = AWSMPIOptimizer.get_mpi_performance_tuning()

    return {
        "spack_packages": all_packages,
        "aws_instance_config": optimized_instance_config,
        "mpi_performance_tuning": performance_config,
        "setup_commands": AWSMPIOptimizer.get_efa_setup_commands(),
        "parallelcluster_config": AWSMPIOptimizer.get_aws_parallelcluster_config_snippet(),
        "optimization_notes": [
            "EFA (Elastic Fabric Adapter) enabled for ultra-low latency MPI communication",
            "Cluster placement groups ensure minimal network latency between nodes",
            "UCX and libfabric optimized for AWS infrastructure",
            "Multi-GPU NCCL communication optimized for AWS instances",
            "Slurm scheduler configured for optimal resource allocation",
            "Up to 32-node scaling with linear performance for most MPI workloads"
        ]
    }


if __name__ == "__main__":
    # Example usage
    optimizer = AWSMPIOptimizer()

    # Example instance configuration
    base_config = {
        "instance_type": "hpc6a.48xlarge",
        "vcpus": 96,
        "memory_gb": 384,
        "cost_per_hour": 2.88
    }

    # Apply optimizations
    result = apply_aws_mpi_optimizations(
        spack_packages=["python@3.11.5", "numpy@1.25.2"],
        instance_config=base_config
    )

    print("AWS MPI Optimized Configuration:")
    print(f"Packages: {len(result['spack_packages'])}")
    print(f"Performance tuning options: {len(result['mpi_performance_tuning']['environment_variables'])}")
    print(f"Setup commands: {len(result['setup_commands'])}")
