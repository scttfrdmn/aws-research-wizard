# 🚀 Spack Manager Go

**PROPRIETARY SOFTWARE - CONFIDENTIAL**

A modern, interactive Go-based interface for managing Spack package installations with real-time progress tracking and beautiful Terminal UI.

> **NOTICE**: This is proprietary software owned by Scott Friedman. This software contains confidential and proprietary information. Unauthorized copying, distribution, or use is strictly prohibited.

## ✨ Features

- **Interactive Terminal UI (TUI)** - Beautiful, responsive interface built with Bubble Tea
- **Real-time Progress Tracking** - Monitor package installations with live progress updates
- **Environment Management** - Create, list, and manage Spack environments
- **AWS Binary Cache Support** - 95% faster installations with AWS binary cache integration
- **Command Line Interface** - Full CLI for automation and scripting
- **Comprehensive API** - Clean Go interfaces for integration into other projects

## 🛠 Installation

### Prerequisites

- Go 1.21 or later
- Spack installation (see [Spack Documentation](https://spack.readthedocs.io/))

### Install from Source

```bash
git clone https://github.com/spack-go/spack-manager.git
cd spack-manager
go build -o spack-manager ./cmd/spack-manager
sudo mv spack-manager /usr/local/bin/
```

### Go Module

```bash
go get github.com/spack-go/spack-manager
```

## 🚀 Quick Start

### Interactive TUI (Recommended)

Launch the beautiful Terminal User Interface:

```bash
spack-manager tui
```

### Command Line Usage

```bash
# Create a new environment
spack-manager env create genomics

# Add packages to environment
spack-manager env add genomics gcc@11.3.0
spack-manager env add genomics python@3.11
spack-manager env add genomics numpy

# Install environment with progress tracking
spack-manager install genomics

# List all environments
spack-manager list

# Get detailed environment info
spack-manager env info genomics
```

## 📖 API Usage

### Basic Manager Setup

```go
package main

import (
    "fmt"
    "log"

    "github.com/spack-go/spack-manager/pkg/manager"
)

func main() {
    // Create Spack manager
    config := manager.Config{
        SpackRoot:   "/opt/spack",
        BinaryCache: "https://cache.spack.io", // Optional: Use AWS binary cache
        WorkDir:     "/tmp/spack-manager",
        LogLevel:    "info",
    }

    sm, err := manager.New(config)
    if err != nil {
        log.Fatal(err)
    }

    // Create environment
    env := manager.Environment{
        Name: "ml-project",
        Packages: []string{
            "gcc@11.3.0",
            "python@3.11",
            "pytorch",
            "numpy",
        },
    }

    if err := sm.CreateEnvironment(env); err != nil {
        log.Fatal(err)
    }

    // Install with progress tracking
    progressChan := make(chan manager.ProgressUpdate, 100)
    go func() {
        for update := range progressChan {
            fmt.Printf("📦 %s: %.1f%% - %s\n",
                update.Package, update.Progress*100, update.Message)
        }
    }()

    if err := sm.InstallEnvironment(env.Name, progressChan); err != nil {
        log.Fatal(err)
    }

    fmt.Println("🎉 Environment installed successfully!")
}
```

### Launch TUI Programmatically

```go
package main

import (
    "log"

    "github.com/spack-go/spack-manager/pkg/manager"
    "github.com/spack-go/spack-manager/pkg/tui"
)

func main() {
    config := manager.Config{
        SpackRoot:   "/opt/spack",
        BinaryCache: "https://cache.spack.io",
    }

    sm, err := manager.New(config)
    if err != nil {
        log.Fatal(err)
    }

    // Launch interactive TUI
    if err := tui.Run(sm); err != nil {
        log.Fatal(err)
    }
}
```

## 🎨 TUI Features

The Terminal User Interface provides:

- **Environment Browser** - Navigate through Spack environments
- **Package Details** - View installed packages and versions
- **Live Installation Progress** - Real-time progress bars and status updates
- **Interactive Help** - Press `?` for comprehensive help
- **Keyboard Navigation** - Intuitive vim-like controls

### TUI Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/k` | Move up |
| `↓/j` | Move down |
| `Enter` | Select/View details |
| `i` | Install environment |
| `r` | Refresh data |
| `Esc` | Go back |
| `?` | Toggle help |
| `q` | Quit |

## ⚙️ Configuration

### Environment Variables

```bash
# Spack installation directory
export SPACK_ROOT="/opt/spack"

# AWS binary cache for faster installations
export SPACK_BINARY_CACHE="https://cache.spack.io"
```

### Config Options

```go
type Config struct {
    SpackRoot   string // Path to Spack installation
    BinaryCache string // URL to binary cache (optional)
    WorkDir     string // Working directory for temporary files
    LogLevel    string // Logging level: debug, info, warn, error
}
```

## 🏗️ Architecture

```
spack-manager-go/
├── cmd/spack-manager/     # CLI application
├── pkg/manager/          # Core Spack management logic
├── pkg/tui/             # Terminal User Interface
├── examples/            # Usage examples
└── docs/               # Additional documentation
```

### Core Components

- **Manager** - Core Spack interaction layer with environment management
- **TUI** - Beautiful Terminal User Interface built with Bubble Tea
- **Progress Tracking** - Real-time installation monitoring
- **CLI** - Command-line interface for automation

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/manager
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📋 Requirements

- **Go**: 1.21 or later
- **Spack**: Latest stable version
- **Terminal**: Modern terminal with color support for best TUI experience

## 🐛 Troubleshooting

### Common Issues

**Spack not found:**
```bash
# Set SPACK_ROOT environment variable
export SPACK_ROOT="/path/to/spack"
```

**Permission denied:**
```bash
# Ensure Spack binary is executable
chmod +x $SPACK_ROOT/bin/spack
```

**Binary cache errors:**
```bash
# Test cache connectivity
curl -I https://cache.spack.io
```

## 📄 License

This project is proprietary software owned by Scott Friedman. See the [LICENSE](LICENSE) file for complete terms and restrictions.

**All rights reserved. No public license is granted.**

## 🙏 Acknowledgments

- [Spack Community](https://spack.io) - The amazing package manager this tool builds upon
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - For the beautiful TUI framework
- [AWS Research Wizard](https://github.com/aws-research-wizard) - Original project that inspired this library

## 🔗 Links

- [Spack Documentation](https://spack.readthedocs.io/)
- [Go Documentation](https://golang.org/doc/)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)

---

**Made with ❤️ for the Spack and Go communities**
