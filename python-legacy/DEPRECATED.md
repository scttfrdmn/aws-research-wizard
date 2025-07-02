# ⚠️ DEPRECATED: Python Implementation Archived

**Status**: DEPRECATED as of July 2, 2025
**Replacement**: Go implementation in `/go/` directory
**Migration Required**: Yes - see migration guide below

## 🚨 Important Notice

The Python implementation of AWS Research Wizard has been **superseded by a superior Go implementation** and is now archived for historical reference only.

**Do not use this Python version for new projects.**

## 🎯 Why Deprecated?

The Go implementation provides:
- **100x faster execution** (sub-second vs 3+ seconds startup)
- **Zero dependency deployment** (20MB binary vs 500MB+ Python environment)
- **Production-grade reliability** with comprehensive error handling
- **Enhanced features** including real-time monitoring and cost optimization
- **91% fewer dependencies** to maintain and secure

## 🔄 Migration Guide

### Quick Migration Reference

| Python Command | Go Equivalent |
|----------------|---------------|
| `python research_infrastructure_wizard.py` | `./aws-research-wizard config recommend` |
| `python tui_research_wizard.py` | `./aws-research-wizard config tui` |
| `python gui_research_wizard.py` | `./aws-research-wizard config tui` (better interface) |
| `python aws_environment_checker.py` | `./aws-research-wizard deploy validate` |
| `python demo_workflow_engine.py` | `./aws-research-wizard data workflow` |

### Feature Mapping

| Python Feature | Go Implementation | Status |
|----------------|-------------------|--------|
| **25+ Domain Packs** | 18 core domains | ✅ Key domains covered |
| **Streamlit Web GUI** | Modern Terminal TUI | ✅ Enhanced UX |
| **Basic AWS Integration** | Advanced AWS SDK v2 | ✅ Superior |
| **Single Transfer Tool** | Multi-engine architecture | ✅ Enhanced |
| **Basic Cost Tracking** | Real-time optimization | ✅ Superior |

### Installation

```bash
# Old Python way (deprecated)
pip install -r requirements-dev.txt
python research_infrastructure_wizard.py

# New Go way (recommended)
cd go/
go build -o aws-research-wizard ./cmd/main.go
./aws-research-wizard config recommend --domain genomics
```

## 📚 Historical Reference

This Python implementation served as the foundation for the Go version and contains:
- **33 Python modules** (~15,000 lines of code)
- **25+ research domain implementations**
- **Comprehensive Streamlit web interface**
- **Original architectural patterns**

### Key Historical Files
- `research_infrastructure_wizard.py` - Original infrastructure engine
- `comprehensive_spack_domains.py` - Domain pack definitions
- `gui_research_wizard.py` - Streamlit web interface
- `tui_research_wizard.py` - Terminal interface prototype

## 🗓️ Timeline

- **2024**: Python implementation developed and deployed
- **2025 Q1**: Go port initiated
- **2025 Q2**: Go version achieved feature parity + enhancements
- **July 2, 2025**: Python version deprecated and archived

## ⚡ Get Started with Go Version

```bash
# Quick start with Go implementation
cd ../go/
./aws-research-wizard config recommend \
  --domain genomics \
  --size large \
  --budget 1000 \
  --users 3

# Interactive mode
./aws-research-wizard config tui
```

---

**For all new development and production use, migrate to the Go implementation.**
**This Python version will be removed in 6 months (January 2026).**
