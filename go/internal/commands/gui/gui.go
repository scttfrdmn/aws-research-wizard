package gui

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/aws-research-wizard/go/internal/web"
)

// GuiCmd represents the gui command for Enhanced GUI
var GuiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch Enhanced GUI web interface",
	Long: `Launch the Enhanced GUI web interface for AWS Research Wizard.

The Enhanced GUI provides a modern, comprehensive web-based interface that
surpasses the deprecated Python Streamlit interface while leveraging Go's
performance advantages.

Features:
- Interactive domain selection and configuration
- Real-time cost calculation and optimization
- Visual deployment management
- Live monitoring dashboards
- Responsive design for desktop, tablet, and mobile

This is Phase 1 of the 17-week Enhanced GUI development plan.`,
	Example: `  # Launch GUI on default port (8080)
  aws-research-wizard gui

  # Launch on custom port
  aws-research-wizard gui --port 3000

  # Launch in development mode with verbose logging
  aws-research-wizard gui --dev --port 8080

  # Launch with TLS enabled
  aws-research-wizard gui --tls --cert server.crt --key server.key`,
	RunE: runGUI,
}

var (
	port        int
	host        string
	development bool
	enableTLS   bool
	tlsCertPath string
	tlsKeyPath  string
)

func init() {
	GuiCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the web server on")
	GuiCmd.Flags().StringVar(&host, "host", "localhost", "Host to bind the web server to")
	GuiCmd.Flags().BoolVar(&development, "dev", false, "Enable development mode with verbose logging")
	GuiCmd.Flags().BoolVar(&enableTLS, "tls", false, "Enable TLS/HTTPS")
	GuiCmd.Flags().StringVar(&tlsCertPath, "cert", "", "Path to TLS certificate file")
	GuiCmd.Flags().StringVar(&tlsKeyPath, "key", "", "Path to TLS private key file")
}

func runGUI(cmd *cobra.Command, args []string) error {
	// Validate TLS configuration
	if enableTLS && (tlsCertPath == "" || tlsKeyPath == "") {
		return fmt.Errorf("TLS enabled but certificate or key path not provided")
	}

	// Create server configuration
	config := web.ServerConfig{
		Port:        port,
		Host:        host,
		Development: development,
		EnableTLS:   enableTLS,
		TLSCertPath: tlsCertPath,
		TLSKeyPath:  tlsKeyPath,
	}

	// Create and start server
	server, err := web.NewServer(config)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Setup graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to listen for interrupt signal to trigger shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.Start()
	}()

	// Display startup information
	fmt.Println("🚀 AWS Research Wizard Enhanced GUI")
	fmt.Println("====================================")
	fmt.Printf("📡 Server: http://%s:%d\n", host, port)
	fmt.Printf("🔧 Development Mode: %v\n", development)
	fmt.Printf("🔒 TLS Enabled: %v\n", enableTLS)
	if development {
		fmt.Printf("🌐 API Endpoints:\n")
		fmt.Printf("   • Domains: http://%s:%d/api/domains\n", host, port)
		fmt.Printf("   • Health: http://%s:%d/api/health\n", host, port)
		fmt.Printf("   • Version: http://%s:%d/api/version\n", host, port)
	}
	fmt.Println("✋ Press Ctrl+C to stop")
	fmt.Println()

	// Wait for shutdown signal or server error
	select {
	case <-sigChan:
		fmt.Println("\n🛑 Shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			fmt.Printf("\n❌ Server error: %v\n", err)
			return err
		}
	}

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	fmt.Println("🔄 Gracefully shutting down server...")
	if err := server.Stop(shutdownCtx); err != nil {
		fmt.Printf("⚠️ Server shutdown error: %v\n", err)
		return err
	}

	fmt.Println("✅ Server stopped gracefully")
	return nil
}
