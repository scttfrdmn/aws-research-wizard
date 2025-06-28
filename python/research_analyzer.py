#!/usr/bin/env python3
"""
Research Guide Document Analyzer
Scans researcher documents and extracts AWS-deployable solutions from sections 1, 2, and 7
"""

import os
import re
import json
import glob
from typing import Dict, List, Tuple, Optional
from pathlib import Path
from dataclasses import dataclass, asdict
from collections import defaultdict, Counter
import argparse
import logging

try:
    import anthropic
except ImportError:
    print("Please install the anthropic package: pip install anthropic")
    exit(1)

@dataclass
class ResearcherProfile:
    name: str
    score: float
    email: str
    file_path: str
    research_summary: str
    computational_needs: str
    tools_analysis: str
    aws_services: List[str]
    research_domain: str

@dataclass
class AWSolution:
    service_name: str
    description: str
    use_cases: List[str]
    frequency: int
    research_domains: List[str]

class ResearchDocumentAnalyzer:
    def __init__(self, api_key: str, base_directory: str = "/Users/scttfrdmn/src/award"):
        self.client = anthropic.Anthropic(api_key=api_key)
        self.base_directory = base_directory
        self.researchers = []
        self.aws_solutions = defaultdict(lambda: {"frequency": 0, "use_cases": set(), "domains": set()})
        
        # Setup logging
        logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
        self.logger = logging.getLogger(__name__)

    def find_researcher_documents(self) -> List[str]:
        """Find all researcher markdown files in the directory structure"""
        pattern = os.path.join(self.base_directory, "**", "researchers", "*.md")
        files = glob.glob(pattern, recursive=True)
        self.logger.info(f"Found {len(files)} researcher documents")
        return files

    def extract_sections(self, content: str) -> Tuple[str, str, str]:
        """Extract sections 1, 2, and 7 from the markdown content"""
        
        # Section 1: Research Summary and Institutional Context
        section1_match = re.search(r'## 1\. Research Summary and Institutional Context\n(.*?)(?=## 2\.)', content, re.DOTALL)
        section1 = section1_match.group(1).strip() if section1_match else ""
        
        # Section 2: Computational & Data Needs Assessment
        section2_match = re.search(r'## 2\. Computational & Data Needs Assessment\n(.*?)(?=## 3\.)', content, re.DOTALL)
        section2 = section2_match.group(1).strip() if section2_match else ""
        
        # Section 7: Research Applications and Tools Analysis
        section7_match = re.search(r'## 7\. Research Applications and Tools Analysis\n(.*?)(?=## 8\.)', content, re.DOTALL)
        section7 = section7_match.group(1).strip() if section7_match else ""
        
        return section1, section2, section7

    def extract_aws_services(self, content: str) -> List[str]:
        """Extract AWS service names from the document"""
        aws_services = set()
        
        # Specific AWS services to look for
        known_services = [
            'AWS Batch', 'Amazon S3', 'Amazon SageMaker', 'AWS ParallelCluster',
            'Amazon WorkSpaces', 'AWS Lambda', 'Amazon AppStream', 'Amazon QuickSight',
            'Amazon EC2', 'Amazon ECS', 'Amazon EKS', 'AWS Glue', 'Amazon EMR',
            'AWS Step Functions', 'Amazon CloudWatch', 'AWS IAM', 'Amazon VPC',
            'AWS CloudFormation', 'Amazon RDS', 'Amazon DynamoDB', 'Amazon ElastiCache',
            'AWS CodePipeline', 'AWS CodeBuild', 'Amazon ECR', 'AWS Fargate',
            'Amazon Redshift', 'AWS Data Lake', 'Amazon Kinesis', 'AWS Athena',
            'Amazon Comprehend', 'Amazon Textract', 'Amazon Rekognition',
            'AWS Ground Truth', 'Amazon Forecast', 'Amazon Personalize',
            'AWS DeepRacer', 'Amazon Bedrock', 'AWS HealthLake'
        ]
        
        # Look for exact matches (case insensitive)
        content_lower = content.lower()
        for service in known_services:
            if service.lower() in content_lower:
                aws_services.add(service)
        
        # Also look for generic patterns but filter more carefully
        generic_patterns = [
            r'AWS\s+(Batch|Lambda|Glue|IAM|CodePipeline|CodeBuild|Fargate|Athena|Ground\s+Truth|DeepRacer|HealthLake)',
            r'Amazon\s+(S3|EC2|ECS|EKS|EMR|RDS|DynamoDB|ElastiCache|ECR|Redshift|Kinesis|SageMaker|WorkSpaces|AppStream|QuickSight|CloudWatch|VPC|CloudFormation|Comprehend|Textract|Rekognition|Forecast|Personalize|Bedrock)'
        ]
        
        for pattern in generic_patterns:
            matches = re.findall(pattern, content, re.IGNORECASE)
            for match in matches:
                if isinstance(match, tuple):
                    match = ' '.join(match)
                service_name = f"AWS {match}" if not match.startswith(('Amazon', 'AWS')) else match
                aws_services.add(service_name)
        
        return list(aws_services)

    def parse_researcher_file(self, file_path: str) -> Optional[ResearcherProfile]:
        """Parse a single researcher file and extract relevant information"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Extract basic info from header
            name_match = re.search(r'# (.+?) \(Score: ([\d.]+)\)', content)
            if not name_match:
                self.logger.warning(f"Could not extract name/score from {file_path}")
                return None
            
            name = name_match.group(1)
            score = float(name_match.group(2))
            
            # Extract email
            email_match = re.search(r'\*\*Email:\*\* (.+)', content)
            email = email_match.group(1) if email_match else ""
            
            # Extract sections
            section1, section2, section7 = self.extract_sections(content)
            
            # Extract AWS services
            aws_services = self.extract_aws_services(content)
            
            # Determine research domain from file path
            path_parts = file_path.split('/')
            research_domain = ""
            for part in path_parts:
                if 'university' in part.lower() or 'college' in part.lower():
                    research_domain = part.replace('_', ' ').title()
                    break
            
            return ResearcherProfile(
                name=name,
                score=score,
                email=email,
                file_path=file_path,
                research_summary=section1,
                computational_needs=section2,
                tools_analysis=section7,
                aws_services=aws_services,
                research_domain=research_domain
            )
            
        except Exception as e:
            self.logger.error(f"Error parsing {file_path}: {e}")
            return None

    def analyze_with_claude(self, text: str, analysis_type: str) -> str:
        """Use Claude API to analyze text and extract AWS solutions"""
        
        prompts = {
            "computational_needs": """
            Analyze this computational needs assessment and identify:
            1. Key computational bottlenecks
            2. Data management challenges
            3. Collaboration requirements
            4. AWS services that could address these needs
            
            Focus on extracting specific, deployable AWS solutions. Return as structured text.
            """,
            
            "tools_analysis": """
            Analyze this research tools section and identify:
            1. Software tools mentioned
            2. Infrastructure requirements
            3. Performance bottlenecks
            4. AWS services that could enhance or replace these tools
            
            Focus on specific AWS services and their applications. Return as structured text.
            """,
            
            "solution_synthesis": """
            Based on this collection of researcher needs, synthesize:
            1. Top 10 most common AWS solutions for research computing
            2. Deployment patterns for each solution
            3. Cost optimization strategies
            4. Implementation priorities
            
            Focus on practical, deployable solutions. Return as structured JSON.
            """
        }
        
        try:
            response = self.client.messages.create(
                model="claude-3-5-sonnet-20241022",
                max_tokens=2000,
                messages=[
                    {
                        "role": "user", 
                        "content": f"{prompts.get(analysis_type, prompts['computational_needs'])}\n\nText to analyze:\n{text[:4000]}"
                    }
                ]
            )
            return response.content[0].text
        except Exception as e:
            self.logger.error(f"Claude API error: {e}")
            return f"Error analyzing with Claude: {e}"

    def aggregate_solutions(self) -> Dict[str, AWSolution]:
        """Aggregate and rank AWS solutions across all researchers"""
        
        # Count service frequencies and collect use cases
        service_counter = Counter()
        service_domains = defaultdict(set)
        service_use_cases = defaultdict(set)
        
        for researcher in self.researchers:
            domain = researcher.research_domain
            
            # Extract services from AWS services list
            for service in researcher.aws_services:
                service_counter[service] += 1
                service_domains[service].add(domain)
            
            # Analyze computational needs for additional services
            if researcher.computational_needs:
                analysis = self.analyze_with_claude(researcher.computational_needs, "computational_needs")
                # Extract services mentioned in analysis
                for service in self.extract_aws_services(analysis):
                    service_counter[service] += 1
                    service_domains[service].add(domain)
                    service_use_cases[service].add("Computational needs")
            
            # Analyze tools for additional services
            if researcher.tools_analysis:
                analysis = self.analyze_with_claude(researcher.tools_analysis, "tools_analysis")
                for service in self.extract_aws_services(analysis):
                    service_counter[service] += 1
                    service_domains[service].add(domain)
                    service_use_cases[service].add("Tools enhancement")
        
        # Create AWSolution objects
        solutions = {}
        for service, frequency in service_counter.most_common():
            solutions[service] = AWSolution(
                service_name=service,
                description=f"AWS service used across {frequency} researchers",
                use_cases=list(service_use_cases[service]),
                frequency=frequency,
                research_domains=list(service_domains[service])
            )
        
        return solutions

    def generate_report(self, solutions: Dict[str, AWSolution]) -> str:
        """Generate a comprehensive report of AWS solutions"""
        
        report = []
        report.append("# AWS Research Computing Solutions Analysis")
        report.append(f"## Summary: {len(self.researchers)} researchers analyzed")
        report.append("")
        
        # Top solutions
        report.append("## Top 10 AWS Solutions for Research Computing")
        report.append("")
        
        top_solutions = sorted(solutions.values(), key=lambda x: x.frequency, reverse=True)[:10]
        
        for i, solution in enumerate(top_solutions, 1):
            report.append(f"### {i}. {solution.service_name}")
            report.append(f"**Frequency:** {solution.frequency} researchers")
            report.append(f"**Domains:** {', '.join(solution.research_domains[:5])}")
            report.append(f"**Use Cases:** {', '.join(solution.use_cases)}")
            report.append("")
        
        # Research domains analysis
        domains = defaultdict(list)
        for researcher in self.researchers:
            domains[researcher.research_domain].append(researcher)
        
        report.append("## Research Domains Analysis")
        report.append("")
        
        for domain, researchers in sorted(domains.items(), key=lambda x: len(x[1]), reverse=True)[:10]:
            report.append(f"### {domain}")
            report.append(f"**Researchers:** {len(researchers)}")
            
            # Most common services in this domain
            domain_services = Counter()
            for r in researchers:
                domain_services.update(r.aws_services)
            
            if domain_services:
                top_services = domain_services.most_common(3)
                report.append(f"**Top Services:** {', '.join([s[0] for s in top_services])}")
            report.append("")
        
        # Generate deployment recommendations using Claude
        all_needs = "\n".join([r.computational_needs for r in self.researchers[:50] if r.computational_needs])
        if all_needs:
            synthesis = self.analyze_with_claude(all_needs[:8000], "solution_synthesis")
            report.append("## Deployment Recommendations (AI Analysis)")
            report.append("")
            report.append(synthesis)
        
        return "\n".join(report)

    def run_analysis(self, max_files: Optional[int] = None) -> str:
        """Run the complete analysis pipeline"""
        
        self.logger.info("Starting research document analysis...")
        
        # Find all researcher files
        files = self.find_researcher_documents()
        
        if max_files:
            files = files[:max_files]
            self.logger.info(f"Limiting analysis to {max_files} files")
        
        # Parse each file
        for i, file_path in enumerate(files):
            if i % 50 == 0:
                self.logger.info(f"Processed {i}/{len(files)} files")
            
            researcher = self.parse_researcher_file(file_path)
            if researcher:
                self.researchers.append(researcher)
        
        self.logger.info(f"Successfully parsed {len(self.researchers)} researcher profiles")
        
        # Aggregate solutions
        solutions = self.aggregate_solutions()
        self.logger.info(f"Identified {len(solutions)} unique AWS solutions")
        
        # Generate report
        report = self.generate_report(solutions)
        
        return report

def main():
    parser = argparse.ArgumentParser(description='Analyze research documents for AWS solutions')
    parser.add_argument('--api-key', required=True, help='Claude API key')
    parser.add_argument('--base-dir', default='/Users/scttfrdmn/src/award', help='Base directory to search')
    parser.add_argument('--max-files', type=int, help='Maximum number of files to analyze')
    parser.add_argument('--output', default='aws_research_solutions.md', help='Output report file')
    
    args = parser.parse_args()
    
    # Initialize analyzer
    analyzer = ResearchDocumentAnalyzer(args.api_key, args.base_dir)
    
    # Run analysis
    report = analyzer.run_analysis(args.max_files)
    
    # Save report
    with open(args.output, 'w', encoding='utf-8') as f:
        f.write(report)
    
    print(f"Analysis complete! Report saved to {args.output}")
    print(f"Analyzed {len(analyzer.researchers)} researchers")
    print(f"Report length: {len(report)} characters")

if __name__ == "__main__":
    main()