#!/bin/bash
"""
Demo Script for Spack Manager TUI
Showcases the interactive Terminal User Interface for Spack environment management
"""

echo "🚀 Spack Manager TUI Demo"
echo "========================="
echo ""

echo "📦 Available Features:"
echo "  • Interactive environment browser"
echo "  • Real-time installation progress"  
echo "  • Package details and status"
echo "  • Live logs and debugging"
echo "  • Keyboard shortcuts:"
echo "    - ↑/↓ or j/k: Navigate"
echo "    - Enter: Select environment"
echo "    - i: Install environment"
echo "    - v: View logs"
echo "    - Tab: Switch views"
echo "    - q: Quit"
echo ""

echo "🔧 Available Environments:"
echo "  • genomics_lab (Bioinformatics tools)"
echo "  • climate_modeling (Weather/climate simulation)"
echo "  • ai_research_studio (ML/AI frameworks)"
echo "  • astronomy_lab (Astronomical analysis)"
echo "  • materials_science (Quantum chemistry)"
echo "  • neuroscience_lab (Brain simulation)"
echo "  • physics_simulation (HEP tools)"
echo ""

echo "💡 TUI Views:"
echo "  1. Environment List - Browse and select environments"
echo "  2. Environment Details - View packages and status"
echo "  3. Installation Progress - Real-time progress bars"
echo "  4. Logs - Live installation logs and debugging"
echo ""

echo "🎯 Try it out:"
echo "  ./aws-research-wizard deploy spack tui"
echo ""

echo "📊 Example TUI Layout:"
echo "┌─────────────────────────────────────────────────────────────────┐"
echo "│                        Spack Manager TUI                       │"
echo "├─────────────────────────────────────────────────────────────────┤"
echo "│ Spack Environments                                              │"
echo "│ > genomics_lab                     Ready • 25 packages          │"
echo "│   climate_modeling                 Installed • 15 packages      │"
echo "│   ai_research_studio              Ready • 20 packages          │"
echo "│   astronomy_lab                   Ready • 18 packages          │"
echo "│                                                                 │"
echo "│                                                                 │"
echo "│                                                                 │"
echo "├─────────────────────────────────────────────────────────────────┤"
echo "│ i: install • enter: details • v: logs • q: quit               │"
echo "└─────────────────────────────────────────────────────────────────┘"
echo ""

echo "✨ Advanced Features:"
echo "  • Progress monitoring with real-time updates"
echo "  • Package dependency visualization"
echo "  • Build time estimates and optimization"
echo "  • Binary cache status and optimization"
echo "  • Error handling with actionable suggestions"
echo ""

echo "🚀 Ready to launch TUI? (Press Enter to continue or Ctrl+C to exit)"
read -r

echo "Launching Spack Manager TUI..."
cd "$(dirname "$0")/../go" || exit 1
./main deploy spack tui