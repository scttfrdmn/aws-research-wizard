# Mathematical Modeling Research Environment - Getting Started

> **Time to Complete**: 20 minutes
> **Cost**: $9-15 for tutorial
> **Skill Level**: Beginner (no cloud experience needed)

## What You'll Build

By the end of this guide, you'll have a working mathematical modeling environment that can:
- Solve differential equations and numerical optimization problems
- Run Monte Carlo simulations and stochastic modeling
- Perform linear algebra computations and matrix operations
- Handle large-scale mathematical computations and algorithms

### Meet Dr. Ahmed Hassan

Dr. Ahmed Hassan is an applied mathematician at MIT. He develops climate prediction models but waits days for supercomputer access. Each model run requires solving millions of differential equations simultaneously.

**Before**: 3-day waits + 12-hour computation = 4 days per model run
**After**: 15-minute setup + 2-hour computation = same day results
**Time Saved**: 96% faster mathematical modeling cycle
**Cost Savings**: $500/month vs $2,000 supercomputer allocation

## Before You Start

### What You Need
- [ ] AWS account (free to create)
- [ ] Credit card for AWS billing (charged only for what you use)
- [ ] Computer with internet connection
- [ ] 20 minutes of uninterrupted time

### Cost Expectations
- **Tutorial cost**: $9-15 (we'll clean up resources when done)
- **Daily research cost**: $18-45 per day when actively computing
- **Monthly estimate**: $200-550 per month for typical usage
- **Free tier**: Some compute included free for first 12 months

### Skills Needed
- Basic computer use (creating folders, installing software)
- Copy and paste commands
- No mathematics or programming experience required

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
- **Region**: Choose `us-east-1` (recommended for mathematical modeling with good CPU performance)

**What this does**: Connects the research wizard to your AWS account.

**Expected result**: "✅ AWS credentials configured successfully"

**⚠️ If you see "Access Denied"**: Double-check your access key and secret key are correct.

## Step 4: Validate Your Setup

```bash
aws-research-wizard deploy validate --domain mathematical_modeling --region us-east-1
```

**What this does**: Checks that everything is working before we spend money.

**Expected result**:
```
✅ AWS credentials valid
✅ Domain configuration valid: mathematical_modeling
✅ Region valid: us-east-1 (6 availability zones)
🎉 All validations passed!
```

## Step 5: Deploy Your Mathematical Modeling Environment

```bash
aws-research-wizard deploy start --domain mathematical_modeling --region us-east-1 --instance c6i.xlarge
```

**What this does**: Creates your mathematical modeling environment optimized for numerical computations.

**This will take**: 5-7 minutes

**Expected result**:
```
🎉 Deployment completed successfully!

Deployment Details:
  Instance ID: i-1234567890abcdef0
  Public IP: 12.34.56.78
  SSH Command: ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
  CPU: 4 cores for parallel mathematical computing
  Memory: 8GB RAM for large matrix operations
```

**💰 Billing starts now**: Your environment costs about $0.34 per hour while running.

## Step 6: Connect to Your Environment

Use the SSH command from the previous step:

```bash
ssh -i ~/.ssh/id_rsa ubuntu@12.34.56.78
```

**What this does**: Connects you to your mathematical modeling computer in the cloud.

**Expected result**: You see a command prompt like `ubuntu@ip-10-0-1-123:~$`

**⚠️ If connection fails**: Your computer might block SSH. Try adding `-o StrictHostKeyChecking=no` to the command.

## Step 7: Explore Your Mathematical Tools

Your environment comes pre-installed with:

### Core Mathematical Software
- **Python Scientific Stack**: NumPy, SciPy, SymPy - Type `python -c "import numpy; print(numpy.__version__)"` to check
- **Jupyter Notebooks**: Interactive computing - Type `jupyter --version` to check
- **BLAS/LAPACK**: Optimized linear algebra - Type `python -c "import numpy; numpy.show_config()"` to check
- **Matplotlib**: Mathematical plotting - Type `python -c "import matplotlib; print(matplotlib.__version__)"` to check
- **Pandas**: Data manipulation - Type `python -c "import pandas; print(pandas.__version__)"` to check

### Try Your First Command
```bash
python -c "import numpy; print('NumPy version:', numpy.__version__)"
```

**What this does**: Shows NumPy version and confirms mathematical computing tools are installed.

**Expected result**: You see NumPy version info confirming mathematical libraries are ready.

## Step 8: Solve Differential Equations

Let's solve mathematical problems to test everything works:

### Numerical Integration and ODEs
```bash
# Create working directory
mkdir ~/mathematical-modeling-tutorial
cd ~/mathematical-modeling-tutorial

# Create differential equations script
cat > differential_equations.py << 'EOF'
import numpy as np
import scipy.integrate as integrate
from scipy.integrate import odeint, solve_ivp
import matplotlib.pyplot as plt

print("Starting differential equation analysis...")

def solve_first_order_ode():
    """Solve first-order ordinary differential equation"""
    print("\n=== First-Order ODE Solution ===")

    # Example: dy/dt = -2y + 1, y(0) = 0
    # Analytical solution: y(t) = 0.5(1 - e^(-2t))

    def dydt(t, y):
        return -2*y + 1

    # Time points
    t_span = (0, 5)
    t_eval = np.linspace(0, 5, 100)
    y0 = [0]  # Initial condition

    # Solve ODE
    solution = solve_ivp(dydt, t_span, y0, t_eval=t_eval, method='RK45')

    # Analytical solution for comparison
    y_analytical = 0.5 * (1 - np.exp(-2 * t_eval))

    print(f"ODE: dy/dt = -2y + 1, y(0) = 0")
    print(f"Time span: {t_span}")
    print(f"Integration points: {len(t_eval)}")
    print(f"Final value (numerical): {solution.y[0][-1]:.6f}")
    print(f"Final value (analytical): {y_analytical[-1]:.6f}")
    print(f"Error: {abs(solution.y[0][-1] - y_analytical[-1]):.2e}")

    # Calculate maximum error across all points
    max_error = np.max(np.abs(solution.y[0] - y_analytical))
    print(f"Maximum error: {max_error:.2e}")

    return solution, y_analytical, t_eval

def solve_system_of_odes():
    """Solve system of coupled ODEs (Lotka-Volterra equations)"""
    print("\n=== System of ODEs: Predator-Prey Model ===")

    # Lotka-Volterra equations
    # dx/dt = ax - bxy  (prey)
    # dy/dt = -cy + dxy (predator)

    def lotka_volterra(t, z):
        x, y = z
        a, b, c, d = 1.0, 0.1, 1.5, 0.075  # Parameters

        dxdt = a*x - b*x*y
        dydt = -c*y + d*x*y

        return [dxdt, dydt]

    # Initial conditions: x0 = 10 (prey), y0 = 5 (predator)
    z0 = [10, 5]
    t_span = (0, 15)
    t_eval = np.linspace(0, 15, 1000)

    # Solve system
    solution = solve_ivp(lotka_volterra, t_span, z0, t_eval=t_eval, method='RK45')

    x_solution = solution.y[0]  # Prey population
    y_solution = solution.y[1]  # Predator population

    print(f"Lotka-Volterra predator-prey system")
    print(f"Parameters: a=1.0, b=0.1, c=1.5, d=0.075")
    print(f"Initial conditions: prey=10, predator=5")
    print(f"Time span: {t_span}")

    # Analyze oscillations
    prey_max = np.max(x_solution)
    prey_min = np.min(x_solution)
    pred_max = np.max(y_solution)
    pred_min = np.min(y_solution)

    print(f"Prey oscillation: {prey_min:.2f} to {prey_max:.2f}")
    print(f"Predator oscillation: {pred_min:.2f} to {pred_max:.2f}")

    # Find period (approximate)
    # Look for peaks in prey population
    peaks = []
    for i in range(1, len(x_solution)-1):
        if x_solution[i] > x_solution[i-1] and x_solution[i] > x_solution[i+1]:
            if x_solution[i] > 0.8 * prey_max:  # Only major peaks
                peaks.append(t_eval[i])

    if len(peaks) > 1:
        periods = np.diff(peaks)
        avg_period = np.mean(periods)
        print(f"Approximate period: {avg_period:.2f} time units")

    return solution, t_eval

def solve_boundary_value_problem():
    """Solve boundary value problem"""
    print("\n=== Boundary Value Problem ===")

    # Second-order BVP: d²y/dx² + y = 0, y(0) = 0, y(π) = 1
    # This is a simple harmonic oscillator with boundary conditions

    # Convert to system of first-order ODEs
    # Let u1 = y, u2 = dy/dx
    # Then du1/dx = u2, du2/dx = -u1

    def bvp_system(x, u):
        u1, u2 = u
        return [u2, -u1]

    def boundary_conditions(ua, ub):
        # ua = [y(0), y'(0)], ub = [y(π), y'(π)]
        return np.array([ua[0] - 0, ub[0] - 1])  # y(0) = 0, y(π) = 1

    # Set up grid
    x = np.linspace(0, np.pi, 100)

    # Initial guess for solution
    y_init = np.zeros((2, x.size))
    y_init[0] = np.sin(x)  # Guess that satisfies boundary conditions approximately

    # Solve BVP
    from scipy.integrate import solve_bvp
    solution = solve_bvp(bvp_system, boundary_conditions, x, y_init)

    print(f"Boundary Value Problem: d²y/dx² + y = 0")
    print(f"Boundary conditions: y(0) = 0, y(π) = 1")
    print(f"Grid points: {len(x)}")
    print(f"Solution at x=0: {solution.y[0][0]:.6f} (should be 0)")
    print(f"Solution at x=π: {solution.y[0][-1]:.6f} (should be 1)")

    # The analytical solution is y(x) = sin(x)/sin(π) * 1, but sin(π) = 0
    # This BVP actually has no solution with these exact boundary conditions
    # Let's check the residual
    residual_0 = abs(solution.y[0][0] - 0)
    residual_pi = abs(solution.y[0][-1] - 1)

    print(f"Boundary condition residuals:")
    print(f"  At x=0: {residual_0:.2e}")
    print(f"  At x=π: {residual_pi:.2e}")

    return solution, x

def numerical_integration():
    """Demonstrate numerical integration techniques"""
    print("\n=== Numerical Integration ===")

    # Example 1: Simple definite integral ∫₀¹ x² dx = 1/3
    def f1(x):
        return x**2

    result1, error1 = integrate.quad(f1, 0, 1)
    analytical1 = 1/3

    print(f"∫₀¹ x² dx")
    print(f"  Numerical: {result1:.8f}")
    print(f"  Analytical: {analytical1:.8f}")
    print(f"  Error: {abs(result1 - analytical1):.2e}")
    print(f"  Estimated error: {error1:.2e}")

    # Example 2: Gaussian integral ∫₋∞^∞ e^(-x²) dx = √π
    def f2(x):
        return np.exp(-x**2)

    result2, error2 = integrate.quad(f2, -np.inf, np.inf)
    analytical2 = np.sqrt(np.pi)

    print(f"\n∫₋∞^∞ e^(-x²) dx")
    print(f"  Numerical: {result2:.8f}")
    print(f"  Analytical: {analytical2:.8f}")
    print(f"  Error: {abs(result2 - analytical2):.2e}")
    print(f"  Estimated error: {error2:.2e}")

    # Example 3: Double integral ∫∫ xy dA over [0,1]×[0,1]
    def f3(y, x):
        return x * y

    result3, error3 = integrate.dblquad(f3, 0, 1, lambda x: 0, lambda x: 1)
    analytical3 = 1/4  # ∫₀¹∫₀¹ xy dx dy = 1/4

    print(f"\n∫∫ xy dA over [0,1]×[0,1]")
    print(f"  Numerical: {result3:.8f}")
    print(f"  Analytical: {analytical3:.8f}")
    print(f"  Error: {abs(result3 - analytical3):.2e}")
    print(f"  Estimated error: {error3:.2e}")

    return result1, result2, result3

# Run differential equation analysis
first_order_sol, analytical_sol, t_points = solve_first_order_ode()
system_sol, system_t = solve_system_of_odes()
bvp_sol, bvp_x = solve_boundary_value_problem()
int_results = numerical_integration()

print("\n✅ Differential equation analysis completed!")
print("Mathematical modeling environment ready for advanced computations")
EOF

python3 differential_equations.py
```

**What this does**: Solves ODEs, systems of equations, and boundary value problems with numerical integration.

**This will take**: 2-3 minutes

### Linear Algebra and Optimization
```bash
# Create linear algebra and optimization script
cat > linear_algebra_optimization.py << 'EOF'
import numpy as np
from scipy import linalg, optimize
import numpy.linalg as la

print("Starting linear algebra and optimization analysis...")

def matrix_operations():
    """Demonstrate advanced matrix operations"""
    print("\n=== Matrix Operations ===")

    # Create test matrices
    np.random.seed(42)
    A = np.random.rand(5, 5)
    B = np.random.rand(5, 5)
    C = np.random.rand(5, 3)

    print(f"Matrix dimensions: A={A.shape}, B={B.shape}, C={C.shape}")

    # Basic operations
    matrix_sum = A + B
    matrix_product = A @ B

    print(f"Matrix sum shape: {matrix_sum.shape}")
    print(f"Matrix product shape: {matrix_product.shape}")

    # Matrix properties
    det_A = la.det(A)
    trace_A = np.trace(A)
    rank_A = la.matrix_rank(A)

    print(f"Determinant of A: {det_A:.6f}")
    print(f"Trace of A: {trace_A:.6f}")
    print(f"Rank of A: {rank_A}")

    # Eigenvalue decomposition
    eigenvalues, eigenvectors = la.eig(A)

    print(f"Eigenvalues (first 3): {eigenvalues[:3]}")
    print(f"Largest eigenvalue: {np.max(np.real(eigenvalues)):.6f}")
    print(f"Condition number: {la.cond(A):.2f}")

    # SVD decomposition
    U, s, Vt = la.svd(A)

    print(f"SVD shapes: U={U.shape}, s={s.shape}, Vt={Vt.shape}")
    print(f"Singular values: {s}")
    print(f"Matrix rank (SVD): {np.sum(s > 1e-10)}")

    # Reconstruction check
    A_reconstructed = U @ np.diag(s) @ Vt
    reconstruction_error = la.norm(A - A_reconstructed)

    print(f"SVD reconstruction error: {reconstruction_error:.2e}")

    return A, eigenvalues, s

def solve_linear_systems():
    """Solve various linear systems"""
    print("\n=== Linear System Solutions ===")

    # Well-conditioned system
    np.random.seed(42)
    n = 100
    A = np.random.rand(n, n) + n * np.eye(n)  # Diagonally dominant
    x_true = np.random.rand(n)
    b = A @ x_true

    # Solve using different methods
    x_solve = la.solve(A, b)
    x_lstsq = la.lstsq(A, b, rcond=None)[0]

    # Calculate errors
    error_solve = la.norm(x_solve - x_true)
    error_lstsq = la.norm(x_lstsq - x_true)

    print(f"Linear system size: {n}×{n}")
    print(f"Condition number: {la.cond(A):.2f}")
    print(f"Direct solver error: {error_solve:.2e}")
    print(f"Least squares error: {error_lstsq:.2e}")

    # Ill-conditioned system (Hilbert matrix)
    n_hilbert = 8
    H = np.array([[1/(i+j+1) for j in range(n_hilbert)] for i in range(n_hilbert)])

    x_true_hilbert = np.ones(n_hilbert)
    b_hilbert = H @ x_true_hilbert

    try:
        x_hilbert = la.solve(H, b_hilbert)
        error_hilbert = la.norm(x_hilbert - x_true_hilbert)

        print(f"\nHilbert matrix {n_hilbert}×{n_hilbert}:")
        print(f"Condition number: {la.cond(H):.2e}")
        print(f"Solution error: {error_hilbert:.2e}")

    except la.LinAlgError:
        print(f"\nHilbert matrix {n_hilbert}×{n_hilbert}: Singular (cannot solve)")

    # Overdetermined system
    m, n = 200, 100
    A_over = np.random.rand(m, n)
    x_true_over = np.random.rand(n)
    b_over = A_over @ x_true_over + 0.1 * np.random.rand(m)  # Add noise

    x_over, residuals, rank, s = la.lstsq(A_over, b_over, rcond=None)
    error_over = la.norm(x_over - x_true_over)

    print(f"\nOverdetermined system {m}×{n}:")
    print(f"Rank: {rank}")
    print(f"Residual norm: {residuals[0]:.6f}")
    print(f"Solution error: {error_over:.6f}")

    return A, H, A_over

def optimization_problems():
    """Solve various optimization problems"""
    print("\n=== Optimization Problems ===")

    # 1. Minimize a simple quadratic function
    def quadratic(x):
        return (x[0] - 2)**2 + (x[1] - 1)**2

    def quadratic_grad(x):
        return np.array([2*(x[0] - 2), 2*(x[1] - 1)])

    # Starting point
    x0 = np.array([0, 0])

    # Minimize using different methods
    result_bfgs = optimize.minimize(quadratic, x0, method='BFGS', jac=quadratic_grad)
    result_nelder = optimize.minimize(quadratic, x0, method='Nelder-Mead')

    print(f"Quadratic optimization f(x,y) = (x-2)² + (y-1)²")
    print(f"True minimum: [2, 1], f = 0")
    print(f"BFGS result: {result_bfgs.x}, f = {result_bfgs.fun:.6f}")
    print(f"Nelder-Mead result: {result_nelder.x}, f = {result_nelder.fun:.6f}")
    print(f"BFGS iterations: {result_bfgs.nit}")
    print(f"Nelder-Mead iterations: {result_nelder.nit}")

    # 2. Constrained optimization
    def objective(x):
        return x[0]**2 + x[1]**2

    def constraint1(x):
        return x[0] + x[1] - 1  # x + y = 1

    def constraint2(x):
        return x[0] - 2*x[1] + 1  # x - 2y >= -1

    constraints = [
        {'type': 'eq', 'fun': constraint1},
        {'type': 'ineq', 'fun': constraint2}
    ]

    x0_constrained = np.array([0.5, 0.5])
    result_constrained = optimize.minimize(objective, x0_constrained,
                                         method='SLSQP', constraints=constraints)

    print(f"\nConstrained optimization: minimize x² + y²")
    print(f"Subject to: x + y = 1, x - 2y ≥ -1")
    print(f"Result: {result_constrained.x}")
    print(f"Objective value: {result_constrained.fun:.6f}")
    print(f"Constraint satisfaction:")
    print(f"  x + y - 1 = {constraint1(result_constrained.x):.6f}")
    print(f"  x - 2y + 1 = {constraint2(result_constrained.x):.6f}")

    # 3. Global optimization (find global minimum of a multi-modal function)
    def rastrigin(x):
        A = 10
        n = len(x)
        return A * n + sum(xi**2 - A * np.cos(2 * np.pi * xi) for xi in x)

    # Rastrigin function has global minimum at origin
    bounds = [(-5.12, 5.12), (-5.12, 5.12)]

    # Try differential evolution for global optimization
    result_global = optimize.differential_evolution(rastrigin, bounds, seed=42)

    print(f"\nGlobal optimization: Rastrigin function")
    print(f"True global minimum: [0, 0], f = 0")
    print(f"Differential evolution result: {result_global.x}")
    print(f"Function value: {result_global.fun:.6f}")
    print(f"Function evaluations: {result_global.nfev}")

    # 4. Root finding
    def equation_system(x):
        return [x[0]**2 + x[1]**2 - 1,  # Circle: x² + y² = 1
                x[0] - x[1]]             # Line: x = y

    x0_root = np.array([0.5, 0.5])
    root_result = optimize.fsolve(equation_system, x0_root)

    print(f"\nRoot finding: intersection of x² + y² = 1 and x = y")
    print(f"Solution: {root_result}")
    print(f"Residual: {equation_system(root_result)}")

    return result_bfgs, result_constrained, result_global

def numerical_methods():
    """Demonstrate various numerical methods"""
    print("\n=== Numerical Methods ===")

    # 1. Interpolation
    x_data = np.array([0, 1, 2, 3, 4])
    y_data = np.array([1, 3, 2, 5, 4])

    from scipy import interpolate

    # Different interpolation methods
    f_linear = interpolate.interp1d(x_data, y_data, kind='linear')
    f_cubic = interpolate.interp1d(x_data, y_data, kind='cubic')

    x_new = np.linspace(0, 4, 20)
    y_linear = f_linear(x_new)
    y_cubic = f_cubic(x_new)

    print(f"Interpolation with {len(x_data)} data points")
    print(f"Evaluation at {len(x_new)} new points")
    print(f"Linear interpolation range: [{y_linear.min():.3f}, {y_linear.max():.3f}]")
    print(f"Cubic interpolation range: [{y_cubic.min():.3f}, {y_cubic.max():.3f}]")

    # 2. Numerical differentiation
    def f_diff(x):
        return x**3 - 2*x**2 + x + 1

    def f_diff_analytical(x):
        return 3*x**2 - 4*x + 1

    # Central difference approximation
    h = 1e-5
    x_test = 2.0

    derivative_numerical = (f_diff(x_test + h) - f_diff(x_test - h)) / (2 * h)
    derivative_analytical = f_diff_analytical(x_test)

    print(f"\nNumerical differentiation at x = {x_test}")
    print(f"Function: f(x) = x³ - 2x² + x + 1")
    print(f"Analytical derivative: {derivative_analytical}")
    print(f"Numerical derivative: {derivative_numerical:.8f}")
    print(f"Error: {abs(derivative_numerical - derivative_analytical):.2e}")

    # 3. Fourier transform
    N = 1024
    t = np.linspace(0, 1, N, endpoint=False)

    # Signal: sum of sinusoids
    freq1, freq2 = 50, 120  # Hz
    signal = np.sin(2 * np.pi * freq1 * t) + 0.5 * np.sin(2 * np.pi * freq2 * t)

    # Add noise
    signal_noisy = signal + 0.2 * np.random.rand(N)

    # FFT
    fft_result = np.fft.fft(signal_noisy)
    frequencies = np.fft.fftfreq(N, t[1] - t[0])

    # Find peaks in frequency domain
    magnitude = np.abs(fft_result)
    peak_indices = np.where(magnitude > 0.1 * np.max(magnitude))[0]
    detected_frequencies = np.abs(frequencies[peak_indices])
    detected_frequencies = detected_frequencies[detected_frequencies > 0]
    detected_frequencies = np.sort(detected_frequencies)[:2]  # Take first two

    print(f"\nFourier analysis of signal with frequencies {freq1} Hz and {freq2} Hz")
    print(f"Signal length: {N} samples")
    print(f"Detected frequencies: {detected_frequencies}")
    print(f"Detection error: {np.abs(detected_frequencies - [freq1, freq2])}")

    return x_new, y_linear, y_cubic, signal, fft_result

# Run linear algebra and optimization analysis
matrix_A, eigenvals, singular_vals = matrix_operations()
matrices = solve_linear_systems()
opt_results = optimization_problems()
numerical_results = numerical_methods()

print("\n✅ Linear algebra and optimization analysis completed!")
print("Advanced mathematical computation capabilities demonstrated")
EOF

python3 linear_algebra_optimization.py
```

**What this does**: Demonstrates matrix operations, linear systems, optimization, and numerical methods.

**Expected result**: Shows comprehensive mathematical computing results and algorithm performance.

## Step 9: Monte Carlo and Stochastic Modeling

Test advanced mathematical modeling capabilities:

```bash
# Create Monte Carlo and stochastic modeling script
cat > monte_carlo_stochastic.py << 'EOF'
import numpy as np
import pandas as pd
from scipy import stats

print("Starting Monte Carlo and stochastic modeling...")

def monte_carlo_integration():
    """Use Monte Carlo methods for numerical integration"""
    print("\n=== Monte Carlo Integration ===")

    np.random.seed(42)

    # Example 1: Estimate π using Monte Carlo
    def estimate_pi(n_samples):
        # Generate random points in [0,1] × [0,1]
        x = np.random.uniform(0, 1, n_samples)
        y = np.random.uniform(0, 1, n_samples)

        # Count points inside unit circle
        inside_circle = (x**2 + y**2) <= 1
        pi_estimate = 4 * np.sum(inside_circle) / n_samples

        return pi_estimate

    sample_sizes = [1000, 10000, 100000, 1000000]
    pi_estimates = []

    print("Estimating π using Monte Carlo:")
    for n in sample_sizes:
        pi_est = estimate_pi(n)
        pi_estimates.append(pi_est)
        error = abs(pi_est - np.pi)
        print(f"  n = {n:7d}: π ≈ {pi_est:.6f}, error = {error:.6f}")

    # Example 2: High-dimensional integration
    # Estimate ∫∫∫ e^(-(x²+y²+z²)) dx dy dz over [-1,1]³
    def high_dim_integral(n_samples):
        # Generate random points in [-1,1]³
        points = np.random.uniform(-1, 1, (n_samples, 3))

        # Evaluate function at random points
        values = np.exp(-np.sum(points**2, axis=1))

        # Monte Carlo estimate
        volume = 8  # Volume of [-1,1]³
        integral_estimate = volume * np.mean(values)

        return integral_estimate

    n_samples = 100000
    integral_result = high_dim_integral(n_samples)

    # Analytical result: (√π * erf(1))³ ≈ 5.568
    analytical_result = (np.sqrt(np.pi) * 2 * stats.norm.cdf(np.sqrt(2)) - np.sqrt(np.pi))**3

    print(f"\nHigh-dimensional integration:")
    print(f"Monte Carlo estimate: {integral_result:.6f}")
    print(f"Analytical result: {analytical_result:.6f}")
    print(f"Error: {abs(integral_result - analytical_result):.6f}")

    return pi_estimates, integral_result

def stochastic_processes():
    """Model various stochastic processes"""
    print("\n=== Stochastic Processes ===")

    np.random.seed(42)

    # 1. Random Walk
    n_steps = 1000
    steps = np.random.choice([-1, 1], n_steps)
    position = np.cumsum(steps)

    print(f"Random Walk ({n_steps} steps):")
    print(f"  Final position: {position[-1]}")
    print(f"  Maximum position: {np.max(position)}")
    print(f"  Minimum position: {np.min(position)}")
    print(f"  Standard deviation: {np.std(position):.2f}")
    print(f"  Theoretical std: {np.sqrt(n_steps):.2f}")

    # 2. Brownian Motion (Wiener Process)
    dt = 0.01
    T = 10.0
    n_points = int(T / dt)
    t = np.linspace(0, T, n_points)

    dW = np.random.normal(0, np.sqrt(dt), n_points-1)
    W = np.concatenate([[0], np.cumsum(dW)])

    print(f"\nBrownian Motion (T = {T}, dt = {dt}):")
    print(f"  Final value: {W[-1]:.4f}")
    print(f"  Theoretical variance: {T:.2f}")
    print(f"  Observed variance: {np.var(W):.4f}")
    print(f"  Maximum value: {np.max(W):.4f}")
    print(f"  Minimum value: {np.min(W):.4f}")

    # 3. Geometric Brownian Motion (Stock Price Model)
    # dS = μS dt + σS dW
    S0 = 100  # Initial stock price
    mu = 0.05  # Drift (expected return)
    sigma = 0.2  # Volatility

    # Exact solution: S(t) = S0 * exp((μ - σ²/2)t + σW(t))
    S = S0 * np.exp((mu - sigma**2/2) * t + sigma * W)

    print(f"\nGeometric Brownian Motion (Stock Price):")
    print(f"  Initial price: ${S0}")
    print(f"  Parameters: μ = {mu:.3f}, σ = {sigma:.3f}")
    print(f"  Final price: ${S[-1]:.2f}")
    print(f"  Maximum price: ${np.max(S):.2f}")
    print(f"  Minimum price: ${np.min(S):.2f}")
    print(f"  Total return: {(S[-1]/S0 - 1)*100:.2f}%")

    # 4. Ornstein-Uhlenbeck Process (Mean Reverting)
    # dX = θ(μ - X)dt + σ dW
    theta = 0.5  # Mean reversion speed
    mu_ou = 0.0  # Long-term mean
    sigma_ou = 1.0  # Volatility

    X = np.zeros(n_points)
    for i in range(1, n_points):
        dX = theta * (mu_ou - X[i-1]) * dt + sigma_ou * dW[i-1]
        X[i] = X[i-1] + dX

    print(f"\nOrnstein-Uhlenbeck Process:")
    print(f"  Parameters: θ = {theta}, μ = {mu_ou}, σ = {sigma_ou}")
    print(f"  Final value: {X[-1]:.4f}")
    print(f"  Mean value: {np.mean(X):.4f}")
    print(f"  Theoretical variance: {sigma_ou**2/(2*theta):.4f}")
    print(f"  Observed variance: {np.var(X):.4f}")

    return position, W, S, X, t

def monte_carlo_option_pricing():
    """Price financial options using Monte Carlo"""
    print("\n=== Monte Carlo Option Pricing ===")

    np.random.seed(42)

    # Black-Scholes parameters
    S0 = 100    # Initial stock price
    K = 105     # Strike price
    T = 1.0     # Time to maturity (1 year)
    r = 0.05    # Risk-free rate
    sigma = 0.2 # Volatility

    n_simulations = 100000

    # Generate stock price paths
    Z = np.random.standard_normal(n_simulations)
    ST = S0 * np.exp((r - 0.5 * sigma**2) * T + sigma * np.sqrt(T) * Z)

    # European Call Option
    call_payoffs = np.maximum(ST - K, 0)
    call_price_mc = np.exp(-r * T) * np.mean(call_payoffs)
    call_std_error = np.exp(-r * T) * np.std(call_payoffs) / np.sqrt(n_simulations)

    # European Put Option
    put_payoffs = np.maximum(K - ST, 0)
    put_price_mc = np.exp(-r * T) * np.mean(put_payoffs)
    put_std_error = np.exp(-r * T) * np.std(put_payoffs) / np.sqrt(n_simulations)

    # Black-Scholes analytical prices for comparison
    from scipy.stats import norm

    d1 = (np.log(S0/K) + (r + 0.5*sigma**2)*T) / (sigma*np.sqrt(T))
    d2 = d1 - sigma*np.sqrt(T)

    call_price_bs = S0*norm.cdf(d1) - K*np.exp(-r*T)*norm.cdf(d2)
    put_price_bs = K*np.exp(-r*T)*norm.cdf(-d2) - S0*norm.cdf(-d1)

    print(f"European Option Pricing:")
    print(f"  Parameters: S0=${S0}, K=${K}, T={T}, r={r:.3f}, σ={sigma:.3f}")
    print(f"  Simulations: {n_simulations:,}")

    print(f"\nCall Option:")
    print(f"  Monte Carlo: ${call_price_mc:.4f} ± ${call_std_error:.4f}")
    print(f"  Black-Scholes: ${call_price_bs:.4f}")
    print(f"  Error: ${abs(call_price_mc - call_price_bs):.4f}")

    print(f"\nPut Option:")
    print(f"  Monte Carlo: ${put_price_mc:.4f} ± ${put_std_error:.4f}")
    print(f"  Black-Scholes: ${put_price_bs:.4f}")
    print(f"  Error: ${abs(put_price_mc - put_price_bs):.4f}")

    # Asian Option (Path-dependent)
    n_time_steps = 252  # Daily observations for 1 year
    dt = T / n_time_steps

    # Generate multiple paths
    n_paths = 10000

    asian_payoffs = []
    for _ in range(n_paths):
        # Generate one path
        dW = np.random.normal(0, np.sqrt(dt), n_time_steps)
        W = np.cumsum(dW)
        S_path = S0 * np.exp((r - 0.5*sigma**2)*np.arange(1, n_time_steps+1)*dt + sigma*W)

        # Asian option payoff (arithmetic average)
        average_price = np.mean(S_path)
        asian_payoffs.append(max(average_price - K, 0))

    asian_price = np.exp(-r * T) * np.mean(asian_payoffs)
    asian_std_error = np.exp(-r * T) * np.std(asian_payoffs) / np.sqrt(n_paths)

    print(f"\nAsian Call Option (Arithmetic Average):")
    print(f"  Price: ${asian_price:.4f} ± ${asian_std_error:.4f}")
    print(f"  Paths: {n_paths}")
    print(f"  Time steps per path: {n_time_steps}")

    return call_price_mc, put_price_mc, asian_price

def markov_chain_simulation():
    """Simulate discrete-time Markov chains"""
    print("\n=== Markov Chain Simulation ===")

    np.random.seed(42)

    # Weather model: Sunny, Cloudy, Rainy
    states = ['Sunny', 'Cloudy', 'Rainy']
    n_states = len(states)

    # Transition matrix
    # P[i,j] = probability of transitioning from state i to state j
    P = np.array([
        [0.7, 0.2, 0.1],  # From Sunny
        [0.3, 0.5, 0.2],  # From Cloudy
        [0.2, 0.3, 0.5]   # From Rainy
    ])

    print(f"Weather Markov Chain:")
    print(f"States: {states}")
    print(f"Transition Matrix:")
    for i, state in enumerate(states):
        print(f"  {state:6}: {P[i]}")

    # Simulate chain
    n_days = 1000
    current_state = 0  # Start with Sunny
    state_sequence = [current_state]

    for day in range(n_days - 1):
        # Choose next state based on transition probabilities
        next_state = np.random.choice(n_states, p=P[current_state])
        state_sequence.append(next_state)
        current_state = next_state

    # Analyze simulation
    state_counts = np.bincount(state_sequence, minlength=n_states)
    state_frequencies = state_counts / n_days

    print(f"\nSimulation Results ({n_days} days):")
    for i, state in enumerate(states):
        print(f"  {state}: {state_counts[i]} days ({state_frequencies[i]:.3f})")

    # Steady-state distribution (theoretical)
    eigenvals, eigenvecs = np.linalg.eig(P.T)
    steady_state_idx = np.argmax(np.real(eigenvals))
    steady_state = np.real(eigenvecs[:, steady_state_idx])
    steady_state = steady_state / np.sum(steady_state)

    print(f"\nSteady-State Distribution:")
    for i, state in enumerate(states):
        print(f"  {state}: {steady_state[i]:.3f} (observed: {state_frequencies[i]:.3f})")

    # Calculate expected return time to each state
    fundamental_matrix = np.linalg.inv(np.eye(n_states) - P + np.ones((n_states, n_states)) / n_states)
    return_times = np.diag(fundamental_matrix) / steady_state

    print(f"\nExpected Return Times:")
    for i, state in enumerate(states):
        print(f"  {state}: {return_times[i]:.2f} days")

    return state_sequence, P, steady_state

def queueing_theory_simulation():
    """Simulate queueing systems"""
    print("\n=== Queueing Theory Simulation ===")

    np.random.seed(42)

    # M/M/1 queue: Poisson arrivals, exponential service, 1 server
    arrival_rate = 0.8  # λ (customers per time unit)
    service_rate = 1.0  # μ (customers served per time unit)
    simulation_time = 1000.0

    # Generate arrival times
    n_arrivals = np.random.poisson(arrival_rate * simulation_time)
    arrival_times = np.sort(np.random.uniform(0, simulation_time, n_arrivals))

    # Generate service times
    service_times = np.random.exponential(1/service_rate, n_arrivals)

    # Simulate queue
    departure_times = np.zeros(n_arrivals)
    wait_times = np.zeros(n_arrivals)

    departure_times[0] = arrival_times[0] + service_times[0]
    wait_times[0] = 0

    for i in range(1, n_arrivals):
        # Customer waits if previous customer hasn't finished
        service_start = max(arrival_times[i], departure_times[i-1])
        wait_times[i] = service_start - arrival_times[i]
        departure_times[i] = service_start + service_times[i]

    # Calculate metrics
    avg_wait_time = np.mean(wait_times)
    avg_system_time = np.mean(departure_times - arrival_times)
    utilization = arrival_rate / service_rate

    # Theoretical results (M/M/1)
    theoretical_wait = (arrival_rate / service_rate) / (service_rate - arrival_rate)
    theoretical_system = 1 / (service_rate - arrival_rate)

    print(f"M/M/1 Queue Simulation:")
    print(f"  Parameters: λ = {arrival_rate}, μ = {service_rate}")
    print(f"  Utilization (ρ): {utilization:.3f}")
    print(f"  Simulation time: {simulation_time}")
    print(f"  Total arrivals: {n_arrivals}")

    print(f"\nSimulation Results:")
    print(f"  Average wait time: {avg_wait_time:.4f}")
    print(f"  Average system time: {avg_system_time:.4f}")

    print(f"\nTheoretical Results:")
    print(f"  Average wait time: {theoretical_wait:.4f}")
    print(f"  Average system time: {theoretical_system:.4f}")

    print(f"\nErrors:")
    print(f"  Wait time error: {abs(avg_wait_time - theoretical_wait):.4f}")
    print(f"  System time error: {abs(avg_system_time - theoretical_system):.4f}")

    return arrival_times, departure_times, wait_times

# Run Monte Carlo and stochastic modeling
pi_est, integral_est = monte_carlo_integration()
random_walk, brownian, stock_prices, ou_process, time_points = stochastic_processes()
call_mc, put_mc, asian_mc = monte_carlo_option_pricing()
markov_sequence, transition_matrix, steady_dist = markov_chain_simulation()
arrivals, departures, waits = queueing_theory_simulation()

print("\n✅ Monte Carlo and stochastic modeling completed!")
print("Advanced probabilistic and stochastic analysis capabilities demonstrated")
EOF

python3 monte_carlo_stochastic.py
```

**What this does**: Demonstrates Monte Carlo integration, stochastic processes, option pricing, and queueing theory.

**Expected result**: Shows comprehensive stochastic modeling and probabilistic analysis results.

## Step 10: Monitor Your Costs

Check your current spending:

```bash
exit  # Exit SSH session first
aws-research-wizard monitor costs --region us-east-1
```

**Expected result**: Shows costs so far (should be under $6 for this tutorial)

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
- **Compute**: $0.34 per hour for computing instance while environment is running
- **Storage**: $0.10 per GB per month for mathematical data you save
- **Data Transfer**: Usually free for mathematical modeling amounts

### Cost Control Tips
- Always delete environments when not needed
- Use spot instances for 60% savings (advanced)
- Store large datasets in S3, not on the instance
- Optimize algorithms to minimize compute time

### Typical Monthly Costs by Usage
- **Light use** (15 hours/week): $125-250
- **Medium use** (4 hours/day): $250-500
- **Heavy use** (8 hours/day): $500-1000

## What's Next?

Now that you have a working mathematical modeling environment, you can:

### Learn More About Mathematical Computing
- [High-Performance Computing Tutorial](high-performance-computing.md)
- [Advanced Numerical Methods Guide](advanced-numerical-methods.md)
- [Cost Optimization for Mathematical Modeling](mathematical-cost-optimization.md)

### Explore Advanced Features
- [Parallel computing with multiple cores](parallel-mathematical-computing.md)
- [Team collaboration with computational notebooks](team-mathematical-collaboration.md)
- [Automated mathematical pipelines](automated-mathematical-pipelines.md)

### Join the Mathematical Modeling Community
- [Mathematical Modeling Forum](https://forum.researchwizard.app/mathematical-modeling)
- [GitHub Mathematical Examples](https://github.com/aws-research-wizard/mathematical-examples)
- [Monthly Mathematical Computing Office Hours](https://calendar.researchwizard.app/mathematical-office-hours)

## Troubleshooting

### Common Issues

**Problem**: "Memory error" during large matrix operations
**Solution**: Use a larger instance type or implement chunked processing
**Prevention**: Monitor memory usage with `htop` during computations

**Problem**: "Convergence error" in optimization algorithms
**Solution**: Try different starting points or optimization methods
**Prevention**: Check function properties and use appropriate solvers

**Problem**: "Numerical instability" in differential equations
**Solution**: Reduce step size or use more stable integration methods
**Prevention**: Analyze equation stability and choose appropriate solvers

**Problem**: "Linear algebra error" for singular matrices
**Solution**: Check matrix condition number and use regularization if needed
**Prevention**: Verify matrix properties before attempting operations

### Getting Help
- Check the [mathematical modeling troubleshooting guide](troubleshooting-mathematical.md)
- Ask in [community forum](https://forum.researchwizard.app)
- File an issue on [GitHub](https://github.com/aws-research-wizard/aws-research-wizard/issues)

### Emergency: Stop All Billing
If something goes wrong and you want to stop all charges immediately:
```bash
aws-research-wizard emergency-stop --region us-east-1 --confirm
```

## Feedback

This guide should take 20 minutes and cost under $15. Help us improve:

**Was this guide helpful?** [Yes/No feedback buttons]

**What was confusing?** [Text box for feedback]

**What would you add?** [Text box for suggestions]

**Rate the clarity (1-5)**: ⭐⭐⭐⭐⭐

---

*Last updated: January 2025 | Reading level: 8th grade | Tutorial tested: January 15, 2025*
