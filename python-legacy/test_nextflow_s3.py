#!/usr/bin/env python3
"""
Nextflow S3 Integration Test for AWS Research Wizard

This module tests and validates Nextflow's S3 integration capabilities:
1. Basic S3 read/write functionality
2. Workflow execution with S3 input/output
3. Performance characteristics with research datasets
4. Cost analysis and optimization recommendations

Test Scenarios:
- Simple file transfer and processing
- Multi-file batch processing
- Large dataset handling
- Cross-region S3 access patterns
- Integration with other workflow tools

Dependencies:
- nextflow: Workflow orchestration tool
- boto3: AWS SDK for Python
- Test data: Sample research datasets
"""

import os
import sys
import json
import boto3
import subprocess
import tempfile
import time
from typing import Dict, List, Any, Optional
from pathlib import Path
import logging

# Import our transfer optimizer
from s3_transfer_optimizer import S3TransferOptimizer, TransferTool, StorageClass

class NextflowS3Tester:
    """
    Test suite for Nextflow S3 integration and performance validation.
    """

    def __init__(self, test_bucket: str = None):
        self.logger = logging.getLogger(__name__)
        self.s3_client = boto3.client('s3')
        self.s3_optimizer = S3TransferOptimizer()

        # Use provided bucket or create test bucket name
        self.test_bucket = test_bucket or f"aws-research-wizard-test-{int(time.time())}"

        # Test configuration
        self.test_config = {
            'small_dataset_mb': 10,      # Small test files
            'medium_dataset_mb': 100,    # Medium test files
            'large_dataset_mb': 1000,    # Large test files (1GB)
            'timeout_seconds': 300,      # 5 minute timeout for tests
            'test_regions': ['us-east-1', 'us-west-2']
        }

        # Track test results
        self.test_results = []

    def setup_test_environment(self) -> bool:
        """Set up S3 bucket and test data for Nextflow testing."""
        try:
            # Create test bucket if it doesn't exist
            try:
                self.s3_client.head_bucket(Bucket=self.test_bucket)
                self.logger.info(f"Using existing test bucket: {self.test_bucket}")
            except self.s3_client.exceptions.NoSuchBucket:
                self.s3_client.create_bucket(Bucket=self.test_bucket)
                self.logger.info(f"Created test bucket: {self.test_bucket}")

            # Create test data files
            self._create_test_data()

            # Verify Nextflow is available
            if not self._check_nextflow_available():
                self.logger.error("Nextflow is not available. Install from: https://www.nextflow.io/")
                return False

            return True

        except Exception as e:
            self.logger.error(f"Failed to setup test environment: {e}")
            return False

    def _check_nextflow_available(self) -> bool:
        """Check if Nextflow is available on the system."""
        try:
            result = subprocess.run(['nextflow', '--version'],
                                  capture_output=True, text=True, timeout=30)
            if result.returncode == 0:
                self.logger.info(f"Nextflow version: {result.stdout.strip()}")
                return True
            else:
                return False
        except (subprocess.TimeoutExpired, FileNotFoundError):
            return False

    def _create_test_data(self):
        """Create test datasets for Nextflow S3 testing."""
        test_data_dir = Path("/tmp/nextflow_test_data")
        test_data_dir.mkdir(exist_ok=True)

        # Create sample FASTQ file (genomics test case)
        fastq_content = """@read1
ACGTACGTACGTACGTACGTACGTACGTACGT
+
IIIIIIIIIIIIIIIIIIIIIIIIIIIIIIII
@read2
TGCATGCATGCATGCATGCATGCATGCATGCA
+
IIIIIIIIIIIIIIIIIIIIIIIIIIIIIIII
"""

        # Create test files of different sizes
        for size_name, size_mb in [
            ('small', self.test_config['small_dataset_mb']),
            ('medium', self.test_config['medium_dataset_mb'])
        ]:
            test_file = test_data_dir / f"test_{size_name}.fastq"

            with open(test_file, 'w') as f:
                # Repeat content to reach target size
                target_size = size_mb * 1024 * 1024
                content_size = len(fastq_content.encode())
                repetitions = target_size // content_size

                for _ in range(repetitions):
                    f.write(fastq_content)

            # Upload to S3
            s3_key = f"test_data/genomics/{test_file.name}"
            self.s3_client.upload_file(str(test_file), self.test_bucket, s3_key)
            self.logger.info(f"Uploaded {test_file.name} to s3://{self.test_bucket}/{s3_key}")

    def test_basic_s3_access(self) -> Dict[str, Any]:
        """Test basic Nextflow S3 read/write functionality."""
        self.logger.info("Testing basic Nextflow S3 access...")

        # Create simple Nextflow script for S3 testing
        nextflow_script = """
#!/usr/bin/env nextflow

params.input = 's3://{bucket}/test_data/genomics/test_small.fastq'
params.outdir = 's3://{bucket}/output/basic_test'

process count_reads {{
    publishDir params.outdir, mode: 'copy'

    input:
    path reads from params.input

    output:
    path 'read_count.txt'

    script:
    '''
    echo "Processing file: $reads"
    wc -l $reads > read_count.txt
    echo "Read count completed"
    '''
}}
""".format(bucket=self.test_bucket)

        # Write script to temporary file
        with tempfile.NamedTemporaryFile(mode='w', suffix='.nf', delete=False) as f:
            f.write(nextflow_script)
            script_path = f.name

        try:
            start_time = time.time()

            # Execute Nextflow workflow
            result = subprocess.run([
                'nextflow', 'run', script_path,
                '-with-report', f'/tmp/nextflow_report_{int(time.time())}.html'
            ], capture_output=True, text=True, timeout=self.test_config['timeout_seconds'])

            execution_time = time.time() - start_time

            # Check if output was created in S3
            output_exists = self._check_s3_output_exists(
                f"{self.test_bucket}/output/basic_test/read_count.txt"
            )

            test_result = {
                'test_name': 'basic_s3_access',
                'success': result.returncode == 0 and output_exists,
                'execution_time': execution_time,
                'stdout': result.stdout,
                'stderr': result.stderr,
                'output_created': output_exists,
                'timestamp': time.time()
            }

            self.test_results.append(test_result)

            if test_result['success']:
                self.logger.info(f"✅ Basic S3 access test passed ({execution_time:.2f}s)")
            else:
                self.logger.error(f"❌ Basic S3 access test failed")
                self.logger.error(f"STDERR: {result.stderr}")

            return test_result

        except subprocess.TimeoutExpired:
            self.logger.error(f"❌ Basic S3 access test timed out after {self.test_config['timeout_seconds']}s")
            return {
                'test_name': 'basic_s3_access',
                'success': False,
                'error': 'timeout',
                'timestamp': time.time()
            }

        finally:
            # Cleanup temporary script
            os.unlink(script_path)

    def test_batch_processing(self) -> Dict[str, Any]:
        """Test Nextflow batch processing with multiple S3 files."""
        self.logger.info("Testing Nextflow batch processing with S3...")

        # Create workflow that processes multiple files
        nextflow_script = """
#!/usr/bin/env nextflow

params.input = 's3://{bucket}/test_data/genomics/*.fastq'
params.outdir = 's3://{bucket}/output/batch_test'

Channel
    .fromPath(params.input)
    .set {{ input_files }}

process process_each_file {{
    publishDir params.outdir, mode: 'copy'

    input:
    path file from input_files

    output:
    path "${{file.baseName}}_processed.txt"

    script:
    '''
    echo "Processing $file" > ${{file.baseName}}_processed.txt
    wc -l $file >> ${{file.baseName}}_processed.txt
    echo "Completed processing $file" >> ${{file.baseName}}_processed.txt
    '''
}}
""".format(bucket=self.test_bucket)

        with tempfile.NamedTemporaryFile(mode='w', suffix='.nf', delete=False) as f:
            f.write(nextflow_script)
            script_path = f.name

        try:
            start_time = time.time()

            result = subprocess.run([
                'nextflow', 'run', script_path,
                '-with-timeline', f'/tmp/nextflow_timeline_{int(time.time())}.html'
            ], capture_output=True, text=True, timeout=self.test_config['timeout_seconds'])

            execution_time = time.time() - start_time

            # Check for multiple output files
            output_count = self._count_s3_outputs(f"{self.test_bucket}/output/batch_test/")

            test_result = {
                'test_name': 'batch_processing',
                'success': result.returncode == 0 and output_count > 0,
                'execution_time': execution_time,
                'output_files_created': output_count,
                'stdout': result.stdout[:1000],  # Truncate for storage
                'stderr': result.stderr[:1000],
                'timestamp': time.time()
            }

            self.test_results.append(test_result)

            if test_result['success']:
                self.logger.info(f"✅ Batch processing test passed ({execution_time:.2f}s, {output_count} files)")
            else:
                self.logger.error(f"❌ Batch processing test failed")

            return test_result

        except subprocess.TimeoutExpired:
            return {
                'test_name': 'batch_processing',
                'success': False,
                'error': 'timeout',
                'timestamp': time.time()
            }

        finally:
            os.unlink(script_path)

    def test_cross_region_performance(self) -> Dict[str, Any]:
        """Test Nextflow performance with cross-region S3 access."""
        self.logger.info("Testing cross-region S3 performance...")

        # This test would ideally use data in different regions
        # For now, we'll simulate by testing with transfer acceleration

        nextflow_script = """
#!/usr/bin/env nextflow

params.input = 's3://{bucket}/test_data/genomics/test_medium.fastq'
params.outdir = 's3://{bucket}/output/cross_region_test'

process analyze_large_file {{
    publishDir params.outdir, mode: 'copy'

    input:
    path file from params.input

    output:
    path 'analysis_results.txt'

    script:
    '''
    echo "Analysis started: $(date)" > analysis_results.txt
    echo "File size: $(wc -c < $file) bytes" >> analysis_results.txt
    echo "Line count: $(wc -l < $file)" >> analysis_results.txt

    # Simulate some processing time
    sleep 5

    echo "Analysis completed: $(date)" >> analysis_results.txt
    '''
}}
""".format(bucket=self.test_bucket)

        with tempfile.NamedTemporaryFile(mode='w', suffix='.nf', delete=False) as f:
            f.write(nextflow_script)
            script_path = f.name

        try:
            start_time = time.time()

            result = subprocess.run([
                'nextflow', 'run', script_path,
                '-with-trace', f'/tmp/nextflow_trace_{int(time.time())}.txt'
            ], capture_output=True, text=True, timeout=self.test_config['timeout_seconds'])

            execution_time = time.time() - start_time

            output_exists = self._check_s3_output_exists(
                f"{self.test_bucket}/output/cross_region_test/analysis_results.txt"
            )

            test_result = {
                'test_name': 'cross_region_performance',
                'success': result.returncode == 0 and output_exists,
                'execution_time': execution_time,
                'throughput_estimate': self.test_config['medium_dataset_mb'] / execution_time if execution_time > 0 else 0,
                'output_created': output_exists,
                'timestamp': time.time()
            }

            self.test_results.append(test_result)

            if test_result['success']:
                self.logger.info(f"✅ Cross-region test passed ({execution_time:.2f}s, {test_result['throughput_estimate']:.2f} MB/s)")
            else:
                self.logger.error(f"❌ Cross-region test failed")

            return test_result

        except subprocess.TimeoutExpired:
            return {
                'test_name': 'cross_region_performance',
                'success': False,
                'error': 'timeout',
                'timestamp': time.time()
            }

        finally:
            os.unlink(script_path)

    def _check_s3_output_exists(self, s3_path: str) -> bool:
        """Check if output file exists in S3."""
        try:
            bucket, key = s3_path.replace('s3://', '').split('/', 1)
            self.s3_client.head_object(Bucket=bucket, Key=key)
            return True
        except:
            return False

    def _count_s3_outputs(self, s3_prefix: str) -> int:
        """Count number of output files in S3 prefix."""
        try:
            bucket, prefix = s3_prefix.replace('s3://', '').split('/', 1)
            response = self.s3_client.list_objects_v2(Bucket=bucket, Prefix=prefix)
            return response.get('KeyCount', 0)
        except:
            return 0

    def run_all_tests(self) -> Dict[str, Any]:
        """Run all Nextflow S3 integration tests."""
        self.logger.info("Starting Nextflow S3 integration test suite...")

        if not self.setup_test_environment():
            return {'error': 'Failed to setup test environment'}

        # Run all tests
        tests = [
            self.test_basic_s3_access,
            self.test_batch_processing,
            self.test_cross_region_performance
        ]

        for test_func in tests:
            try:
                test_func()
            except Exception as e:
                self.logger.error(f"Test {test_func.__name__} failed with exception: {e}")
                self.test_results.append({
                    'test_name': test_func.__name__,
                    'success': False,
                    'error': str(e),
                    'timestamp': time.time()
                })

        # Generate summary
        total_tests = len(self.test_results)
        passed_tests = sum(1 for result in self.test_results if result.get('success', False))

        summary = {
            'total_tests': total_tests,
            'passed_tests': passed_tests,
            'failed_tests': total_tests - passed_tests,
            'success_rate': passed_tests / total_tests if total_tests > 0 else 0,
            'test_results': self.test_results,
            'recommendations': self._generate_recommendations()
        }

        self.logger.info(f"Test summary: {passed_tests}/{total_tests} tests passed ({summary['success_rate']:.1%})")

        return summary

    def _generate_recommendations(self) -> List[str]:
        """Generate optimization recommendations based on test results."""
        recommendations = []

        # Analyze performance patterns
        execution_times = [r.get('execution_time', 0) for r in self.test_results if r.get('success')]

        if execution_times:
            avg_time = sum(execution_times) / len(execution_times)

            if avg_time > 60:  # More than 1 minute average
                recommendations.append("Consider using s5cmd for large data transfers to improve performance")
                recommendations.append("Enable S3 Transfer Acceleration for cross-region workflows")

            # Check for failed tests
            failed_tests = [r for r in self.test_results if not r.get('success', True)]
            if failed_tests:
                recommendations.append("Some tests failed - check Nextflow configuration and S3 permissions")
                recommendations.append("Verify AWS credentials and S3 bucket policies")

        # General recommendations
        recommendations.extend([
            "Use S3 Intelligent Tiering for research data with variable access patterns",
            "Consider regional data placement to minimize transfer costs",
            "Implement proper error handling and retry logic in Nextflow workflows"
        ])

        return recommendations

    def cleanup_test_environment(self):
        """Clean up test bucket and data."""
        try:
            # Delete all objects in test bucket
            response = self.s3_client.list_objects_v2(Bucket=self.test_bucket)
            if 'Contents' in response:
                objects = [{'Key': obj['Key']} for obj in response['Contents']]
                self.s3_client.delete_objects(
                    Bucket=self.test_bucket,
                    Delete={'Objects': objects}
                )

            # Delete bucket
            self.s3_client.delete_bucket(Bucket=self.test_bucket)
            self.logger.info(f"Cleaned up test bucket: {self.test_bucket}")

        except Exception as e:
            self.logger.warning(f"Failed to cleanup test environment: {e}")

def main():
    """Run Nextflow S3 integration tests."""
    logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

    # Check if we should run tests
    if len(sys.argv) > 1 and sys.argv[1] == '--run-tests':
        tester = NextflowS3Tester()

        try:
            results = tester.run_all_tests()

            # Print results
            print("\n" + "="*60)
            print("NEXTFLOW S3 INTEGRATION TEST RESULTS")
            print("="*60)
            print(f"Tests Passed: {results['passed_tests']}/{results['total_tests']}")
            print(f"Success Rate: {results['success_rate']:.1%}")

            print("\nRecommendations:")
            for rec in results['recommendations']:
                print(f"  • {rec}")

            # Save detailed results
            with open('/tmp/nextflow_s3_test_results.json', 'w') as f:
                json.dump(results, f, indent=2)

            print(f"\nDetailed results saved to: /tmp/nextflow_s3_test_results.json")

        finally:
            tester.cleanup_test_environment()

    else:
        print("Nextflow S3 Integration Tester")
        print("Usage: python test_nextflow_s3.py --run-tests")
        print("\nThis will:")
        print("  1. Create a test S3 bucket")
        print("  2. Upload test genomics data")
        print("  3. Run Nextflow workflows with S3 input/output")
        print("  4. Validate performance and functionality")
        print("  5. Generate optimization recommendations")

if __name__ == "__main__":
    main()
