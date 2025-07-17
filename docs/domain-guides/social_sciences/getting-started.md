# Social Sciences Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $8-14 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working social sciences research environment that can:
- Analyze survey data and social network structures
- Process large-scale demographic and sociological datasets
- Run statistical models for social research
- Handle census data, social media data, and survey responses

### Meet Dr. Sarah Kim

Dr. Sarah Kim is a sociologist at University of Chicago. She studies social inequality but waits weeks for university computing resources. Each analysis requires processing millions of survey responses and census records.

**Before**: 2-week waits + 5-day analysis = 3 weeks per study
**After**: 15-minute setup + 3-hour analysis = same day results
**Time Saved**: 95% faster social research cycle
**Cost Savings**: $300/month vs $1,200 university allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $8-14 (we'll clean up resources when done)
- **Daily research cost**: $12-30 per day when actively analyzing
- **Monthly estimate**: $150-400 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No social science or programming experience required

## Step 1: Install AWS Research Wizard

Choose your operating system:

### macOS/Linux
```bash
curl -fsSL https://install.aws-research-wizard.com | sh
```

### Windows
Download from: https://github.com/aws-research-wizard/releases/latest

**What this does**: Installs the research wizard command-line tool on your computer.

**Expected result**: You should see "Installation successful" message.

**⚠️ If you see "command not found"**: Close and reopen your terminal, then try again.

## Step 2: Set Up AWS Account

If you don't have an AWS account:

1. Go to [aws.amazon.com](https://aws.amazon.com)
2. Click "Create an AWS Account"
3. Follow the signup process
4. **Important**: Choose the free tier options

**What this does**: Creates your personal cloud computing account.

**Expected result**: You receive email confirmation from AWS.

**💰 Cost note**: Account creation is free. You only pay for resources you use.

## Step 3: Configure Your Credentials

```bash
aws-research-wizard config setup
```

The wizard will ask for:
- **AWS Access Key**: Found in AWS Console → Security Credentials
- **Secret Key**: Created with your access key
- **Region**: Choose `us-east-1` (recommended for social sciences with good data access)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain social_sciences --region us-east-1
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: social_sciences
✅ Region valid: us-east-1 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Social Sciences Environment

```bash
aws-research-wizard deploy start --domain social_sciences --region us-east-1 --instance m6i.large
```

**What this does**: Creates your social sciences environment optimized for statistical analysis and data processing.

**This will take**: 5-7 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
  CPU: 2 cores for statistical computing
  Memory: 8GB RAM for large datasets
```

**💰 Billing starts now**: Your environment costs about $0.19 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
```

**What this does**: Connects you to your social sciences computer in the cloud.

**Expected result**: You see a command prompt like `ubuntu@ip-10-0-1-123:~$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Social Sciences Tools

Your environment comes pre-installed with:

### Core Research Tools
- **R Statistical Software**: Statistical analysis - Type `R --version` to check
- **Python Scientific Stack**: NumPy, Pandas, SciPy - Type `python -c "import pandas; print(pandas.__version__)"` to check
- **SPSS Syntax Support**: Statistical package compatibility - Type `python -c "import pyreadstat; print('SPSS support available')"` to check
- **Jupyter Notebooks**: Interactive analysis - Type `jupyter --version` to check
- **NetworkX**: Social network analysis - Type `python -c "import networkx; print(networkx.__version__)"` to check

### Try Your First Command
```bash
python -c "import pandas; print('Pandas version:', pandas.__version__)"
```

**What this does**: Shows Pandas version and confirms data analysis tools are installed.

**Expected result**: You see Pandas version info confirming social science libraries are ready.

## Step 8: Analyze Survey Data

Let's analyze social research data to test everything works:

### Create Sample Survey Data
```bash
# Create working directory
mkdir ~/social-science-tutorial
cd ~/social-science-tutorial

# Create survey data analysis script
cat > survey_analysis.py << 'EOF'
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from scipy import stats

print("Starting social sciences survey analysis...")

def generate_survey_data():
    """Generate synthetic survey data for analysis"""
    print("\n=== Survey Data Generation ===")

    np.random.seed(42)
    n_respondents = 2500

    # Demographic variables
    ages = np.random.normal(45, 15, n_respondents)
    ages = np.clip(ages, 18, 85).astype(int)

    # Gender (binary for simplicity)
    genders = np.random.choice(['Male', 'Female'], n_respondents, p=[0.48, 0.52])

    # Education levels
    education_levels = np.random.choice([
        'High School', 'Some College', 'Bachelor\'s', 'Master\'s', 'PhD'
    ], n_respondents, p=[0.25, 0.30, 0.25, 0.15, 0.05])

    # Income (correlated with education and age)
    base_income = 35000
    education_multipliers = {
        'High School': 1.0,
        'Some College': 1.2,
        'Bachelor\'s': 1.6,
        'Master\'s': 2.0,
        'PhD': 2.4
    }

    incomes = []
    for edu, age in zip(education_levels, ages):
        multiplier = education_multipliers[edu]
        age_factor = 1 + (age - 25) * 0.01  # Income increases with age
        income = base_income * multiplier * age_factor * np.random.lognormal(0, 0.3)
        incomes.append(max(15000, income))  # Minimum wage floor

    # Likert scale responses (1-5)
    # Job satisfaction (correlated with income)
    job_satisfaction = []
    for income in incomes:
        base_satisfaction = 2.5 + (income - 35000) / 100000  # Higher income = higher satisfaction
        satisfaction = np.random.normal(base_satisfaction, 0.8)
        job_satisfaction.append(np.clip(satisfaction, 1, 5))

    # Life satisfaction (correlated with multiple factors)
    life_satisfaction = []
    for i, (income, job_sat, age) in enumerate(zip(incomes, job_satisfaction, ages)):
        base_life_sat = 2.8 + (job_sat - 3) * 0.4 + (income - 50000) / 150000
        # Adjust for age (U-shaped curve)
        age_adjustment = -0.3 * ((age - 50) / 20) ** 2
        life_sat = np.random.normal(base_life_sat + age_adjustment, 0.7)
        life_satisfaction.append(np.clip(life_sat, 1, 5))

    # Political views (1 = very liberal, 5 = very conservative)
    political_views = []
    for age, edu in zip(ages, education_levels):
        base_political = 3.0  # Center
        # Age effect (older = more conservative)
        age_effect = (age - 40) * 0.01
        # Education effect (higher education = more liberal)
        edu_effects = {
            'High School': 0.2,
            'Some College': 0.1,
            'Bachelor\'s': -0.1,
            'Master\'s': -0.2,
            'PhD': -0.3
        }
        political = np.random.normal(base_political + age_effect + edu_effects[edu], 0.8)
        political_views.append(np.clip(political, 1, 5))

    # Create DataFrame
    survey_data = pd.DataFrame({
        'respondent_id': range(1, n_respondents + 1),
        'age': ages,
        'gender': genders,
        'education': education_levels,
        'income': incomes,
        'job_satisfaction': job_satisfaction,
        'life_satisfaction': life_satisfaction,
        'political_views': political_views
    })

    print(f"Generated survey data: {len(survey_data)} respondents")
    print(f"Age range: {survey_data['age'].min()}-{survey_data['age'].max()}")
    print(f"Income range: ${survey_data['income'].min():,.0f}-${survey_data['income'].max():,.0f}")

    return survey_data

def descriptive_analysis(survey_data):
    """Perform descriptive analysis of survey data"""
    print("\n=== Descriptive Analysis ===")

    # Basic demographics
    print("Demographics Summary:")
    print(f"  Total respondents: {len(survey_data)}")
    print(f"  Mean age: {survey_data['age'].mean():.1f} years")
    print(f"  Gender distribution:")
    gender_counts = survey_data['gender'].value_counts()
    for gender, count in gender_counts.items():
        percentage = (count / len(survey_data)) * 100
        print(f"    {gender}: {count} ({percentage:.1f}%)")

    # Education distribution
    print(f"  Education distribution:")
    edu_counts = survey_data['education'].value_counts()
    for edu, count in edu_counts.items():
        percentage = (count / len(survey_data)) * 100
        print(f"    {edu}: {count} ({percentage:.1f}%)")

    # Income statistics
    print(f"  Income statistics:")
    print(f"    Mean: ${survey_data['income'].mean():,.0f}")
    print(f"    Median: ${survey_data['income'].median():,.0f}")
    print(f"    Standard deviation: ${survey_data['income'].std():,.0f}")

    # Likert scale variables
    likert_vars = ['job_satisfaction', 'life_satisfaction', 'political_views']
    print(f"  Likert scale variables (1-5):")
    for var in likert_vars:
        mean_score = survey_data[var].mean()
        print(f"    {var.replace('_', ' ').title()}: {mean_score:.2f}")

    return survey_data.describe()

def correlation_analysis(survey_data):
    """Analyze correlations between variables"""
    print("\n=== Correlation Analysis ===")

    # Select numeric variables
    numeric_vars = ['age', 'income', 'job_satisfaction', 'life_satisfaction', 'political_views']
    correlation_matrix = survey_data[numeric_vars].corr()

    print("Correlation Matrix:")
    print(correlation_matrix.round(3))

    # Identify significant correlations
    significant_correlations = []
    for i in range(len(numeric_vars)):
        for j in range(i+1, len(numeric_vars)):
            var1, var2 = numeric_vars[i], numeric_vars[j]
            corr_value = correlation_matrix.loc[var1, var2]

            # Calculate p-value for correlation
            r, p_value = stats.pearsonr(survey_data[var1], survey_data[var2])

            if abs(corr_value) > 0.1 and p_value < 0.05:
                significant_correlations.append((var1, var2, corr_value, p_value))

    print(f"\nSignificant correlations (|r| > 0.1, p < 0.05):")
    for var1, var2, r, p in significant_correlations:
        strength = "strong" if abs(r) > 0.5 else "moderate" if abs(r) > 0.3 else "weak"
        direction = "positive" if r > 0 else "negative"
        print(f"  {var1} ↔ {var2}: r = {r:.3f} (p = {p:.3f}) - {strength} {direction}")

    return correlation_matrix

def demographic_analysis(survey_data):
    """Analyze differences across demographic groups"""
    print("\n=== Demographic Group Analysis ===")

    # Age group analysis
    survey_data['age_group'] = pd.cut(survey_data['age'],
                                    bins=[18, 30, 45, 60, 85],
                                    labels=['18-30', '31-45', '46-60', '61+'])

    age_group_analysis = survey_data.groupby('age_group')[
        ['income', 'job_satisfaction', 'life_satisfaction']
    ].mean()

    print("Analysis by Age Group:")
    print(age_group_analysis.round(2))

    # Gender analysis
    gender_analysis = survey_data.groupby('gender')[
        ['income', 'job_satisfaction', 'life_satisfaction', 'political_views']
    ].mean()

    print(f"\nAnalysis by Gender:")
    print(gender_analysis.round(2))

    # Education analysis
    education_analysis = survey_data.groupby('education')[
        ['income', 'job_satisfaction', 'life_satisfaction']
    ].mean().sort_values('income')

    print(f"\nAnalysis by Education Level:")
    print(education_analysis.round(2))

    # Statistical tests
    print(f"\nStatistical Tests:")

    # T-test for gender differences in income
    male_income = survey_data[survey_data['gender'] == 'Male']['income']
    female_income = survey_data[survey_data['gender'] == 'Female']['income']
    t_stat, p_value = stats.ttest_ind(male_income, female_income)

    print(f"  Gender income difference (t-test): t = {t_stat:.3f}, p = {p_value:.3f}")
    if p_value < 0.05:
        mean_diff = male_income.mean() - female_income.mean()
        print(f"    Significant difference: ${mean_diff:,.0f}")

    # ANOVA for education level differences in satisfaction
    education_groups = [group['life_satisfaction'].values for name, group in survey_data.groupby('education')]
    f_stat, p_value = stats.f_oneway(*education_groups)

    print(f"  Education satisfaction difference (ANOVA): F = {f_stat:.3f}, p = {p_value:.3f}")

    return age_group_analysis, gender_analysis, education_analysis

# Run survey analysis
survey_data = generate_survey_data()
descriptive_stats = descriptive_analysis(survey_data)
correlation_matrix = correlation_analysis(survey_data)
age_analysis, gender_analysis, education_analysis = demographic_analysis(survey_data)

print("\n✅ Survey analysis completed!")
print("Social sciences research environment is ready for advanced analysis")
EOF

python3 survey_analysis.py
```

**What this does**: Analyzes survey data with demographics, correlations, and statistical tests.

**This will take**: 2-3 minutes

### Social Network Analysis
```bash
# Create social network analysis script
cat > network_analysis.py << 'EOF'
import networkx as nx
import numpy as np
import pandas as pd

print("Starting social network analysis...")

def create_social_network():
    """Create a synthetic social network for analysis"""
    print("\n=== Social Network Generation ===")

    np.random.seed(42)

    # Create a scale-free network (common in social networks)
    n_nodes = 500
    G = nx.barabasi_albert_graph(n_nodes, 3)

    # Add node attributes (demographic information)
    for node in G.nodes():
        G.nodes[node]['age'] = np.random.randint(18, 70)
        G.nodes[node]['gender'] = np.random.choice(['M', 'F'])
        G.nodes[node]['education'] = np.random.choice(['HS', 'College', 'Graduate'], p=[0.4, 0.4, 0.2])
        G.nodes[node]['income'] = np.random.lognormal(10.5, 0.5)  # Log-normal income distribution

    # Add edge attributes (relationship strength)
    for edge in G.edges():
        G.edges[edge]['weight'] = np.random.uniform(0.1, 1.0)
        G.edges[edge]['relationship_type'] = np.random.choice(
            ['friend', 'colleague', 'family'], p=[0.6, 0.3, 0.1]
        )

    print(f"Created social network:")
    print(f"  Nodes (people): {G.number_of_nodes()}")
    print(f"  Edges (connections): {G.number_of_edges()}")
    print(f"  Density: {nx.density(G):.4f}")

    return G

def analyze_network_structure(G):
    """Analyze the structure of the social network"""
    print("\n=== Network Structure Analysis ===")

    # Basic network metrics
    print("Basic Network Metrics:")
    print(f"  Number of nodes: {G.number_of_nodes()}")
    print(f"  Number of edges: {G.number_of_edges()}")
    print(f"  Density: {nx.density(G):.4f}")
    print(f"  Is connected: {nx.is_connected(G)}")

    if nx.is_connected(G):
        print(f"  Average shortest path: {nx.average_shortest_path_length(G):.2f}")
        print(f"  Diameter: {nx.diameter(G)}")

    print(f"  Clustering coefficient: {nx.average_clustering(G):.4f}")

    # Degree distribution
    degrees = [d for n, d in G.degree()]
    print(f"\nDegree Distribution:")
    print(f"  Mean degree: {np.mean(degrees):.2f}")
    print(f"  Median degree: {np.median(degrees):.0f}")
    print(f"  Max degree: {max(degrees)}")
    print(f"  Min degree: {min(degrees)}")

    # Components analysis
    if not nx.is_connected(G):
        components = list(nx.connected_components(G))
        print(f"\nConnected Components:")
        print(f"  Number of components: {len(components)}")
        component_sizes = [len(c) for c in components]
        print(f"  Largest component size: {max(component_sizes)}")
        print(f"  Average component size: {np.mean(component_sizes):.1f}")

    return degrees

def centrality_analysis(G):
    """Analyze centrality measures to identify important nodes"""
    print("\n=== Centrality Analysis ===")

    # Calculate different centrality measures
    degree_centrality = nx.degree_centrality(G)
    betweenness_centrality = nx.betweenness_centrality(G, k=100)  # Sample for speed
    closeness_centrality = nx.closeness_centrality(G)
    eigenvector_centrality = nx.eigenvector_centrality(G, max_iter=1000)

    # Convert to DataFrame for analysis
    centrality_df = pd.DataFrame({
        'node': list(G.nodes()),
        'degree_centrality': [degree_centrality[n] for n in G.nodes()],
        'betweenness_centrality': [betweenness_centrality[n] for n in G.nodes()],
        'closeness_centrality': [closeness_centrality[n] for n in G.nodes()],
        'eigenvector_centrality': [eigenvector_centrality[n] for n in G.nodes()]
    })

    print("Centrality Measures Summary:")
    centrality_measures = ['degree_centrality', 'betweenness_centrality',
                          'closeness_centrality', 'eigenvector_centrality']

    for measure in centrality_measures:
        values = centrality_df[measure]
        print(f"  {measure.replace('_', ' ').title()}:")
        print(f"    Mean: {values.mean():.4f}")
        print(f"    Std: {values.std():.4f}")
        print(f"    Max: {values.max():.4f}")

    # Identify top central nodes
    print(f"\nTop 5 Most Central Nodes:")
    for measure in centrality_measures:
        top_nodes = centrality_df.nlargest(5, measure)
        print(f"  {measure.replace('_', ' ').title()}:")
        for _, row in top_nodes.iterrows():
            print(f"    Node {row['node']}: {row[measure]:.4f}")

    # Correlation between centrality measures
    centrality_corr = centrality_df[centrality_measures].corr()
    print(f"\nCentrality Measure Correlations:")
    print(centrality_corr.round(3))

    return centrality_df

def community_detection(G):
    """Detect communities in the social network"""
    print("\n=== Community Detection ===")

    # Use Louvain method for community detection
    try:
        import community as community_louvain
        partition = community_louvain.best_partition(G)
        modularity = community_louvain.modularity(partition, G)
    except ImportError:
        # Fallback to basic community detection
        communities = list(nx.community.greedy_modularity_communities(G))
        partition = {}
        for i, community in enumerate(communities):
            for node in community:
                partition[node] = i
        modularity = nx.community.modularity(G, communities)

    # Analyze communities
    community_sizes = {}
    for node, comm_id in partition.items():
        if comm_id not in community_sizes:
            community_sizes[comm_id] = 0
        community_sizes[comm_id] += 1

    print(f"Community Detection Results:")
    print(f"  Number of communities: {len(community_sizes)}")
    print(f"  Modularity: {modularity:.4f}")
    print(f"  Largest community size: {max(community_sizes.values())}")
    print(f"  Smallest community size: {min(community_sizes.values())}")
    print(f"  Average community size: {np.mean(list(community_sizes.values())):.1f}")

    # Community size distribution
    size_distribution = {}
    for size in community_sizes.values():
        if size not in size_distribution:
            size_distribution[size] = 0
        size_distribution[size] += 1

    print(f"\nCommunity Size Distribution:")
    for size in sorted(size_distribution.keys()):
        count = size_distribution[size]
        print(f"  Size {size}: {count} communities")

    return partition, modularity

def homophily_analysis(G):
    """Analyze homophily (tendency to connect with similar others)"""
    print("\n=== Homophily Analysis ===")

    # Analyze gender homophily
    gender_homophily = 0
    total_edges = 0

    for edge in G.edges():
        node1, node2 = edge
        if G.nodes[node1]['gender'] == G.nodes[node2]['gender']:
            gender_homophily += 1
        total_edges += 1

    gender_homophily_rate = gender_homophily / total_edges
    print(f"Gender Homophily:")
    print(f"  Same-gender connections: {gender_homophily}/{total_edges} ({gender_homophily_rate:.3f})")

    # Expected rate if connections were random
    gender_counts = {'M': 0, 'F': 0}
    for node in G.nodes():
        gender_counts[G.nodes[node]['gender']] += 1

    p_male = gender_counts['M'] / G.number_of_nodes()
    expected_same_gender = p_male**2 + (1-p_male)**2

    print(f"  Expected random rate: {expected_same_gender:.3f}")
    print(f"  Homophily index: {(gender_homophily_rate - expected_same_gender) / (1 - expected_same_gender):.3f}")

    # Analyze education homophily
    education_homophily = 0
    for edge in G.edges():
        node1, node2 = edge
        if G.nodes[node1]['education'] == G.nodes[node2]['education']:
            education_homophily += 1

    education_homophily_rate = education_homophily / total_edges
    print(f"\nEducation Homophily:")
    print(f"  Same-education connections: {education_homophily}/{total_edges} ({education_homophily_rate:.3f})")

    # Age homophily (similar ages)
    age_homophily = 0
    for edge in G.edges():
        node1, node2 = edge
        age_diff = abs(G.nodes[node1]['age'] - G.nodes[node2]['age'])
        if age_diff <= 10:  # Within 10 years
            age_homophily += 1

    age_homophily_rate = age_homophily / total_edges
    print(f"\nAge Homophily (within 10 years):")
    print(f"  Similar-age connections: {age_homophily}/{total_edges} ({age_homophily_rate:.3f})")

    return {
        'gender_homophily': gender_homophily_rate,
        'education_homophily': education_homophily_rate,
        'age_homophily': age_homophily_rate
    }

# Run social network analysis
social_network = create_social_network()
degrees = analyze_network_structure(social_network)
centrality_data = centrality_analysis(social_network)
communities, modularity = community_detection(social_network)
homophily_results = homophily_analysis(social_network)

print("\n✅ Social network analysis completed!")
print("Advanced social science network analysis capabilities demonstrated")
EOF

python3 network_analysis.py
```

**What this does**: Analyzes social networks including centrality, communities, and homophily patterns.

**Expected result**: Shows network structure analysis and social connection patterns.

## Step 9: Statistical Modeling

Test advanced social science capabilities:

```bash
# Create statistical modeling script
cat > statistical_modeling.py << 'EOF'
import pandas as pd
import numpy as np
from scipy import stats
import matplotlib.pyplot as plt

print("Running advanced statistical modeling for social sciences...")

def regression_analysis():
    """Perform multiple regression analysis"""
    print("\n=== Multiple Regression Analysis ===")

    np.random.seed(42)
    n = 1000

    # Generate synthetic data for regression
    education_years = np.random.normal(14, 3, n)  # Years of education
    experience_years = np.random.normal(15, 8, n)  # Years of work experience
    gender = np.random.choice([0, 1], n)  # 0 = female, 1 = male

    # Generate income with realistic relationships
    # Income increases with education and experience, with gender gap
    income = (2000 * education_years +
              800 * experience_years +
              5000 * gender +  # Gender wage gap
              np.random.normal(0, 8000, n) +
              25000)  # Base income

    income = np.maximum(income, 15000)  # Minimum wage floor

    # Create DataFrame
    regression_data = pd.DataFrame({
        'income': income,
        'education_years': education_years,
        'experience_years': experience_years,
        'gender_male': gender
    })

    print("Regression Data Summary:")
    print(regression_data.describe().round(2))

    # Calculate correlation matrix
    correlation_matrix = regression_data.corr()
    print(f"\nCorrelation Matrix:")
    print(correlation_matrix.round(3))

    # Simple linear regression (income ~ education)
    from scipy.stats import linregress
    slope, intercept, r_value, p_value, std_err = linregress(
        regression_data['education_years'], regression_data['income']
    )

    print(f"\nSimple Linear Regression (Income ~ Education):")
    print(f"  Slope: ${slope:.0f} per year of education")
    print(f"  Intercept: ${intercept:.0f}")
    print(f"  R-squared: {r_value**2:.3f}")
    print(f"  P-value: {p_value:.3e}")

    # Multiple regression simulation (simplified)
    # Calculate partial correlations manually

    # Control for education when looking at experience effect
    education_residuals = regression_data['experience_years'] - (
        np.mean(regression_data['experience_years']) +
        correlation_matrix.loc['experience_years', 'education_years'] *
        (regression_data['education_years'] - np.mean(regression_data['education_years']))
    )

    income_residuals = regression_data['income'] - (
        np.mean(regression_data['income']) +
        correlation_matrix.loc['income', 'education_years'] *
        (regression_data['education_years'] - np.mean(regression_data['education_years']))
    )

    partial_corr = np.corrcoef(education_residuals, income_residuals)[0, 1]

    print(f"\nPartial correlation (Experience-Income, controlling for Education): {partial_corr:.3f}")

    return regression_data

def hypothesis_testing():
    """Perform various hypothesis tests"""
    print("\n=== Hypothesis Testing ===")

    np.random.seed(42)

    # Generate data for hypothesis testing
    # Research question: Do men and women have different job satisfaction scores?
    male_satisfaction = np.random.normal(3.2, 0.8, 300)
    female_satisfaction = np.random.normal(3.0, 0.9, 350)

    # Ensure scores are within 1-5 range
    male_satisfaction = np.clip(male_satisfaction, 1, 5)
    female_satisfaction = np.clip(female_satisfaction, 1, 5)

    print("Hypothesis Test: Gender Differences in Job Satisfaction")
    print(f"Male satisfaction (n={len(male_satisfaction)}): M = {np.mean(male_satisfaction):.3f}, SD = {np.std(male_satisfaction):.3f}")
    print(f"Female satisfaction (n={len(female_satisfaction)}): M = {np.mean(female_satisfaction):.3f}, SD = {np.std(female_satisfaction):.3f}")

    # Independent samples t-test
    t_stat, p_value = stats.ttest_ind(male_satisfaction, female_satisfaction)

    print(f"\nIndependent Samples T-Test:")
    print(f"  t-statistic: {t_stat:.3f}")
    print(f"  p-value: {p_value:.3f}")
    print(f"  Effect size (Cohen's d): {(np.mean(male_satisfaction) - np.mean(female_satisfaction)) / np.sqrt(((len(male_satisfaction)-1)*np.var(male_satisfaction) + (len(female_satisfaction)-1)*np.var(female_satisfaction)) / (len(male_satisfaction) + len(female_satisfaction) - 2)):.3f}")

    if p_value < 0.05:
        print("  Result: Significant difference (p < 0.05)")
    else:
        print("  Result: No significant difference (p ≥ 0.05)")

    # Chi-square test of independence
    # Research question: Is political affiliation related to education level?
    np.random.seed(42)

    # Create contingency table
    education_levels = ['High School', 'College', 'Graduate']
    political_affiliations = ['Liberal', 'Moderate', 'Conservative']

    # Simulate data with some relationship
    contingency_table = np.array([
        [120, 80, 45],   # High School
        [90, 110, 85],   # College
        [65, 70, 40]     # Graduate
    ])

    print(f"\nChi-Square Test: Education Level vs Political Affiliation")
    print("Contingency Table:")
    print(f"{'':12} {'Liberal':>8} {'Moderate':>8} {'Conservative':>8}")
    for i, edu in enumerate(education_levels):
        print(f"{edu:12} {contingency_table[i,0]:8} {contingency_table[i,1]:8} {contingency_table[i,2]:8}")

    chi2, chi2_p, dof, expected = stats.chi2_contingency(contingency_table)

    print(f"\nChi-Square Test Results:")
    print(f"  Chi-square statistic: {chi2:.3f}")
    print(f"  Degrees of freedom: {dof}")
    print(f"  p-value: {chi2_p:.3f}")

    # Calculate Cramér's V (effect size for chi-square)
    n = np.sum(contingency_table)
    cramers_v = np.sqrt(chi2 / (n * (min(contingency_table.shape) - 1)))
    print(f"  Cramér's V (effect size): {cramers_v:.3f}")

    if chi2_p < 0.05:
        print("  Result: Significant association (p < 0.05)")
    else:
        print("  Result: No significant association (p ≥ 0.05)")

    # One-way ANOVA
    # Research question: Do different age groups have different life satisfaction?
    young_adults = np.random.normal(3.4, 0.9, 200)  # 18-30
    middle_aged = np.random.normal(3.1, 1.0, 250)   # 31-50
    older_adults = np.random.normal(3.6, 0.8, 180)  # 51+

    # Clip to valid range
    young_adults = np.clip(young_adults, 1, 5)
    middle_aged = np.clip(middle_aged, 1, 5)
    older_adults = np.clip(older_adults, 1, 5)

    print(f"\nOne-Way ANOVA: Age Group Differences in Life Satisfaction")
    print(f"Young adults (18-30): M = {np.mean(young_adults):.3f}, SD = {np.std(young_adults):.3f}")
    print(f"Middle-aged (31-50): M = {np.mean(middle_aged):.3f}, SD = {np.std(middle_aged):.3f}")
    print(f"Older adults (51+): M = {np.mean(older_adults):.3f}, SD = {np.std(older_adults):.3f}")

    f_stat, anova_p = stats.f_oneway(young_adults, middle_aged, older_adults)

    print(f"\nANOVA Results:")
    print(f"  F-statistic: {f_stat:.3f}")
    print(f"  p-value: {anova_p:.3f}")

    if anova_p < 0.05:
        print("  Result: Significant group differences (p < 0.05)")
    else:
        print("  Result: No significant group differences (p ≥ 0.05)")

    return {
        't_test': {'t': t_stat, 'p': p_value},
        'chi_square': {'chi2': chi2, 'p': chi2_p},
        'anova': {'f': f_stat, 'p': anova_p}
    }

def longitudinal_analysis():
    """Analyze longitudinal data (repeated measures)"""
    print("\n=== Longitudinal Data Analysis ===")

    np.random.seed(42)

    # Simulate 3-wave longitudinal study
    n_participants = 200
    participant_ids = range(1, n_participants + 1)

    # Generate baseline individual differences
    baseline_wellbeing = np.random.normal(3.0, 0.8, n_participants)

    # Simulate 3 time points with some change over time
    wellbeing_data = []

    for wave in range(1, 4):  # Time 1, 2, 3
        # Overall population trend (slight increase over time)
        time_effect = 0.1 * wave

        # Individual variation around trend
        individual_change = np.random.normal(0, 0.3, n_participants)

        # Regression to the mean effect
        regression_effect = -0.2 * (baseline_wellbeing - 3.0)

        wellbeing_wave = baseline_wellbeing + time_effect + individual_change + regression_effect
        wellbeing_wave = np.clip(wellbeing_wave, 1, 5)

        for i, participant_id in enumerate(participant_ids):
            wellbeing_data.append({
                'participant_id': participant_id,
                'wave': wave,
                'wellbeing': wellbeing_wave[i],
                'age_baseline': np.random.randint(25, 65),
                'gender': np.random.choice(['M', 'F'])
            })

    longitudinal_df = pd.DataFrame(wellbeing_data)

    print("Longitudinal Study Summary:")
    print(f"  Participants: {n_participants}")
    print(f"  Time points: 3 waves")
    print(f"  Total observations: {len(longitudinal_df)}")

    # Calculate descriptive statistics by wave
    wave_summary = longitudinal_df.groupby('wave')['wellbeing'].agg(['mean', 'std', 'count'])
    print(f"\nWellbeing by Wave:")
    print(wave_summary.round(3))

    # Repeated measures analysis (simplified)
    # Calculate change scores
    wide_data = longitudinal_df.pivot(index='participant_id', columns='wave', values='wellbeing')
    wide_data.columns = ['wave1', 'wave2', 'wave3']

    # Remove participants with missing data
    complete_data = wide_data.dropna()

    print(f"\nComplete cases for analysis: {len(complete_data)}")

    # Calculate change scores
    change_1_to_2 = complete_data['wave2'] - complete_data['wave1']
    change_2_to_3 = complete_data['wave3'] - complete_data['wave2']
    change_1_to_3 = complete_data['wave3'] - complete_data['wave1']

    print(f"\nChange Score Analysis:")
    print(f"  Wave 1 to 2: M = {change_1_to_2.mean():.3f}, SD = {change_1_to_2.std():.3f}")
    print(f"  Wave 2 to 3: M = {change_2_to_3.mean():.3f}, SD = {change_2_to_3.std():.3f}")
    print(f"  Wave 1 to 3: M = {change_1_to_3.mean():.3f}, SD = {change_1_to_3.std():.3f}")

    # Test if changes are significant
    t_stat_12, p_val_12 = stats.ttest_1samp(change_1_to_2, 0)
    t_stat_23, p_val_23 = stats.ttest_1samp(change_2_to_3, 0)
    t_stat_13, p_val_13 = stats.ttest_1samp(change_1_to_3, 0)

    print(f"\nSignificance Tests (one-sample t-tests against 0):")
    print(f"  Wave 1-2 change: t = {t_stat_12:.3f}, p = {p_val_12:.3f}")
    print(f"  Wave 2-3 change: t = {t_stat_23:.3f}, p = {p_val_23:.3f}")
    print(f"  Wave 1-3 change: t = {t_stat_13:.3f}, p = {p_val_13:.3f}")

    # Stability analysis (test-retest correlation)
    corr_12 = complete_data['wave1'].corr(complete_data['wave2'])
    corr_23 = complete_data['wave2'].corr(complete_data['wave3'])
    corr_13 = complete_data['wave1'].corr(complete_data['wave3'])

    print(f"\nStability Correlations:")
    print(f"  Wave 1-2: r = {corr_12:.3f}")
    print(f"  Wave 2-3: r = {corr_23:.3f}")
    print(f"  Wave 1-3: r = {corr_13:.3f}")

    return longitudinal_df, complete_data

# Run statistical modeling
regression_data = regression_analysis()
hypothesis_results = hypothesis_testing()
longitudinal_df, longitudinal_complete = longitudinal_analysis()

print("\n✅ Statistical modeling completed!")
print("Advanced social science statistical analysis capabilities demonstrated")
EOF

python3 statistical_modeling.py
```

**What this does**: Demonstrates regression analysis, hypothesis testing, and longitudinal data analysis.

**Expected result**: Shows comprehensive statistical modeling results for social science research.

## Step 10: Monitor Your Costs

Check your current spending:

```bash
exit  # Exit SSH session first
aws-research-wizard monitor costs --region us-east-1
```

**Expected result**: Shows costs so far (should be under $5 for this tutorial)

## Step 11: Clean Up (Important!)

When you're done experimenting:

```bash
aws-research-wizard deploy delete --region us-east-1
```

Type `y` when prompted.

**What this does**: Stops billing by removing your cloud resources.

**💰 Important**: Always clean up to avoid ongoing charges.

**Expected result**: "🗑️ Deletion completed successfully"

## Understanding Your Costs

### What You're Paying For
- **Compute**: $0.19 per hour for general-purpose instance while environment is running
- **Storage**: $0.10 per GB per month for research datasets you save
- **Data Transfer**: Usually free for social science data amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 60% savings (advanced)
- Store large datasets in S3, not on the instance
- Process data efficiently to minimize compute time

### Typical Monthly Costs by Usage
- **Light use** (10 hours/week): $75-150
- **Medium use** (3 hours/day): $150-300
- **Heavy use** (6 hours/day): $300-600

## What's Next?

Now that you have a working social sciences environment, you can:

### Learn More About Social Research
- [Large-scale Survey Analysis Tutorial](large-scale-survey-analysis.md)
- [Advanced Social Network Analysis Guide](advanced-social-networks.md)
- [Cost Optimization for Social Sciences](social-sciences-cost-optimization.md)

### Explore Advanced Features
- [Multi-method research integration](multi-method-research.md)
- [Team collaboration with research databases](team-social-science-collaboration.md)
- [Automated research pipelines](automated-social-science-pipelines.md)

### Join the Social Sciences Community
- [Social Sciences Research Forum](https://forum.researchwizard.app/social-sciences)
- [GitHub Social Sciences Examples](https://github.com/aws-research-wizard/social-sciences-examples)
- [Monthly Social Research Office Hours](https://calendar.researchwizard.app/social-sciences-office-hours)

## Troubleshooting

### Common Issues

**Problem**: "R package not found" during statistical analysis
**Solution**: Install missing packages: `R -e "install.packages('package_name')"` or use Python alternatives
**Prevention**: Wait 5-7 minutes after deployment for all statistical packages to initialize

**Problem**: "Memory error" during large survey processing
**Solution**: Process data in smaller chunks or use a larger instance type
**Prevention**: Monitor memory usage with `htop` during analysis

**Problem**: "Statistical test assumption violation"
**Solution**: Check data distributions and consider non-parametric alternatives
**Prevention**: Always examine data with descriptive statistics before testing

**Problem**: "Network analysis import error"
**Solution**: Check NetworkX installation: `python -c "import networkx"` and reinstall if needed
**Prevention**: Verify all required packages are available before starting analysis

### Getting Help
- Check the [social sciences troubleshooting guide](troubleshooting-social-sciences.md)
- Ask in [community forum](https://forum.researchwizard.app)
- File an issue on [GitHub](https://github.com/aws-research-wizard/aws-research-wizard/issues)

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-east-1 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $14. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
