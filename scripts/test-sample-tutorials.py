#!/usr/bin/env python3
"""
Test a representative sample of tutorials to collect cost and performance data
"""

import os
import time
import json
import subprocess
from datetime import datetime

class TutorialTester:
    def __init__(self):
        self.aws_profile = "aws"
        self.test_results = []
        self.current_test = {}

        # Selected diverse sample of tutorials
        self.sample_tutorials = [
            {
                'domain': 'machine_learning',
                'reason': 'Popular domain, GPU instance, complex data processing',
                'expected_cost': 15,
                'expected_time': 25
            },
            {
                'domain': 'climate_modeling',
                'reason': 'Large data downloads, MPI/HPC processing',
                'expected_cost': 18,
                'expected_time': 20
            },
            {
                'domain': 'genomics',
                'reason': 'Bioinformatics tools, large datasets, CPU intensive',
                'expected_cost': 10,
                'expected_time': 20
            },
            {
                'domain': 'quantum_computing',
                'reason': 'Specialized frameworks, moderate compute',
                'expected_cost': 14,
                'expected_time': 20
            },
            {
                'domain': 'geospatial_research',
                'reason': 'GIS tools, satellite data processing',
                'expected_cost': 12,
                'expected_time': 20
            }
        ]

    def log_message(self, message, level="INFO"):
        """Log a timestamped message"""
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        print(f"[{timestamp}] [{level}] {message}")

    def run_aws_command(self, command, description=""):
        """Run an AWS CLI command and capture output"""
        full_command = f"aws --profile {self.aws_profile} {command}"
        self.log_message(f"Running: {full_command}")

        try:
            result = subprocess.run(
                full_command.split(),
                capture_output=True,
                text=True,
                timeout=300  # 5 minute timeout
            )

            if result.returncode == 0:
                self.log_message(f"✅ {description} successful")
                return result.stdout.strip()
            else:
                self.log_message(f"❌ {description} failed: {result.stderr}", "ERROR")
                return None
        except subprocess.TimeoutExpired:
            self.log_message(f"⏰ {description} timed out", "ERROR")
            return None
        except Exception as e:
            self.log_message(f"❌ {description} error: {e}", "ERROR")
            return None

    def get_account_costs(self):
        """Get current AWS costs for the account"""
        # Note: This would require AWS Cost Explorer API or billing data
        # For now, return placeholder
        self.log_message("📊 Getting cost data (placeholder implementation)")
        return {"placeholder": "cost_data_would_go_here"}

    def test_tutorial_validation(self, domain):
        """Test tutorial validation steps (without deployment)"""
        self.log_message(f"🔍 Testing {domain} validation steps...")

        validation_commands = [
            f"sts get-caller-identity",  # Test AWS connectivity
            f"ec2 describe-regions --query 'Regions[0].RegionName'",  # Test region access
        ]

        validation_results = {}
        for cmd in validation_commands:
            result = self.run_aws_command(cmd, f"Validation: {cmd}")
            validation_results[cmd] = result is not None

        return validation_results

    def test_data_downloads(self, domain):
        """Test data download accessibility (without actually downloading)"""
        self.log_message(f"🔍 Testing {domain} data download accessibility...")

        # Common S3 URLs used in tutorials
        test_urls = {
            'machine_learning': [
                's3://amazon-reviews-pds/tsv/',
                's3://commoncrawl/crawl-data/'
            ],
            'climate_modeling': [
                's3://era5-pds/2023/01/data/',
                's3://noaa-gfs-bdp-pds/'
            ],
            'genomics': [
                's3://1000genomes/phase3/data/',
                's3://tcga-2-open/'
            ],
            'quantum_computing': [
                's3://ibm-quantum-network/',
                's3://google-quantum-ai/'
            ],
            'geospatial_research': [
                's3://landsat-pds/',
                's3://sentinel-s2-l1c/'
            ]
        }

        domain_urls = test_urls.get(domain, [])
        accessibility_results = {}

        for url in domain_urls:
            try:
                cmd = f"s3 ls {url} --no-sign-request"
                result = self.run_aws_command(cmd, f"Data access: {url}")
                accessibility_results[url] = result is not None
            except Exception as e:
                self.log_message(f"❌ Could not test {url}: {e}", "ERROR")
                accessibility_results[url] = False

        return accessibility_results

    def estimate_instance_costs(self, domain):
        """Estimate costs for tutorial instance types"""
        # Instance type recommendations from tutorials
        instance_costs_per_hour = {
            'machine_learning': {'type': 'g5.xlarge', 'cost': 1.20},
            'climate_modeling': {'type': 'c6i.2xlarge', 'cost': 0.68},
            'genomics': {'type': 'r6i.large', 'cost': 0.24},
            'quantum_computing': {'type': 'c6i.xlarge', 'cost': 0.34},
            'geospatial_research': {'type': 'c5.xlarge', 'cost': 0.17}
        }

        return instance_costs_per_hour.get(domain, {'type': 'unknown', 'cost': 0.0})

    def run_tutorial_test(self, tutorial_info):
        """Run a complete test of a tutorial (validation only for now)"""
        domain = tutorial_info['domain']
        self.current_test = {
            'domain': domain,
            'start_time': datetime.now(),
            'reason': tutorial_info['reason'],
            'expected_cost': tutorial_info['expected_cost'],
            'expected_time': tutorial_info['expected_time']
        }

        self.log_message(f"🧪 Starting test for {domain}")
        self.log_message(f"   Reason: {tutorial_info['reason']}")

        # Phase 1: Validation tests
        validation_results = self.test_tutorial_validation(domain)
        self.current_test['validation'] = validation_results

        # Phase 2: Data accessibility tests
        data_results = self.test_data_downloads(domain)
        self.current_test['data_accessibility'] = data_results

        # Phase 3: Cost estimation
        cost_estimate = self.estimate_instance_costs(domain)
        self.current_test['cost_estimate'] = cost_estimate

        # Phase 4: Calculate readiness score
        validation_score = sum(validation_results.values()) / len(validation_results) if validation_results else 0
        data_score = sum(data_results.values()) / len(data_results) if data_results else 0
        readiness_score = (validation_score + data_score) / 2

        self.current_test['readiness_score'] = readiness_score
        self.current_test['end_time'] = datetime.now()
        self.current_test['test_duration'] = (self.current_test['end_time'] - self.current_test['start_time']).total_seconds()

        # Log results
        self.log_message(f"📊 {domain} test completed:")
        self.log_message(f"   Readiness Score: {readiness_score:.2%}")
        self.log_message(f"   Validation: {validation_score:.2%}")
        self.log_message(f"   Data Access: {data_score:.2%}")
        self.log_message(f"   Estimated Cost: ${cost_estimate['cost']:.2f}/hour ({cost_estimate['type']})")

        self.test_results.append(self.current_test.copy())
        return self.current_test

    def run_sample_tests(self):
        """Run tests on the sample tutorials"""
        self.log_message("🚀 Starting Phase 2: Sample Tutorial Testing")
        self.log_message(f"Testing {len(self.sample_tutorials)} representative tutorials")

        overall_start = datetime.now()

        for tutorial_info in self.sample_tutorials:
            try:
                self.run_tutorial_test(tutorial_info)
                self.log_message("")  # Blank line for readability
            except Exception as e:
                self.log_message(f"❌ Test failed for {tutorial_info['domain']}: {e}", "ERROR")

        overall_end = datetime.now()
        overall_duration = (overall_end - overall_start).total_seconds()

        # Generate summary report
        self.generate_summary_report(overall_duration)

    def generate_summary_report(self, overall_duration):
        """Generate a summary report of all tests"""
        self.log_message("=" * 60)
        self.log_message("📊 PHASE 2 TESTING SUMMARY")
        self.log_message("=" * 60)

        total_tests = len(self.test_results)
        avg_readiness = sum(t['readiness_score'] for t in self.test_results) / total_tests if total_tests > 0 else 0
        total_estimated_cost = sum(t['cost_estimate']['cost'] for t in self.test_results)

        self.log_message(f"Tests completed: {total_tests}")
        self.log_message(f"Overall test time: {overall_duration:.1f} seconds")
        self.log_message(f"Average readiness score: {avg_readiness:.2%}")
        self.log_message(f"Total estimated hourly cost: ${total_estimated_cost:.2f}")

        self.log_message("\n📋 Individual Results:")
        for result in self.test_results:
            domain = result['domain']
            score = result['readiness_score']
            cost = result['cost_estimate']['cost']
            instance = result['cost_estimate']['type']

            status = "✅ Ready" if score >= 0.8 else "⚠️  Issues" if score >= 0.5 else "❌ Not Ready"
            self.log_message(f"  {domain:20} | {status} | {score:.2%} | ${cost:.2f}/hr ({instance})")

        # Recommendations
        self.log_message("\n💡 RECOMMENDATIONS:")

        ready_count = sum(1 for t in self.test_results if t['readiness_score'] >= 0.8)
        if ready_count >= 3:
            self.log_message(f"✅ {ready_count} tutorials are ready for full testing")
            self.log_message("   Proceed with Phase 3: Full tutorial testing")
        else:
            self.log_message(f"⚠️  Only {ready_count} tutorials are fully ready")
            self.log_message("   Fix validation issues before full testing")

        # Save detailed results
        results_file = f"/Users/scttfrdmn/src/aws-research-wizard/tutorial_test_results_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
        with open(results_file, 'w') as f:
            json.dump(self.test_results, f, indent=2, default=str)

        self.log_message(f"\n📁 Detailed results saved to: {results_file}")

def main():
    """Run the tutorial testing"""
    tester = TutorialTester()
    tester.run_sample_tests()

if __name__ == "__main__":
    main()
