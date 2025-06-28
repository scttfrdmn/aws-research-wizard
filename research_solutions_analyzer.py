#!/usr/bin/env python3
"""
Research Solutions Analyzer
Identifies domain-specific deployable solutions that can serve multiple researchers
"""

import os
import re
import json
import glob
from typing import Dict, List, Tuple, Optional, Set
from pathlib import Path
from dataclasses import dataclass, asdict
from collections import defaultdict, Counter
import argparse
import logging

@dataclass
class ResearchSolution:
    name: str
    description: str
    applicable_domains: List[str]
    applicable_researchers: List[str]
    computational_patterns: List[str]
    aws_architecture: Dict[str, str]
    implementation_time: str
    estimated_cost: str
    roi_metrics: List[str]
    use_cases: List[str]

@dataclass
class ResearcherProfile:
    name: str
    score: float
    domain: str
    computational_needs: str
    research_summary: str
    tools_analysis: str
    file_path: str

class ResearchSolutionsAnalyzer:
    def __init__(self, base_directory: str = "/Users/scttfrdmn/src/award"):
        self.base_directory = base_directory
        self.researchers = []
        self.solutions = {}
        
        # Define solution patterns to look for
        self.solution_patterns = {
            "burst_computing": {
                "keywords": ["parallel", "simulation", "modeling", "ensemble", "parameter sweep", "monte carlo", "batch processing"],
                "indicators": ["queue", "wait time", "scaling", "peak demand", "computational bursts"]
            },
            "data_pipeline": {
                "keywords": ["pipeline", "workflow", "etl", "data processing", "ingestion", "transformation"],
                "indicators": ["automated", "orchestration", "batch", "streaming", "data flow"]
            },
            "ml_analytics": {
                "keywords": ["machine learning", "neural network", "deep learning", "ai", "analytics", "prediction"],
                "indicators": ["training", "model", "gpu", "tensorflow", "pytorch", "jupyter"]
            },
            "collaborative_platform": {
                "keywords": ["collaboration", "sharing", "multi-user", "team", "consortium"],
                "indicators": ["access control", "shared storage", "remote access", "distributed team"]
            },
            "hpc_modernization": {
                "keywords": ["high performance", "cluster", "mpi", "parallel computing", "supercomputing"],
                "indicators": ["legacy", "modernize", "migration", "hybrid cloud", "burst"]
            },
            "data_lake": {
                "keywords": ["data lake", "storage", "archive", "repository", "dataset management"],
                "indicators": ["petabyte", "lifecycle", "governance", "metadata", "discovery"]
            },
            "secure_research": {
                "keywords": ["sensitive data", "compliance", "security", "privacy", "protected"],
                "indicators": ["hipaa", "ferpa", "export control", "classified", "encryption"]
            },
            "visualization_platform": {
                "keywords": ["visualization", "graphics", "rendering", "vr", "interactive"],
                "indicators": ["gpu", "workstation", "remote", "streaming", "desktop"]
            }
        }
        
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
        section1_match = re.search(r'## 1\. Research Summary and Institutional Context\n(.*?)(?=## 2\.)', content, re.DOTALL)
        section1 = section1_match.group(1).strip() if section1_match else ""
        
        section2_match = re.search(r'## 2\. Computational & Data Needs Assessment\n(.*?)(?=## 3\.)', content, re.DOTALL)
        section2 = section2_match.group(1).strip() if section2_match else ""
        
        section7_match = re.search(r'## 7\. Research Applications and Tools Analysis\n(.*?)(?=## 8\.)', content, re.DOTALL)
        section7 = section7_match.group(1).strip() if section7_match else ""
        
        return section1, section2, section7

    def identify_solution_patterns(self, text: str) -> Dict[str, int]:
        """Identify which solution patterns match the given text"""
        text_lower = text.lower()
        pattern_scores = {}
        
        for pattern_name, pattern_data in self.solution_patterns.items():
            score = 0
            
            # Score based on keywords
            for keyword in pattern_data["keywords"]:
                if keyword in text_lower:
                    score += 2
            
            # Score based on indicators
            for indicator in pattern_data["indicators"]:
                if indicator in text_lower:
                    score += 1
            
            pattern_scores[pattern_name] = score
        
        return pattern_scores

    def extract_research_domain(self, file_path: str, research_summary: str) -> str:
        """Extract research domain from file path and content"""
        # First try to get domain from file path
        path_parts = file_path.split('/')
        institution = ""
        for part in path_parts:
            if 'university' in part.lower() or 'college' in part.lower():
                institution = part.replace('_v0_9_3', '').replace('_', ' ').title()
                break
        
        # Try to identify research field from content
        research_fields = {
            "Computer Science": ["computer science", "artificial intelligence", "machine learning", "cybersecurity", "software"],
            "Climate Science": ["climate", "atmospheric", "weather", "meteorology", "oceanography"],
            "Biomedical": ["biomedical", "biology", "genetics", "medical", "health", "bioinformatics"],
            "Physics": ["physics", "particle", "quantum", "astronomy", "astrophysics"],
            "Chemistry": ["chemistry", "chemical", "molecular", "materials science"],
            "Engineering": ["engineering", "mechanical", "electrical", "civil", "aerospace"],
            "Mathematics": ["mathematics", "statistics", "computational", "numerical"],
            "Earth Sciences": ["geology", "geophysics", "environmental", "earth science"],
            "Social Sciences": ["social", "economics", "psychology", "linguistics", "education"]
        }
        
        research_summary_lower = research_summary.lower()
        detected_fields = []
        
        for field, keywords in research_fields.items():
            for keyword in keywords:
                if keyword in research_summary_lower:
                    detected_fields.append(field)
                    break
        
        if detected_fields:
            return f"{institution} - {', '.join(detected_fields[:2])}"
        
        return institution

    def parse_researcher_file(self, file_path: str) -> Optional[ResearcherProfile]:
        """Parse a single researcher file and extract relevant information"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            # Extract basic info from header
            name_match = re.search(r'# (.+?) \(Score: ([\d.]+)\)', content)
            if not name_match:
                return None
            
            name = name_match.group(1)
            score = float(name_match.group(2))
            
            # Extract sections
            section1, section2, section7 = self.extract_sections(content)
            
            # Extract research domain
            domain = self.extract_research_domain(file_path, section1)
            
            return ResearcherProfile(
                name=name,
                score=score,
                domain=domain,
                computational_needs=section2,
                research_summary=section1,
                tools_analysis=section7,
                file_path=file_path
            )
            
        except Exception as e:
            self.logger.error(f"Error parsing {file_path}: {e}")
            return None

    def create_solution_definitions(self) -> Dict[str, ResearchSolution]:
        """Create deployable solution definitions based on patterns"""
        solutions = {
            "elastic_hpc_platform": ResearchSolution(
                name="Elastic High-Performance Computing Platform",
                description="Auto-scaling HPC environment that eliminates queue wait times and provides burst computing capabilities",
                applicable_domains=[],
                applicable_researchers=[],
                computational_patterns=["burst_computing", "hpc_modernization"],
                aws_architecture={
                    "compute": "AWS ParallelCluster + EC2 HPC instances + Spot instances",
                    "storage": "FSx for Lustre + S3 for data staging",
                    "networking": "Elastic Fabric Adapter for HPC workloads",
                    "management": "AWS Batch for job scheduling"
                },
                implementation_time="3-4 weeks",
                estimated_cost="$2,000-15,000/month (scales with usage)",
                roi_metrics=["70-90% reduction in queue wait times", "30-60% cost savings vs on-premises"],
                use_cases=[]
            ),
            
            "research_data_platform": ResearchSolution(
                name="Scalable Research Data Lake and Analytics Platform",
                description="Unified data platform for ingestion, processing, and analysis of research datasets",
                applicable_domains=[],
                applicable_researchers=[],
                computational_patterns=["data_pipeline", "data_lake"],
                aws_architecture={
                    "storage": "S3 with intelligent tiering + Glacier for archival",
                    "processing": "AWS Glue for ETL + EMR for big data processing",
                    "analytics": "Athena for queries + QuickSight for visualization",
                    "governance": "Lake Formation for data governance"
                },
                implementation_time="4-6 weeks",
                estimated_cost="$1,000-8,000/month (based on data volume)",
                roi_metrics=["50-80% faster data processing", "90% reduction in data management overhead"],
                use_cases=[]
            ),
            
            "ml_research_accelerator": ResearchSolution(
                name="AI/ML Research Acceleration Platform",
                description="Comprehensive platform for machine learning research, training, and deployment",
                applicable_domains=[],
                applicable_researchers=[],
                computational_patterns=["ml_analytics"],
                aws_architecture={
                    "compute": "SageMaker + EC2 GPU instances (P4/P5)",
                    "storage": "S3 for datasets + EFS for shared notebooks",
                    "orchestration": "Step Functions for ML pipelines",
                    "deployment": "SageMaker endpoints + Lambda for inference"
                },
                implementation_time="2-3 weeks",
                estimated_cost="$1,500-10,000/month (varies with GPU usage)",
                roi_metrics=["10x faster model training", "60% reduction in time-to-publication"],
                use_cases=[]
            ),
            
            "secure_collaboration_hub": ResearchSolution(
                name="Secure Multi-Institution Research Collaboration Hub",
                description="Platform enabling secure data sharing and collaboration across institutional boundaries",
                applicable_domains=[],
                applicable_researchers=[],
                computational_patterns=["collaborative_platform", "secure_research"],
                aws_architecture={
                    "security": "Organizations + IAM Identity Center + KMS",
                    "networking": "Transit Gateway + Direct Connect",
                    "storage": "S3 with cross-account access + DataSync",
                    "compute": "WorkSpaces + AppStream for remote access"
                },
                implementation_time="6-8 weeks",
                estimated_cost="$500-3,000/month per institution",
                roi_metrics=["10x improvement in data sharing speed", "100% compliance with security requirements"],
                use_cases=[]
            ),
            
            "hybrid_burst_architecture": ResearchSolution(
                name="Hybrid Cloud Burst Computing Architecture",
                description="Seamless integration with existing HPC infrastructure for overflow computing",
                applicable_domains=[],
                applicable_researchers=[],
                computational_patterns=["burst_computing", "hpc_modernization"],
                aws_architecture={
                    "connectivity": "Direct Connect + VPN for hybrid connectivity",
                    "compute": "ParallelCluster + Batch for burst workloads",
                    "orchestration": "Custom job schedulers + CloudWatch monitoring",
                    "storage": "Storage Gateway + S3 for data synchronization"
                },
                implementation_time="8-12 weeks",
                estimated_cost="$3,000-20,000/month (usage-based)",
                roi_metrics=["50% increase in computational throughput", "Eliminates resource contention"],
                use_cases=[]
            ),
            
            "interactive_research_environment": ResearchSolution(
                name="Cloud-Based Interactive Research Environment",
                description="Virtual research environment with specialized software and collaborative tools",
                applicable_domains=[],
                applicable_researchers=[],
                computational_patterns=["visualization_platform", "collaborative_platform"],
                aws_architecture={
                    "desktop": "WorkSpaces + AppStream 2.0 for applications",
                    "compute": "EC2 with GPU support for visualization",
                    "storage": "EFS for shared project storage",
                    "software": "Custom AMIs with research software stack"
                },
                implementation_time="3-4 weeks",
                estimated_cost="$200-1,500/month per user",
                roi_metrics=["Anywhere access to research tools", "50% reduction in software licensing costs"],
                use_cases=[]
            )
        }
        
        return solutions

    def match_researchers_to_solutions(self):
        """Match researchers to applicable solutions based on their computational patterns"""
        for researcher in self.researchers:
            # Analyze all text content
            all_text = f"{researcher.computational_needs} {researcher.research_summary} {researcher.tools_analysis}"
            
            # Get pattern scores
            pattern_scores = self.identify_solution_patterns(all_text)
            
            # Match to solutions with threshold
            threshold = 3  # Minimum score to consider a match
            
            for solution_name, solution in self.solutions.items():
                match_score = 0
                for pattern in solution.computational_patterns:
                    if pattern in pattern_scores:
                        match_score += pattern_scores[pattern]
                
                if match_score >= threshold:
                    solution.applicable_researchers.append(researcher.name)
                    if researcher.domain not in solution.applicable_domains:
                        solution.applicable_domains.append(researcher.domain)
                    
                    # Extract specific use cases from researcher's needs
                    use_case = self.extract_use_case(researcher, all_text)
                    if use_case and use_case not in solution.use_cases:
                        solution.use_cases.append(use_case)

    def extract_use_case(self, researcher: ResearcherProfile, text: str) -> str:
        """Extract specific use case description for this researcher"""
        # Look for key phrases that describe computational needs
        use_case_patterns = [
            r"(.*processing.*datasets.*)",
            r"(.*simulations?.*require.*)",
            r"(.*analysis.*large.*)",
            r"(.*computational.*challenges.*)",
            r"(.*modeling.*complex.*)"
        ]
        
        for pattern in use_case_patterns:
            matches = re.findall(pattern, text, re.IGNORECASE)
            if matches:
                # Clean up and return first meaningful match
                use_case = matches[0].strip()
                if len(use_case) > 20 and len(use_case) < 200:
                    return f"{researcher.domain.split(' - ')[1] if ' - ' in researcher.domain else researcher.domain}: {use_case}"
        
        return None

    def generate_solutions_report(self) -> str:
        """Generate comprehensive solutions-focused report"""
        report = []
        report.append("# Research Computing Solution Architecture Report")
        report.append(f"## Analysis Summary: {len(self.researchers)} researchers analyzed")
        report.append("")
        
        # Sort solutions by number of applicable researchers
        sorted_solutions = sorted(
            self.solutions.items(), 
            key=lambda x: len(x[1].applicable_researchers), 
            reverse=True
        )
        
        report.append("## Deployable Solutions Ranked by Applicability")
        report.append("")
        
        for i, (solution_id, solution) in enumerate(sorted_solutions, 1):
            if len(solution.applicable_researchers) == 0:
                continue
                
            report.append(f"### {i}. {solution.name}")
            report.append(f"**Applicable to {len(solution.applicable_researchers)} researchers** across {len(solution.applicable_domains)} domains")
            report.append("")
            report.append(f"**Description:** {solution.description}")
            report.append("")
            
            # Architecture
            report.append("**AWS Architecture:**")
            for component, description in solution.aws_architecture.items():
                report.append(f"- **{component.title()}:** {description}")
            report.append("")
            
            # Implementation details
            report.append(f"**Implementation Time:** {solution.implementation_time}")
            report.append(f"**Estimated Cost:** {solution.estimated_cost}")
            report.append("")
            
            # ROI metrics
            if solution.roi_metrics:
                report.append("**Expected ROI:**")
                for metric in solution.roi_metrics:
                    report.append(f"- {metric}")
                report.append("")
            
            # Applicable domains (top 5)
            if solution.applicable_domains:
                top_domains = solution.applicable_domains[:5]
                report.append(f"**Key Domains:** {', '.join(top_domains)}")
                if len(solution.applicable_domains) > 5:
                    report.append(f" (and {len(solution.applicable_domains) - 5} others)")
                report.append("")
            
            # Sample use cases
            if solution.use_cases:
                report.append("**Sample Use Cases:**")
                for use_case in solution.use_cases[:3]:
                    report.append(f"- {use_case}")
                report.append("")
            
            report.append("---")
            report.append("")
        
        # Domain analysis
        report.append("## Cross-Domain Solution Mapping")
        report.append("")
        
        domain_solution_map = defaultdict(list)
        for solution_id, solution in self.solutions.items():
            for domain in solution.applicable_domains:
                if len(solution.applicable_researchers) > 0:
                    domain_solution_map[domain].append((solution.name, len(solution.applicable_researchers)))
        
        for domain, solutions in sorted(domain_solution_map.items(), key=lambda x: len(x[1]), reverse=True)[:10]:
            report.append(f"### {domain}")
            solutions.sort(key=lambda x: x[1], reverse=True)
            report.append("**Recommended Solutions:**")
            for solution_name, researcher_count in solutions[:3]:
                report.append(f"- {solution_name} ({researcher_count} researchers)")
            report.append("")
        
        return "\n".join(report)

    def run_analysis(self, max_files: Optional[int] = None) -> str:
        """Run the complete solutions analysis"""
        self.logger.info("Starting research solutions analysis...")
        
        # Find and parse researcher files
        files = self.find_researcher_documents()
        if max_files:
            files = files[:max_files]
        
        for i, file_path in enumerate(files):
            if i % 100 == 0:
                self.logger.info(f"Processed {i}/{len(files)} files")
            
            researcher = self.parse_researcher_file(file_path)
            if researcher:
                self.researchers.append(researcher)
        
        self.logger.info(f"Parsed {len(self.researchers)} researcher profiles")
        
        # Create solution definitions
        self.solutions = self.create_solution_definitions()
        
        # Match researchers to solutions
        self.match_researchers_to_solutions()
        
        # Generate report
        report = self.generate_solutions_report()
        
        return report

def main():
    parser = argparse.ArgumentParser(description='Analyze research documents for deployable solutions')
    parser.add_argument('--base-dir', default='/Users/scttfrdmn/src/award', help='Base directory to search')
    parser.add_argument('--max-files', type=int, help='Maximum number of files to analyze')
    parser.add_argument('--output', default='research_solutions_analysis.md', help='Output report file')
    
    args = parser.parse_args()
    
    # Initialize analyzer
    analyzer = ResearchSolutionsAnalyzer(args.base_dir)
    
    # Run analysis
    report = analyzer.run_analysis(args.max_files)
    
    # Save report
    with open(args.output, 'w', encoding='utf-8') as f:
        f.write(report)
    
    print(f"Solutions analysis complete! Report saved to {args.output}")
    print(f"Analyzed {len(analyzer.researchers)} researchers")

if __name__ == "__main__":
    main()