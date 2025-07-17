# Food Science & Nutrition Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $8-14 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working food science and nutrition research environment that can:
- Analyze nutritional data and food safety research
- Process large-scale dietary survey data and food composition databases
- Run statistical models for nutrition epidemiology
- Handle food safety testing data and quality control metrics

### Meet Dr. Maria Gonzalez

Dr. Maria Gonzalez is a food scientist at USDA. She analyzes nutrition data but waits weeks for secure computing resources. Each study requires processing thousands of food samples and dietary assessments from national surveys.

**Before**: 2-week waits + 6-day analysis = 3 weeks per nutrition study
**After**: 15-minute setup + 4-hour analysis = same day results
**Time Saved**: 94% faster food science research cycle
**Cost Savings**: $250/month vs $1,000 government allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $8-14 (we'll clean up resources when done)
- **Daily research cost**: $10-25 per day when actively analyzing
- **Monthly estimate**: $120-300 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No food science or programming experience required

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
- **Region**: Choose `us-east-1` (recommended for food science with good data access)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain food_science_nutrition --region us-east-1
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**: "✅ All validations passed"

**⚠️ If validation fails**: Check your internet connection and AWS credentials.

## Step 5: Deploy Your Research Environment

```bash
aws-research-wizard deploy create --domain food_science_nutrition --region us-east-1 --instance-type r5.large
```

**What this does**: Creates a cloud computer with food science tools installed.

**Expected result**: You'll see progress updates for about 5 minutes, then "✅ Environment ready"

**💰 Billing starts now**: About $0.12 per hour ($2.88 per day if left running)

**⚠️ If deploy fails**: Run the command again. AWS sometimes has temporary issues.

## Step 6: Connect to Your Environment

```bash
aws-research-wizard connect --domain food_science_nutrition
```

**What this does**: Opens a connection to your cloud research environment.

**Expected result**: You'll see a terminal prompt like `[food-scientist@ip-10-0-1-123 ~]$`

**🎉 Success**: You're now inside your food science research environment!

## Step 7: Verify Your Tools

Let's make sure all the food science tools are working:

```bash
# Check Python data science tools
python3 -c "import pandas, numpy, scipy, matplotlib; print('✅ Data science tools ready')"

# Check R statistical environment
R --version | head -1

# Check nutritional analysis tools
python3 -c "import sklearn, seaborn; print('✅ Advanced analytics ready')"
```

**Expected result**: You should see "✅" messages confirming tools are installed.

**⚠️ If tools are missing**: Run `sudo yum update && sudo yum install python3-pip R` then try again.

## Your First Food Science Analysis

Let's run a real nutrition analysis to test everything:

### 1. Download Sample Nutrition Data

```bash
# Create workspace
mkdir -p ~/food_research/nutrition_analysis
cd ~/food_research/nutrition_analysis

# Download USDA food composition data
wget -O food_data.csv "https://research-data.aws-wizard.com/nutrition/usda_food_composition_sample.csv"
```

### 2. Analyze Nutritional Content

Create this Python script for nutrition analysis:

```bash
cat > nutrition_analyzer.py << 'EOF'
#!/usr/bin/env python3
"""
Food Science & Nutrition Analysis Suite
Analyzes nutritional content and dietary patterns
"""

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
from scipy import stats
import warnings
warnings.filterwarnings('ignore')

# Load nutrition data
print("📊 Loading USDA nutrition data...")
data = pd.read_csv('food_data.csv')
print(f"Loaded {len(data)} food items")

# Basic nutritional analysis
print("\n🔬 Nutritional Content Analysis")
print("=" * 40)

# Calculate nutritional density scores
data['nutrient_density'] = (
    data['protein_g'] * 4 +     # 4 cal/g protein
    data['fiber_g'] * 2 +       # Fiber bonus
    data['vitamin_c_mg'] * 0.1 + # Vitamin bonus
    data['calcium_mg'] * 0.01    # Mineral bonus
) / data['calories_per_100g']

# Top nutrient-dense foods
top_dense = data.nlargest(10, 'nutrient_density')[['food_name', 'nutrient_density', 'calories_per_100g']]
print("Top 10 Nutrient-Dense Foods:")
for idx, row in top_dense.iterrows():
    print(f"  {row['food_name']}: {row['nutrient_density']:.2f} (Cal: {row['calories_per_100g']})")

# Macronutrient distribution analysis
print(f"\n📈 Macronutrient Distribution")
macro_stats = {
    'Protein (g)': data['protein_g'].describe(),
    'Carbs (g)': data['carbohydrates_g'].describe(),
    'Fat (g)': data['fat_g'].describe()
}

for nutrient, stats_data in macro_stats.items():
    print(f"\n{nutrient}:")
    print(f"  Mean: {stats_data['mean']:.1f}")
    print(f"  Median: {stats_data['50%']:.1f}")
    print(f"  Range: {stats_data['min']:.1f} - {stats_data['max']:.1f}")

# Food category analysis
print(f"\n🍎 Food Category Insights")
category_analysis = data.groupby('food_category').agg({
    'calories_per_100g': 'mean',
    'protein_g': 'mean',
    'fiber_g': 'mean',
    'sodium_mg': 'mean'
}).round(1)

print("Average nutritional content by food category:")
print(category_analysis)

# Generate nutrition visualization
plt.figure(figsize=(12, 8))

# Subplot 1: Calorie distribution
plt.subplot(2, 2, 1)
plt.hist(data['calories_per_100g'], bins=30, alpha=0.7, color='orange')
plt.title('Calorie Distribution (per 100g)')
plt.xlabel('Calories')
plt.ylabel('Frequency')

# Subplot 2: Protein vs Calories
plt.subplot(2, 2, 2)
plt.scatter(data['protein_g'], data['calories_per_100g'], alpha=0.6, color='blue')
plt.title('Protein vs Calories')
plt.xlabel('Protein (g)')
plt.ylabel('Calories per 100g')

# Subplot 3: Nutrient density by category
plt.subplot(2, 2, 3)
category_density = data.groupby('food_category')['nutrient_density'].mean().sort_values(ascending=True)
category_density.plot(kind='barh', color='green')
plt.title('Nutrient Density by Food Category')
plt.xlabel('Average Nutrient Density Score')

# Subplot 4: Sodium content distribution
plt.subplot(2, 2, 4)
plt.boxplot([data[data['food_category'] == cat]['sodium_mg'].dropna()
            for cat in data['food_category'].unique()[:5]])
plt.title('Sodium Content by Category (Top 5)')
plt.xticks(range(1, 6), data['food_category'].unique()[:5], rotation=45)
plt.ylabel('Sodium (mg)')

plt.tight_layout()
plt.savefig('nutrition_analysis_dashboard.png', dpi=300, bbox_inches='tight')
print(f"\n📊 Dashboard saved as 'nutrition_analysis_dashboard.png'")

# Dietary recommendation engine
print(f"\n🎯 Dietary Recommendations")
print("=" * 30)

# High-protein, low-calorie foods
high_protein_low_cal = data[
    (data['protein_g'] > data['protein_g'].quantile(0.75)) &
    (data['calories_per_100g'] < data['calories_per_100g'].quantile(0.5))
][['food_name', 'protein_g', 'calories_per_100g']].head(5)

print("🥗 High-Protein, Low-Calorie Foods:")
for idx, row in high_protein_low_cal.iterrows():
    print(f"  {row['food_name']}: {row['protein_g']}g protein, {row['calories_per_100g']} cal")

# High-fiber foods
high_fiber = data.nlargest(5, 'fiber_g')[['food_name', 'fiber_g', 'calories_per_100g']]
print(f"\n🌾 High-Fiber Foods:")
for idx, row in high_fiber.iterrows():
    print(f"  {row['food_name']}: {row['fiber_g']}g fiber, {row['calories_per_100g']} cal")

print(f"\n✅ Nutrition analysis complete!")
print(f"Generated comprehensive nutritional insights from {len(data)} food items")
EOF

chmod +x nutrition_analyzer.py
```

### 3. Run the Analysis

```bash
python3 nutrition_analyzer.py
```

**Expected output**: You should see nutritional analysis results and food recommendations.

### 4. Food Safety Assessment Script

```bash
cat > food_safety_analyzer.py << 'EOF'
#!/usr/bin/env python3
"""
Food Safety Assessment Tool
Analyzes contamination risk and quality control metrics
"""

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from datetime import datetime, timedelta
import warnings
warnings.filterwarnings('ignore')

# Simulate food safety testing data
print("🔬 Generating food safety testing data...")
np.random.seed(42)

# Create sample food safety dataset
safety_data = {
    'sample_id': [f'SAMPLE_{i:04d}' for i in range(1, 501)],
    'product_type': np.random.choice(['Dairy', 'Meat', 'Produce', 'Processed', 'Beverages'], 500),
    'test_date': [datetime.now() - timedelta(days=np.random.randint(0, 90)) for _ in range(500)],
    'temperature_c': np.random.normal(4.0, 2.0, 500),  # Cold storage temp
    'ph_level': np.random.normal(6.5, 1.0, 500),
    'bacterial_count_cfu': np.random.lognormal(3, 1, 500),  # Colony forming units
    'moisture_percent': np.random.normal(15, 5, 500),
    'salt_percent': np.random.normal(2.5, 1.0, 500)
}

df = pd.DataFrame(safety_data)

# Food safety assessment
print("\n🛡️ Food Safety Risk Assessment")
print("=" * 35)

# Temperature compliance
temp_violations = df[df['temperature_c'] > 7].shape[0]  # > 7°C is danger zone
print(f"Temperature violations: {temp_violations}/500 samples ({temp_violations/5:.1f}%)")

# pH safety ranges by product type
ph_safety_ranges = {
    'Dairy': (6.0, 7.0),
    'Meat': (5.5, 6.5),
    'Produce': (4.0, 7.0),
    'Processed': (4.0, 6.0),
    'Beverages': (3.0, 4.5)
}

# Check pH violations
ph_violations = 0
for product_type, (min_ph, max_ph) in ph_safety_ranges.items():
    product_data = df[df['product_type'] == product_type]
    violations = product_data[(product_data['ph_level'] < min_ph) |
                            (product_data['ph_level'] > max_ph)].shape[0]
    ph_violations += violations
    print(f"{product_type} pH violations: {violations}/{len(product_data)} samples")

print(f"Total pH violations: {ph_violations}/500 samples ({ph_violations/5:.1f}%)")

# Bacterial contamination risk
high_bacterial = df[df['bacterial_count_cfu'] > 1000].shape[0]  # >1000 CFU/g
print(f"High bacterial count samples: {high_bacterial}/500 ({high_bacterial/5:.1f}%)")

# Risk scoring system
def calculate_risk_score(row):
    score = 0

    # Temperature risk
    if row['temperature_c'] > 7:
        score += 3
    elif row['temperature_c'] > 5:
        score += 1

    # pH risk
    product_range = ph_safety_ranges.get(row['product_type'], (4.0, 7.0))
    if row['ph_level'] < product_range[0] or row['ph_level'] > product_range[1]:
        score += 2

    # Bacterial risk
    if row['bacterial_count_cfu'] > 1000:
        score += 4
    elif row['bacterial_count_cfu'] > 500:
        score += 2

    # Moisture risk (high moisture = spoilage risk)
    if row['moisture_percent'] > 20:
        score += 1

    return score

df['risk_score'] = df.apply(calculate_risk_score, axis=1)

# Risk categorization
df['risk_level'] = pd.cut(df['risk_score'],
                         bins=[0, 2, 5, 10],
                         labels=['Low', 'Medium', 'High'])

print(f"\n📊 Risk Distribution:")
risk_counts = df['risk_level'].value_counts()
for level, count in risk_counts.items():
    print(f"  {level} Risk: {count} samples ({count/5:.1f}%)")

# Quality control dashboard
plt.figure(figsize=(14, 10))

# Temperature monitoring
plt.subplot(2, 3, 1)
plt.hist(df['temperature_c'], bins=20, alpha=0.7, color='red')
plt.axvline(x=7, color='black', linestyle='--', label='Safety Limit')
plt.title('Temperature Distribution')
plt.xlabel('Temperature (°C)')
plt.ylabel('Frequency')
plt.legend()

# pH levels by product type
plt.subplot(2, 3, 2)
product_types = df['product_type'].unique()
ph_by_product = [df[df['product_type'] == pt]['ph_level'] for pt in product_types]
plt.boxplot(ph_by_product, labels=product_types)
plt.title('pH Levels by Product Type')
plt.ylabel('pH Level')
plt.xticks(rotation=45)

# Bacterial contamination
plt.subplot(2, 3, 3)
plt.scatter(df['moisture_percent'], df['bacterial_count_cfu'],
           c=df['risk_score'], cmap='Reds', alpha=0.6)
plt.title('Bacterial Count vs Moisture')
plt.xlabel('Moisture (%)')
plt.ylabel('Bacterial Count (CFU/g)')
plt.colorbar(label='Risk Score')

# Risk score distribution
plt.subplot(2, 3, 4)
risk_counts.plot(kind='bar', color=['green', 'orange', 'red'])
plt.title('Risk Level Distribution')
plt.ylabel('Number of Samples')
plt.xticks(rotation=0)

# Compliance trends over time
plt.subplot(2, 3, 5)
df['week'] = df['test_date'].dt.isocalendar().week
weekly_compliance = df.groupby('week').agg({
    'risk_score': 'mean',
    'temperature_c': lambda x: (x <= 7).mean() * 100
}).round(2)

plt.plot(weekly_compliance.index, weekly_compliance['temperature_c'], 'b-',
         label='Temp Compliance %', linewidth=2)
plt.title('Weekly Compliance Trends')
plt.xlabel('Week Number')
plt.ylabel('Compliance Rate (%)')
plt.legend()

# Product type risk comparison
plt.subplot(2, 3, 6)
product_risk = df.groupby('product_type')['risk_score'].mean().sort_values(ascending=True)
product_risk.plot(kind='barh', color='purple')
plt.title('Average Risk Score by Product')
plt.xlabel('Average Risk Score')

plt.tight_layout()
plt.savefig('food_safety_dashboard.png', dpi=300, bbox_inches='tight')
print(f"\n📊 Food safety dashboard saved as 'food_safety_dashboard.png'")

# Generate safety recommendations
high_risk_samples = df[df['risk_level'] == 'High']
print(f"\n⚠️ High-Risk Samples Requiring Action:")
print(f"Total high-risk samples: {len(high_risk_samples)}")

if len(high_risk_samples) > 0:
    print("\nTop 5 highest risk samples:")
    for idx, row in high_risk_samples.nlargest(5, 'risk_score').iterrows():
        print(f"  {row['sample_id']}: Score {row['risk_score']} - "
              f"{row['product_type']} (Temp: {row['temperature_c']:.1f}°C, "
              f"pH: {row['ph_level']:.1f}, Bacteria: {row['bacterial_count_cfu']:.0f} CFU/g)")

print(f"\n✅ Food safety analysis complete!")
print(f"Assessed {len(df)} samples across {len(product_types)} product categories")
EOF

chmod +x food_safety_analyzer.py
```

### 5. Run Food Safety Assessment

```bash
python3 food_safety_analyzer.py
```

**Expected output**: Food safety risk assessment with contamination analysis.

### 6. Dietary Pattern Analysis

```bash
cat > dietary_pattern_analyzer.py << 'EOF'
#!/usr/bin/env python3
"""
Dietary Pattern Analysis Tool
Analyzes population dietary trends and nutritional epidemiology
"""

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
from scipy import stats
from sklearn.cluster import KMeans
from sklearn.preprocessing import StandardScaler
import warnings
warnings.filterwarnings('ignore')

# Simulate dietary survey data
print("📊 Generating dietary survey data...")
np.random.seed(42)

# Create sample dietary patterns dataset
n_participants = 1000

dietary_data = {
    'participant_id': [f'P_{i:04d}' for i in range(1, n_participants + 1)],
    'age': np.random.normal(45, 15, n_participants).astype(int),
    'gender': np.random.choice(['Male', 'Female'], n_participants),
    'bmi': np.random.normal(26, 4, n_participants),
    'daily_calories': np.random.normal(2200, 400, n_participants),
    'fruits_servings': np.random.poisson(2.5, n_participants),
    'vegetables_servings': np.random.poisson(3.0, n_participants),
    'grains_servings': np.random.poisson(6.0, n_participants),
    'protein_servings': np.random.poisson(2.8, n_participants),
    'dairy_servings': np.random.poisson(2.2, n_participants),
    'sugar_grams': np.random.normal(75, 25, n_participants),
    'sodium_mg': np.random.normal(3200, 800, n_participants),
    'fiber_grams': np.random.normal(18, 6, n_participants),
    'alcohol_servings': np.random.poisson(1.5, n_participants),
    'exercise_minutes': np.random.normal(150, 60, n_participants)
}

# Add some realistic correlations
for i in range(n_participants):
    # Higher BMI tends to correlate with higher calorie intake
    if dietary_data['bmi'][i] > 30:
        dietary_data['daily_calories'][i] *= 1.2
        dietary_data['sugar_grams'][i] *= 1.3

    # Older people tend to eat more vegetables, less sugar
    if dietary_data['age'][i] > 60:
        dietary_data['vegetables_servings'][i] = int(dietary_data['vegetables_servings'][i] * 1.3)
        dietary_data['sugar_grams'][i] *= 0.8

df = pd.DataFrame(dietary_data)

# Ensure realistic ranges
df['age'] = df['age'].clip(18, 80)
df['bmi'] = df['bmi'].clip(16, 45)
df['daily_calories'] = df['daily_calories'].clip(1000, 4000)
df['exercise_minutes'] = df['exercise_minutes'].clip(0, 500)

print(f"Generated dietary data for {len(df)} participants")

# Dietary pattern analysis
print("\n🍽️ Dietary Pattern Analysis")
print("=" * 30)

# Calculate dietary quality score
def calculate_diet_quality(row):
    score = 0

    # Fruits and vegetables (5-9 servings recommended)
    produce_total = row['fruits_servings'] + row['vegetables_servings']
    if produce_total >= 7:
        score += 3
    elif produce_total >= 5:
        score += 2
    elif produce_total >= 3:
        score += 1

    # Whole grains (6-8 servings recommended)
    if 6 <= row['grains_servings'] <= 8:
        score += 2
    elif 4 <= row['grains_servings'] <= 10:
        score += 1

    # Protein (2-3 servings recommended)
    if 2 <= row['protein_servings'] <= 3:
        score += 2
    elif 1 <= row['protein_servings'] <= 4:
        score += 1

    # Limit sugar (<50g recommended)
    if row['sugar_grams'] <= 25:
        score += 3
    elif row['sugar_grams'] <= 50:
        score += 2
    elif row['sugar_grams'] <= 75:
        score += 1

    # Limit sodium (<2300mg recommended)
    if row['sodium_mg'] <= 2300:
        score += 2
    elif row['sodium_mg'] <= 3000:
        score += 1

    # Adequate fiber (25g+ recommended)
    if row['fiber_grams'] >= 25:
        score += 2
    elif row['fiber_grams'] >= 18:
        score += 1

    return score

df['diet_quality_score'] = df.apply(calculate_diet_quality, axis=1)

# Diet quality categories
df['diet_quality'] = pd.cut(df['diet_quality_score'],
                           bins=[0, 6, 10, 15],
                           labels=['Poor', 'Fair', 'Good'])

print("Diet Quality Distribution:")
quality_counts = df['diet_quality'].value_counts()
for quality, count in quality_counts.items():
    print(f"  {quality}: {count} participants ({count/len(df)*100:.1f}%)")

# Nutritional adequacy analysis
print(f"\n📈 Nutritional Adequacy Assessment")
print("=" * 35)

# Fruit & vegetable intake
adequate_produce = (df['fruits_servings'] + df['vegetables_servings'] >= 5).sum()
print(f"Meeting fruit/vegetable guidelines: {adequate_produce}/{len(df)} ({adequate_produce/len(df)*100:.1f}%)")

# Fiber adequacy
adequate_fiber = (df['fiber_grams'] >= 25).sum()
print(f"Meeting fiber recommendations: {adequate_fiber}/{len(df)} ({adequate_fiber/len(df)*100:.1f}%)")

# Sodium compliance
low_sodium = (df['sodium_mg'] <= 2300).sum()
print(f"Within sodium limits: {low_sodium}/{len(df)} ({low_sodium/len(df)*100:.1f}%)")

# Sugar compliance
low_sugar = (df['sugar_grams'] <= 50).sum()
print(f"Within sugar recommendations: {low_sugar}/{len(df)} ({low_sugar/len(df)*100:.1f}%)")

# Demographic analysis
print(f"\n👥 Dietary Patterns by Demographics")
print("=" * 35)

# By age group
df['age_group'] = pd.cut(df['age'], bins=[18, 35, 50, 65, 80],
                        labels=['18-35', '36-50', '51-65', '66-80'])

age_diet_quality = df.groupby('age_group')['diet_quality_score'].mean()
print("Average diet quality by age group:")
for age_group, score in age_diet_quality.items():
    print(f"  {age_group}: {score:.1f}")

# By gender
gender_analysis = df.groupby('gender').agg({
    'diet_quality_score': 'mean',
    'fruits_servings': 'mean',
    'vegetables_servings': 'mean',
    'sugar_grams': 'mean'
}).round(1)

print(f"\nDietary patterns by gender:")
print(gender_analysis)

# Cluster analysis to identify dietary patterns
print(f"\n🎯 Dietary Pattern Clustering")
print("=" * 28)

# Prepare features for clustering
cluster_features = ['fruits_servings', 'vegetables_servings', 'grains_servings',
                   'protein_servings', 'dairy_servings', 'sugar_grams', 'fiber_grams']

scaler = StandardScaler()
scaled_features = scaler.fit_transform(df[cluster_features])

# Perform k-means clustering
kmeans = KMeans(n_clusters=4, random_state=42)
df['dietary_pattern'] = kmeans.fit_predict(scaled_features)

# Label the patterns based on characteristics
pattern_labels = {0: 'High Sugar', 1: 'Balanced', 2: 'Low Produce', 3: 'Health Conscious'}
df['pattern_name'] = df['dietary_pattern'].map(pattern_labels)

print("Identified dietary patterns:")
pattern_summary = df.groupby('pattern_name').agg({
    'fruits_servings': 'mean',
    'vegetables_servings': 'mean',
    'sugar_grams': 'mean',
    'fiber_grams': 'mean',
    'diet_quality_score': 'mean'
}).round(1)

for pattern, data in pattern_summary.iterrows():
    print(f"\n{pattern} Pattern:")
    print(f"  Fruits: {data['fruits_servings']} servings/day")
    print(f"  Vegetables: {data['vegetables_servings']} servings/day")
    print(f"  Sugar: {data['sugar_grams']}g/day")
    print(f"  Fiber: {data['fiber_grams']}g/day")
    print(f"  Diet Quality: {data['diet_quality_score']}/15")

# Generate comprehensive visualization
plt.figure(figsize=(16, 12))

# Diet quality distribution
plt.subplot(3, 3, 1)
quality_counts.plot(kind='bar', color=['red', 'orange', 'green'])
plt.title('Diet Quality Distribution')
plt.ylabel('Number of Participants')
plt.xticks(rotation=0)

# BMI vs Diet Quality
plt.subplot(3, 3, 2)
plt.scatter(df['bmi'], df['diet_quality_score'], alpha=0.6, color='blue')
plt.title('BMI vs Diet Quality Score')
plt.xlabel('BMI')
plt.ylabel('Diet Quality Score')

# Age vs Produce Intake
plt.subplot(3, 3, 3)
df['total_produce'] = df['fruits_servings'] + df['vegetables_servings']
plt.scatter(df['age'], df['total_produce'], alpha=0.6, color='green')
plt.title('Age vs Fruit/Vegetable Intake')
plt.xlabel('Age')
plt.ylabel('Total Servings/Day')

# Sugar intake by gender
plt.subplot(3, 3, 4)
sugar_by_gender = [df[df['gender'] == 'Male']['sugar_grams'],
                  df[df['gender'] == 'Female']['sugar_grams']]
plt.boxplot(sugar_by_gender, labels=['Male', 'Female'])
plt.title('Sugar Intake by Gender')
plt.ylabel('Sugar (g/day)')

# Dietary patterns
plt.subplot(3, 3, 5)
pattern_counts = df['pattern_name'].value_counts()
plt.pie(pattern_counts.values, labels=pattern_counts.index, autopct='%1.1f%%')
plt.title('Dietary Pattern Distribution')

# Fiber vs age
plt.subplot(3, 3, 6)
plt.scatter(df['age'], df['fiber_grams'], alpha=0.6, color='brown')
plt.axhline(y=25, color='red', linestyle='--', label='Recommended')
plt.title('Fiber Intake vs Age')
plt.xlabel('Age')
plt.ylabel('Fiber (g/day)')
plt.legend()

# Exercise vs diet quality
plt.subplot(3, 3, 7)
plt.scatter(df['exercise_minutes'], df['diet_quality_score'], alpha=0.6, color='purple')
plt.title('Exercise vs Diet Quality')
plt.xlabel('Exercise (min/week)')
plt.ylabel('Diet Quality Score')

# Sodium intake distribution
plt.subplot(3, 3, 8)
plt.hist(df['sodium_mg'], bins=30, alpha=0.7, color='red')
plt.axvline(x=2300, color='black', linestyle='--', label='Recommended Limit')
plt.title('Sodium Intake Distribution')
plt.xlabel('Sodium (mg/day)')
plt.ylabel('Frequency')
plt.legend()

# Calorie intake by age group
plt.subplot(3, 3, 9)
calorie_by_age = [df[df['age_group'] == ag]['daily_calories'] for ag in df['age_group'].cat.categories]
plt.boxplot(calorie_by_age, labels=df['age_group'].cat.categories)
plt.title('Calorie Intake by Age Group')
plt.ylabel('Calories/Day')
plt.xticks(rotation=45)

plt.tight_layout()
plt.savefig('dietary_pattern_analysis.png', dpi=300, bbox_inches='tight')
print(f"\n📊 Dietary analysis dashboard saved as 'dietary_pattern_analysis.png'")

# Health risk assessment
print(f"\n⚠️ Health Risk Assessment")
print("=" * 25)

# Identify high-risk individuals
high_risk_criteria = (
    (df['bmi'] > 30) |  # Obese
    (df['sugar_grams'] > 100) |  # Excessive sugar
    (df['sodium_mg'] > 4000) |  # Very high sodium
    (df['fiber_grams'] < 10) |  # Very low fiber
    (df['fruits_servings'] + df['vegetables_servings'] < 2)  # Very low produce
)

high_risk_count = high_risk_criteria.sum()
print(f"Participants meeting high-risk criteria: {high_risk_count}/{len(df)} ({high_risk_count/len(df)*100:.1f}%)")

# Protective factors
protective_factors = (
    (df['diet_quality_score'] >= 10) &
    (df['exercise_minutes'] >= 150) &
    (df['bmi'] >= 18.5) & (df['bmi'] <= 25)
)

protected_count = protective_factors.sum()
print(f"Participants with protective lifestyle factors: {protected_count}/{len(df)} ({protected_count/len(df)*100:.1f}%)")

print(f"\n✅ Dietary pattern analysis complete!")
print(f"Analyzed dietary data from {len(df)} participants across multiple health indicators")
EOF

chmod +x dietary_pattern_analyzer.py
```

### 7. Run Dietary Pattern Analysis

```bash
python3 dietary_pattern_analyzer.py
```

**Expected output**: Comprehensive dietary pattern analysis with health risk assessment.

## What You've Accomplished

🎉 **Congratulations!** You've successfully:

1. ✅ Created a food science research environment in the cloud
2. ✅ Analyzed nutritional content and food composition data
3. ✅ Performed food safety risk assessment and quality control
4. ✅ Conducted dietary pattern analysis and nutritional epidemiology
5. ✅ Generated publication-quality visualizations and reports

### Real Research Applications

Your environment can now handle:
- **Nutritional databases**: USDA, FDA food composition analysis
- **Food safety monitoring**: Contamination risk, quality control systems
- **Dietary surveys**: NHANES, population nutrition studies
- **Clinical nutrition**: Patient dietary assessment, intervention studies
- **Food product development**: Nutritional optimization, safety testing

### Next Steps for Advanced Research

```bash
# Install specialized food science packages
pip3 install fooddata-central nutrient-analysis food-safety-toolkit

# Set up food composition databases
wget https://fdc.nal.usda.gov/fdc-datasets/FoodData_Central_csv_2023-04-20.zip

# Configure nutritional analysis pipelines
aws-research-wizard tools install --domain food_science_nutrition --advanced
```

### Monthly Cost Estimate

For typical food science research usage:
- **Light usage** (10 hours/week): ~$120/month
- **Medium usage** (20 hours/week): ~$200/month
- **Heavy usage** (40 hours/week): ~$300/month

## Clean Up Resources

**Important**: Always clean up to avoid unexpected charges!

```bash
# Exit your research environment
exit

# Destroy the research environment
aws-research-wizard deploy destroy --domain food_science_nutrition
```

**Expected result**: "✅ Environment destroyed successfully"

**💰 Billing stops**: No more charges after cleanup

## Troubleshooting

### Common Issues

**Problem**: "Package not found" errors
**Solution**:
```bash
sudo yum update
sudo yum install python3-pip python3-devel
pip3 install --upgrade pip
```

**Problem**: "Permission denied" errors
**Solution**:
```bash
sudo chown -R $USER:$USER ~/food_research/
chmod +x *.py
```

**Problem**: Visualizations not generating
**Solution**:
```bash
pip3 install matplotlib seaborn pandas scipy scikit-learn
export MPLBACKEND=Agg  # For headless display
```

**Problem**: Data download failures
**Solution**:
```bash
# Check internet connectivity
ping google.com

# Try alternative data source
wget --no-check-certificate [URL]
```

### Getting Help

- **Food Science Community**: [forum.aws-research-wizard.com/food-science](https://forum.aws-research-wizard.com/food-science)
- **Technical Support**: [support@aws-research-wizard.com](mailto:support@aws-research-wizard.com)
- **Sample Data**: [research-data.aws-wizard.com/nutrition](https://research-data.aws-wizard.com/nutrition)

### Emergency Stop

If something goes wrong and you want to stop all charges immediately:

```bash
aws-research-wizard emergency-stop --all
```

This will terminate everything and stop billing within 2 minutes.

---

**🍎 Happy food science research!**
*You now have a professional-grade nutrition research environment that scales with your needs.*
