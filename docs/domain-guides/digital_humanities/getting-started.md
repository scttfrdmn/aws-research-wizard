# Digital Humanities Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $7-12 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working digital humanities research environment that can:
- Analyze large collections of text and historical documents
- Perform natural language processing and sentiment analysis
- Create text visualizations and topic modeling
- Handle datasets with millions of documents

### Meet Dr. Maria Santos
Dr. Maria Santos is a digital historian at Columbia University. She analyzes 19th-century newspaper archives but waits weeks for university computing resources. Each text analysis project takes months to complete manually.

**Before**: 3-week waits + 2-month analysis = 2.7 months per project
**After**: 15-minute setup + 4-hour analysis = same day results
**Time Saved**: 99% faster research cycle
**Cost Savings**: $300/month vs $1,200 research budget

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $7-12 (we'll clean up resources when done)
- **Daily research cost**: $12-35 per day when actively analyzing
- **Monthly estimate**: $150-400 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No cloud or programming experience required

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
- **Region**: Choose `us-west-2` (recommended for digital humanities with good text processing performance)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain digital_humanities --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: digital_humanities
✅ Region valid: us-west-2 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Digital Humanities Environment

```bash
aws-research-wizard deploy start --domain digital_humanities --region us-west-2 --instance t3.large
```

**What this does**: Creates your digital humanities computing environment optimized for text analysis.

**This will take**: 3-5 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
  CPU: 2 cores optimized for text processing
  Memory: 8GB RAM for document analysis
```

**💰 Billing starts now**: Your environment costs about $0.17 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
```

**What this does**: Connects you to your digital humanities computer in the cloud.

**Expected result**: You see a command prompt like `ubuntu@ip-10-0-1-123:~$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Digital Humanities Tools

Your environment comes pre-installed with:

### Core Text Analysis Tools
- **spaCy**: Advanced NLP library - Type `python -c "import spacy; print(spacy.__version__)"` to check
- **NLTK**: Natural Language Toolkit - Type `python -c "import nltk; print(nltk.__version__)"` to check
- **Gensim**: Topic modeling - Type `python -c "import gensim; print(gensim.__version__)"` to check
- **Pandas**: Data manipulation - Type `python -c "import pandas; print(pandas.__version__)"` to check
- **Jupyter**: Interactive notebooks - Type `jupyter --version` to check

### Try Your First Command
```bash
python -c "import spacy; print('spaCy version:', spacy.__version__)"
```

**What this does**: Shows spaCy version and confirms natural language processing tools are installed.

**Expected result**: You see spaCy version info confirming NLP libraries are ready.

## Step 8: Analyze Historical Text Data

Let's analyze some historical documents to test everything works:

### Download Sample Historical Texts
```bash
# Create working directory
mkdir ~/humanities-tutorial
cd ~/humanities-tutorial

# Download sample historical documents (Project Gutenberg texts)
wget -O shakespeare_hamlet.txt "https://www.gutenberg.org/files/1524/1524-0.txt"
wget -O dickens_tale_of_two_cities.txt "https://www.gutenberg.org/files/98/98-0.txt"
wget -O austen_pride_prejudice.txt "https://www.gutenberg.org/files/1342/1342-0.txt"

echo "Sample historical texts downloaded successfully!"
```

### Basic Text Analysis
```bash
# Create text analysis script
cat > text_analysis.py << 'EOF'
import nltk
import spacy
from collections import Counter
import matplotlib.pyplot as plt
import pandas as pd

# Download required NLTK data
print("Downloading NLTK data...")
nltk.download('punkt', quiet=True)
nltk.download('stopwords', quiet=True)
nltk.download('vader_lexicon', quiet=True)

# Load spaCy model
print("Loading spaCy English model...")
nlp = spacy.load('en_core_web_sm')

def analyze_text(filename, title):
    print(f"\n=== Analyzing {title} ===")

    # Read text file
    with open(filename, 'r', encoding='utf-8') as f:
        text = f.read()

    # Basic statistics
    words = nltk.word_tokenize(text.lower())
    sentences = nltk.sent_tokenize(text)

    print(f"Total characters: {len(text):,}")
    print(f"Total words: {len(words):,}")
    print(f"Total sentences: {len(sentences):,}")
    print(f"Average words per sentence: {len(words)/len(sentences):.1f}")

    # Remove stopwords and analyze vocabulary
    from nltk.corpus import stopwords
    stop_words = set(stopwords.words('english'))
    content_words = [word for word in words if word.isalpha() and word not in stop_words]

    print(f"Content words (no stopwords): {len(content_words):,}")
    print(f"Unique vocabulary: {len(set(content_words)):,}")
    print(f"Vocabulary richness: {len(set(content_words))/len(content_words):.3f}")

    # Most common words
    word_freq = Counter(content_words)
    print(f"Most common words: {word_freq.most_common(10)}")

    # Named entity recognition
    doc = nlp(text[:100000])  # Process first 100k characters
    entities = [(ent.text, ent.label_) for ent in doc.ents]
    entity_types = Counter([ent[1] for ent in entities])

    print(f"Named entities found: {len(entities)}")
    print(f"Entity types: {dict(entity_types.most_common(5))}")

    return {
        'title': title,
        'characters': len(text),
        'words': len(words),
        'sentences': len(sentences),
        'vocabulary': len(set(content_words)),
        'top_words': word_freq.most_common(5)
    }

# Analyze each text
results = []
texts = [
    ('shakespeare_hamlet.txt', 'Shakespeare - Hamlet'),
    ('dickens_tale_of_two_cities.txt', 'Dickens - A Tale of Two Cities'),
    ('austen_pride_prejudice.txt', 'Austen - Pride and Prejudice')
]

for filename, title in texts:
    try:
        result = analyze_text(filename, title)
        results.append(result)
    except FileNotFoundError:
        print(f"File {filename} not found - skipping")

# Create comparison
if results:
    df = pd.DataFrame(results)
    print("\n=== Comparative Analysis ===")
    print(df[['title', 'words', 'sentences', 'vocabulary']])

print("\n✅ Historical text analysis completed!")
EOF

python3 text_analysis.py
```

**What this does**: Analyzes vocabulary, sentence structure, and named entities in historical literature.

**This will take**: 2-3 minutes

### Sentiment Analysis
```bash
# Create sentiment analysis script
cat > sentiment_analysis.py << 'EOF'
import nltk
from nltk.sentiment import SentimentIntensityAnalyzer
import matplotlib.pyplot as plt
import numpy as np

print("Performing sentiment analysis on historical texts...")

# Initialize sentiment analyzer
sia = SentimentIntensityAnalyzer()

def analyze_sentiment(filename, title):
    print(f"\n=== Sentiment Analysis: {title} ===")

    try:
        with open(filename, 'r', encoding='utf-8') as f:
            text = f.read()

        # Split into chapters or sections (every 5000 characters)
        sections = [text[i:i+5000] for i in range(0, len(text), 5000)][:20]  # First 20 sections

        sentiments = []
        for i, section in enumerate(sections):
            # Clean section text
            clean_section = ' '.join(section.split()[:500])  # First 500 words

            # Get sentiment scores
            scores = sia.polarity_scores(clean_section)
            sentiments.append({
                'section': i+1,
                'positive': scores['pos'],
                'negative': scores['neg'],
                'neutral': scores['neu'],
                'compound': scores['compound']
            })

        # Calculate overall sentiment
        avg_compound = np.mean([s['compound'] for s in sentiments])
        avg_positive = np.mean([s['positive'] for s in sentiments])
        avg_negative = np.mean([s['negative'] for s in sentiments])

        print(f"Overall sentiment (compound): {avg_compound:.3f}")
        print(f"Average positive sentiment: {avg_positive:.3f}")
        print(f"Average negative sentiment: {avg_negative:.3f}")

        # Classify overall tone
        if avg_compound >= 0.05:
            tone = "Positive"
        elif avg_compound <= -0.05:
            tone = "Negative"
        else:
            tone = "Neutral"

        print(f"Overall tone: {tone}")

        return {
            'title': title,
            'compound': avg_compound,
            'positive': avg_positive,
            'negative': avg_negative,
            'tone': tone,
            'sections': len(sentiments)
        }

    except FileNotFoundError:
        print(f"File {filename} not found")
        return None

# Analyze sentiment for each text
texts = [
    ('shakespeare_hamlet.txt', 'Hamlet'),
    ('dickens_tale_of_two_cities.txt', 'Tale of Two Cities'),
    ('austen_pride_prejudice.txt', 'Pride and Prejudice')
]

sentiment_results = []
for filename, title in texts:
    result = analyze_sentiment(filename, title)
    if result:
        sentiment_results.append(result)

# Summary comparison
if sentiment_results:
    print("\n=== Sentiment Comparison ===")
    for result in sentiment_results:
        print(f"{result['title']}: {result['tone']} (compound: {result['compound']:.3f})")

print("\n✅ Sentiment analysis completed!")
EOF

python3 sentiment_analysis.py
```

**What this does**: Analyzes emotional tone and sentiment patterns across different historical literary works.

**Expected result**: Shows sentiment scores and comparative emotional analysis of the texts.

**🎉 Success!** You've analyzed real historical documents in the cloud.

## Step 9: Topic Modeling

Test advanced digital humanities capabilities:

```bash
# Create topic modeling script
cat > topic_modeling.py << 'EOF'
import gensim
from gensim import corpora
import nltk
from nltk.corpus import stopwords
import re

print("Performing topic modeling on historical texts...")

def preprocess_text(filename):
    """Clean and prepare text for topic modeling"""
    with open(filename, 'r', encoding='utf-8') as f:
        text = f.read()

    # Remove Project Gutenberg header/footer
    start_marker = "*** START OF"
    end_marker = "*** END OF"

    if start_marker in text:
        text = text.split(start_marker)[1] if start_marker in text else text
    if end_marker in text:
        text = text.split(end_marker)[0] if end_marker in text else text

    # Clean text
    text = re.sub(r'[^a-zA-Z\s]', '', text)  # Remove punctuation
    text = text.lower()

    # Tokenize and remove stopwords
    words = nltk.word_tokenize(text)
    stop_words = set(stopwords.words('english'))

    # Add common words that aren't meaningful for topic modeling
    stop_words.update(['said', 'say', 'come', 'came', 'go', 'went', 'one', 'two', 'would', 'could'])

    # Filter words
    filtered_words = [word for word in words if len(word) > 3 and word not in stop_words]

    return filtered_words

# Process all texts
texts = [
    ('shakespeare_hamlet.txt', 'Hamlet'),
    ('dickens_tale_of_two_cities.txt', 'Tale of Two Cities'),
    ('austen_pride_prejudice.txt', 'Pride and Prejudice')
]

processed_texts = []
for filename, title in texts:
    try:
        words = preprocess_text(filename)
        if words:
            # Split into chunks for better topic analysis
            chunk_size = 1000
            chunks = [words[i:i+chunk_size] for i in range(0, len(words), chunk_size)]
            processed_texts.extend(chunks[:5])  # Take first 5 chunks per text
            print(f"Processed {title}: {len(words)} words in {len(chunks)} chunks")
    except FileNotFoundError:
        print(f"File {filename} not found - skipping")

if processed_texts:
    print(f"\nTotal documents for topic modeling: {len(processed_texts)}")

    # Create dictionary and corpus
    dictionary = corpora.Dictionary(processed_texts)
    dictionary.filter_extremes(no_below=2, no_above=0.7)  # Remove rare and common words
    corpus = [dictionary.doc2bow(text) for text in processed_texts]

    print(f"Dictionary size: {len(dictionary)}")
    print(f"Corpus size: {len(corpus)}")

    # Train LDA model
    print("Training topic model...")
    lda_model = gensim.models.LdaModel(
        corpus=corpus,
        id2word=dictionary,
        num_topics=5,
        random_state=42,
        passes=10,
        alpha='auto',
        per_word_topics=True
    )

    # Show topics
    print("\n=== Discovered Topics ===")
    for idx, topic in lda_model.print_topics():
        print(f"Topic {idx}: {topic}")

    # Model coherence
    coherence_model = gensim.models.CoherenceModel(
        model=lda_model, texts=processed_texts, dictionary=dictionary, coherence='c_v'
    )
    coherence_score = coherence_model.get_coherence()
    print(f"\nModel coherence score: {coherence_score:.3f}")

else:
    print("No texts available for topic modeling")

print("\n✅ Topic modeling completed!")
EOF

python3 topic_modeling.py
```

**What this does**: Discovers hidden themes and topics across the collection of historical texts.

**Expected result**: Shows discovered topics and their associated keywords.

## Step 10: Monitor Your Costs

Check your current spending:

```bash
exit  # Exit SSH session first
aws-research-wizard monitor costs --region us-west-2
```

**Expected result**: Shows costs so far (should be under $3 for this tutorial)

## Step 11: Clean Up (Important!)

When you're done experimenting:

```bash
aws-research-wizard deploy delete --region us-west-2
```

Type `y` when prompted.

**What this does**: Stops billing by removing your cloud resources.

**💰 Important**: Always clean up to avoid ongoing charges.

**Expected result**: "🗑️ Deletion completed successfully"

## Understanding Your Costs

### What You're Paying For
- **Compute**: $0.17 per hour for text processing instance while environment is running
- **Storage**: $0.10 per GB per month for document collections you save
- **Data Transfer**: Usually free for digital humanities data amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 60% savings (advanced)
- Store large document collections in S3, not on the instance
- Process texts efficiently to minimize compute time

### Typical Monthly Costs by Usage
- **Light use** (10 hours/week): $30-75
- **Medium use** (3 hours/day): $75-150
- **Heavy use** (6 hours/day): $150-300

## What's Next?

Now that you have a working digital humanities environment, you can:

### Learn More About Digital Text Analysis
- [Large Corpus Analysis Tutorial](large-corpus-analysis.md)
- [Advanced Topic Modeling Guide](advanced-topic-modeling.md)
- [Cost Optimization for Digital Humanities](humanities-cost-optimization.md)

### Explore Advanced Features
- [Network analysis of historical documents](network-analysis-humanities.md)
- [Team collaboration with text databases](team-humanities-collaboration.md)
- [Automated text processing pipelines](automated-humanities-pipelines.md)

### Join the Digital Humanities Community
- [Digital Humanities Forum](https://forum.researchwizard.app/humanities)
- [GitHub Humanities Examples](https://github.com/aws-research-wizard/humanities-examples)
- [Monthly Digital Humanities Office Hours](https://calendar.researchwizard.app/humanities-office-hours)

## Troubleshooting

### Common Issues

**Problem**: "spaCy model not found" error during analysis
**Solution**: Download language model: `python -m spacy download en_core_web_sm`
**Prevention**: Wait 3-5 minutes after deployment for all NLP models to initialize

**Problem**: "NLTK data not found" error
**Solution**: Download required data: `python -c "import nltk; nltk.download('all')"`
**Prevention**: The analysis scripts automatically download required data

**Problem**: "Memory error" during large text processing
**Solution**: Process texts in smaller chunks or use a larger instance type
**Prevention**: Monitor memory usage with `htop` during analysis

**Problem**: "Unicode decode error" when reading text files
**Solution**: Try different encoding: `open(filename, 'r', encoding='latin-1')`
**Prevention**: Check file encoding with `file -bi filename.txt` before processing

### Getting Help
- Check the [digital humanities troubleshooting guide](troubleshooting-humanities.md)
- Ask in [community forum](https://forum.researchwizard.app)
- File an issue on [GitHub](https://github.com/aws-research-wizard/aws-research-wizard/issues)

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-west-2 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $12. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
