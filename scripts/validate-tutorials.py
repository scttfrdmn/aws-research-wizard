#!/usr/bin/env python3
"""
Tutorial Validation Script
Validates structure, syntax, and consistency of all domain guide tutorials
"""

import os
import re
import glob
import subprocess
import json
from pathlib import Path

class TutorialValidator:
    def __init__(self):
        self.base_path = "/Users/scttfrdmn/src/aws-research-wizard/docs/domain-guides"
        self.issues = []
        self.validated_guides = 0

    def log_issue(self, domain, severity, issue):
        """Log a validation issue"""
        self.issues.append({
            'domain': domain,
            'severity': severity,
            'issue': issue
        })
        print(f"[{severity.upper()}] {domain}: {issue}")

    def validate_tutorial_structure(self, guide_path):
        """Validate basic tutorial structure"""
        domain = os.path.basename(os.path.dirname(guide_path))

        try:
            with open(guide_path, 'r', encoding='utf-8') as f:
                content = f.read()
        except Exception as e:
            self.log_issue(domain, 'error', f"Could not read file: {e}")
            return False

        # Check required sections
        required_sections = [
            r'# .+ Research Environment - Getting Started',
            r'## What You\'ll Build',
            r'## Before You Start',
            r'## Step 1: Install AWS Research Wizard',
            r'## Step 2: Set Up AWS Account',
            r'## Step 3: Configure Your Credentials',
            r'## Step 4: Validate Your Setup',
            r'## Step 5: Deploy',
            r'## Step 6: Connect',
            r'## Step 7: .+ Tools',
            r'## Step 8: .+ Data',
            r'## Step 9: Using Your Own .+ Data',
        ]

        for section in required_sections:
            if not re.search(section, content):
                self.log_issue(domain, 'error', f"Missing required section: {section}")

        # Check for cost information
        if not re.search(r'\$\d+', content):
            self.log_issue(domain, 'warning', "No cost information found")

        # Check for data download summary
        if not re.search(r'📊 Data Download Summary:', content):
            self.log_issue(domain, 'error', "Missing data download summary")

        # Check for extend and contribute section
        if not re.search(r'### Extend and Contribute', content):
            self.log_issue(domain, 'error', "Missing 'Extend and Contribute' section")

        return True

    def validate_aws_commands(self, guide_path):
        """Validate AWS CLI commands in tutorial"""
        domain = os.path.basename(os.path.dirname(guide_path))

        try:
            with open(guide_path, 'r', encoding='utf-8') as f:
                content = f.read()
        except Exception as e:
            return False

        # Extract code blocks
        code_blocks = re.findall(r'```bash\n(.*?)\n```', content, re.DOTALL)

        for i, block in enumerate(code_blocks):
            lines = block.split('\n')
            for line_num, line in enumerate(lines):
                line = line.strip()
                if not line or line.startswith('#'):
                    continue

                # Check aws-research-wizard commands
                if line.startswith('aws-research-wizard'):
                    self.validate_wizard_command(domain, line, i+1, line_num+1)

                # Check AWS CLI commands
                elif line.startswith('aws '):
                    self.validate_aws_cli_command(domain, line, i+1, line_num+1)

    def validate_wizard_command(self, domain, command, block_num, line_num):
        """Validate aws-research-wizard command syntax"""
        valid_commands = [
            'config setup',
            'deploy validate',
            'deploy start',
            'deploy delete',
            'deploy destroy',
            'connect',
            'monitor costs',
            'emergency-stop'
        ]

        # Extract the subcommand
        parts = command.split()
        if len(parts) < 2:
            self.log_issue(domain, 'error', f"Invalid wizard command (block {block_num}): {command}")
            return

        subcommand = ' '.join(parts[1:3]) if len(parts) > 2 else parts[1]

        if not any(subcommand.startswith(valid) for valid in valid_commands):
            self.log_issue(domain, 'warning', f"Unrecognized wizard command (block {block_num}): {subcommand}")

    def validate_aws_cli_command(self, domain, command, block_num, line_num):
        """Validate AWS CLI command syntax"""
        # Basic syntax check for common AWS CLI patterns
        if ' s3 cp ' in command:
            # Check S3 copy command format
            if not re.search(r's3://[\w\-\.]+/', command):
                self.log_issue(domain, 'warning', f"Potentially invalid S3 URL (block {block_num}): {command}")

        # Check for common AWS CLI options
        if '--no-sign-request' in command and 's3 cp' not in command:
            self.log_issue(domain, 'warning', f"--no-sign-request used outside S3 cp (block {block_num}): {command}")

    def validate_data_downloads(self, guide_path):
        """Validate data download information and URLs"""
        domain = os.path.basename(os.path.dirname(guide_path))

        try:
            with open(guide_path, 'r', encoding='utf-8') as f:
                content = f.read()
        except Exception as e:
            return False

        # Extract data download summary
        download_match = re.search(r'📊 Data Download Summary:\*\*(.*?)\*\*Total download\*\*: ~([0-9.]+) GB', content, re.DOTALL)
        if download_match:
            total_size = float(download_match.group(2))

            # Check if size is in target range (4-7 GB as we standardized)
            if total_size < 4.0 or total_size > 7.0:
                self.log_issue(domain, 'warning', f"Data download size {total_size}GB outside target range 4-7GB")

        # Extract S3 URLs for basic format validation
        s3_urls = re.findall(r's3://[\w\-\.]+/[^\s]*', content)
        for url in s3_urls:
            if not re.match(r's3://[\w\-\.]+/.*', url):
                self.log_issue(domain, 'error', f"Invalid S3 URL format: {url}")

    def validate_step_numbering(self, guide_path):
        """Validate step numbering consistency"""
        domain = os.path.basename(os.path.dirname(guide_path))

        try:
            with open(guide_path, 'r', encoding='utf-8') as f:
                content = f.read()
        except Exception as e:
            return False

        # Extract step numbers
        steps = re.findall(r'## Step (\d+):', content)
        expected_steps = list(range(1, len(steps) + 1))
        actual_steps = [int(step) for step in steps]

        if actual_steps != expected_steps:
            self.log_issue(domain, 'error', f"Step numbering inconsistent. Expected: {expected_steps}, Found: {actual_steps}")

    def validate_code_syntax(self, guide_path):
        """Validate embedded code syntax"""
        domain = os.path.basename(os.path.dirname(guide_path))

        try:
            with open(guide_path, 'r', encoding='utf-8') as f:
                content = f.read()
        except Exception as e:
            return False

        # Extract Python code blocks
        python_blocks = re.findall(r'```python\n(.*?)\n```', content, re.DOTALL)
        for i, block in enumerate(python_blocks):
            try:
                compile(block, f'<tutorial-{domain}-block-{i}>', 'exec')
            except SyntaxError as e:
                self.log_issue(domain, 'error', f"Python syntax error in block {i+1}: {e}")

        # Check for common bash syntax issues
        bash_blocks = re.findall(r'```bash\n(.*?)\n```', content, re.DOTALL)
        for i, block in enumerate(bash_blocks):
            lines = block.split('\n')
            for line_num, line in enumerate(lines):
                line = line.strip()
                if not line or line.startswith('#'):
                    continue

                # Check for unmatched quotes
                if line.count('"') % 2 != 0:
                    self.log_issue(domain, 'warning', f"Unmatched quotes in bash block {i+1}, line {line_num+1}: {line}")

                # Check for dangerous commands
                dangerous_patterns = ['rm -rf /', 'sudo rm', 'chmod 777']
                for pattern in dangerous_patterns:
                    if pattern in line:
                        self.log_issue(domain, 'warning', f"Potentially dangerous command in block {i+1}: {line}")

    def run_validation(self):
        """Run complete validation on all tutorials"""
        print("🔍 Starting tutorial validation...")
        print("=" * 50)

        guide_files = glob.glob(f"{self.base_path}/*/getting-started.md")

        for guide_file in sorted(guide_files):
            domain = os.path.basename(os.path.dirname(guide_file))
            print(f"\n📋 Validating {domain}...")

            self.validate_tutorial_structure(guide_file)
            self.validate_aws_commands(guide_file)
            self.validate_data_downloads(guide_file)
            self.validate_step_numbering(guide_file)
            self.validate_code_syntax(guide_file)

            self.validated_guides += 1

        self.print_summary()

    def print_summary(self):
        """Print validation summary"""
        print("\n" + "=" * 50)
        print("📊 VALIDATION SUMMARY")
        print("=" * 50)

        print(f"Guides validated: {self.validated_guides}")
        print(f"Total issues found: {len(self.issues)}")

        # Group issues by severity
        errors = [i for i in self.issues if i['severity'] == 'error']
        warnings = [i for i in self.issues if i['severity'] == 'warning']

        print(f"Errors: {len(errors)}")
        print(f"Warnings: {len(warnings)}")

        if errors:
            print("\n🚨 CRITICAL ERRORS:")
            for error in errors:
                print(f"  {error['domain']}: {error['issue']}")

        if warnings:
            print("\n⚠️  WARNINGS:")
            for warning in warnings:
                print(f"  {warning['domain']}: {warning['issue']}")

        if not self.issues:
            print("\n✅ All tutorials passed validation!")

        # Export issues to JSON for further analysis
        with open('/Users/scttfrdmn/src/aws-research-wizard/tutorial_validation_results.json', 'w') as f:
            json.dump(self.issues, f, indent=2)

        print(f"\n📁 Detailed results saved to: tutorial_validation_results.json")

if __name__ == "__main__":
    validator = TutorialValidator()
    validator.run_validation()
