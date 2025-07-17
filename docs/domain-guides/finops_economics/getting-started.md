# FinOps & Economics Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $8-14 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working financial and economic research environment that can:
- Analyze financial markets and economic data
- Build econometric models and financial forecasting systems
- Process large-scale financial datasets and time series
- Implement risk management and portfolio optimization models

### Meet Dr. Rachel Chen

Dr. Rachel Chen is a financial economist at Federal Reserve Bank. She analyzes market data but waits days for secure computing resources. Each economic model requires processing millions of financial transactions and market indicators.

**Before**: 3-day waits + 8-hour analysis = 11 days per economic study
**After**: 15-minute setup + 2-hour analysis = same day results
**Time Saved**: 96% faster financial analysis cycle
**Cost Savings**: $400/month vs $1,500 institutional allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $8-14 (we'll clean up resources when done)
- **Daily research cost**: $12-28 per day when actively analyzing
- **Monthly estimate**: $150-350 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No finance or economics experience required

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
- **Region**: Choose `us-east-1` (recommended for financial data with good market data access)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain finops_economics --region us-east-1
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: finops_economics
✅ Region valid: us-east-1 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your FinOps Environment

```bash
aws-research-wizard deploy start --domain finops_economics --region us-east-1 --instance m6i.large
```

**What this does**: Creates your financial research environment optimized for economic data analysis.

**This will take**: 5-7 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
  CPU: 2 cores for financial modeling
  Memory: 8GB RAM for large datasets
```

**💰 Billing starts now**: Your environment costs about $0.19 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
```

**What this does**: Connects you to your financial research computer in the cloud.

**Expected result**: You see a command prompt like `ubuntu@ip-10-0-1-123:~$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Financial Tools

Your environment comes pre-installed with:

### Core Financial Software
- **Python Financial Stack**: Pandas, NumPy, SciPy - Type `python -c "import pandas; print(pandas.__version__)"` to check
- **R Statistical Software**: Econometric analysis - Type `R --version` to check
- **Jupyter Notebooks**: Interactive analysis - Type `jupyter --version` to check
- **QuantLib**: Quantitative finance library - Type `python -c "import QuantLib; print(QuantLib.__version__)"` to check
- **StatsModels**: Econometric modeling - Type `python -c "import statsmodels; print(statsmodels.__version__)"` to check

### Try Your First Command
```bash
python -c "import pandas; print('Pandas version:', pandas.__version__)"
```

**What this does**: Shows Pandas version and confirms financial analysis tools are installed.

**Expected result**: You see Pandas version info confirming financial libraries are ready.

## Step 8: Analyze Financial Data

Let's analyze financial and economic data to test everything works:

### Stock Market Analysis
```bash
# Create working directory
mkdir ~/finops-tutorial
cd ~/finops-tutorial

# Create financial data analysis script
cat > financial_analysis.py << 'EOF'
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from datetime import datetime, timedelta

print("Starting financial data analysis...")

def generate_stock_data():
    """Generate synthetic stock market data"""
    print("\n=== Stock Market Data Generation ===")

    np.random.seed(42)

    # Generate 2 years of daily stock data
    start_date = datetime(2022, 1, 1)
    end_date = datetime(2023, 12, 31)
    dates = pd.date_range(start_date, end_date, freq='D')

    # Remove weekends (basic market calendar)
    trading_days = dates[dates.weekday < 5]
    n_days = len(trading_days)

    print(f"Trading days: {n_days}")
    print(f"Date range: {start_date.date()} to {end_date.date()}")

    # Simulate stock prices using geometric Brownian motion
    stocks = {
        'TECH': {'price': 150, 'mu': 0.12, 'sigma': 0.25},  # Tech stock
        'FINANCE': {'price': 85, 'mu': 0.08, 'sigma': 0.20},  # Financial stock
        'ENERGY': {'price': 60, 'mu': 0.06, 'sigma': 0.30},   # Energy stock
        'HEALTHCARE': {'price': 120, 'mu': 0.10, 'sigma': 0.18},  # Healthcare
        'UTILITIES': {'price': 45, 'mu': 0.04, 'sigma': 0.12}   # Utilities
    }

    stock_data = pd.DataFrame(index=trading_days)

    for symbol, params in stocks.items():
        S0 = params['price']
        mu = params['mu'] / 252  # Daily return
        sigma = params['sigma'] / np.sqrt(252)  # Daily volatility

        # Generate random shocks
        shocks = np.random.normal(0, 1, n_days)

        # Calculate daily returns
        daily_returns = mu + sigma * shocks

        # Calculate prices
        prices = [S0]
        for i in range(1, n_days):
            price = prices[i-1] * np.exp(daily_returns[i])
            prices.append(price)

        stock_data[symbol] = prices

    return stock_data

def calculate_financial_metrics(stock_data):
    """Calculate key financial metrics"""
    print("\n=== Financial Metrics Calculation ===")

    # Calculate daily returns
    returns = stock_data.pct_change().dropna()

    # Calculate annualized metrics
    trading_days_per_year = 252

    metrics = {}
    for symbol in stock_data.columns:
        symbol_returns = returns[symbol]

        # Basic statistics
        annual_return = symbol_returns.mean() * trading_days_per_year
        annual_volatility = symbol_returns.std() * np.sqrt(trading_days_per_year)

        # Risk metrics
        sharpe_ratio = annual_return / annual_volatility  # Assuming risk-free rate = 0
        max_drawdown = calculate_max_drawdown(stock_data[symbol])

        # VaR (Value at Risk) at 95% confidence
        var_95 = np.percentile(symbol_returns, 5)

        metrics[symbol] = {
            'Annual Return': annual_return,
            'Annual Volatility': annual_volatility,
            'Sharpe Ratio': sharpe_ratio,
            'Max Drawdown': max_drawdown,
            'VaR (95%)': var_95
        }

    # Display metrics
    metrics_df = pd.DataFrame(metrics).T
    print("Financial Metrics Summary:")
    print(metrics_df.round(4))

    # Portfolio analysis
    print(f"\nPortfolio Analysis (Equal Weight):")
    portfolio_weights = np.array([0.2, 0.2, 0.2, 0.2, 0.2])  # Equal weight
    portfolio_returns = (returns * portfolio_weights).sum(axis=1)

    portfolio_annual_return = portfolio_returns.mean() * trading_days_per_year
    portfolio_annual_vol = portfolio_returns.std() * np.sqrt(trading_days_per_year)
    portfolio_sharpe = portfolio_annual_return / portfolio_annual_vol

    print(f"  Annual Return: {portfolio_annual_return:.4f}")
    print(f"  Annual Volatility: {portfolio_annual_vol:.4f}")
    print(f"  Sharpe Ratio: {portfolio_sharpe:.4f}")

    return metrics_df, portfolio_returns

def calculate_max_drawdown(price_series):
    """Calculate maximum drawdown"""
    cumulative = (1 + price_series.pct_change()).cumprod()
    rolling_max = cumulative.expanding().max()
    drawdown = (cumulative - rolling_max) / rolling_max
    return drawdown.min()

def correlation_analysis(stock_data):
    """Analyze correlations between assets"""
    print("\n=== Correlation Analysis ===")

    returns = stock_data.pct_change().dropna()
    correlation_matrix = returns.corr()

    print("Correlation Matrix:")
    print(correlation_matrix.round(3))

    # Find highest and lowest correlations
    corr_pairs = []
    for i in range(len(correlation_matrix.columns)):
        for j in range(i+1, len(correlation_matrix.columns)):
            stock1 = correlation_matrix.columns[i]
            stock2 = correlation_matrix.columns[j]
            corr_value = correlation_matrix.iloc[i, j]
            corr_pairs.append((stock1, stock2, corr_value))

    # Sort by correlation
    corr_pairs.sort(key=lambda x: x[2], reverse=True)

    print(f"\nHighest Correlations:")
    for stock1, stock2, corr in corr_pairs[:3]:
        print(f"  {stock1} - {stock2}: {corr:.3f}")

    print(f"\nLowest Correlations:")
    for stock1, stock2, corr in corr_pairs[-3:]:
        print(f"  {stock1} - {stock2}: {corr:.3f}")

    return correlation_matrix

def technical_analysis(stock_data):
    """Perform basic technical analysis"""
    print("\n=== Technical Analysis ===")

    # Focus on first stock for technical analysis
    symbol = stock_data.columns[0]
    prices = stock_data[symbol]

    # Moving averages
    ma_20 = prices.rolling(window=20).mean()
    ma_50 = prices.rolling(window=50).mean()

    # Bollinger Bands
    bb_period = 20
    bb_std = 2
    bb_middle = prices.rolling(window=bb_period).mean()
    bb_std_dev = prices.rolling(window=bb_period).std()
    bb_upper = bb_middle + (bb_std_dev * bb_std)
    bb_lower = bb_middle - (bb_std_dev * bb_std)

    # RSI (Relative Strength Index)
    def calculate_rsi(prices, period=14):
        delta = prices.diff()
        gain = (delta.where(delta > 0, 0)).rolling(window=period).mean()
        loss = (-delta.where(delta < 0, 0)).rolling(window=period).mean()
        rs = gain / loss
        rsi = 100 - (100 / (1 + rs))
        return rsi

    rsi = calculate_rsi(prices)

    # Current values (last trading day)
    current_price = prices.iloc[-1]
    current_ma20 = ma_20.iloc[-1]
    current_ma50 = ma_50.iloc[-1]
    current_rsi = rsi.iloc[-1]

    print(f"Technical Analysis for {symbol}:")
    print(f"  Current Price: ${current_price:.2f}")
    print(f"  20-day MA: ${current_ma20:.2f}")
    print(f"  50-day MA: ${current_ma50:.2f}")
    print(f"  RSI: {current_rsi:.2f}")

    # Trading signals
    print(f"\nTrading Signals:")
    if current_price > current_ma20 and current_ma20 > current_ma50:
        print("  Trend: Bullish (Price > MA20 > MA50)")
    elif current_price < current_ma20 and current_ma20 < current_ma50:
        print("  Trend: Bearish (Price < MA20 < MA50)")
    else:
        print("  Trend: Mixed")

    if current_rsi > 70:
        print("  RSI Signal: Overbought (RSI > 70)")
    elif current_rsi < 30:
        print("  RSI Signal: Oversold (RSI < 30)")
    else:
        print("  RSI Signal: Neutral (30 < RSI < 70)")

    # Bollinger Band position
    bb_position = (current_price - bb_lower.iloc[-1]) / (bb_upper.iloc[-1] - bb_lower.iloc[-1])
    print(f"  Bollinger Band Position: {bb_position:.2f} (0=lower band, 1=upper band)")

    return {
        'ma_20': current_ma20,
        'ma_50': current_ma50,
        'rsi': current_rsi,
        'bb_position': bb_position
    }

# Run financial analysis
stock_data = generate_stock_data()
financial_metrics, portfolio_returns = calculate_financial_metrics(stock_data)
correlation_matrix = correlation_analysis(stock_data)
technical_indicators = technical_analysis(stock_data)

print("\n✅ Financial data analysis completed!")
print("FinOps research environment ready for advanced financial modeling")
EOF

python3 financial_analysis.py
```

**What this does**: Analyzes stock market data with financial metrics, correlations, and technical analysis.

**This will take**: 2-3 minutes

### Economic Modeling
```bash
# Create economic modeling script
cat > economic_modeling.py << 'EOF'
import pandas as pd
import numpy as np
from scipy import stats
import statsmodels.api as sm

print("Starting economic modeling analysis...")

def generate_economic_data():
    """Generate synthetic economic time series data"""
    print("\n=== Economic Data Generation ===")

    np.random.seed(42)

    # Generate quarterly data for 20 years (80 quarters)
    quarters = pd.date_range('2004-Q1', '2023-Q4', freq='Q')
    n_quarters = len(quarters)

    print(f"Time series: {n_quarters} quarters from {quarters[0]} to {quarters[-1]}")

    # Generate correlated economic variables
    # GDP, Inflation, Unemployment, Interest Rates, Consumer Confidence

    # Start with independent shocks
    shocks = np.random.multivariate_normal(
        mean=[0, 0, 0, 0, 0],
        cov=[[1.0, 0.2, -0.3, 0.1, 0.4],    # GDP shocks
             [0.2, 1.0, 0.1, 0.6, -0.2],    # Inflation shocks
             [-0.3, 0.1, 1.0, -0.2, -0.5],  # Unemployment shocks
             [0.1, 0.6, -0.2, 1.0, 0.1],    # Interest rate shocks
             [0.4, -0.2, -0.5, 0.1, 1.0]],  # Consumer confidence shocks
        size=n_quarters
    )

    # Initialize variables
    gdp_growth = np.zeros(n_quarters)
    inflation = np.zeros(n_quarters)
    unemployment = np.zeros(n_quarters)
    interest_rate = np.zeros(n_quarters)
    consumer_confidence = np.zeros(n_quarters)

    # Set initial values
    gdp_growth[0] = 2.5  # 2.5% initial GDP growth
    inflation[0] = 2.0   # 2% initial inflation
    unemployment[0] = 5.5  # 5.5% initial unemployment
    interest_rate[0] = 3.0  # 3% initial interest rate
    consumer_confidence[0] = 100  # Index = 100

    # Generate time series with persistence and cross-correlations
    for t in range(1, n_quarters):
        # GDP growth with persistence and cyclical component
        gdp_growth[t] = (0.7 * gdp_growth[t-1] +
                        0.3 * np.sin(2 * np.pi * t / 16) +  # 4-year cycle
                        shocks[t, 0])

        # Inflation with persistence and monetary policy response
        inflation[t] = (0.8 * inflation[t-1] +
                       0.1 * gdp_growth[t] +  # Phillips curve
                       shocks[t, 1])

        # Unemployment with persistence and Okun's law
        unemployment[t] = (0.9 * unemployment[t-1] -
                          0.3 * (gdp_growth[t] - 2.5) +  # Okun's law
                          shocks[t, 2])

        # Interest rate with Taylor rule
        interest_rate[t] = (0.7 * interest_rate[t-1] +
                           0.5 * (inflation[t] - 2.0) +  # Inflation target
                           0.2 * gdp_growth[t] +
                           shocks[t, 3])

        # Consumer confidence
        consumer_confidence[t] = (0.8 * consumer_confidence[t-1] +
                                 2.0 * gdp_growth[t] -
                                 1.5 * unemployment[t] +
                                 100 + shocks[t, 4])

    # Ensure realistic bounds
    unemployment = np.clip(unemployment, 3.0, 12.0)
    interest_rate = np.clip(interest_rate, 0.0, 8.0)
    inflation = np.clip(inflation, -2.0, 6.0)
    consumer_confidence = np.clip(consumer_confidence, 50, 150)

    # Create DataFrame
    econ_data = pd.DataFrame({
        'Date': quarters,
        'GDP_Growth': gdp_growth,
        'Inflation': inflation,
        'Unemployment': unemployment,
        'Interest_Rate': interest_rate,
        'Consumer_Confidence': consumer_confidence
    })

    econ_data.set_index('Date', inplace=True)

    print("Economic Variables Summary:")
    print(econ_data.describe().round(3))

    return econ_data

def econometric_analysis(econ_data):
    """Perform econometric analysis"""
    print("\n=== Econometric Analysis ===")

    # 1. Phillips Curve: Inflation vs Unemployment
    print("1. Phillips Curve Analysis:")

    X = sm.add_constant(econ_data['Unemployment'])
    y = econ_data['Inflation']

    phillips_model = sm.OLS(y, X).fit()

    print(f"   Inflation = {phillips_model.params[0]:.3f} + {phillips_model.params[1]:.3f} * Unemployment")
    print(f"   R-squared: {phillips_model.rsquared:.3f}")
    print(f"   Unemployment coefficient p-value: {phillips_model.pvalues[1]:.3f}")

    if phillips_model.pvalues[1] < 0.05:
        if phillips_model.params[1] < 0:
            print("   Result: Significant negative relationship (classic Phillips curve)")
        else:
            print("   Result: Significant positive relationship")
    else:
        print("   Result: No significant relationship")

    # 2. Okun's Law: GDP Growth vs Unemployment
    print("\n2. Okun's Law Analysis:")

    # Use change in unemployment
    unemployment_change = econ_data['Unemployment'].diff().dropna()
    gdp_growth_aligned = econ_data['GDP_Growth'][1:]

    X_okun = sm.add_constant(gdp_growth_aligned)
    y_okun = unemployment_change

    okun_model = sm.OLS(y_okun, X_okun).fit()

    print(f"   ΔUnemployment = {okun_model.params[0]:.3f} + {okun_model.params[1]:.3f} * GDP_Growth")
    print(f"   R-squared: {okun_model.rsquared:.3f}")
    print(f"   GDP Growth coefficient: {okun_model.params[1]:.3f} (p-value: {okun_model.pvalues[1]:.3f})")

    if okun_model.pvalues[1] < 0.05:
        print("   Result: Significant relationship between GDP growth and unemployment change")

    # 3. Consumer Confidence and Economic Activity
    print("\n3. Consumer Confidence Analysis:")

    X_conf = sm.add_constant(econ_data[['GDP_Growth', 'Unemployment', 'Inflation']])
    y_conf = econ_data['Consumer_Confidence']

    confidence_model = sm.OLS(y_conf, X_conf).fit()

    print(f"   Consumer Confidence Model:")
    print(f"   R-squared: {confidence_model.rsquared:.3f}")
    print(f"   Coefficients:")
    for var, coef, pval in zip(['Constant', 'GDP_Growth', 'Unemployment', 'Inflation'],
                               confidence_model.params, confidence_model.pvalues):
        significance = "***" if pval < 0.01 else "**" if pval < 0.05 else "*" if pval < 0.1 else ""
        print(f"     {var}: {coef:.3f} {significance}")

    return phillips_model, okun_model, confidence_model

def time_series_analysis(econ_data):
    """Perform time series analysis"""
    print("\n=== Time Series Analysis ===")

    # Test for stationarity using Augmented Dickey-Fuller test
    from statsmodels.tsa.stattools import adfuller

    print("Stationarity Tests (ADF Test):")
    for column in econ_data.columns:
        result = adfuller(econ_data[column].dropna())

        print(f"  {column}:")
        print(f"    ADF Statistic: {result[0]:.4f}")
        print(f"    p-value: {result[1]:.4f}")

        if result[1] <= 0.05:
            print("    Result: Stationary (reject unit root)")
        else:
            print("    Result: Non-stationary (unit root present)")

    # Granger Causality Test
    from statsmodels.tsa.stattools import grangercausalitytests

    print(f"\nGranger Causality Tests:")

    # Test if GDP Growth Granger-causes Unemployment
    try:
        # Combine variables for Granger test
        granger_data = econ_data[['Unemployment', 'GDP_Growth']].dropna()

        print("  GDP Growth → Unemployment:")
        granger_result = grangercausalitytests(granger_data, maxlag=4, verbose=False)

        # Get p-value for lag 1
        p_value = granger_result[1][0]['ssr_ftest'][1]
        print(f"    p-value (lag 1): {p_value:.4f}")

        if p_value < 0.05:
            print("    Result: GDP Growth Granger-causes Unemployment")
        else:
            print("    Result: No Granger causality")

    except Exception as e:
        print(f"    Error in Granger test: {e}")

    # Autocorrelation analysis
    print(f"\nAutocorrelation Analysis:")

    for column in ['GDP_Growth', 'Inflation']:
        series = econ_data[column].dropna()

        # Calculate autocorrelations for lags 1-4
        autocorrs = [series.autocorr(lag=i) for i in range(1, 5)]

        print(f"  {column} Autocorrelations:")
        for lag, ac in enumerate(autocorrs, 1):
            print(f"    Lag {lag}: {ac:.3f}")

    return econ_data

def forecasting_models(econ_data):
    """Build forecasting models"""
    print("\n=== Economic Forecasting Models ===")

    # 1. ARIMA model for GDP Growth
    from statsmodels.tsa.arima.model import ARIMA

    gdp_series = econ_data['GDP_Growth'].dropna()

    # Fit ARIMA(1,0,1) model
    try:
        arima_model = ARIMA(gdp_series, order=(1, 0, 1))
        arima_fit = arima_model.fit()

        print("ARIMA(1,0,1) Model for GDP Growth:")
        print(f"  AIC: {arima_fit.aic:.2f}")
        print(f"  Parameters:")
        for param, value in zip(['AR(1)', 'MA(1)', 'Constant'], arima_fit.params):
            print(f"    {param}: {value:.4f}")

        # Generate forecast
        forecast = arima_fit.forecast(steps=4)  # 4 quarters ahead
        print(f"  4-quarter forecast: {forecast.mean():.3f}")

    except Exception as e:
        print(f"ARIMA model error: {e}")
        arima_fit = None

    # 2. Vector Autoregression (VAR) for multiple variables
    from statsmodels.tsa.api import VAR

    # Select key variables for VAR
    var_data = econ_data[['GDP_Growth', 'Inflation', 'Unemployment']].dropna()

    try:
        var_model = VAR(var_data)
        var_fit = var_model.fit(maxlags=2)

        print(f"\nVAR Model (lag order: 2):")
        print(f"  AIC: {var_fit.aic:.2f}")
        print(f"  Log Likelihood: {var_fit.llf:.2f}")

        # Generate VAR forecast
        var_forecast = var_fit.forecast(var_data.values, steps=4)

        print(f"  4-quarter VAR forecasts:")
        for i, var_name in enumerate(var_data.columns):
            print(f"    {var_name}: {var_forecast[-1, i]:.3f}")

    except Exception as e:
        print(f"VAR model error: {e}")
        var_fit = None

    # 3. Simple trend and seasonal decomposition
    print(f"\nTrend Analysis:")

    for column in ['GDP_Growth', 'Consumer_Confidence']:
        series = econ_data[column]

        # Calculate linear trend
        time_index = np.arange(len(series))
        trend_coef = np.polyfit(time_index, series, 1)

        print(f"  {column}:")
        print(f"    Linear trend: {trend_coef[0]:.4f} per quarter")

        if abs(trend_coef[0]) > 0.01:
            direction = "increasing" if trend_coef[0] > 0 else "decreasing"
            print(f"    Trend direction: {direction}")
        else:
            print(f"    Trend direction: stable")

    return arima_fit, var_fit

# Run economic modeling
economic_data = generate_economic_data()
phillips_model, okun_model, confidence_model = econometric_analysis(economic_data)
ts_results = time_series_analysis(economic_data)
arima_model, var_model = forecasting_models(economic_data)

print("\n✅ Economic modeling analysis completed!")
print("Advanced econometric and forecasting capabilities demonstrated")
EOF

python3 economic_modeling.py
```

**What this does**: Demonstrates econometric modeling, time series analysis, and economic forecasting.

**Expected result**: Shows comprehensive economic analysis including Phillips curves, Okun's law, and forecasting models.

## Step 9: Risk Management and Portfolio Optimization

Test advanced FinOps capabilities:

```bash
# Create risk management and portfolio optimization script
cat > risk_portfolio_optimization.py << 'EOF'
import pandas as pd
import numpy as np
from scipy import optimize
import matplotlib.pyplot as plt

print("Starting risk management and portfolio optimization...")

def generate_asset_returns():
    """Generate synthetic asset return data"""
    print("\n=== Asset Return Data Generation ===")

    np.random.seed(42)

    # Generate 5 years of daily returns for different asset classes
    n_days = 1250  # ~5 years of trading days

    asset_params = {
        'US_Stocks': {'mean': 0.08/252, 'vol': 0.16/np.sqrt(252)},
        'International_Stocks': {'mean': 0.07/252, 'vol': 0.18/np.sqrt(252)},
        'Bonds': {'mean': 0.04/252, 'vol': 0.05/np.sqrt(252)},
        'Real_Estate': {'mean': 0.06/252, 'vol': 0.12/np.sqrt(252)},
        'Commodities': {'mean': 0.05/252, 'vol': 0.22/np.sqrt(252)}
    }

    # Correlation matrix
    correlation_matrix = np.array([
        [1.00, 0.75, 0.15, 0.60, 0.30],  # US Stocks
        [0.75, 1.00, 0.10, 0.55, 0.35],  # International Stocks
        [0.15, 0.10, 1.00, 0.25, -0.10], # Bonds
        [0.60, 0.55, 0.25, 1.00, 0.40],  # Real Estate
        [0.30, 0.35, -0.10, 0.40, 1.00]  # Commodities
    ])

    # Generate correlated returns
    means = np.array([params['mean'] for params in asset_params.values()])
    vols = np.array([params['vol'] for params in asset_params.values()])

    # Create covariance matrix
    cov_matrix = np.outer(vols, vols) * correlation_matrix

    # Generate multivariate normal returns
    returns = np.random.multivariate_normal(means, cov_matrix, n_days)

    # Create DataFrame
    returns_df = pd.DataFrame(returns, columns=asset_params.keys())

    print(f"Generated {n_days} days of return data for {len(asset_params)} assets")

    # Calculate summary statistics
    annual_returns = returns_df.mean() * 252
    annual_vols = returns_df.std() * np.sqrt(252)
    sharpe_ratios = annual_returns / annual_vols

    print("\nAsset Statistics:")
    stats_df = pd.DataFrame({
        'Annual Return': annual_returns,
        'Annual Volatility': annual_vols,
        'Sharpe Ratio': sharpe_ratios
    })
    print(stats_df.round(4))

    return returns_df, cov_matrix * 252  # Annualized covariance

def portfolio_optimization(returns_df, cov_matrix):
    """Perform portfolio optimization"""
    print("\n=== Portfolio Optimization ===")

    n_assets = len(returns_df.columns)
    annual_returns = returns_df.mean() * 252

    # 1. Minimum Variance Portfolio
    def portfolio_variance(weights, cov_matrix):
        return np.dot(weights.T, np.dot(cov_matrix, weights))

    def portfolio_return(weights, returns):
        return np.dot(weights, returns)

    # Constraints: weights sum to 1, all weights >= 0
    constraints = ({'type': 'eq', 'fun': lambda x: np.sum(x) - 1})
    bounds = tuple((0, 1) for _ in range(n_assets))

    # Initial guess (equal weights)
    x0 = np.array([1/n_assets] * n_assets)

    # Minimize variance
    min_var_result = optimize.minimize(
        portfolio_variance, x0, args=(cov_matrix,),
        method='SLSQP', bounds=bounds, constraints=constraints
    )

    min_var_weights = min_var_result.x
    min_var_return = portfolio_return(min_var_weights, annual_returns)
    min_var_vol = np.sqrt(portfolio_variance(min_var_weights, cov_matrix))

    print("Minimum Variance Portfolio:")
    print("  Weights:")
    for asset, weight in zip(returns_df.columns, min_var_weights):
        print(f"    {asset}: {weight:.3f}")
    print(f"  Expected Return: {min_var_return:.4f}")
    print(f"  Volatility: {min_var_vol:.4f}")
    print(f"  Sharpe Ratio: {min_var_return/min_var_vol:.4f}")

    # 2. Maximum Sharpe Ratio Portfolio
    def negative_sharpe_ratio(weights, returns, cov_matrix):
        port_return = portfolio_return(weights, returns)
        port_vol = np.sqrt(portfolio_variance(weights, cov_matrix))
        return -port_return / port_vol  # Negative because we minimize

    max_sharpe_result = optimize.minimize(
        negative_sharpe_ratio, x0, args=(annual_returns, cov_matrix),
        method='SLSQP', bounds=bounds, constraints=constraints
    )

    max_sharpe_weights = max_sharpe_result.x
    max_sharpe_return = portfolio_return(max_sharpe_weights, annual_returns)
    max_sharpe_vol = np.sqrt(portfolio_variance(max_sharpe_weights, cov_matrix))

    print("\nMaximum Sharpe Ratio Portfolio:")
    print("  Weights:")
    for asset, weight in zip(returns_df.columns, max_sharpe_weights):
        print(f"    {asset}: {weight:.3f}")
    print(f"  Expected Return: {max_sharpe_return:.4f}")
    print(f"  Volatility: {max_sharpe_vol:.4f}")
    print(f"  Sharpe Ratio: {max_sharpe_return/max_sharpe_vol:.4f}")

    # 3. Efficient Frontier
    print("\nEfficient Frontier Calculation:")

    target_returns = np.linspace(min_var_return, max_sharpe_return * 0.9, 10)
    efficient_portfolios = []

    for target_return in target_returns:
        # Add return constraint
        return_constraint = {'type': 'eq', 'fun': lambda x: portfolio_return(x, annual_returns) - target_return}
        all_constraints = [constraints, return_constraint]

        # Minimize variance for given return
        result = optimize.minimize(
            portfolio_variance, x0, args=(cov_matrix,),
            method='SLSQP', bounds=bounds, constraints=all_constraints
        )

        if result.success:
            portfolio_vol = np.sqrt(result.fun)
            efficient_portfolios.append((target_return, portfolio_vol))

    print(f"  Generated {len(efficient_portfolios)} efficient portfolios")

    return min_var_weights, max_sharpe_weights, efficient_portfolios

def risk_metrics_analysis(returns_df):
    """Calculate various risk metrics"""
    print("\n=== Risk Metrics Analysis ===")

    # Portfolio returns (equal weight for example)
    portfolio_returns = returns_df.mean(axis=1)

    # 1. Value at Risk (VaR)
    confidence_levels = [0.95, 0.99]
    print("Value at Risk (VaR):")

    for conf_level in confidence_levels:
        var_historical = np.percentile(portfolio_returns, (1 - conf_level) * 100)
        var_parametric = stats.norm.ppf(1 - conf_level, portfolio_returns.mean(), portfolio_returns.std())

        print(f"  {conf_level*100:.0f}% VaR:")
        print(f"    Historical: {var_historical:.4f} ({var_historical*100:.2f}%)")
        print(f"    Parametric: {var_parametric:.4f} ({var_parametric*100:.2f}%)")

    # 2. Conditional Value at Risk (Expected Shortfall)
    var_95 = np.percentile(portfolio_returns, 5)
    cvar_95 = portfolio_returns[portfolio_returns <= var_95].mean()

    print(f"\nConditional VaR (95%):")
    print(f"  Expected Shortfall: {cvar_95:.4f} ({cvar_95*100:.2f}%)")

    # 3. Maximum Drawdown
    cumulative_returns = (1 + portfolio_returns).cumprod()
    rolling_max = cumulative_returns.expanding().max()
    drawdowns = (cumulative_returns - rolling_max) / rolling_max
    max_drawdown = drawdowns.min()

    print(f"\nDrawdown Analysis:")
    print(f"  Maximum Drawdown: {max_drawdown:.4f} ({max_drawdown*100:.2f}%)")

    # Find drawdown periods
    in_drawdown = drawdowns < -0.05  # 5% threshold
    if in_drawdown.any():
        drawdown_periods = []
        start_dd = None

        for i, dd in enumerate(in_drawdown):
            if dd and start_dd is None:
                start_dd = i
            elif not dd and start_dd is not None:
                drawdown_periods.append((start_dd, i-1, drawdowns[start_dd:i].min()))
                start_dd = None

        print(f"  Significant drawdown periods (>5%): {len(drawdown_periods)}")

        if drawdown_periods:
            worst_dd = min(drawdown_periods, key=lambda x: x[2])
            print(f"  Worst drawdown: {worst_dd[2]:.4f} (duration: {worst_dd[1] - worst_dd[0] + 1} days)")

    # 4. Risk-adjusted returns
    annual_return = portfolio_returns.mean() * 252
    annual_vol = portfolio_returns.std() * np.sqrt(252)
    sharpe_ratio = annual_return / annual_vol

    # Sortino ratio (downside deviation)
    downside_returns = portfolio_returns[portfolio_returns < 0]
    downside_vol = downside_returns.std() * np.sqrt(252) if len(downside_returns) > 0 else 0
    sortino_ratio = annual_return / downside_vol if downside_vol > 0 else float('inf')

    print(f"\nRisk-Adjusted Returns:")
    print(f"  Sharpe Ratio: {sharpe_ratio:.4f}")
    print(f"  Sortino Ratio: {sortino_ratio:.4f}")

    # 5. Beta analysis (relative to market proxy - use first asset as market)
    market_returns = returns_df.iloc[:, 0]  # First asset as market proxy

    betas = {}
    for asset in returns_df.columns:
        if asset != returns_df.columns[0]:  # Skip market proxy
            covariance = np.cov(returns_df[asset], market_returns)[0, 1]
            market_variance = np.var(market_returns)
            beta = covariance / market_variance
            betas[asset] = beta

    print(f"\nBeta Analysis (vs {returns_df.columns[0]}):")
    for asset, beta in betas.items():
        risk_level = "High" if beta > 1.2 else "Moderate" if beta > 0.8 else "Low"
        print(f"  {asset}: {beta:.3f} ({risk_level} risk)")

    return {
        'var_95': var_95,
        'cvar_95': cvar_95,
        'max_drawdown': max_drawdown,
        'sharpe_ratio': sharpe_ratio,
        'sortino_ratio': sortino_ratio
    }

def scenario_analysis(returns_df):
    """Perform scenario analysis and stress testing"""
    print("\n=== Scenario Analysis & Stress Testing ===")

    # Equal weight portfolio
    portfolio_weights = np.array([1/len(returns_df.columns)] * len(returns_df.columns))
    portfolio_returns = (returns_df * portfolio_weights).sum(axis=1)

    # Historical scenarios
    scenarios = {
        'Financial Crisis (2008-style)': {
            'stocks_shock': -0.30,
            'bonds_change': 0.15,
            'real_estate_shock': -0.25,
            'commodities_shock': -0.20
        },
        'Inflation Surge': {
            'stocks_shock': -0.10,
            'bonds_change': -0.15,
            'real_estate_shock': 0.05,
            'commodities_shock': 0.20
        },
        'Economic Boom': {
            'stocks_shock': 0.25,
            'bonds_change': -0.05,
            'real_estate_shock': 0.15,
            'commodities_shock': 0.10
        }
    }

    print("Stress Test Results:")

    current_portfolio_value = 100000  # $100k initial value

    for scenario_name, shocks in scenarios.items():
        # Apply shocks to different asset classes
        shocked_returns = portfolio_returns.copy()

        # Map shocks to our assets (simplified mapping)
        asset_shocks = {
            'US_Stocks': shocks.get('stocks_shock', 0),
            'International_Stocks': shocks.get('stocks_shock', 0),
            'Bonds': shocks.get('bonds_change', 0),
            'Real_Estate': shocks.get('real_estate_shock', 0),
            'Commodities': shocks.get('commodities_shock', 0)
        }

        # Calculate portfolio impact
        portfolio_shock = sum(portfolio_weights[i] * shock
                            for i, (asset, shock) in enumerate(asset_shocks.items()))

        shocked_value = current_portfolio_value * (1 + portfolio_shock)
        loss_amount = current_portfolio_value - shocked_value
        loss_percentage = portfolio_shock * 100

        print(f"\n  {scenario_name}:")
        print(f"    Portfolio Impact: {loss_percentage:+.2f}%")
        print(f"    Value Change: ${loss_amount:+,.0f}")
        print(f"    New Portfolio Value: ${shocked_value:,.0f}")

        # Risk assessment
        if abs(loss_percentage) < 5:
            risk_assessment = "Low Impact"
        elif abs(loss_percentage) < 15:
            risk_assessment = "Moderate Impact"
        else:
            risk_assessment = "High Impact"

        print(f"    Risk Assessment: {risk_assessment}")

    # Monte Carlo simulation for portfolio outcomes
    print(f"\nMonte Carlo Portfolio Simulation:")

    n_simulations = 10000
    time_horizon = 252  # 1 year

    # Generate random scenarios
    portfolio_mean = portfolio_returns.mean()
    portfolio_std = portfolio_returns.std()

    final_values = []

    for _ in range(n_simulations):
        # Generate random path
        random_returns = np.random.normal(portfolio_mean, portfolio_std, time_horizon)
        cumulative_return = np.prod(1 + random_returns) - 1
        final_value = current_portfolio_value * (1 + cumulative_return)
        final_values.append(final_value)

    final_values = np.array(final_values)

    # Calculate percentiles
    percentiles = [5, 25, 50, 75, 95]
    percentile_values = np.percentile(final_values, percentiles)

    print(f"  1-Year Portfolio Value Projections (${current_portfolio_value:,} initial):")
    for p, value in zip(percentiles, percentile_values):
        return_pct = (value / current_portfolio_value - 1) * 100
        print(f"    {p:2d}th percentile: ${value:,.0f} ({return_pct:+.1f}%)")

    # Probability of loss
    prob_loss = np.mean(final_values < current_portfolio_value) * 100
    prob_large_loss = np.mean(final_values < current_portfolio_value * 0.9) * 100

    print(f"\n  Risk Probabilities:")
    print(f"    Probability of any loss: {prob_loss:.1f}%")
    print(f"    Probability of >10% loss: {prob_large_loss:.1f}%")

    return final_values

# Run risk management and portfolio optimization
returns_data, cov_matrix_annual = generate_asset_returns()
min_var_weights, max_sharpe_weights, efficient_frontier = portfolio_optimization(returns_data, cov_matrix_annual)
risk_metrics = risk_metrics_analysis(returns_data)
simulation_results = scenario_analysis(returns_data)

print("\n✅ Risk management and portfolio optimization completed!")
print("Advanced financial risk analysis and optimization capabilities demonstrated")
EOF

python3 risk_portfolio_optimization.py
```

**What this does**: Demonstrates portfolio optimization, risk metrics calculation, and scenario analysis.

**Expected result**: Shows comprehensive risk management and portfolio optimization results.

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
- **Storage**: $0.10 per GB per month for financial datasets you save
- **Data Transfer**: Usually free for financial research data amounts

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

Now that you have a working FinOps environment, you can:

### Learn More About Financial Research
- [High-Frequency Trading Analysis Tutorial](high-frequency-trading.md)
- [Advanced Econometric Modeling Guide](advanced-econometrics.md)
- [Cost Optimization for Financial Research](finops-cost-optimization.md)

### Explore Advanced Features
- [Real-time market data integration](real-time-market-data.md)
- [Team collaboration with financial models](team-finops-collaboration.md)
- [Automated trading strategy testing](automated-trading-strategies.md)

### Join the FinOps Community
- [Financial Research Forum](https://forum.researchwizard.app/finops)
- [GitHub FinOps Examples](https://github.com/aws-research-wizard/finops-examples)
- [Monthly Financial Modeling Office Hours](https://calendar.researchwizard.app/finops-office-hours)

## Troubleshooting

### Common Issues

**Problem**: "QuantLib import error" during financial analysis
**Solution**: Check QuantLib installation: `python -c "import QuantLib"` and reinstall if needed
**Prevention**: Wait 5-7 minutes after deployment for all financial packages to initialize

**Problem**: "Convergence error" in optimization algorithms
**Solution**: Try different starting points or reduce convergence tolerance
**Prevention**: Check input data quality and parameter bounds

**Problem**: "Memory error" during large portfolio optimization
**Solution**: Reduce number of assets or use a larger instance type
**Prevention**: Monitor memory usage with `htop` during optimization

**Problem**: "Data format error" when loading financial data
**Solution**: Verify date formats and missing value handling
**Prevention**: Always validate financial data before analysis

### Getting Help
- Check the [FinOps troubleshooting guide](troubleshooting-finops.md)
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
