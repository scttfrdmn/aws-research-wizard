#!/usr/bin/env python3
"""
Quick tutorial readiness assessment
"""

import os
import subprocess
import json
from datetime import datetime

def run_quick_aws_test():
    """Quick AWS connectivity test"""
    try:
        result = subprocess.run(
            ["aws", "--profile", "aws", "sts", "get-caller-identity"],
            capture_output=True,
            text=True,
            timeout=10
        )
        return result.returncode == 0
    except:
        return False

def assess_tutorial_readiness():
    """Assess tutorial readiness for testing"""

    # Sample tutorials with expected characteristics
    tutorials = [
        {
            'domain': 'machine_learning',
            'instance_type': 'g5.xlarge',
            'cost_per_hour': 1.20,
            'expected_duration': 25,
            'data_size_gb': 0.53,
            'complexity': 'high'
        },
        {
            'domain': 'climate_modeling',
            'instance_type': 'c6i.2xlarge',
            'cost_per_hour': 0.68,
            'expected_duration': 20,
            'data_size_gb': 5.3,
            'complexity': 'high'
        },
        {
            'domain': 'genomics',
            'instance_type': 'r6i.large',
            'cost_per_hour': 0.24,
            'expected_duration': 20,
            'data_size_gb': 5.5,
            'complexity': 'medium'
        },
        {
            'domain': 'quantum_computing',
            'instance_type': 'c6i.xlarge',
            'cost_per_hour': 0.34,
            'expected_duration': 20,
            'data_size_gb': 4.0,
            'complexity': 'medium'
        },
        {
            'domain': 'geospatial_research',
            'instance_type': 'c5.xlarge',
            'cost_per_hour': 0.17,
            'expected_duration': 20,
            'data_size_gb': 4.8,
            'complexity': 'medium'
        }
    ]

    print("🔍 Quick Tutorial Readiness Assessment")
    print("=" * 50)

    # Test AWS connectivity
    aws_ready = run_quick_aws_test()
    print(f"AWS Connectivity: {'✅ Ready' if aws_ready else '❌ Failed'}")

    if not aws_ready:
        print("❌ Cannot proceed without AWS access")
        return

    print("\n📊 Tutorial Cost & Time Analysis:")
    print(f"{'Domain':<20} | {'Instance':<12} | {'$/hr':<6} | {'Duration':<8} | {'Est Cost':<8} | {'Data GB':<7}")
    print("-" * 80)

    total_cost = 0
    total_time = 0
    total_data = 0

    for tutorial in tutorials:
        domain = tutorial['domain']
        instance = tutorial['instance_type']
        cost_per_hour = tutorial['cost_per_hour']
        duration_min = tutorial['expected_duration']
        data_gb = tutorial['data_size_gb']

        est_cost = (cost_per_hour * duration_min) / 60
        total_cost += est_cost
        total_time += duration_min
        total_data += data_gb

        print(f"{domain:<20} | {instance:<12} | ${cost_per_hour:<5.2f} | {duration_min:<8}min | ${est_cost:<7.2f} | {data_gb:<7.1f}")

    print("-" * 80)
    print(f"{'TOTALS':<20} | {'5 tutorials':<12} | {'':>6} | {total_time:<8}min | ${total_cost:<7.2f} | {total_data:<7.1f}")

    print(f"\n📈 Testing Scenarios:")
    print(f"Sample Testing (5 tutorials): ${total_cost:.2f}, {total_time} minutes, {total_data:.1f} GB")
    print(f"Full Testing (27 tutorials):  ${total_cost * 5.4:.2f}, {total_time * 5.4:.0f} minutes, {total_data * 5.4:.1f} GB")

    print(f"\n💡 Recommendations:")
    if total_cost < 20:
        print("✅ Sample testing is cost-effective - proceed with Phase 2")
        print("✅ Full testing estimated at ${:.2f} - manageable cost".format(total_cost * 5.4))
    else:
        print("⚠️  Sample testing costs ${:.2f} - consider subset".format(total_cost))

    # Check for any obvious readiness issues
    print(f"\n🔧 Readiness Check:")

    # Check if tutorial files exist
    base_path = "/Users/scttfrdmn/src/aws-research-wizard/docs/domain-guides"
    ready_count = 0

    for tutorial in tutorials:
        domain = tutorial['domain']
        guide_path = f"{base_path}/{domain}/getting-started.md"

        if os.path.exists(guide_path):
            ready_count += 1
            print(f"✅ {domain}: Tutorial file exists")
        else:
            print(f"❌ {domain}: Tutorial file missing")

    print(f"\nReadiness Summary: {ready_count}/{len(tutorials)} tutorials ready")

    if ready_count >= 4:
        print("🚀 Proceed with actual tutorial testing!")
        return True
    else:
        print("⚠️  Fix missing tutorials before testing")
        return False

if __name__ == "__main__":
    assess_tutorial_readiness()
