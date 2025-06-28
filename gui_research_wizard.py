#!/usr/bin/env python3
"""
AWS Research Wizard - GUI Interface
Interactive web-based interface for the AWS Research Wizard using Streamlit
"""

import streamlit as st
import json
import pandas as pd
import plotly.express as px
import plotly.graph_objects as go
from typing import Dict, List, Any, Optional
import sys
import os
from pathlib import Path

# Add the current directory to Python path for imports
current_dir = Path(__file__).parent
sys.path.append(str(current_dir))

# Import our research pack modules
try:
    from research_infrastructure_wizard import ResearchInfrastructureWizard, WorkloadCharacteristics, Priority, WorkloadSize
    from geospatial_research_pack import GeospatialResearchPack, GeospatialDomain, GeospatialWorkload
    from atmospheric_chemistry_pack import AtmosphericChemistryPack, AtmosphericDomain, AtmosphericWorkload
    from cybersecurity_research_pack import CyberSecurityResearchPack, CyberSecurityDomain, CyberSecurityWorkload
    from renewable_energy_systems_pack import RenewableEnergyResearchPack, RenewableEnergyDomain, RenewableEnergyWorkload
    from visualization_studio_pack import VisualizationStudioPack
except ImportError as e:
    st.error(f"Error importing research pack modules: {e}")
    st.error("Please ensure all research pack files are in the same directory as this GUI script.")

# Page configuration
st.set_page_config(
    page_title="AWS Research Wizard",
    page_icon="🧙‍♂️",
    layout="wide",
    initial_sidebar_state="expanded"
)

# Custom CSS for better styling
st.markdown("""
<style>
    .main-header {
        font-size: 3rem;
        color: #FF9900;
        text-align: center;
        margin-bottom: 2rem;
    }
    .domain-card {
        background-color: #f0f2f6;
        padding: 1rem;
        border-radius: 0.5rem;
        margin: 0.5rem 0;
    }
    .cost-box {
        background-color: #e6f3ff;
        padding: 1rem;
        border-radius: 0.5rem;
        border-left: 4px solid #0066cc;
    }
    .recommendation-box {
        background-color: #f0f8f0;
        padding: 1rem;
        border-radius: 0.5rem;
        border-left: 4px solid #00cc66;
    }
    .stSelectbox > div > div > select {
        background-color: white;
    }
</style>
""", unsafe_allow_html=True)

def initialize_session_state():
    """Initialize session state variables"""
    if 'current_step' not in st.session_state:
        st.session_state.current_step = 'domain_selection'
    if 'selected_domain' not in st.session_state:
        st.session_state.selected_domain = None
    if 'workload_config' not in st.session_state:
        st.session_state.workload_config = {}
    if 'recommendation' not in st.session_state:
        st.session_state.recommendation = None

def get_domain_info():
    """Get information about all available research domains"""
    return {
        "core_domains": {
            "Genomics & Bioinformatics": {
                "description": "DNA/RNA analysis, variant calling, population genomics",
                "tools": ["GATK", "BWA", "STAR", "SAMtools"],
                "cost_range": "$150-900/month",
                "icon": "🧬"
            },
            "Climate Science": {
                "description": "Weather modeling, climate simulation, atmospheric physics",
                "tools": ["WRF", "CESM", "NCO/CDO"],
                "cost_range": "$300-1500/month",
                "icon": "🌍"
            },
            "Materials Science": {
                "description": "DFT calculations, molecular dynamics, materials screening",
                "tools": ["VASP", "Quantum ESPRESSO", "LAMMPS"],
                "cost_range": "$400-2000/month",
                "icon": "⚛️"
            },
            "Machine Learning": {
                "description": "Deep learning, computer vision, NLP, model training",
                "tools": ["PyTorch", "TensorFlow", "scikit-learn"],
                "cost_range": "$200-1000/month",
                "icon": "🤖"
            }
        },
        "specialized_packs": {
            "Geospatial Research": {
                "description": "Remote sensing, GIS analysis, earth observation",
                "tools": ["QGIS", "GDAL", "Satellite imagery"],
                "cost_range": "$150-10000/month",
                "icon": "🗺️"
            },
            "Atmospheric Chemistry": {
                "description": "Air quality modeling, chemical transport, GEOS-Chem",
                "tools": ["GEOS-Chem", "CMAQ", "WRF-Chem"],
                "cost_range": "$300-25000/month",
                "icon": "🌬️"
            },
            "Cybersecurity Research": {
                "description": "Threat analysis, malware research, digital forensics",
                "tools": ["MISP", "Ghidra", "Wireshark"],
                "cost_range": "$500-40000/month",
                "icon": "🔒"
            },
            "Renewable Energy": {
                "description": "Solar/wind analysis, grid integration, energy storage",
                "tools": ["SAM", "OpenFAST", "GridLAB-D"],
                "cost_range": "$300-30000/month",
                "icon": "⚡"
            },
            "Visualization Studio": {
                "description": "Scientific visualization, 3D rendering, collaborative analysis",
                "tools": ["ParaView", "VTK", "Jupyter"],
                "cost_range": "$150-1200/month",
                "icon": "📊"
            }
        }
    }

def render_domain_selection():
    """Render the domain selection interface"""
    st.markdown('<h1 class="main-header">🧙‍♂️ AWS Research Wizard</h1>', unsafe_allow_html=True)
    st.markdown("### Choose your research domain to get started")
    
    domain_info = get_domain_info()
    
    # Core domains
    st.subheader("🎯 Core Research Domains")
    core_cols = st.columns(2)
    
    for i, (domain, info) in enumerate(domain_info["core_domains"].items()):
        col = core_cols[i % 2]
        with col:
            with st.container():
                st.markdown(f"""
                <div class="domain-card">
                    <h4>{info['icon']} {domain}</h4>
                    <p>{info['description']}</p>
                    <p><strong>Tools:</strong> {', '.join(info['tools'][:3])}</p>
                    <p><strong>Cost:</strong> {info['cost_range']}</p>
                </div>
                """, unsafe_allow_html=True)
                if st.button(f"Select {domain}", key=f"core_{domain}"):
                    st.session_state.selected_domain = domain
                    st.session_state.current_step = 'workload_config'
                    st.rerun()
    
    # Specialized packs
    st.subheader("🔬 Specialized Research Packs")
    spec_cols = st.columns(3)
    
    for i, (domain, info) in enumerate(domain_info["specialized_packs"].items()):
        col = spec_cols[i % 3]
        with col:
            with st.container():
                st.markdown(f"""
                <div class="domain-card">
                    <h4>{info['icon']} {domain}</h4>
                    <p>{info['description']}</p>
                    <p><strong>Tools:</strong> {', '.join(info['tools'][:3])}</p>
                    <p><strong>Cost:</strong> {info['cost_range']}</p>
                </div>
                """, unsafe_allow_html=True)
                if st.button(f"Select {domain}", key=f"spec_{domain}"):
                    st.session_state.selected_domain = domain
                    st.session_state.current_step = 'workload_config'
                    st.rerun()

def render_workload_configuration():
    """Render the workload configuration interface"""
    st.markdown(f"# Configure Your {st.session_state.selected_domain} Workload")
    
    if st.button("← Back to Domain Selection"):
        st.session_state.current_step = 'domain_selection'
        st.rerun()
    
    # Create configuration form based on selected domain
    with st.form("workload_config_form"):
        col1, col2 = st.columns(2)
        
        with col1:
            st.subheader("🔧 Basic Configuration")
            
            # Common configuration options
            problem_size = st.selectbox(
                "Problem Size",
                ["Small", "Medium", "Large", "Massive"],
                help="Small: <50GB data, Medium: 50GB-500GB, Large: 500GB-5TB, Massive: >5TB"
            )
            
            priority = st.selectbox(
                "Optimization Priority",
                ["Cost", "Performance", "Balanced", "Deadline"],
                index=2,
                help="Cost: Minimize expenses, Performance: Fastest results, Balanced: Good compromise"
            )
            
            team_size = st.number_input(
                "Team Size",
                min_value=1,
                max_value=50,
                value=3,
                help="Number of researchers who will use this environment"
            )
            
            data_size_gb = st.number_input(
                "Data Size (GB)",
                min_value=1,
                max_value=100000,
                value=100,
                help="Expected dataset size in gigabytes"
            )
        
        with col2:
            st.subheader("⚙️ Advanced Configuration")
            
            # Domain-specific configuration
            if st.session_state.selected_domain == "Geospatial Research":
                spatial_resolution = st.selectbox(
                    "Spatial Resolution",
                    ["High (1-10m)", "Medium (10-100m)", "Low (>100m)"]
                )
                
                coverage_area = st.selectbox(
                    "Coverage Area",
                    ["Local", "Regional", "National", "Global"]
                )
                
                analysis_type = st.selectbox(
                    "Analysis Type",
                    ["Remote Sensing", "GIS Analysis", "Environmental Modeling", "Precision Agriculture"]
                )
            
            elif st.session_state.selected_domain == "Atmospheric Chemistry":
                model_type = st.selectbox(
                    "Model Type",
                    ["GEOS-Chem", "CMAQ", "WRF-Chem", "Chemical Transport"]
                )
                
                spatial_scale = st.selectbox(
                    "Spatial Scale",
                    ["Global", "Regional", "Urban", "Local"]
                )
                
                chemical_complexity = st.selectbox(
                    "Chemical Complexity",
                    ["Basic", "Standard", "Full", "Custom"]
                )
            
            elif st.session_state.selected_domain == "Cybersecurity Research":
                research_type = st.selectbox(
                    "Research Type",
                    ["Academic", "Industry", "Government", "Red Team"]
                )
                
                data_sensitivity = st.selectbox(
                    "Data Sensitivity",
                    ["Public", "Confidential", "Restricted", "Top Secret"]
                )
                
                analysis_scale = st.selectbox(
                    "Analysis Scale",
                    ["Individual", "Enterprise", "National", "Global"]
                )
            
            elif st.session_state.selected_domain == "Renewable Energy":
                energy_type = st.selectbox(
                    "Energy Type",
                    ["Solar", "Wind", "Energy Storage", "Grid Integration", "Hydroelectric", "Geothermal"]
                )
                
                analysis_focus = st.selectbox(
                    "Analysis Focus",
                    ["Resource Assessment", "System Design", "Economic Analysis", "Grid Impact"]
                )
                
                temporal_scale = st.selectbox(
                    "Temporal Scale",
                    ["Real-time", "Hourly", "Daily", "Monthly", "Annual"]
                )
            
            # Common advanced options
            gpu_required = st.checkbox(
                "GPU Acceleration Required",
                help="Check if your workload requires GPU processing"
            )
            
            real_time_req = st.checkbox(
                "Real-time Processing",
                help="Check if you need real-time data processing capabilities"
            )
        
        # Submit button
        submitted = st.form_submit_button(
            "🚀 Generate Recommendation",
            type="primary",
            use_container_width=True
        )
        
        if submitted:
            # Store configuration in session state
            st.session_state.workload_config = {
                'domain': st.session_state.selected_domain,
                'problem_size': problem_size,
                'priority': priority,
                'team_size': team_size,
                'data_size_gb': data_size_gb,
                'gpu_required': gpu_required,
                'real_time_req': real_time_req
            }
            
            # Add domain-specific config
            if st.session_state.selected_domain == "Geospatial Research":
                st.session_state.workload_config.update({
                    'spatial_resolution': spatial_resolution,
                    'coverage_area': coverage_area,
                    'analysis_type': analysis_type
                })
            elif st.session_state.selected_domain == "Atmospheric Chemistry":
                st.session_state.workload_config.update({
                    'model_type': model_type,
                    'spatial_scale': spatial_scale,
                    'chemical_complexity': chemical_complexity
                })
            elif st.session_state.selected_domain == "Cybersecurity Research":
                st.session_state.workload_config.update({
                    'research_type': research_type,
                    'data_sensitivity': data_sensitivity,
                    'analysis_scale': analysis_scale
                })
            elif st.session_state.selected_domain == "Renewable Energy":
                st.session_state.workload_config.update({
                    'energy_type': energy_type,
                    'analysis_focus': analysis_focus,
                    'temporal_scale': temporal_scale
                })
            
            # Generate recommendation
            generate_recommendation()
            st.session_state.current_step = 'recommendation'
            st.rerun()

def generate_recommendation():
    """Generate infrastructure recommendation based on configuration"""
    config = st.session_state.workload_config
    
    try:
        if config['domain'] == "Geospatial Research":
            pack = GeospatialResearchPack()
            # Create workload object - simplified for demo
            workload = GeospatialWorkload(
                domain=GeospatialDomain.REMOTE_SENSING,  # Default, could map from analysis_type
                data_sources=["satellite"],
                processing_type="analysis",
                spatial_resolution="Medium",
                temporal_frequency="Weekly",
                coverage_area=config.get('coverage_area', 'Regional'),
                data_volume_tb=config['data_size_gb'] / 1000,
                analysis_complexity="Moderate"
            )
            recommendation = pack.generate_deployment_recommendation(workload)
        
        elif config['domain'] == "Atmospheric Chemistry":
            pack = AtmosphericChemistryPack()
            workload = AtmosphericWorkload(
                domain=AtmosphericDomain.CHEMICAL_TRANSPORT,
                model_type=config.get('model_type', 'GEOS-Chem'),
                spatial_resolution=config.get('spatial_scale', 'Regional'),
                temporal_scale="Daily",
                chemical_complexity=config.get('chemical_complexity', 'Standard'),
                emission_sources=["anthropogenic", "biogenic"],
                data_volume_tb=config['data_size_gb'] / 1000,
                computational_intensity="Moderate"
            )
            recommendation = pack.generate_atmospheric_recommendation(workload)
        
        elif config['domain'] == "Cybersecurity Research":
            pack = CyberSecurityResearchPack()
            workload = CyberSecurityWorkload(
                domain=CyberSecurityDomain.THREAT_INTELLIGENCE,
                research_type=config.get('research_type', 'Academic'),
                data_sensitivity=config.get('data_sensitivity', 'Confidential'),
                analysis_scale=config.get('analysis_scale', 'Individual'),
                real_time_req=config.get('real_time_req', False),
                compliance_level="Basic",
                data_volume_tb=config['data_size_gb'] / 1000,
                computational_intensity="Moderate"
            )
            recommendation = pack.generate_cybersecurity_recommendation(workload)
        
        elif config['domain'] == "Renewable Energy":
            pack = RenewableEnergyResearchPack()
            workload = RenewableEnergyWorkload(
                domain=RenewableEnergyDomain.SOLAR_ENERGY,  # Default, could map from energy_type
                analysis_type=config.get('analysis_focus', 'System Design'),
                temporal_scale=config.get('temporal_scale', 'Daily'),
                spatial_scale="Regional",
                modeling_complexity="Intermediate",
                data_sources=["Weather", "Market"],
                optimization_focus="Cost",
                data_volume_tb=config['data_size_gb'] / 1000
            )
            recommendation = pack.generate_renewable_energy_recommendation(workload)
        
        else:
            # For core domains, use the main wizard
            wizard = ResearchInfrastructureWizard()
            
            # Map domain to wizard domain
            domain_mapping = {
                "Genomics & Bioinformatics": "genomics",
                "Climate Science": "climate_modeling",
                "Materials Science": "materials_science",
                "Machine Learning": "machine_learning"
            }
            
            wizard_domain = domain_mapping.get(config['domain'], 'genomics')
            
            workload = WorkloadCharacteristics(
                domain=wizard_domain,
                primary_tools=["default"],
                problem_size=WorkloadSize(config['problem_size'].lower()),
                priority=Priority(config['priority'].lower()),
                deadline_hours=None,
                budget_limit=None,
                data_size_gb=config['data_size_gb'],
                parallel_scaling="linear",
                gpu_requirement="required" if config['gpu_required'] else "none",
                memory_intensity="medium",
                io_pattern="sequential",
                collaboration_users=config['team_size']
            )
            
            recommendation = wizard.generate_recommendation(workload)
        
        st.session_state.recommendation = recommendation
    
    except Exception as e:
        st.error(f"Error generating recommendation: {e}")
        st.session_state.recommendation = None

def render_recommendation():
    """Render the recommendation results"""
    if not st.session_state.recommendation:
        st.error("No recommendation available. Please go back and reconfigure.")
        return
    
    st.markdown(f"# 📋 Recommendation for {st.session_state.selected_domain}")
    
    if st.button("← Back to Configuration"):
        st.session_state.current_step = 'workload_config'
        st.rerun()
    
    if st.button("🔄 Start Over"):
        st.session_state.current_step = 'domain_selection'
        st.session_state.selected_domain = None
        st.session_state.workload_config = {}
        st.session_state.recommendation = None
        st.rerun()
    
    recommendation = st.session_state.recommendation
    
    # Cost Summary
    col1, col2, col3 = st.columns(3)
    
    if 'estimated_cost' in recommendation:
        cost_data = recommendation['estimated_cost']
        
        with col1:
            st.metric(
                "Monthly Compute Cost",
                f"${cost_data.get('compute', 0):,.0f}",
                help="Estimated monthly compute costs"
            )
        
        with col2:
            st.metric(
                "Monthly Storage Cost",
                f"${cost_data.get('storage', 0):,.0f}",
                help="Estimated monthly storage costs"
            )
        
        with col3:
            st.metric(
                "Total Monthly Cost",
                f"${cost_data.get('total', 0):,.0f}",
                help="Total estimated monthly costs"
            )
    
    # Instance Recommendations
    st.subheader("💻 Recommended AWS Infrastructure")
    
    if 'configuration' in recommendation:
        config = recommendation['configuration']
        
        if 'aws_instance_recommendations' in config:
            instances = config['aws_instance_recommendations']
            
            # Create a DataFrame for instance comparison
            instance_data = []
            for name, details in instances.items():
                instance_data.append({
                    'Configuration': name.replace('_', ' ').title(),
                    'Instance Type': details.get('instance_type', 'N/A'),
                    'vCPUs': details.get('vcpus', 'N/A'),
                    'Memory (GB)': details.get('memory_gb', 'N/A'),
                    'Cost/Hour': f"${details.get('cost_per_hour', 0):.3f}",
                    'Use Case': details.get('use_case', 'N/A')
                })
            
            if instance_data:
                df = pd.DataFrame(instance_data)
                st.dataframe(df, use_container_width=True)
        
        # Tools and Software
        if 'spack_packages' in config:
            st.subheader("🛠️ Included Software and Tools")
            
            packages = config['spack_packages']
            
            # Group packages by category (simplified)
            core_tools = [p for p in packages if any(tool in p.lower() for tool in ['python', 'gcc', 'openmpi'])]
            domain_tools = [p for p in packages if p not in core_tools][:10]  # Show first 10 domain-specific tools
            
            col1, col2 = st.columns(2)
            
            with col1:
                st.write("**Core Infrastructure:**")
                for tool in core_tools[:5]:
                    st.write(f"• {tool.split('@')[0]}")
            
            with col2:
                st.write("**Domain-Specific Tools:**")
                for tool in domain_tools:
                    st.write(f"• {tool.split('@')[0]}")
        
        # Research Capabilities
        if 'research_capabilities' in config:
            st.subheader("🔬 Research Capabilities")
            capabilities = config['research_capabilities']
            
            # Display in columns
            cols = st.columns(2)
            for i, capability in enumerate(capabilities):
                col = cols[i % 2]
                col.write(f"• {capability}")
    
    # Cost Breakdown Chart
    if 'estimated_cost' in recommendation:
        st.subheader("💰 Cost Breakdown")
        
        cost_data = recommendation['estimated_cost']
        
        # Create pie chart
        labels = []
        values = []
        
        for key, value in cost_data.items():
            if key != 'total' and value > 0:
                labels.append(key.title())
                values.append(value)
        
        if labels and values:
            fig = px.pie(
                values=values,
                names=labels,
                title="Monthly Cost Distribution"
            )
            st.plotly_chart(fig, use_container_width=True)
    
    # Optimization Recommendations
    if 'optimization_recommendations' in recommendation:
        st.subheader("⚡ Optimization Recommendations")
        
        recommendations = recommendation['optimization_recommendations']
        for rec in recommendations:
            st.markdown(f"• {rec}")
    
    # Deployment Instructions
    st.subheader("🚀 Next Steps")
    
    st.markdown("""
    **To deploy this configuration:**
    
    1. **Download Configuration**: Click the button below to download your configuration
    2. **Review Costs**: Verify the cost estimates meet your budget requirements
    3. **Deploy Infrastructure**: Use the AWS Research Wizard CLI or Terraform templates
    4. **Monitor Usage**: Set up CloudWatch monitoring and cost alerts
    
    **Need Help?**
    - 📖 Read the [AWS Research Wizard Documentation](https://docs.aws-research-wizard.com)
    - 💬 Join our [Community Forum](https://community.aws-research-wizard.com)
    - 📧 Contact [Support](mailto:support@aws-research-wizard.com)
    """)
    
    # Download configuration
    if st.button("📥 Download Configuration", type="primary"):
        config_json = json.dumps(recommendation, indent=2)
        st.download_button(
            label="Download JSON Configuration",
            data=config_json,
            file_name=f"aws_research_config_{st.session_state.selected_domain.lower().replace(' ', '_')}.json",
            mime="application/json"
        )

def render_about():
    """Render the about page"""
    st.markdown("# About AWS Research Wizard")
    
    st.markdown("""
    ## 🎯 Overview
    
    The AWS Research Wizard is a comprehensive suite of tools designed to bridge the gap between research computing needs and optimal AWS infrastructure deployment. It provides domain-specific, cost-optimized, and performance-tuned solutions for researchers across 25+ scientific disciplines.
    
    ## ✨ Key Features
    
    ### 🧙‍♂️ Intelligent Infrastructure Wizard
    - **Domain-aware recommendations**: Research-specific questions get optimal AWS infrastructure
    - **Cost/performance/deadline optimization**: Balance competing priorities automatically
    - **Workload-aware instance selection**: Choose from 400+ AWS instance types intelligently
    - **Alternative configurations**: Compare cost-optimized vs performance-optimized deployments
    
    ### 🔬 Domain-Specific Solutions
    - **25+ Research Domains**: From genomics to digital humanities, each with tailored toolstacks
    - **Spack-powered environments**: Optimized, reproducible software deployment
    - **Ready-to-run workflows**: Pre-configured pipelines for immediate productivity
    - **Transparent pricing**: Clear cost estimates from $0 idle to $3000+/day for massive simulations
    
    ### 💰 FinOps-First Architecture
    - **Ephemeral computing**: Pay only for active compute, $0 when idle
    - **Intelligent storage tiering**: Hot → Warm → Cold storage optimization
    - **Auto-scaling**: Dynamic resource allocation based on workload
    - **Cost monitoring**: Real-time spend tracking and budget alerts
    
    ## 📊 Supported Research Domains
    
    **Core Domains:**
    - Genomics & Bioinformatics
    - Climate Science & Atmospheric Physics
    - Materials Science & Computational Chemistry
    - Machine Learning & AI
    - Physics Simulation
    
    **Specialized Packs:**
    - Geospatial Research & Remote Sensing
    - Atmospheric Chemistry & Air Quality
    - Cybersecurity Research
    - Renewable Energy Systems
    - Visualization Studio
    
    ## 🏗️ Architecture
    
    **Core Components:**
    - Interactive GUI (this interface)
    - Command-line interface
    - Spack environment capture
    - Custom configuration generator
    - Cost optimization engine
    
    ## 🤝 Contributing
    
    We welcome contributions from the research computing community! Visit our [GitHub repository](https://github.com/aws-research-wizard) to get involved.
    """)

def main():
    """Main application logic"""
    initialize_session_state()
    
    # Sidebar navigation
    with st.sidebar:
        st.title("🧙‍♂️ AWS Research Wizard")
        st.markdown("---")
        
        # Navigation
        if st.button("🏠 Home", use_container_width=True):
            st.session_state.current_step = 'domain_selection'
            st.rerun()
        
        if st.button("ℹ️ About", use_container_width=True):
            st.session_state.current_step = 'about'
            st.rerun()
        
        st.markdown("---")
        
        # Current step indicator
        if st.session_state.current_step == 'domain_selection':
            st.markdown("**📍 Current Step:** Domain Selection")
        elif st.session_state.current_step == 'workload_config':
            st.markdown(f"**📍 Current Step:** Configuring {st.session_state.selected_domain}")
        elif st.session_state.current_step == 'recommendation':
            st.markdown(f"**📍 Current Step:** Recommendation Ready")
        elif st.session_state.current_step == 'about':
            st.markdown("**📍 Current Page:** About")
        
        # Configuration summary
        if st.session_state.workload_config:
            st.markdown("---")
            st.markdown("**🔧 Configuration Summary:**")
            config = st.session_state.workload_config
            st.write(f"Domain: {config.get('domain', 'N/A')}")
            st.write(f"Size: {config.get('problem_size', 'N/A')}")
            st.write(f"Team: {config.get('team_size', 'N/A')} users")
            st.write(f"Data: {config.get('data_size_gb', 'N/A')} GB")
        
        st.markdown("---")
        st.markdown("**🆘 Need Help?**")
        st.markdown("[📖 Documentation](https://docs.aws-research-wizard.com)")
        st.markdown("[💬 Community](https://community.aws-research-wizard.com)")
        st.markdown("[🐛 Report Issue](https://github.com/aws-research-wizard/issues)")
    
    # Main content area
    if st.session_state.current_step == 'domain_selection':
        render_domain_selection()
    elif st.session_state.current_step == 'workload_config':
        render_workload_configuration()
    elif st.session_state.current_step == 'recommendation':
        render_recommendation()
    elif st.session_state.current_step == 'about':
        render_about()

if __name__ == "__main__":
    main()