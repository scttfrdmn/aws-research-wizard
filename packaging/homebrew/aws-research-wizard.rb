class AwsResearchWizard < Formula
  desc "Easily run research workloads in the cloud - 27 research domains"
  homepage "https://researchwizard.app"
  version "0.3.1"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/scttfrdmn/aws-research-wizard/releases/download/v#{version}/aws-research-wizard-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_ARM64_SHA256"
    else
      url "https://github.com/scttfrdmn/aws-research-wizard/releases/download/v#{version}/aws-research-wizard-darwin-amd64.tar.gz"
      sha256 "PLACEHOLDER_AMD64_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/scttfrdmn/aws-research-wizard/releases/download/v#{version}/aws-research-wizard-linux-arm64.tar.gz"
      sha256 "PLACEHOLDER_LINUX_ARM64_SHA256"
    else
      url "https://github.com/scttfrdmn/aws-research-wizard/releases/download/v#{version}/aws-research-wizard-linux-amd64.tar.gz"
      sha256 "PLACEHOLDER_LINUX_AMD64_SHA256"
    end
  end

  depends_on "awscli"
  depends_on "terraform"

  def install
    bin.install "aws-research-wizard"

    # Generate shell completions if supported
    if respond_to?(:generate_completions_from_executable)
      generate_completions_from_executable(bin/"aws-research-wizard", "completion")
    end
  end

  test do
    assert_match "0.3.0", shell_output("#{bin}/aws-research-wizard --version")

    # Test that help commands work
    assert_match "Usage:", shell_output("#{bin}/aws-research-wizard --help")
    assert_match "config", shell_output("#{bin}/aws-research-wizard --help")
    assert_match "deploy", shell_output("#{bin}/aws-research-wizard --help")
    assert_match "monitor", shell_output("#{bin}/aws-research-wizard --help")
  end

  def caveats
    <<~EOS
      🔬 AWS Research Wizard has been installed!

      Getting Started:
        1. Configure AWS credentials:
           aws configure

        2. See available research domains:
           aws-research-wizard config list

        3. Deploy a research environment:
           aws-research-wizard deploy --domain genomics --region us-west-2

      Quick Start Guide:
        https://researchwizard.app/getting-started/

      Available Research Domains:
        • Genomics & DNA Analysis      • Climate Modeling
        • Machine Learning & AI        • Materials Science
        • Neuroscience                 • Astronomy & Astrophysics
        • Cybersecurity Research       • Digital Humanities
        • Drug Discovery               • Economics & Finance
        • And 17 more domains...

      Documentation: https://researchwizard.app

      Note: Requires valid AWS account with appropriate permissions.
    EOS
  end
end
