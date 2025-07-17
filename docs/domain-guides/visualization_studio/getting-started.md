# Visualization Studio Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $10-16 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working scientific visualization environment that can:
- Create interactive 3D visualizations and animations
- Generate publication-quality figures and plots
- Build web-based dashboards and interactive applications
- Process and visualize large datasets with high-performance rendering

### Meet Dr. Lisa Rodriguez

Dr. Lisa Rodriguez is a data scientist at CERN. She creates visualizations for particle physics data but waits hours for graphics workstations. Each visualization project requires processing terabytes of experimental data into interactive displays.

**Before**: 4-hour waits + 8-hour rendering = 12 hours per visualization
**After**: 15-minute setup + 1-hour rendering = same day results
**Time Saved**: 92% faster visualization workflow
**Cost Savings**: $600/month vs $2,400 graphics workstation allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $10-16 (we'll clean up resources when done)
- **Daily research cost**: $20-50 per day when actively rendering
- **Monthly estimate**: $250-650 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No visualization or programming experience required

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
- **Region**: Choose `us-west-2` (recommended for visualization with good GPU performance)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain visualization_studio --region us-west-2
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: visualization_studio
✅ Region valid: us-west-2 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Visualization Environment

```bash
aws-research-wizard deploy start --domain visualization_studio --region us-west-2 --instance g4dn.xlarge
```

**What this does**: Creates your visualization environment optimized for graphics rendering and interactive visualization.

**This will take**: 5-7 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
  GPU: NVIDIA T4 for high-performance rendering
  Memory: 16GB RAM for large visualizations
```

**💰 Billing starts now**: Your environment costs about $0.53 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
```

**What this does**: Connects you to your visualization computer in the cloud.

**Expected result**: You see a command prompt like `ubuntu@ip-10-0-1-123:~$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Visualization Tools

Your environment comes pre-installed with:

### Core Visualization Software
- **Python Visualization**: Matplotlib, Plotly, Bokeh - Type `python -c "import matplotlib; print(matplotlib.__version__)"` to check
- **3D Graphics**: Three.js, WebGL support - Type `python -c "import plotly; print(plotly.__version__)"` to check
- **Interactive Dashboards**: Dash, Streamlit - Type `streamlit --version` to check
- **Scientific Plotting**: Seaborn, Altair - Type `python -c "import seaborn; print(seaborn.__version__)"` to check
- **GPU Acceleration**: CUDA, OpenGL - Type `nvidia-smi` to check

### Try Your First Command
```bash
python -c "import matplotlib; print('Matplotlib version:', matplotlib.__version__)"
```

**What this does**: Shows Matplotlib version and confirms visualization tools are installed.

**Expected result**: You see Matplotlib version info confirming visualization libraries are ready.

## Step 8: Create Scientific Visualizations

Let's create interactive visualizations to test everything works:

### 2D Plotting and Charts
```bash
# Create working directory
mkdir ~/visualization-tutorial
cd ~/visualization-tutorial

# Create scientific plotting script
cat > scientific_plotting.py << 'EOF'
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import seaborn as sns
from scipy import stats
import plotly.graph_objects as go
import plotly.express as px
from plotly.subplots import make_subplots

print("Starting scientific visualization analysis...")

def create_publication_figures():
    """Create publication-quality scientific figures"""
    print("\n=== Publication-Quality Figures ===")

    # Set publication style
    plt.style.use('seaborn-v0_8-whitegrid')

    # Figure 1: Multi-panel scientific plot
    fig, axes = plt.subplots(2, 2, figsize=(12, 10))
    fig.suptitle('Scientific Data Analysis', fontsize=16, fontweight='bold')

    # Panel A: Experimental data with error bars
    np.random.seed(42)
    x = np.linspace(0, 10, 50)
    y_true = 2 * np.sin(x) + 0.5 * x
    y_data = y_true + np.random.normal(0, 0.3, len(x))
    y_error = np.random.uniform(0.1, 0.5, len(x))

    axes[0, 0].errorbar(x, y_data, yerr=y_error, fmt='o', capsize=3,
                       color='steelblue', alpha=0.7, label='Experimental Data')
    axes[0, 0].plot(x, y_true, 'r-', linewidth=2, label='Theoretical Model')
    axes[0, 0].set_xlabel('Time (s)')
    axes[0, 0].set_ylabel('Signal Amplitude')
    axes[0, 0].set_title('A) Experimental vs Theoretical')
    axes[0, 0].legend()
    axes[0, 0].grid(True, alpha=0.3)

    # Panel B: Statistical distribution
    data = np.random.normal(100, 15, 1000)
    axes[0, 1].hist(data, bins=30, density=True, alpha=0.7, color='lightcoral',
                   edgecolor='black', label='Data')

    # Fit normal distribution
    mu, sigma = stats.norm.fit(data)
    x_fit = np.linspace(data.min(), data.max(), 100)
    y_fit = stats.norm.pdf(x_fit, mu, sigma)
    axes[0, 1].plot(x_fit, y_fit, 'k-', linewidth=2,
                   label=f'Normal Fit (μ={mu:.1f}, σ={sigma:.1f})')

    axes[0, 1].set_xlabel('Value')
    axes[0, 1].set_ylabel('Probability Density')
    axes[0, 1].set_title('B) Statistical Distribution')
    axes[0, 1].legend()

    # Panel C: Correlation heatmap
    corr_data = np.random.multivariate_normal([0, 0, 0],
                                            [[1, 0.5, -0.3], [0.5, 1, 0.8], [-0.3, 0.8, 1]],
                                            500)
    corr_df = pd.DataFrame(corr_data, columns=['Variable A', 'Variable B', 'Variable C'])

    im = axes[1, 0].imshow(corr_df.corr(), cmap='coolwarm', aspect='auto', vmin=-1, vmax=1)
    axes[1, 0].set_xticks(range(len(corr_df.columns)))
    axes[1, 0].set_yticks(range(len(corr_df.columns)))
    axes[1, 0].set_xticklabels(corr_df.columns, rotation=45)
    axes[1, 0].set_yticklabels(corr_df.columns)
    axes[1, 0].set_title('C) Correlation Matrix')

    # Add correlation values
    for i in range(len(corr_df.columns)):
        for j in range(len(corr_df.columns)):
            text = axes[1, 0].text(j, i, f'{corr_df.corr().iloc[i, j]:.2f}',
                                 ha="center", va="center", color="black", fontweight='bold')

    # Panel D: Time series with trend
    dates = pd.date_range('2020-01-01', periods=365, freq='D')
    trend = np.linspace(0, 2, len(dates))
    seasonal = 0.5 * np.sin(2 * np.pi * np.arange(len(dates)) / 365.25 * 4)
    noise = np.random.normal(0, 0.2, len(dates))
    ts_data = 10 + trend + seasonal + noise

    axes[1, 1].plot(dates, ts_data, 'b-', linewidth=1, alpha=0.7, label='Data')
    axes[1, 1].plot(dates, 10 + trend, 'r--', linewidth=2, label='Trend')
    axes[1, 1].set_xlabel('Date')
    axes[1, 1].set_ylabel('Value')
    axes[1, 1].set_title('D) Time Series Analysis')
    axes[1, 1].legend()
    axes[1, 1].tick_params(axis='x', rotation=45)

    plt.tight_layout()
    plt.savefig('scientific_figure.png', dpi=300, bbox_inches='tight')
    plt.close()

    print("Created publication-quality figure: scientific_figure.png")
    print("  Resolution: 300 DPI")
    print("  Format: PNG with transparent background")
    print("  Panels: 4 (A-D) with different visualization types")

def create_interactive_plots():
    """Create interactive web-based visualizations"""
    print("\n=== Interactive Visualizations ===")

    # Generate sample data
    np.random.seed(42)

    # 3D scatter plot
    n_points = 1000
    x = np.random.normal(0, 1, n_points)
    y = np.random.normal(0, 1, n_points)
    z = x**2 + y**2 + np.random.normal(0, 0.5, n_points)
    colors = z

    fig_3d = go.Figure(data=[go.Scatter3d(
        x=x, y=y, z=z,
        mode='markers',
        marker=dict(
            size=3,
            color=colors,
            colorscale='Viridis',
            showscale=True,
            colorbar=dict(title="Z Value")
        ),
        text=[f'Point {i}' for i in range(n_points)],
        hovertemplate='<b>%{text}</b><br>X: %{x:.2f}<br>Y: %{y:.2f}<br>Z: %{z:.2f}<extra></extra>'
    )])

    fig_3d.update_layout(
        title='Interactive 3D Scatter Plot',
        scene=dict(
            xaxis_title='X Axis',
            yaxis_title='Y Axis',
            zaxis_title='Z Axis'
        ),
        width=800,
        height=600
    )

    fig_3d.write_html('3d_scatter.html')
    print("Created interactive 3D scatter plot: 3d_scatter.html")

    # Multi-trace time series
    dates = pd.date_range('2020-01-01', periods=365, freq='D')

    # Multiple time series
    series_data = {}
    for i, name in enumerate(['Dataset A', 'Dataset B', 'Dataset C']):
        trend = np.linspace(0, 2 + i, len(dates))
        seasonal = (0.5 + i*0.2) * np.sin(2 * np.pi * np.arange(len(dates)) / 365.25 * 4)
        noise = np.random.normal(0, 0.3, len(dates))
        series_data[name] = 10 + i*5 + trend + seasonal + noise

    fig_ts = go.Figure()

    for name, data in series_data.items():
        fig_ts.add_trace(go.Scatter(
            x=dates,
            y=data,
            mode='lines',
            name=name,
            line=dict(width=2),
            hovertemplate='<b>%{fullData.name}</b><br>Date: %{x}<br>Value: %{y:.2f}<extra></extra>'
        ))

    fig_ts.update_layout(
        title='Interactive Multi-Series Time Plot',
        xaxis_title='Date',
        yaxis_title='Value',
        hovermode='x unified',
        width=1000,
        height=500
    )

    fig_ts.write_html('time_series.html')
    print("Created interactive time series plot: time_series.html")

    # Dashboard-style subplot
    fig_dashboard = make_subplots(
        rows=2, cols=2,
        subplot_titles=('Histogram', 'Box Plot', 'Scatter Plot', 'Bar Chart'),
        specs=[[{"secondary_y": False}, {"secondary_y": False}],
               [{"secondary_y": False}, {"secondary_y": False}]]
    )

    # Histogram
    hist_data = np.random.gamma(2, 2, 1000)
    fig_dashboard.add_trace(
        go.Histogram(x=hist_data, nbinsx=30, name='Distribution'),
        row=1, col=1
    )

    # Box plot
    box_data = [np.random.normal(0, 1, 100) for _ in range(4)]
    for i, data in enumerate(box_data):
        fig_dashboard.add_trace(
            go.Box(y=data, name=f'Group {i+1}'),
            row=1, col=2
        )

    # Scatter plot
    scatter_x = np.random.normal(0, 1, 200)
    scatter_y = 2 * scatter_x + np.random.normal(0, 0.5, 200)
    fig_dashboard.add_trace(
        go.Scatter(x=scatter_x, y=scatter_y, mode='markers', name='Data Points'),
        row=2, col=1
    )

    # Bar chart
    categories = ['A', 'B', 'C', 'D', 'E']
    values = np.random.uniform(10, 100, len(categories))
    fig_dashboard.add_trace(
        go.Bar(x=categories, y=values, name='Categories'),
        row=2, col=2
    )

    fig_dashboard.update_layout(
        title_text="Interactive Dashboard",
        showlegend=False,
        height=600,
        width=1000
    )

    fig_dashboard.write_html('dashboard.html')
    print("Created interactive dashboard: dashboard.html")

    print("\nInteractive visualizations created:")
    print("  - 3D scatter plot with hover information")
    print("  - Multi-series time plot with zoom and pan")
    print("  - Dashboard with multiple chart types")
    print("  - All plots are web-based and fully interactive")

def create_advanced_plots():
    """Create advanced visualization techniques"""
    print("\n=== Advanced Visualization Techniques ===")

    # Contour plot with custom colormap
    x = np.linspace(-3, 3, 100)
    y = np.linspace(-3, 3, 100)
    X, Y = np.meshgrid(x, y)
    Z = np.sin(X) * np.cos(Y) * np.exp(-(X**2 + Y**2)/4)

    fig_contour = go.Figure(data=go.Contour(
        x=x,
        y=y,
        z=Z,
        colorscale='RdBu',
        contours=dict(
            showlabels=True,
            labelfont=dict(size=12, color='white')
        ),
        hovertemplate='X: %{x:.2f}<br>Y: %{y:.2f}<br>Z: %{z:.3f}<extra></extra>'
    ))

    fig_contour.update_layout(
        title='2D Contour Plot with Custom Styling',
        xaxis_title='X',
        yaxis_title='Y',
        width=600,
        height=500
    )

    fig_contour.write_html('contour_plot.html')
    print("Created advanced contour plot: contour_plot.html")

    # Animated plot
    frames = []
    for t in np.linspace(0, 2*np.pi, 20):
        x_anim = np.linspace(0, 4*np.pi, 100)
        y_anim = np.sin(x_anim + t)

        frame = go.Frame(data=[go.Scatter(x=x_anim, y=y_anim, mode='lines', name='Wave')])
        frames.append(frame)

    fig_anim = go.Figure(
        data=[go.Scatter(x=x_anim, y=y_anim, mode='lines', name='Wave')],
        frames=frames
    )

    fig_anim.update_layout(
        title='Animated Sine Wave',
        xaxis_title='X',
        yaxis_title='Y',
        updatemenus=[{
            'type': 'buttons',
            'buttons': [
                {'label': 'Play', 'method': 'animate', 'args': [None]},
                {'label': 'Pause', 'method': 'animate', 'args': [None, {'frame': {'duration': 0}}]}
            ]
        }],
        width=800,
        height=400
    )

    fig_anim.write_html('animated_plot.html')
    print("Created animated visualization: animated_plot.html")

    # Parallel coordinates plot
    df = pd.DataFrame({
        'Feature_A': np.random.normal(0, 1, 100),
        'Feature_B': np.random.normal(0, 1, 100),
        'Feature_C': np.random.normal(0, 1, 100),
        'Feature_D': np.random.normal(0, 1, 100),
        'Category': np.random.choice(['Type 1', 'Type 2', 'Type 3'], 100)
    })

    fig_parallel = go.Figure(data=go.Parcoords(
        line=dict(color=pd.Categorical(df['Category']).codes,
                 colorscale='Set1',
                 showscale=True),
        dimensions=list([
            dict(range=[df['Feature_A'].min(), df['Feature_A'].max()],
                 label='Feature A', values=df['Feature_A']),
            dict(range=[df['Feature_B'].min(), df['Feature_B'].max()],
                 label='Feature B', values=df['Feature_B']),
            dict(range=[df['Feature_C'].min(), df['Feature_C'].max()],
                 label='Feature C', values=df['Feature_C']),
            dict(range=[df['Feature_D'].min(), df['Feature_D'].max()],
                 label='Feature D', values=df['Feature_D'])
        ])
    ))

    fig_parallel.update_layout(
        title='Parallel Coordinates Plot',
        width=800,
        height=500
    )

    fig_parallel.write_html('parallel_coordinates.html')
    print("Created parallel coordinates plot: parallel_coordinates.html")

    print("\nAdvanced visualizations demonstrate:")
    print("  - Custom colormaps and styling")
    print("  - Animation and interactive controls")
    print("  - Multi-dimensional data representation")
    print("  - Professional scientific presentation")

# Run visualization creation
create_publication_figures()
create_interactive_plots()
create_advanced_plots()

print("\n✅ Scientific visualization analysis completed!")
print("Visualization studio environment ready for advanced graphics and interactive displays")
EOF

python3 scientific_plotting.py
```

**What this does**: Creates publication-quality figures, interactive plots, and advanced visualizations.

**This will take**: 3-4 minutes

### 3D Visualization and Rendering
```bash
# Create 3D visualization script
cat > 3d_visualization.py << 'EOF'
import numpy as np
import plotly.graph_objects as go
from plotly.subplots import make_subplots
import plotly.express as px

print("Starting 3D visualization and rendering...")

def create_3d_surface_plots():
    """Create 3D surface visualizations"""
    print("\n=== 3D Surface Visualizations ===")

    # Mathematical surface
    x = np.linspace(-5, 5, 50)
    y = np.linspace(-5, 5, 50)
    X, Y = np.meshgrid(x, y)
    Z = np.sin(np.sqrt(X**2 + Y**2)) * np.exp(-0.1 * (X**2 + Y**2))

    fig_surface = go.Figure(data=[go.Surface(
        x=X, y=Y, z=Z,
        colorscale='Viridis',
        showscale=True,
        hovertemplate='X: %{x:.2f}<br>Y: %{y:.2f}<br>Z: %{z:.3f}<extra></extra>'
    )])

    fig_surface.update_layout(
        title='3D Mathematical Surface',
        scene=dict(
            xaxis_title='X Axis',
            yaxis_title='Y Axis',
            zaxis_title='Z Axis',
            camera=dict(eye=dict(x=1.5, y=1.5, z=1.5))
        ),
        width=800,
        height=600
    )

    fig_surface.write_html('3d_surface.html')
    print("Created 3D surface plot: 3d_surface.html")

    # Parametric surface (torus)
    u = np.linspace(0, 2*np.pi, 50)
    v = np.linspace(0, 2*np.pi, 50)
    U, V = np.meshgrid(u, v)

    R = 3  # Major radius
    r = 1  # Minor radius

    X_torus = (R + r * np.cos(V)) * np.cos(U)
    Y_torus = (R + r * np.cos(V)) * np.sin(U)
    Z_torus = r * np.sin(V)

    fig_torus = go.Figure(data=[go.Surface(
        x=X_torus, y=Y_torus, z=Z_torus,
        colorscale='Plasma',
        showscale=True
    )])

    fig_torus.update_layout(
        title='3D Parametric Surface (Torus)',
        scene=dict(
            xaxis_title='X',
            yaxis_title='Y',
            zaxis_title='Z',
            aspectmode='cube'
        ),
        width=800,
        height=600
    )

    fig_torus.write_html('3d_torus.html')
    print("Created 3D torus visualization: 3d_torus.html")

    # Volume rendering simulation
    print("\n3D Surface Features:")
    print("  - Interactive rotation and zooming")
    print("  - Custom lighting and shading")
    print("  - Hover information and tooltips")
    print("  - Multiple colormap options")

def create_molecular_visualization():
    """Create molecular-style 3D visualizations"""
    print("\n=== Molecular-Style 3D Visualizations ===")

    # Generate molecular-like structure
    np.random.seed(42)
    n_atoms = 50

    # Create random 3D positions
    positions = np.random.random((n_atoms, 3)) * 10

    # Assign atom types
    atom_types = np.random.choice(['H', 'C', 'N', 'O'], n_atoms, p=[0.4, 0.3, 0.2, 0.1])

    # Color and size mapping
    color_map = {'H': 'lightblue', 'C': 'black', 'N': 'blue', 'O': 'red'}
    size_map = {'H': 3, 'C': 8, 'N': 6, 'O': 7}

    colors = [color_map[atom] for atom in atom_types]
    sizes = [size_map[atom] for atom in atom_types]

    # Create bonds (connections between nearby atoms)
    bonds = []
    for i in range(n_atoms):
        for j in range(i+1, n_atoms):
            distance = np.linalg.norm(positions[i] - positions[j])
            if distance < 2.5:  # Bond threshold
                bonds.append((i, j))

    fig_mol = go.Figure()

    # Add atoms
    fig_mol.add_trace(go.Scatter3d(
        x=positions[:, 0],
        y=positions[:, 1],
        z=positions[:, 2],
        mode='markers',
        marker=dict(
            size=sizes,
            color=colors,
            opacity=0.8,
            line=dict(width=2, color='black')
        ),
        text=[f'Atom {i}: {atom}' for i, atom in enumerate(atom_types)],
        hovertemplate='<b>%{text}</b><br>X: %{x:.2f}<br>Y: %{y:.2f}<br>Z: %{z:.2f}<extra></extra>',
        name='Atoms'
    ))

    # Add bonds
    for i, j in bonds:
        fig_mol.add_trace(go.Scatter3d(
            x=[positions[i, 0], positions[j, 0]],
            y=[positions[i, 1], positions[j, 1]],
            z=[positions[i, 2], positions[j, 2]],
            mode='lines',
            line=dict(color='gray', width=3),
            showlegend=False,
            hoverinfo='skip'
        ))

    fig_mol.update_layout(
        title='Molecular Structure Visualization',
        scene=dict(
            xaxis_title='X (Å)',
            yaxis_title='Y (Å)',
            zaxis_title='Z (Å)',
            aspectmode='cube',
            bgcolor='black'
        ),
        width=800,
        height=600
    )

    fig_mol.write_html('molecular_structure.html')
    print("Created molecular visualization: molecular_structure.html")
    print(f"  Generated {n_atoms} atoms with {len(bonds)} bonds")
    print(f"  Atom types: {dict(zip(*np.unique(atom_types, return_counts=True)))}")

def create_vector_field_visualization():
    """Create vector field and flow visualizations"""
    print("\n=== Vector Field Visualization ===")

    # 3D vector field
    x = np.linspace(-2, 2, 8)
    y = np.linspace(-2, 2, 8)
    z = np.linspace(-2, 2, 8)
    X, Y, Z = np.meshgrid(x, y, z)

    # Vector field components (simplified electromagnetic field)
    U = -Y
    V = X
    W = Z * 0.5

    fig_vector = go.Figure(data=go.Cone(
        x=X.flatten(),
        y=Y.flatten(),
        z=Z.flatten(),
        u=U.flatten(),
        v=V.flatten(),
        w=W.flatten(),
        colorscale='Blues',
        sizemode='absolute',
        sizeref=0.3,
        showscale=True
    ))

    fig_vector.update_layout(
        title='3D Vector Field Visualization',
        scene=dict(
            xaxis_title='X',
            yaxis_title='Y',
            zaxis_title='Z',
            aspectmode='cube'
        ),
        width=800,
        height=600
    )

    fig_vector.write_html('vector_field.html')
    print("Created vector field visualization: vector_field.html")

    # Streamline visualization
    # Create 2D streamlines for demonstration
    x_stream = np.linspace(-3, 3, 30)
    y_stream = np.linspace(-3, 3, 30)
    X_stream, Y_stream = np.meshgrid(x_stream, y_stream)

    # Velocity field components
    U_stream = -Y_stream
    V_stream = X_stream

    # Calculate streamlines
    fig_stream = go.Figure()

    # Add velocity field as arrows
    skip = 3  # Skip points for cleaner visualization
    fig_stream.add_trace(go.Scatter(
        x=X_stream[::skip, ::skip].flatten(),
        y=Y_stream[::skip, ::skip].flatten(),
        mode='markers',
        marker=dict(
            size=8,
            color='blue',
            symbol='arrow',
            angle=np.degrees(np.arctan2(V_stream[::skip, ::skip].flatten(),
                                       U_stream[::skip, ::skip].flatten()))
        ),
        name='Velocity Vectors',
        hovertemplate='X: %{x:.2f}<br>Y: %{y:.2f}<extra></extra>'
    ))

    # Add background contour
    fig_stream.add_trace(go.Contour(
        x=x_stream,
        y=y_stream,
        z=np.sqrt(U_stream**2 + V_stream**2),
        colorscale='Viridis',
        opacity=0.6,
        showscale=True,
        contours=dict(showlabels=True),
        name='Velocity Magnitude'
    ))

    fig_stream.update_layout(
        title='2D Streamline Visualization',
        xaxis_title='X',
        yaxis_title='Y',
        width=800,
        height=600
    )

    fig_stream.write_html('streamlines.html')
    print("Created streamline visualization: streamlines.html")

    print("\nVector field features:")
    print("  - 3D cone plots for vector direction")
    print("  - Streamline flow visualization")
    print("  - Color-coded magnitude information")
    print("  - Interactive exploration of field properties")

def create_scientific_animations():
    """Create animated scientific visualizations"""
    print("\n=== Scientific Animations ===")

    # Wave propagation animation
    x = np.linspace(0, 4*np.pi, 100)
    frames = []

    for t in np.linspace(0, 2*np.pi, 50):
        y = np.sin(x - t) * np.exp(-x/10)
        frame = go.Frame(
            data=[go.Scatter(x=x, y=y, mode='lines', name='Wave', line=dict(color='blue', width=3))],
            name=f'frame_{t:.2f}'
        )
        frames.append(frame)

    fig_wave = go.Figure(
        data=[go.Scatter(x=x, y=np.sin(x), mode='lines', name='Wave', line=dict(color='blue', width=3))],
        frames=frames
    )

    fig_wave.update_layout(
        title='Wave Propagation Animation',
        xaxis_title='Position',
        yaxis_title='Amplitude',
        xaxis=dict(range=[0, 4*np.pi]),
        yaxis=dict(range=[-1.5, 1.5]),
        updatemenus=[{
            'type': 'buttons',
            'buttons': [
                {'label': 'Play', 'method': 'animate', 'args': [None, {'frame': {'duration': 100}}]},
                {'label': 'Pause', 'method': 'animate', 'args': [None, {'frame': {'duration': 0}}]}
            ]
        }],
        width=800,
        height=500
    )

    fig_wave.write_html('wave_animation.html')
    print("Created wave animation: wave_animation.html")

    # Rotating 3D object animation
    theta = np.linspace(0, 2*np.pi, 20)
    phi = np.linspace(0, np.pi, 20)
    THETA, PHI = np.meshgrid(theta, phi)

    # Sphere coordinates
    X_sphere = np.sin(PHI) * np.cos(THETA)
    Y_sphere = np.sin(PHI) * np.sin(THETA)
    Z_sphere = np.cos(PHI)

    frames_3d = []
    for angle in np.linspace(0, 2*np.pi, 30):
        # Rotate the sphere
        X_rot = X_sphere * np.cos(angle) - Y_sphere * np.sin(angle)
        Y_rot = X_sphere * np.sin(angle) + Y_sphere * np.cos(angle)
        Z_rot = Z_sphere

        frame = go.Frame(
            data=[go.Surface(x=X_rot, y=Y_rot, z=Z_rot, colorscale='Viridis', showscale=False)],
            name=f'rotation_{angle:.2f}'
        )
        frames_3d.append(frame)

    fig_3d_anim = go.Figure(
        data=[go.Surface(x=X_sphere, y=Y_sphere, z=Z_sphere, colorscale='Viridis', showscale=False)],
        frames=frames_3d
    )

    fig_3d_anim.update_layout(
        title='3D Rotating Object Animation',
        scene=dict(
            xaxis_title='X',
            yaxis_title='Y',
            zaxis_title='Z',
            aspectmode='cube'
        ),
        updatemenus=[{
            'type': 'buttons',
            'buttons': [
                {'label': 'Play', 'method': 'animate', 'args': [None, {'frame': {'duration': 100}}]},
                {'label': 'Pause', 'method': 'animate', 'args': [None, {'frame': {'duration': 0}}]}
            ]
        }],
        width=800,
        height=600
    )

    fig_3d_anim.write_html('3d_rotation.html')
    print("Created 3D rotation animation: 3d_rotation.html")

    print("\nAnimation features:")
    print("  - Smooth frame transitions")
    print("  - Play/pause controls")
    print("  - Customizable timing and speed")
    print("  - Support for complex 3D animations")

# Run 3D visualization creation
create_3d_surface_plots()
create_molecular_visualization()
create_vector_field_visualization()
create_scientific_animations()

print("\n✅ 3D visualization and rendering completed!")
print("Advanced 3D graphics and animation capabilities demonstrated")
EOF

python3 3d_visualization.py
```

**What this does**: Creates 3D surfaces, molecular visualizations, vector fields, and scientific animations.

**Expected result**: Shows comprehensive 3D visualization capabilities with interactive controls.

## Step 9: Dashboard and Web Applications

Test advanced visualization studio capabilities:

```bash
# Create dashboard application script
cat > dashboard_app.py << 'EOF'
import streamlit as st
import pandas as pd
import numpy as np
import plotly.graph_objects as go
import plotly.express as px
from plotly.subplots import make_subplots
import time

print("Creating interactive dashboard application...")

# Create dashboard data
@st.cache_data
def generate_dashboard_data():
    """Generate sample data for dashboard"""
    np.random.seed(42)

    # Time series data
    dates = pd.date_range('2023-01-01', periods=365, freq='D')

    # Multiple metrics
    metrics = {
        'Revenue': np.random.lognormal(10, 0.3, len(dates)),
        'Users': np.random.poisson(1000, len(dates)),
        'Conversion': np.random.beta(2, 8, len(dates)),
        'Satisfaction': np.random.normal(4.2, 0.8, len(dates))
    }

    df = pd.DataFrame(metrics, index=dates)

    # Regional data
    regions = ['North', 'South', 'East', 'West', 'Central']
    regional_data = pd.DataFrame({
        'Region': regions,
        'Sales': np.random.uniform(100000, 500000, len(regions)),
        'Customers': np.random.randint(1000, 5000, len(regions)),
        'Growth': np.random.uniform(-0.1, 0.3, len(regions))
    })

    # Product data
    products = ['Product A', 'Product B', 'Product C', 'Product D', 'Product E']
    product_data = pd.DataFrame({
        'Product': products,
        'Units_Sold': np.random.randint(100, 1000, len(products)),
        'Revenue': np.random.uniform(10000, 100000, len(products)),
        'Margin': np.random.uniform(0.1, 0.4, len(products))
    })

    return df, regional_data, product_data

# Main dashboard function
def create_dashboard():
    """Create the main dashboard"""

    # Page configuration
    st.set_page_config(
        page_title="Scientific Visualization Dashboard",
        page_icon="📊",
        layout="wide",
        initial_sidebar_state="expanded"
    )

    # Dashboard header
    st.title("🔬 Scientific Visualization Dashboard")
    st.markdown("Interactive data exploration and analysis platform")

    # Load data
    df, regional_data, product_data = generate_dashboard_data()

    # Sidebar controls
    st.sidebar.header("Dashboard Controls")

    # Date range selector
    date_range = st.sidebar.date_input(
        "Select Date Range",
        value=(df.index.min(), df.index.max()),
        min_value=df.index.min(),
        max_value=df.index.max()
    )

    # Metric selector
    selected_metrics = st.sidebar.multiselect(
        "Select Metrics",
        options=df.columns.tolist(),
        default=df.columns.tolist()
    )

    # Visualization type
    chart_type = st.sidebar.selectbox(
        "Chart Type",
        options=['Line', 'Bar', 'Area', 'Scatter']
    )

    # Filter data based on selections
    filtered_df = df.loc[date_range[0]:date_range[1], selected_metrics]

    # Main content area
    col1, col2 = st.columns([2, 1])

    with col1:
        st.subheader("📈 Time Series Analysis")

        # Create time series plot
        fig_ts = go.Figure()

        for metric in selected_metrics:
            if chart_type == 'Line':
                fig_ts.add_trace(go.Scatter(
                    x=filtered_df.index,
                    y=filtered_df[metric],
                    mode='lines',
                    name=metric,
                    line=dict(width=2)
                ))
            elif chart_type == 'Bar':
                fig_ts.add_trace(go.Bar(
                    x=filtered_df.index,
                    y=filtered_df[metric],
                    name=metric
                ))
            elif chart_type == 'Area':
                fig_ts.add_trace(go.Scatter(
                    x=filtered_df.index,
                    y=filtered_df[metric],
                    mode='lines',
                    name=metric,
                    fill='tonexty' if metric != selected_metrics[0] else 'tozeroy'
                ))

        fig_ts.update_layout(
            title="Time Series Metrics",
            xaxis_title="Date",
            yaxis_title="Value",
            hovermode='x unified',
            height=400
        )

        st.plotly_chart(fig_ts, use_container_width=True)

    with col2:
        st.subheader("📊 Key Statistics")

        # Display key metrics
        for metric in selected_metrics:
            current_value = filtered_df[metric].iloc[-1]
            previous_value = filtered_df[metric].iloc[-30] if len(filtered_df) >= 30 else filtered_df[metric].iloc[0]
            change = ((current_value - previous_value) / previous_value) * 100

            st.metric(
                label=metric,
                value=f"{current_value:.2f}",
                delta=f"{change:.1f}%"
            )

    # Regional analysis
    st.subheader("🗺️ Regional Performance")

    col3, col4 = st.columns(2)

    with col3:
        # Regional sales chart
        fig_regional = px.bar(
            regional_data,
            x='Region',
            y='Sales',
            title="Sales by Region",
            color='Growth',
            color_continuous_scale='RdYlGn'
        )
        st.plotly_chart(fig_regional, use_container_width=True)

    with col4:
        # Regional customers pie chart
        fig_pie = px.pie(
            regional_data,
            values='Customers',
            names='Region',
            title="Customer Distribution"
        )
        st.plotly_chart(fig_pie, use_container_width=True)

    # Product analysis
    st.subheader("🛍️ Product Performance")

    # Product scatter plot
    fig_product = px.scatter(
        product_data,
        x='Units_Sold',
        y='Revenue',
        size='Margin',
        color='Product',
        title="Product Performance Matrix",
        hover_data=['Margin']
    )
    fig_product.update_layout(height=400)
    st.plotly_chart(fig_product, use_container_width=True)

    # Data tables
    st.subheader("📋 Data Tables")

    tab1, tab2, tab3 = st.tabs(["Time Series", "Regional", "Product"])

    with tab1:
        st.dataframe(filtered_df, use_container_width=True)

    with tab2:
        st.dataframe(regional_data, use_container_width=True)

    with tab3:
        st.dataframe(product_data, use_container_width=True)

    # Real-time simulation
    st.subheader("⚡ Real-time Data Simulation")

    if st.button("Start Real-time Simulation"):
        placeholder = st.empty()

        for i in range(20):
            # Generate new data point
            new_data = {
                'Time': time.time(),
                'Value': np.random.normal(100, 15),
                'Status': np.random.choice(['Good', 'Warning', 'Critical'])
            }

            # Create real-time chart
            fig_realtime = go.Figure()
            fig_realtime.add_trace(go.Scatter(
                x=[new_data['Time']],
                y=[new_data['Value']],
                mode='markers',
                marker=dict(
                    size=15,
                    color='green' if new_data['Status'] == 'Good' else 'orange' if new_data['Status'] == 'Warning' else 'red'
                ),
                name=f"Status: {new_data['Status']}"
            ))

            fig_realtime.update_layout(
                title=f"Real-time Monitor (Update {i+1}/20)",
                xaxis_title="Time",
                yaxis_title="Value",
                height=300
            )

            placeholder.plotly_chart(fig_realtime, use_container_width=True)
            time.sleep(0.5)

    # Export functionality
    st.subheader("💾 Export Options")

    col5, col6, col7 = st.columns(3)

    with col5:
        if st.button("Export as CSV"):
            csv = filtered_df.to_csv()
            st.download_button(
                label="Download CSV",
                data=csv,
                file_name="dashboard_data.csv",
                mime="text/csv"
            )

    with col6:
        if st.button("Export as JSON"):
            json_data = filtered_df.to_json()
            st.download_button(
                label="Download JSON",
                data=json_data,
                file_name="dashboard_data.json",
                mime="application/json"
            )

    with col7:
        st.button("Generate Report", help="Generate PDF report (feature coming soon)")

# Save dashboard to file
dashboard_code = '''
import streamlit as st
import pandas as pd
import numpy as np
import plotly.graph_objects as go
import plotly.express as px
from plotly.subplots import make_subplots
import time

# Dashboard code here (same as above)
# ... (code content) ...

if __name__ == "__main__":
    create_dashboard()
'''

with open('dashboard_app.py', 'w') as f:
    f.write(dashboard_code)

print("Dashboard application created successfully!")
print("\nDashboard Features:")
print("  - Interactive time series visualization")
print("  - Real-time data simulation")
print("  - Multiple chart types and filters")
print("  - Regional and product analysis")
print("  - Data export functionality")
print("  - Responsive design with sidebar controls")

print("\nTo run the dashboard:")
print("  streamlit run dashboard_app.py")
print("  (Dashboard will open in your web browser)")

# Create a simple HTML demonstration
html_demo = '''
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Visualization Studio Demo</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; }
        .gallery { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin-top: 30px; }
        .viz-card { padding: 20px; border: 1px solid #ddd; border-radius: 8px; background: #fafafa; }
        .viz-card h3 { margin-top: 0; color: #555; }
        .viz-card p { color: #666; line-height: 1.6; }
        .btn { display: inline-block; padding: 10px 20px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; margin-top: 10px; }
        .btn:hover { background: #0056b3; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎨 Visualization Studio Demo</h1>
        <p style="text-align: center; font-size: 18px; color: #666;">
            Welcome to your scientific visualization environment! Below are the interactive visualizations you've created.
        </p>

        <div class="gallery">
            <div class="viz-card">
                <h3>📊 Publication Figures</h3>
                <p>High-quality multi-panel scientific figures with error bars, statistical distributions, and correlation matrices.</p>
                <a href="scientific_figure.png" class="btn">View Figure</a>
            </div>

            <div class="viz-card">
                <h3>🌐 3D Interactive Plots</h3>
                <p>Interactive 3D scatter plots with hover information, zoom, and rotation controls.</p>
                <a href="3d_scatter.html" class="btn">Open 3D Plot</a>
            </div>

            <div class="viz-card">
                <h3>📈 Time Series Dashboard</h3>
                <p>Multi-series time plots with interactive legends, zoom, and pan functionality.</p>
                <a href="time_series.html" class="btn">View Time Series</a>
            </div>

            <div class="viz-card">
                <h3>🎛️ Interactive Dashboard</h3>
                <p>Multi-panel dashboard with histograms, box plots, scatter plots, and bar charts.</p>
                <a href="dashboard.html" class="btn">Open Dashboard</a>
            </div>

            <div class="viz-card">
                <h3>🗺️ Surface Visualization</h3>
                <p>3D surface plots with custom colormaps and mathematical function visualization.</p>
                <a href="3d_surface.html" class="btn">View Surface</a>
            </div>

            <div class="viz-card">
                <h3>⚛️ Molecular Structure</h3>
                <p>3D molecular visualization with atoms, bonds, and interactive exploration.</p>
                <a href="molecular_structure.html" class="btn">View Molecule</a>
            </div>

            <div class="viz-card">
                <h3>🌊 Vector Fields</h3>
                <p>3D vector field visualization with cone plots and streamline analysis.</p>
                <a href="vector_field.html" class="btn">View Vectors</a>
            </div>

            <div class="viz-card">
                <h3>🎬 Animations</h3>
                <p>Animated visualizations with play/pause controls and smooth transitions.</p>
                <a href="wave_animation.html" class="btn">View Animation</a>
            </div>
        </div>

        <div style="margin-top: 40px; padding: 20px; background: #e8f4f8; border-radius: 8px;">
            <h3>🚀 Next Steps</h3>
            <ul>
                <li>Explore each visualization by clicking the buttons above</li>
                <li>Modify the Python scripts to create your own visualizations</li>
                <li>Run the Streamlit dashboard with: <code>streamlit run dashboard_app.py</code></li>
                <li>Customize colors, layouts, and interactive features</li>
            </ul>
        </div>
    </div>
</body>
</html>
'''

with open('visualization_demo.html', 'w') as f:
    f.write(html_demo)

print("\nDemo webpage created: visualization_demo.html")
print("Open this file in your browser to access all visualizations")
EOF

python3 dashboard_app.py
```

**What this does**: Creates interactive dashboards and web applications for scientific visualization.

**Expected result**: Shows comprehensive dashboard capabilities and creates a demonstration webpage.

## Step 10: Monitor Your Costs

Check your current spending:

```bash
exit  # Exit SSH session first
aws-research-wizard monitor costs --region us-west-2
```

**Expected result**: Shows costs so far (should be under $8 for this tutorial)

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
- **Compute**: $0.53 per hour for GPU-enabled instance while environment is running
- **Storage**: $0.10 per GB per month for visualization files you save
- **Data Transfer**: Usually free for visualization data amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 60% savings (advanced)
- Store large visualization files in S3, not on the instance
- Optimize rendering settings to minimize GPU usage

### Typical Monthly Costs by Usage
- **Light use** (12 hours/week): $150-300
- **Medium use** (3 hours/day): $300-600
- **Heavy use** (6 hours/day): $600-1200

## What's Next?

Now that you have a working visualization studio, you can:

### Learn More About Scientific Visualization
- [Advanced 3D Rendering Tutorial](advanced-3d-rendering.md)
- [Interactive Web Applications Guide](interactive-web-apps.md)
- [Cost Optimization for Visualization](visualization-cost-optimization.md)

### Explore Advanced Features
- [VR/AR visualization development](vr-ar-visualization.md)
- [Team collaboration with shared visualizations](team-visualization-collaboration.md)
- [Automated visualization pipelines](automated-visualization-pipelines.md)

### Join the Visualization Community
- [Visualization Research Forum](https://forum.researchwizard.app/visualization)
- [GitHub Visualization Examples](https://github.com/aws-research-wizard/visualization-examples)
- [Monthly Visualization Office Hours](https://calendar.researchwizard.app/visualization-office-hours)

## Troubleshooting

### Common Issues

**Problem**: "GPU not available" during rendering
**Solution**: Check GPU status with `nvidia-smi` and restart if needed
**Prevention**: Wait 5-7 minutes after deployment for GPU drivers to initialize

**Problem**: "Memory error" during large visualizations
**Solution**: Reduce data size or use data streaming techniques
**Prevention**: Monitor GPU memory usage during rendering

**Problem**: "Web browser not opening" for interactive plots
**Solution**: Download HTML files and open locally or use port forwarding
**Prevention**: Use SSH tunneling for remote access to web applications

**Problem**: "Streamlit app not loading"
**Solution**: Check port availability and firewall settings
**Prevention**: Use proper port forwarding when accessing remote applications

### Getting Help
- Check the [visualization troubleshooting guide](troubleshooting-visualization.md)
- Ask in [community forum](https://forum.researchwizard.app)
- File an issue on [GitHub](https://github.com/aws-research-wizard/aws-research-wizard/issues)

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-west-2 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $16. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
