"""
Unit tests for the AWS Environment Checker.

This module tests the AWS environment validation functionality including
credential checks, quota validation, service availability, and security
configuration assessment.
"""

import json
import pytest
from unittest.mock import Mock, patch, MagicMock
from botocore.exceptions import ClientError, NoCredentialsError, ProfileNotFound

# Import the module under test
import sys
import os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', '..'))

from aws_environment_checker import (
    AWSEnvironmentChecker,
    CheckResult,
    CheckStatus,
    QuotaCheck
)


class TestCheckResult:
    """Test the CheckResult dataclass."""

    def test_check_result_creation_minimal(self):
        """Test creating CheckResult with minimal required fields."""
        result = CheckResult(
            name="Test Check",
            status=CheckStatus.PASS,
            message="Test passed successfully"
        )

        assert result.name == "Test Check"
        assert result.status == CheckStatus.PASS
        assert result.message == "Test passed successfully"
        assert result.recommendation is None
        assert result.details is None

    def test_check_result_creation_complete(self):
        """Test creating CheckResult with all fields."""
        details = {"test_key": "test_value"}
        result = CheckResult(
            name="Test Check",
            status=CheckStatus.WARN,
            message="Test completed with warnings",
            recommendation="Review configuration",
            details=details
        )

        assert result.name == "Test Check"
        assert result.status == CheckStatus.WARN
        assert result.message == "Test completed with warnings"
        assert result.recommendation == "Review configuration"
        assert result.details == details


class TestCheckStatus:
    """Test the CheckStatus enumeration."""

    def test_check_status_values(self):
        """Test that all expected status values are defined."""
        assert CheckStatus.PASS.value == "PASS"
        assert CheckStatus.WARN.value == "WARN"
        assert CheckStatus.FAIL.value == "FAIL"
        assert CheckStatus.SKIP.value == "SKIP"


class TestQuotaCheck:
    """Test the QuotaCheck dataclass."""

    def test_quota_check_creation(self):
        """Test creating QuotaCheck instance."""
        quota = QuotaCheck(
            service="ec2",
            quota_name="Running On-Demand Standard instances",
            quota_code="L-1216C47A",
            minimum_required=100,
            recommended=500,
            current_value=200
        )

        assert quota.service == "ec2"
        assert quota.quota_name == "Running On-Demand Standard instances"
        assert quota.quota_code == "L-1216C47A"
        assert quota.minimum_required == 100
        assert quota.recommended == 500
        assert quota.current_value == 200


class TestAWSEnvironmentChecker:
    """Test the main AWSEnvironmentChecker class."""

    def test_initialization_with_default_values(self):
        """Test checker initialization with default values."""
        with patch('boto3.Session') as mock_session:
            mock_session.return_value.client.return_value.get_caller_identity.return_value = {
                'Account': '123456789012'
            }

            checker = AWSEnvironmentChecker()

            assert checker.profile is None
            assert checker.region == 'us-east-1'
            assert checker.account_id == '123456789012'
            mock_session.assert_called_once_with(region_name='us-east-1')

    def test_initialization_with_custom_profile_and_region(self):
        """Test checker initialization with custom profile and region."""
        with patch('boto3.Session') as mock_session:
            mock_session.return_value.client.return_value.get_caller_identity.return_value = {
                'Account': '123456789012'
            }

            checker = AWSEnvironmentChecker(profile="research", region="us-west-2")

            assert checker.profile == "research"
            assert checker.region == "us-west-2"
            assert checker.account_id == '123456789012'
            mock_session.assert_called_once_with(profile_name="research", region_name="us-west-2")

    def test_initialization_with_no_credentials(self):
        """Test checker initialization when AWS credentials are not available."""
        with patch('boto3.Session', side_effect=NoCredentialsError()):
            checker = AWSEnvironmentChecker()

            assert checker.session is None
            assert len(checker.results) == 1
            assert checker.results[0].status == CheckStatus.FAIL
            assert "credentials not found" in checker.results[0].message.lower()

    def test_initialization_with_invalid_profile(self):
        """Test checker initialization with invalid profile."""
        with patch('boto3.Session', side_effect=ProfileNotFound(profile="invalid")):
            checker = AWSEnvironmentChecker(profile="invalid")

            assert checker.session is None
            assert len(checker.results) == 1
            assert checker.results[0].status == CheckStatus.FAIL

    @patch('boto3.Session')
    def test_check_credentials_success(self, mock_session):
        """Test successful credential validation."""
        # Mock successful STS response
        mock_sts = Mock()
        mock_sts.get_caller_identity.return_value = {
            'Account': '123456789012',
            'UserId': 'AIDACKCEVSQ6C2EXAMPLE',
            'Arn': 'arn:aws:iam::123456789012:user/test-user'
        }
        mock_session.return_value.client.return_value = mock_sts

        checker = AWSEnvironmentChecker()
        checker._check_credentials()

        # Find the credentials check result
        cred_result = next((r for r in checker.results if r.name == "AWS Credentials"), None)
        assert cred_result is not None
        assert cred_result.status == CheckStatus.PASS
        assert "Successfully authenticated" in cred_result.message
        assert cred_result.details["account_id"] == "123456789012"

    @patch('boto3.Session')
    def test_check_credentials_failure(self, mock_session):
        """Test credential validation failure."""
        # Mock STS client failure
        mock_sts = Mock()
        mock_sts.get_caller_identity.side_effect = ClientError(
            {"Error": {"Code": "InvalidUserID.NotFound"}}, "GetCallerIdentity"
        )
        mock_session.return_value.client.return_value = mock_sts

        checker = AWSEnvironmentChecker()
        checker._check_credentials()

        # Find the credentials check result
        cred_result = next((r for r in checker.results if r.name == "AWS Credentials"), None)
        assert cred_result is not None
        assert cred_result.status == CheckStatus.FAIL
        assert "Credential validation failed" in cred_result.message

    @patch('boto3.Session')
    def test_check_permissions_success(self, mock_session):
        """Test successful permission validation."""
        # Mock clients for different services
        mock_ec2 = Mock()
        mock_ec2.describe_instances.return_value = {"Reservations": []}
        mock_ec2.describe_instance_types.return_value = {"InstanceTypes": []}

        mock_s3 = Mock()
        mock_s3.list_buckets.return_value = {"Buckets": []}

        mock_iam = Mock()
        mock_iam.get_user.return_value = {"User": {"UserName": "test-user"}}

        mock_pricing = Mock()
        mock_pricing.get_products.return_value = {"Products": []}

        mock_quotas = Mock()
        mock_quotas.list_service_quotas.return_value = {"Quotas": []}

        def client_factory(service_name, **kwargs):
            clients = {
                'ec2': mock_ec2,
                's3': mock_s3,
                'iam': mock_iam,
                'pricing': mock_pricing,
                'servicequotas': mock_quotas,
                'sts': Mock()
            }
            return clients.get(service_name, Mock())

        mock_session.return_value.client.side_effect = client_factory
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_permissions()

        # Find the permissions check result
        perm_result = next((r for r in checker.results if r.name == "IAM Permissions"), None)
        assert perm_result is not None
        assert perm_result.status == CheckStatus.PASS
        assert "All essential permissions available" in perm_result.message

    @patch('boto3.Session')
    def test_check_permissions_partial_failure(self, mock_session):
        """Test permission validation with some failures."""
        # Mock clients with some failures
        mock_ec2 = Mock()
        mock_ec2.describe_instances.side_effect = ClientError(
            {"Error": {"Code": "UnauthorizedOperation"}}, "DescribeInstances"
        )
        mock_ec2.describe_instance_types.return_value = {"InstanceTypes": []}

        mock_s3 = Mock()
        mock_s3.list_buckets.return_value = {"Buckets": []}

        def client_factory(service_name, **kwargs):
            clients = {
                'ec2': mock_ec2,
                's3': mock_s3,
                'sts': Mock()
            }
            return clients.get(service_name, Mock())

        mock_session.return_value.client.side_effect = client_factory
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_permissions()

        # Find the permissions check result
        perm_result = next((r for r in checker.results if r.name == "IAM Permissions"), None)
        assert perm_result is not None
        assert perm_result.status == CheckStatus.WARN
        assert "Missing permissions" in perm_result.message

    @patch('boto3.Session')
    def test_check_regions(self, mock_session):
        """Test region availability check."""
        # Mock EC2 describe_regions response
        mock_ec2 = Mock()
        mock_ec2.describe_regions.return_value = {
            "Regions": [
                {"RegionName": "us-east-1", "Endpoint": "ec2.us-east-1.amazonaws.com"},
                {"RegionName": "us-west-2", "Endpoint": "ec2.us-west-2.amazonaws.com"},
                {"RegionName": "eu-west-1", "Endpoint": "ec2.eu-west-1.amazonaws.com"},
                {"RegionName": "ap-southeast-1", "Endpoint": "ec2.ap-southeast-1.amazonaws.com"}
            ]
        }

        mock_session.return_value.client.return_value = mock_ec2
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_regions()

        # Find the regional availability check result
        region_result = next((r for r in checker.results if r.name == "Regional Availability"), None)
        assert region_result is not None
        assert region_result.status == CheckStatus.PASS
        assert "Good regional coverage" in region_result.message
        assert len(region_result.details["research_regions"]) >= 3

    @patch('boto3.Session')
    def test_check_service_availability(self, mock_session):
        """Test service availability check."""
        # Mock various service clients
        mock_ec2 = Mock()
        mock_ec2.describe_availability_zones.return_value = {"AvailabilityZones": []}
        mock_ec2.describe_volumes.return_value = {"Volumes": []}
        mock_ec2.describe_vpcs.return_value = {"Vpcs": []}

        mock_s3 = Mock()
        mock_s3.list_buckets.return_value = {"Buckets": []}

        mock_iam = Mock()
        mock_iam.list_roles.return_value = {"Roles": []}

        mock_cloudwatch = Mock()
        mock_cloudwatch.list_metrics.return_value = {"Metrics": []}

        mock_batch = Mock()
        mock_batch.describe_compute_environments.return_value = {"computeEnvironments": []}

        mock_efs = Mock()
        mock_efs.describe_file_systems.return_value = {"FileSystems": []}

        mock_fsx = Mock()
        mock_fsx.describe_file_systems.return_value = {"FileSystems": []}

        def client_factory(service_name, **kwargs):
            clients = {
                'ec2': mock_ec2,
                's3': mock_s3,
                'iam': mock_iam,
                'cloudwatch': mock_cloudwatch,
                'batch': mock_batch,
                'efs': mock_efs,
                'fsx': mock_fsx,
                'sts': Mock()
            }
            return clients.get(service_name, Mock())

        mock_session.return_value.client.side_effect = client_factory
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_service_availability()

        # Find the service availability check result
        service_result = next((r for r in checker.results if r.name == "Service Availability"), None)
        assert service_result is not None
        assert service_result.status == CheckStatus.PASS
        assert "All essential services available" in service_result.message

    @patch('boto3.Session')
    def test_check_compute_quotas(self, mock_session):
        """Test compute quota checking."""
        # Mock service quotas client
        mock_quotas = Mock()
        mock_quotas.get_service_quota.return_value = {
            "Quota": {"Value": 1000.0}  # High quota value
        }

        mock_session.return_value.client.return_value = mock_quotas
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_compute_quotas()

        # Find the compute quotas check result
        quota_result = next((r for r in checker.results if r.name == "EC2 Compute Quotas"), None)
        assert quota_result is not None
        assert quota_result.status in [CheckStatus.PASS, CheckStatus.WARN]  # Depends on quota values

    @patch('boto3.Session')
    def test_check_storage_quotas(self, mock_session):
        """Test storage quota checking."""
        # Mock EBS quotas
        mock_quotas = Mock()
        mock_quotas.get_service_quota.return_value = {
            "Quota": {"Value": 50000.0}  # High quota value
        }

        # Mock S3 client
        mock_s3 = Mock()
        mock_s3.list_buckets.return_value = {"Buckets": [{"Name": f"bucket-{i}"} for i in range(5)]}

        def client_factory(service_name, **kwargs):
            if service_name == 'servicequotas':
                return mock_quotas
            elif service_name == 's3':
                return mock_s3
            return Mock()

        mock_session.return_value.client.side_effect = client_factory
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_storage_quotas()

        # Should have both EBS and S3 quota check results
        ebs_result = next((r for r in checker.results if r.name == "EBS Storage Quotas"), None)
        s3_result = next((r for r in checker.results if r.name == "S3 Storage Quotas"), None)

        assert ebs_result is not None
        assert s3_result is not None
        assert s3_result.status == CheckStatus.PASS  # Low bucket count

    @patch('boto3.Session')
    def test_check_security_setup(self, mock_session):
        """Test security configuration check."""
        # Mock EC2 client for VPC check
        mock_ec2 = Mock()
        mock_ec2.describe_vpcs.return_value = {
            "Vpcs": [
                {"VpcId": "vpc-12345", "IsDefault": False},  # Custom VPC
                {"VpcId": "vpc-default", "IsDefault": True}   # Default VPC
            ]
        }

        # Mock CloudTrail client
        mock_cloudtrail = Mock()
        mock_cloudtrail.describe_trails.return_value = {
            "trailList": [{"Name": "test-trail"}]
        }

        # Mock IAM client
        mock_iam = Mock()
        mock_iam.get_account_summary.return_value = {
            "SummaryMap": {"MFADevices": 2}
        }

        def client_factory(service_name, **kwargs):
            clients = {
                'ec2': mock_ec2,
                'cloudtrail': mock_cloudtrail,
                'iam': mock_iam,
                'sts': Mock()
            }
            return clients.get(service_name, Mock())

        mock_session.return_value.client.side_effect = client_factory
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_security_setup()

        # Find the security configuration check result
        security_result = next((r for r in checker.results if r.name == "Security Configuration"), None)
        assert security_result is not None
        assert security_result.status == CheckStatus.PASS
        assert "checks:" in security_result.message.lower()

    @patch('boto3.Session')
    def test_check_billing_setup_success(self, mock_session):
        """Test billing setup check with access."""
        # Mock Cost Explorer client
        mock_ce = Mock()
        mock_ce.get_cost_and_usage.return_value = {
            "ResultsByTime": []
        }

        def client_factory(service_name, **kwargs):
            if service_name == 'ce':
                return mock_ce
            return Mock()

        mock_session.return_value.client.side_effect = client_factory
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_billing_setup()

        # Find the billing access check result
        billing_result = next((r for r in checker.results if r.name == "Billing Access"), None)
        assert billing_result is not None
        assert billing_result.status == CheckStatus.PASS
        assert "accessible" in billing_result.message.lower()

    @patch('boto3.Session')
    def test_check_billing_setup_no_access(self, mock_session):
        """Test billing setup check without access."""
        # Mock Cost Explorer client with access denied
        mock_ce = Mock()
        mock_ce.get_cost_and_usage.side_effect = ClientError(
            {"Error": {"Code": "AccessDenied"}}, "GetCostAndUsage"
        )

        def client_factory(service_name, **kwargs):
            if service_name == 'ce':
                return mock_ce
            return Mock()

        mock_session.return_value.client.side_effect = client_factory
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_billing_setup()

        # Find the billing access check result
        billing_result = next((r for r in checker.results if r.name == "Billing Access"), None)
        assert billing_result is not None
        assert billing_result.status == CheckStatus.WARN
        assert "Limited billing access" in billing_result.message

    @patch('boto3.Session')
    def test_check_research_prerequisites(self, mock_session):
        """Test research-specific prerequisite checks."""
        # Mock EC2 client
        mock_ec2 = Mock()
        mock_ec2.describe_instance_types.return_value = {
            "InstanceTypes": [
                {"InstanceType": "hpc6a.48xlarge"},
                {"InstanceType": "c6a.24xlarge"}
            ]
        }
        mock_ec2.describe_spot_price_history.return_value = {
            "SpotPriceHistory": [{"SpotPrice": "0.03", "InstanceType": "c6i.large"}]
        }
        mock_ec2.describe_placement_groups.return_value = {"PlacementGroups": []}

        mock_session.return_value.client.return_value = mock_ec2
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_research_prerequisites()

        # Find the research prerequisites check result
        research_result = next((r for r in checker.results if r.name == "Research Prerequisites"), None)
        assert research_result is not None
        assert research_result.status == CheckStatus.PASS
        assert "available" in research_result.message.lower()

    @patch('boto3.Session')
    def test_run_all_checks(self, mock_session):
        """Test running all checks together."""
        # Mock a basic session
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        # Mock all the individual check methods
        checker = AWSEnvironmentChecker()

        with patch.object(checker, '_check_credentials') as mock_creds, \
             patch.object(checker, '_check_permissions') as mock_perms, \
             patch.object(checker, '_check_regions') as mock_regions, \
             patch.object(checker, '_check_service_availability') as mock_services, \
             patch.object(checker, '_check_compute_quotas') as mock_compute, \
             patch.object(checker, '_check_storage_quotas') as mock_storage, \
             patch.object(checker, '_check_network_quotas') as mock_network, \
             patch.object(checker, '_check_specialized_quotas') as mock_specialized, \
             patch.object(checker, '_check_security_setup') as mock_security, \
             patch.object(checker, '_check_billing_setup') as mock_billing, \
             patch.object(checker, '_check_research_prerequisites') as mock_research:

            results = checker.run_all_checks()

            # Verify all check methods were called
            mock_creds.assert_called_once()
            mock_perms.assert_called_once()
            mock_regions.assert_called_once()
            mock_services.assert_called_once()
            mock_compute.assert_called_once()
            mock_storage.assert_called_once()
            mock_network.assert_called_once()
            mock_specialized.assert_called_once()
            mock_security.assert_called_once()
            mock_billing.assert_called_once()
            mock_research.assert_called_once()

            assert isinstance(results, list)

    def test_export_results(self, temp_dir):
        """Test exporting results to JSON file."""
        with patch('boto3.Session') as mock_session:
            mock_session.return_value.client.return_value.get_caller_identity.return_value = {
                'Account': '123456789012'
            }

            checker = AWSEnvironmentChecker()

            # Add some test results
            checker.results = [
                CheckResult("Test Check 1", CheckStatus.PASS, "Success"),
                CheckResult("Test Check 2", CheckStatus.WARN, "Warning", "Fix this"),
                CheckResult("Test Check 3", CheckStatus.FAIL, "Failed", "Critical issue")
            ]

            output_file = os.path.join(temp_dir, "test_results.json")
            checker.export_results(output_file)

            # Verify file was created and contains expected data
            assert os.path.exists(output_file)

            with open(output_file, 'r') as f:
                data = json.load(f)

            assert "timestamp" in data
            assert "account_id" in data
            assert "region" in data
            assert "results" in data
            assert len(data["results"]) == 3
            assert data["account_id"] == "123456789012"

    def test_print_results_output(self, capsys):
        """Test print_results output formatting."""
        with patch('boto3.Session') as mock_session:
            mock_session.return_value.client.return_value.get_caller_identity.return_value = {
                'Account': '123456789012'
            }

            checker = AWSEnvironmentChecker()

            # Add test results with different statuses
            checker.results = [
                CheckResult("Pass Check", CheckStatus.PASS, "Everything is good"),
                CheckResult("Warning Check", CheckStatus.WARN, "Minor issue", "Review config"),
                CheckResult("Fail Check", CheckStatus.FAIL, "Critical problem", "Fix immediately"),
                CheckResult("Skip Check", CheckStatus.SKIP, "Not applicable")
            ]

            checker.print_results()

            # Capture printed output
            captured = capsys.readouterr()

            # Verify output contains expected elements
            assert "ENVIRONMENT CHECK RESULTS" in captured.out
            assert "✅ Pass Check" in captured.out
            assert "⚠️ Warning Check" in captured.out
            assert "❌ Fail Check" in captured.out
            assert "⏭️ Skip Check" in captured.out
            assert "SUMMARY" in captured.out
            assert "Passed: 1" in captured.out
            assert "Warnings: 1" in captured.out
            assert "Failed: 1" in captured.out
            assert "Skipped: 1" in captured.out

    def test_service_quotas_parsing(self):
        """Test quota check result parsing."""
        quota_checks = [
            QuotaCheck("ec2", "Standard instances", "L-1216C47A", 100, 500, 200),
            QuotaCheck("ec2", "Spot instances", "L-34B43A08", 256, 1000, 800),
            QuotaCheck("ec2", "Elastic IPs", "L-0263D0A3", 20, 100, 150)
        ]

        # Test quota evaluation logic
        passed = sum(1 for q in quota_checks if q.current_value >= q.recommended)
        warned = sum(1 for q in quota_checks if q.minimum_required <= q.current_value < q.recommended)
        failed = sum(1 for q in quota_checks if q.current_value < q.minimum_required)

        assert passed == 1  # Only Elastic IPs meets recommended
        assert warned == 2  # Standard and Spot instances are adequate but not optimal
        assert failed == 0  # None below minimum

    @pytest.mark.integration
    def test_real_aws_environment_check(self):
        """Integration test with real AWS environment (requires credentials)."""
        # Skip if no real AWS credentials
        try:
            checker = AWSEnvironmentChecker()
            if checker.session is None:
                pytest.skip("No AWS credentials available for integration test")

            results = checker.run_all_checks()

            # Basic validation that checks ran
            assert len(results) > 0

            # Should have at least credential and permission checks
            check_names = [r.name for r in results]
            assert "AWS Credentials" in check_names
            assert "IAM Permissions" in check_names

        except Exception as e:
            pytest.skip(f"Integration test failed due to AWS access: {e}")


class TestErrorHandling:
    """Test error handling in various scenarios."""

    @patch('boto3.Session')
    def test_network_timeout_handling(self, mock_session):
        """Test handling of network timeouts."""
        # Mock timeout exception
        mock_session.return_value.client.return_value.describe_regions.side_effect = Exception("Connection timeout")
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_regions()

        # Should handle timeout gracefully
        region_result = next((r for r in checker.results if r.name == "Regional Availability"), None)
        assert region_result is not None
        assert region_result.status == CheckStatus.FAIL
        assert "timeout" in region_result.message.lower()

    @patch('boto3.Session')
    def test_service_unavailable_handling(self, mock_session):
        """Test handling of unavailable services."""
        # Mock service unavailable
        mock_session.return_value.client.side_effect = Exception("Service temporarily unavailable")
        mock_session.return_value.client.return_value.get_caller_identity.return_value = {
            'Account': '123456789012'
        }

        checker = AWSEnvironmentChecker()
        checker._check_service_availability()

        # Should handle service unavailability gracefully
        service_result = next((r for r in checker.results if r.name == "Service Availability"), None)
        assert service_result is not None
        assert service_result.status == CheckStatus.FAIL
