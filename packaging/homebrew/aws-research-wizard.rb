class AwsResearchWizard < Formula
  desc "Complete research environment management for AWS"
  homepage "https://github.com/aws-research-wizard/aws-research-wizard"
  version "1.0.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/aws-research-wizard/aws-research-wizard/releases/download/v#{version}/aws-research-wizard-darwin-amd64"
      sha256 "PLACEHOLDER_INTEL_SHA256"
    else
      url "https://github.com/aws-research-wizard/aws-research-wizard/releases/download/v#{version}/aws-research-wizard-darwin-arm64"
      sha256 "PLACEHOLDER_ARM64_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/aws-research-wizard/aws-research-wizard/releases/download/v#{version}/aws-research-wizard-linux-amd64"
      sha256 "PLACEHOLDER_LINUX_INTEL_SHA256"
    else
      url "https://github.com/aws-research-wizard/aws-research-wizard/releases/download/v#{version}/aws-research-wizard-linux-arm64"
      sha256 "PLACEHOLDER_LINUX_ARM64_SHA256"
    end
  end

  def install
    bin.install "aws-research-wizard-#{OS.kernel_name.downcase}-#{Hardware::CPU.arch}" => "aws-research-wizard"
    
    # Generate shell completions
    generate_completions_from_executable(bin/"aws-research-wizard", "completion")
  end

  test do
    assert_match "AWS Research Wizard", shell_output("#{bin}/aws-research-wizard version")
    
    # Test that help commands work
    assert_match "Usage:", shell_output("#{bin}/aws-research-wizard --help")
    assert_match "config", shell_output("#{bin}/aws-research-wizard --help")
    assert_match "deploy", shell_output("#{bin}/aws-research-wizard --help")
    assert_match "monitor", shell_output("#{bin}/aws-research-wizard --help")
    
    # Test subcommand help
    assert_match "Domain configuration", shell_output("#{bin}/aws-research-wizard config --help")
    assert_match "Infrastructure deployment", shell_output("#{bin}/aws-research-wizard deploy --help")
    assert_match "Real-time monitoring", shell_output("#{bin}/aws-research-wizard monitor --help")
  end

  def caveats
    <<~EOS
      🔬 AWS Research Wizard has been installed!
      
      Getting Started:
        aws-research-wizard --help                    # Show all available commands
        aws-research-wizard config list              # List available research domains
        aws-research-wizard deploy --help            # Deploy infrastructure help
        aws-research-wizard monitor --help           # Monitoring dashboard help
      
      Quick Start Example:
        aws-research-wizard config                   # Interactive domain configuration
        aws-research-wizard deploy --domain genomics # Deploy a genomics environment
        aws-research-wizard monitor                  # Launch monitoring dashboard
      
      Documentation:
        https://github.com/aws-research-wizard/aws-research-wizard
      
      Note: AWS credentials are required for deployment and monitoring operations.
      Configure with: aws configure
    EOS
  end
end