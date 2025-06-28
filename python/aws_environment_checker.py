#!/usr/bin/env python3
"""
AWS Environment Checker
Validates AWS account setup, quotas, and prerequisites for AWS Research Wizard
"""

import boto3
import json
import time
from typing import Dict, List, Any, Optional, Tuple
from dataclasses import dataclass, asdict
from enum import Enum
import sys
import argparse
from botocore.exceptions import ClientError, NoCredentialsError, ProfileNotFound

class CheckStatus(Enum):
    PASS = "PASS"
    WARN = "WARN"
    FAIL = "FAIL"
    SKIP = "SKIP"

@dataclass
class CheckResult:
    name: str
    status: CheckStatus
    message: str
    recommendation: Optional[str] = None
    details: Optional[Dict[str, Any]] = None

@dataclass
class QuotaCheck:
    service: str
    quota_name: str
    quota_code: str
    minimum_required: int
    recommended: int
    current_value: Optional[int] = None

class AWSEnvironmentChecker:
    """
    Comprehensive AWS environment checker for AWS Research Wizard
    Validates account setup, permissions, quotas, and prerequisites
    """
    
    def __init__(self, profile: Optional[str] = None, region: Optional[str] = None):
        self.profile = profile
        self.region = region or 'us-east-1'
        self.session = None
        self.account_id = None
        self.results: List[CheckResult] = []
        
        # Initialize AWS session
        try:
            if profile:
                self.session = boto3.Session(profile_name=profile, region_name=self.region)
            else:
                self.session = boto3.Session(region_name=self.region)
            
            # Get account ID
            sts = self.session.client('sts')
            self.account_id = sts.get_caller_identity()['Account']
            
        except (NoCredentialsError, ProfileNotFound) as e:
            self.results.append(CheckResult(
                name="AWS Credentials",
                status=CheckStatus.FAIL,
                message=f"AWS credentials not found or invalid: {e}",
                recommendation="Configure AWS credentials using 'aws configure' or set environment variables"
            ))
            self.session = None
    
    def run_all_checks(self) -> List[CheckResult]:
        """Run all environment checks"""
        if not self.session:
            return self.results
        
        print("🧙‍♂️ AWS Research Wizard - Environment Checker")
        print("=" * 60)
        print(f"Account ID: {self.account_id}")
        print(f"Region: {self.region}")
        print(f"Profile: {self.profile or 'default'}")
        print("=" * 60)
        
        # Core AWS checks
        self._check_credentials()
        self._check_permissions()
        self._check_regions()
        
        # Service availability
        self._check_service_availability()
        
        # Quota checks
        self._check_compute_quotas()
        self._check_storage_quotas()
        self._check_network_quotas()
        self._check_specialized_quotas()
        
        # Security and compliance
        self._check_security_setup()
        
        # Cost and billing
        self._check_billing_setup()
        
        # Research-specific checks
        self._check_research_prerequisites()
        
        return self.results
    
    def _check_credentials(self):
        """Check AWS credentials and basic access"""
        try:
            sts = self.session.client('sts')
            identity = sts.get_caller_identity()
            
            self.results.append(CheckResult(
                name="AWS Credentials",
                status=CheckStatus.PASS,
                message=f"Successfully authenticated as {identity['Arn']}",
                details={"account_id": identity['Account'], "user_id": identity['UserId']}
            ))
            
        except Exception as e:
            self.results.append(CheckResult(
                name="AWS Credentials",
                status=CheckStatus.FAIL,
                message=f"Credential validation failed: {e}",
                recommendation="Check AWS credentials and permissions"
            ))
    
    def _check_permissions(self):
        """Check essential IAM permissions"""
        required_permissions = [
            ('ec2', 'describe_instances'),
            ('ec2', 'describe_instance_types'),
            ('s3', 'list_buckets'),
            ('iam', 'get_user'),
            ('pricing', 'get_products'),
            ('servicequotas', 'list_service_quotas')
        ]
        
        failed_permissions = []
        
        for service, action in required_permissions:
            try:
                client = self.session.client(service)
                
                if service == 'ec2' and action == 'describe_instances':
                    client.describe_instances(MaxResults=5)
                elif service == 'ec2' and action == 'describe_instance_types':
                    client.describe_instance_types(MaxResults=5)
                elif service == 's3' and action == 'list_buckets':
                    client.list_buckets()
                elif service == 'iam' and action == 'get_user':
                    try:
                        client.get_user()
                    except ClientError as e:
                        if 'AccessDenied' not in str(e):
                            continue  # Some roles can't get user info but can perform other actions
                elif service == 'pricing' and action == 'get_products':
                    client.get_products(ServiceCode='AmazonEC2', MaxResults=1)
                elif service == 'servicequotas' and action == 'list_service_quotas':
                    client.list_service_quotas(ServiceCode='ec2', MaxResults=1)
                    
            except ClientError as e:
                if 'AccessDenied' in str(e) or 'UnauthorizedOperation' in str(e):
                    failed_permissions.append(f"{service}:{action}")
            except Exception as e:
                failed_permissions.append(f"{service}:{action} ({e})")
        
        if failed_permissions:
            self.results.append(CheckResult(
                name="IAM Permissions",
                status=CheckStatus.WARN,
                message=f"Missing permissions for: {', '.join(failed_permissions)}",
                recommendation="Ensure IAM user/role has required permissions for research workloads"
            ))
        else:
            self.results.append(CheckResult(
                name="IAM Permissions",
                status=CheckStatus.PASS,
                message="All essential permissions available"
            ))
    
    def _check_regions(self):
        """Check available regions and research-friendly regions"""
        try:
            ec2 = self.session.client('ec2')
            regions = ec2.describe_regions()['Regions']
            available_regions = [r['RegionName'] for r in regions]
            
            # Research-friendly regions (good price/performance, full service availability)
            research_regions = [
                'us-east-1', 'us-east-2', 'us-west-2',
                'eu-west-1', 'eu-central-1',
                'ap-southeast-1', 'ap-northeast-1'
            ]
            
            available_research_regions = [r for r in research_regions if r in available_regions]
            
            if len(available_research_regions) >= 3:
                status = CheckStatus.PASS
                message = f"Good regional coverage: {len(available_research_regions)} research-friendly regions"
            else:
                status = CheckStatus.WARN
                message = f"Limited regional options: {len(available_research_regions)} research-friendly regions"
            
            self.results.append(CheckResult(
                name="Regional Availability",
                status=status,
                message=message,
                details={
                    "total_regions": len(available_regions),
                    "research_regions": available_research_regions,
                    "current_region": self.region
                }
            ))
            
        except Exception as e:
            self.results.append(CheckResult(
                name="Regional Availability",
                status=CheckStatus.FAIL,
                message=f"Could not check regions: {e}"
            ))
    
    def _check_service_availability(self):
        """Check availability of key AWS services"""
        services_to_check = [
            ('EC2', 'ec2'),
            ('S3', 's3'),
            ('EBS', 'ec2'),
            ('VPC', 'ec2'),
            ('IAM', 'iam'),
            ('CloudWatch', 'cloudwatch'),
            ('AWS Batch', 'batch'),
            ('EFS', 'efs'),
            ('FSx', 'fsx'),
            ('ParallelCluster', None)  # CLI-based service
        ]
        
        available_services = []
        unavailable_services = []
        
        for service_name, client_name in services_to_check:
            try:
                if client_name:
                    client = self.session.client(client_name)
                    
                    # Simple API call to test service availability
                    if client_name == 'ec2':
                        if service_name == 'EC2':
                            client.describe_availability_zones()
                        elif service_name == 'EBS':
                            client.describe_volumes(MaxResults=1)
                        elif service_name == 'VPC':
                            client.describe_vpcs(MaxResults=1)
                    elif client_name == 's3':
                        client.list_buckets()
                    elif client_name == 'iam':
                        client.list_roles(MaxItems=1)
                    elif client_name == 'cloudwatch':
                        client.list_metrics(MaxRecords=1)
                    elif client_name == 'batch':
                        client.describe_compute_environments(maxResults=1)
                    elif client_name == 'efs':
                        client.describe_file_systems(MaxItems=1)
                    elif client_name == 'fsx':
                        client.describe_file_systems(MaxResults=1)
                
                available_services.append(service_name)
                
            except ClientError as e:
                if 'AccessDenied' in str(e):
                    available_services.append(f"{service_name} (limited access)")
                else:
                    unavailable_services.append(f"{service_name}: {e}")
            except Exception as e:
                unavailable_services.append(f"{service_name}: {e}")
        
        if len(available_services) >= 8:
            status = CheckStatus.PASS
            message = f"All essential services available ({len(available_services)}/9)"
        elif len(available_services) >= 6:
            status = CheckStatus.WARN
            message = f"Most services available ({len(available_services)}/9)"
        else:
            status = CheckStatus.FAIL
            message = f"Limited service availability ({len(available_services)}/9)"
        
        self.results.append(CheckResult(
            name="Service Availability",
            status=status,
            message=message,
            details={
                "available": available_services,
                "unavailable": unavailable_services
            }
        ))
    
    def _check_compute_quotas(self):
        """Check EC2 compute quotas"""
        quota_checks = [
            QuotaCheck('ec2', 'Running On-Demand Standard instances', 'L-1216C47A', 100, 500),
            QuotaCheck('ec2', 'Running On-Demand HPC instances', 'L-F7808C92', 10, 50),
            QuotaCheck('ec2', 'Running On-Demand High Memory instances', 'L-43DA4232', 10, 50),
            QuotaCheck('ec2', 'Running On-Demand GPU instances', 'L-DB2E81BA', 10, 50),
            QuotaCheck('ec2', 'All Standard Spot Instance Requests', 'L-34B43A08', 256, 1000),
            QuotaCheck('ec2', 'EC2-VPC Elastic IPs', 'L-0263D0A3', 20, 100)
        ]
        
        self._check_service_quotas('EC2 Compute Quotas', quota_checks, 'ec2')
    
    def _check_storage_quotas(self):
        """Check storage-related quotas"""
        quota_checks = [
            QuotaCheck('ebs', 'General Purpose SSD volume storage', 'L-D18FCD1D', 10000, 50000),  # GB
            QuotaCheck('ebs', 'Provisioned IOPS SSD volume storage', 'L-B3A130E6', 1000, 10000),   # GB
            QuotaCheck('ebs', 'Cold HDD volume storage', 'L-9CF3C2EB', 10000, 100000),             # GB
            QuotaCheck('ebs', 'Throughput Optimized HDD volume storage', 'L-7A658B76', 10000, 100000),  # GB
            QuotaCheck('ebs', 'Snapshots per Region', 'L-309BACF6', 10000, 100000),
            QuotaCheck('s3', 'Buckets', 'L-DC2B2D3D', 100, 1000)
        ]
        
        # Check EBS quotas
        ebs_checks = [q for q in quota_checks if q.service == 'ebs']
        self._check_service_quotas('EBS Storage Quotas', ebs_checks, 'ebs')
        
        # S3 has different quota checking mechanism
        try:
            s3 = self.session.client('s3')
            buckets = s3.list_buckets()
            bucket_count = len(buckets['Buckets'])
            
            if bucket_count < 900:  # S3 default limit is 1000
                status = CheckStatus.PASS
                message = f"S3 bucket usage: {bucket_count}/1000"
            else:
                status = CheckStatus.WARN
                message = f"S3 bucket usage high: {bucket_count}/1000"
            
            self.results.append(CheckResult(
                name="S3 Storage Quotas",
                status=status,
                message=message,
                details={"bucket_count": bucket_count, "limit": 1000}
            ))
            
        except Exception as e:
            self.results.append(CheckResult(
                name="S3 Storage Quotas",
                status=CheckStatus.FAIL,
                message=f"Could not check S3 quotas: {e}"
            ))
    
    def _check_network_quotas(self):
        """Check networking quotas"""
        quota_checks = [
            QuotaCheck('vpc', 'VPCs per Region', 'L-F678F1CE', 5, 20),
            QuotaCheck('vpc', 'Subnets per VPC', 'L-407747CB', 200, 500),
            QuotaCheck('vpc', 'Security groups per VPC', 'L-E79EC296', 500, 2000),
            QuotaCheck('vpc', 'Rules per security group', 'L-0EA8095F', 60, 120),
            QuotaCheck('vpc', 'NAT gateways per Availability Zone', 'L-FE5A380F', 5, 20),
            QuotaCheck('ec2', 'EC2-VPC Elastic IPs', 'L-0263D0A3', 20, 100)
        ]
        
        self._check_service_quotas('Network Quotas', quota_checks, 'ec2')
    
    def _check_specialized_quotas(self):
        """Check quotas for specialized research services"""
        # AWS Batch quotas
        try:
            batch = self.session.client('batch')
            compute_envs = batch.describe_compute_environments()
            compute_env_count = len(compute_envs['computeEnvironments'])
            
            if compute_env_count < 40:  # Default limit is 50
                status = CheckStatus.PASS
                message = f"AWS Batch compute environments: {compute_env_count}/50"
            else:
                status = CheckStatus.WARN
                message = f"AWS Batch compute environments high: {compute_env_count}/50"
            
            self.results.append(CheckResult(
                name="AWS Batch Quotas",
                status=status,
                message=message
            ))
            
        except Exception as e:
            self.results.append(CheckResult(
                name="AWS Batch Quotas",
                status=CheckStatus.SKIP,
                message=f"Could not check AWS Batch: {e}"
            ))
        
        # EFS quotas
        try:
            efs = self.session.client('efs')
            filesystems = efs.describe_file_systems()
            fs_count = len(filesystems['FileSystems'])
            
            if fs_count < 900:  # Default limit is 1000
                status = CheckStatus.PASS
                message = f"EFS file systems: {fs_count}/1000"
            else:
                status = CheckStatus.WARN
                message = f"EFS file systems high: {fs_count}/1000"
            
            self.results.append(CheckResult(
                name="EFS Quotas",
                status=status,
                message=message
            ))
            
        except Exception as e:
            self.results.append(CheckResult(
                name="EFS Quotas",
                status=CheckStatus.SKIP,
                message=f"Could not check EFS: {e}"
            ))
    
    def _check_service_quotas(self, check_name: str, quota_checks: List[QuotaCheck], service_code: str):
        """Generic service quota checker"""
        try:
            quotas = self.session.client('servicequotas')
            
            passed = 0
            warned = 0
            failed = 0
            
            for quota_check in quota_checks:
                try:
                    response = quotas.get_service_quota(
                        ServiceCode=service_code,
                        QuotaCode=quota_check.quota_code
                    )
                    
                    current_value = int(response['Quota']['Value'])
                    quota_check.current_value = current_value
                    
                    if current_value >= quota_check.recommended:
                        passed += 1
                    elif current_value >= quota_check.minimum_required:
                        warned += 1
                    else:
                        failed += 1
                        
                except ClientError as e:
                    if 'NoSuchResourceException' in str(e):
                        # Quota doesn't exist in this region/service
                        continue
                    else:
                        failed += 1
                except Exception:
                    failed += 1
            
            total_checks = len(quota_checks)
            if failed == 0 and warned <= total_checks * 0.2:
                status = CheckStatus.PASS
                message = f"Quotas sufficient: {passed} optimal, {warned} adequate, {failed} insufficient"
            elif failed <= total_checks * 0.2:
                status = CheckStatus.WARN
                message = f"Some quotas may limit research: {passed} optimal, {warned} adequate, {failed} insufficient"
            else:
                status = CheckStatus.FAIL
                message = f"Insufficient quotas for research: {passed} optimal, {warned} adequate, {failed} insufficient"
            
            recommendation = None
            if failed > 0 or warned > total_checks * 0.3:
                recommendation = "Consider requesting quota increases for research workloads"
            
            self.results.append(CheckResult(
                name=check_name,
                status=status,
                message=message,
                recommendation=recommendation,
                details={"quota_details": [asdict(q) for q in quota_checks]}
            ))
            
        except Exception as e:
            self.results.append(CheckResult(
                name=check_name,
                status=CheckStatus.FAIL,
                message=f"Could not check quotas: {e}",
                recommendation="Ensure ServiceQuotas API access is available"
            ))
    
    def _check_security_setup(self):
        """Check security configuration"""
        security_checks = []
        
        # Check for default VPC
        try:
            ec2 = self.session.client('ec2')
            vpcs = ec2.describe_vpcs()
            default_vpc = None
            custom_vpcs = []
            
            for vpc in vpcs['Vpcs']:
                if vpc.get('IsDefault', False):
                    default_vpc = vpc
                else:
                    custom_vpcs.append(vpc)
            
            if custom_vpcs:
                security_checks.append(("Custom VPC", CheckStatus.PASS, "Custom VPCs configured"))
            elif default_vpc:
                security_checks.append(("Custom VPC", CheckStatus.WARN, "Only default VPC found"))
            else:
                security_checks.append(("Custom VPC", CheckStatus.FAIL, "No VPC found"))
            
        except Exception as e:
            security_checks.append(("Custom VPC", CheckStatus.FAIL, f"Could not check VPCs: {e}"))
        
        # Check CloudTrail
        try:
            cloudtrail = self.session.client('cloudtrail')
            trails = cloudtrail.describe_trails()
            
            if trails['trailList']:
                security_checks.append(("CloudTrail", CheckStatus.PASS, f"{len(trails['trailList'])} trails configured"))
            else:
                security_checks.append(("CloudTrail", CheckStatus.WARN, "No CloudTrail configured"))
                
        except Exception as e:
            security_checks.append(("CloudTrail", CheckStatus.SKIP, f"Could not check CloudTrail: {e}"))
        
        # Check MFA
        try:
            iam = self.session.client('iam')
            account_summary = iam.get_account_summary()
            
            mfa_devices = account_summary['SummaryMap'].get('MFADevices', 0)
            if mfa_devices > 0:
                security_checks.append(("MFA", CheckStatus.PASS, f"{mfa_devices} MFA devices configured"))
            else:
                security_checks.append(("MFA", CheckStatus.WARN, "No MFA devices found"))
                
        except Exception as e:
            security_checks.append(("MFA", CheckStatus.SKIP, f"Could not check MFA: {e}"))
        
        # Aggregate security results
        passed = sum(1 for _, status, _ in security_checks if status == CheckStatus.PASS)
        total = len(security_checks)
        
        if passed >= total * 0.8:
            overall_status = CheckStatus.PASS
        elif passed >= total * 0.5:
            overall_status = CheckStatus.WARN
        else:
            overall_status = CheckStatus.FAIL
        
        self.results.append(CheckResult(
            name="Security Configuration",
            status=overall_status,
            message=f"Security checks: {passed}/{total} passed",
            details={"checks": security_checks},
            recommendation="Review security best practices for research environments"
        ))
    
    def _check_billing_setup(self):
        """Check billing and cost management setup"""
        try:
            # Check if billing data is accessible
            ce = self.session.client('ce', region_name='us-east-1')  # Cost Explorer only in us-east-1
            
            # Try to get basic cost data
            response = ce.get_cost_and_usage(
                TimePeriod={
                    'Start': '2023-01-01',
                    'End': '2023-01-02'
                },
                Granularity='DAILY',
                Metrics=['BlendedCost']
            )
            
            self.results.append(CheckResult(
                name="Billing Access",
                status=CheckStatus.PASS,
                message="Cost and billing data accessible",
                recommendation="Set up billing alerts and AWS Budgets for cost control"
            ))
            
        except ClientError as e:
            if 'AccessDenied' in str(e):
                self.results.append(CheckResult(
                    name="Billing Access",
                    status=CheckStatus.WARN,
                    message="Limited billing access",
                    recommendation="Enable billing access for cost monitoring"
                ))
            else:
                self.results.append(CheckResult(
                    name="Billing Access",
                    status=CheckStatus.FAIL,
                    message=f"Billing check failed: {e}"
                ))
        except Exception as e:
            self.results.append(CheckResult(
                name="Billing Access",
                status=CheckStatus.FAIL,
                message=f"Could not check billing: {e}"
            ))
    
    def _check_research_prerequisites(self):
        """Check research-specific prerequisites"""
        research_checks = []
        
        # Check for research-optimized instance types
        try:
            ec2 = self.session.client('ec2')
            instance_types = ec2.describe_instance_types(
                Filters=[
                    {'Name': 'instance-type', 'Values': ['hpc6a.*', 'c6a.*', 'r6a.*', 'i4i.*']}
                ]
            )
            
            if instance_types['InstanceTypes']:
                research_checks.append(("Research Instances", CheckStatus.PASS, "Research-optimized instances available"))
            else:
                research_checks.append(("Research Instances", CheckStatus.WARN, "Limited research instance types"))
                
        except Exception as e:
            research_checks.append(("Research Instances", CheckStatus.SKIP, f"Could not check instances: {e}"))
        
        # Check for Spot instance availability
        try:
            ec2 = self.session.client('ec2')
            spot_prices = ec2.describe_spot_price_history(
                InstanceTypes=['c6i.large'],
                ProductDescriptions=['Linux/UNIX'],
                MaxResults=1
            )
            
            if spot_prices['SpotPriceHistory']:
                research_checks.append(("Spot Instances", CheckStatus.PASS, "Spot instances available"))
            else:
                research_checks.append(("Spot Instances", CheckStatus.WARN, "Spot instance data unavailable"))
                
        except Exception as e:
            research_checks.append(("Spot Instances", CheckStatus.SKIP, f"Could not check Spot: {e}"))
        
        # Check for enhanced networking
        try:
            ec2 = self.session.client('ec2')
            placement_groups = ec2.describe_placement_groups()
            
            research_checks.append(("Placement Groups", CheckStatus.PASS, "Placement groups supported"))
            
        except Exception as e:
            research_checks.append(("Placement Groups", CheckStatus.SKIP, f"Could not check placement groups: {e}"))
        
        # Aggregate research prerequisites
        passed = sum(1 for _, status, _ in research_checks if status == CheckStatus.PASS)
        total = len(research_checks)
        
        if passed >= total * 0.8:
            overall_status = CheckStatus.PASS
        elif passed >= total * 0.5:
            overall_status = CheckStatus.WARN
        else:
            overall_status = CheckStatus.FAIL
        
        self.results.append(CheckResult(
            name="Research Prerequisites",
            status=overall_status,
            message=f"Research features: {passed}/{total} available",
            details={"checks": research_checks}
        ))
    
    def print_results(self):
        """Print formatted results"""
        print("\n" + "=" * 60)
        print("🔍 ENVIRONMENT CHECK RESULTS")
        print("=" * 60)
        
        status_counts = {
            CheckStatus.PASS: 0,
            CheckStatus.WARN: 0,
            CheckStatus.FAIL: 0,
            CheckStatus.SKIP: 0
        }
        
        for result in self.results:
            status_counts[result.status] += 1
            
            # Status emoji
            emoji = {
                CheckStatus.PASS: "✅",
                CheckStatus.WARN: "⚠️",
                CheckStatus.FAIL: "❌",
                CheckStatus.SKIP: "⏭️"
            }
            
            print(f"\n{emoji[result.status]} {result.name}")
            print(f"   {result.message}")
            
            if result.recommendation:
                print(f"   💡 {result.recommendation}")
            
            if result.details and any(key in result.details for key in ['checks', 'quota_details']):
                if 'checks' in result.details:
                    for check_name, check_status, check_message in result.details['checks']:
                        check_emoji = emoji.get(check_status, "❓")
                        print(f"      {check_emoji} {check_name}: {check_message}")
        
        # Summary
        print("\n" + "=" * 60)
        print("📊 SUMMARY")
        print("=" * 60)
        print(f"✅ Passed: {status_counts[CheckStatus.PASS]}")
        print(f"⚠️  Warnings: {status_counts[CheckStatus.WARN]}")
        print(f"❌ Failed: {status_counts[CheckStatus.FAIL]}")
        print(f"⏭️  Skipped: {status_counts[CheckStatus.SKIP]}")
        
        # Overall assessment
        total_critical = status_counts[CheckStatus.PASS] + status_counts[CheckStatus.WARN] + status_counts[CheckStatus.FAIL]
        
        if status_counts[CheckStatus.FAIL] == 0 and status_counts[CheckStatus.WARN] <= total_critical * 0.2:
            print("\n🎉 Your AWS environment is ready for research workloads!")
        elif status_counts[CheckStatus.FAIL] <= total_critical * 0.1:
            print("\n✨ Your AWS environment is mostly ready. Address warnings for optimal performance.")
        else:
            print("\n🔧 Your AWS environment needs attention before running research workloads.")
        
        # Next steps
        print("\n" + "=" * 60)
        print("🚀 NEXT STEPS")
        print("=" * 60)
        
        if status_counts[CheckStatus.FAIL] > 0:
            print("1. Address failed checks first (❌)")
            print("2. Fix authentication and permission issues")
            print("3. Request quota increases if needed")
        
        if status_counts[CheckStatus.WARN] > 0:
            print("4. Review warnings for optimization opportunities")
            print("5. Consider additional security configurations")
        
        print("6. Start with the AWS Research Wizard GUI: streamlit run gui_research_wizard.py")
        print("7. Review cost optimization recommendations")
        print("8. Set up billing alerts and budgets")
    
    def export_results(self, filename: str = "aws_environment_check.json"):
        """Export results to JSON file"""
        export_data = {
            "timestamp": time.strftime("%Y-%m-%d %H:%M:%S UTC", time.gmtime()),
            "account_id": self.account_id,
            "region": self.region,
            "profile": self.profile,
            "results": [asdict(result) for result in self.results]
        }
        
        with open(filename, 'w') as f:
            json.dump(export_data, f, indent=2, default=str)
        
        print(f"\n📄 Results exported to {filename}")

def main():
    parser = argparse.ArgumentParser(
        description="AWS Research Wizard Environment Checker",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python aws_environment_checker.py                    # Use default profile and region
  python aws_environment_checker.py --profile research # Use specific profile
  python aws_environment_checker.py --region us-west-2 # Use specific region
  python aws_environment_checker.py --export results.json # Export results to file
        """
    )
    
    parser.add_argument('--profile', '-p', help='AWS profile name')
    parser.add_argument('--region', '-r', help='AWS region (default: us-east-1)')
    parser.add_argument('--export', '-e', help='Export results to JSON file')
    parser.add_argument('--quiet', '-q', action='store_true', help='Minimal output')
    
    args = parser.parse_args()
    
    # Create and run checker
    checker = AWSEnvironmentChecker(profile=args.profile, region=args.region)
    results = checker.run_all_checks()
    
    if not args.quiet:
        checker.print_results()
    
    if args.export:
        checker.export_results(args.export)
    
    # Exit code based on results
    failed_count = sum(1 for r in results if r.status == CheckStatus.FAIL)
    if failed_count > 0:
        sys.exit(1)
    else:
        sys.exit(0)

if __name__ == "__main__":
    main()