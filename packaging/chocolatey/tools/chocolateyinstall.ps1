$ErrorActionPreference = 'Stop'

$packageName = 'aws-research-wizard'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$version = '1.0.0'

# Package parameters
$packageArgs = @{
  packageName    = $packageName
  unzipLocation  = $toolsDir
  fileType       = 'exe'
  url64bit       = "https://github.com/aws-research-wizard/aws-research-wizard/releases/download/v$version/aws-research-wizard-windows-amd64.exe"
  softwareName   = 'AWS Research Wizard*'
  checksum64     = 'PLACEHOLDER_CHECKSUM'
  checksumType64 = 'sha256'
  validExitCodes = @(0)
}

# Download and install
Write-Host "Downloading AWS Research Wizard v$version..." -ForegroundColor Green
Get-ChocolateyWebFile @packageArgs

# Rename the downloaded file to remove architecture suffix
$downloadedFile = Join-Path $toolsDir "aws-research-wizard-windows-amd64.exe"
$targetFile = Join-Path $toolsDir "aws-research-wizard.exe"

if (Test-Path $downloadedFile) {
    Move-Item $downloadedFile $targetFile -Force
    Write-Host "AWS Research Wizard installed to: $targetFile" -ForegroundColor Green
} else {
    throw "Downloaded file not found: $downloadedFile"
}

# Verify installation
try {
    $versionOutput = & $targetFile version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Installation verified: $versionOutput" -ForegroundColor Green
    } else {
        Write-Warning "Installation verification failed, but binary is present"
    }
} catch {
    Write-Warning "Could not verify installation, but binary should be functional"
}

# Display getting started message
Write-Host ""
Write-Host "🎉 AWS Research Wizard has been successfully installed!" -ForegroundColor Green
Write-Host ""
Write-Host "Getting Started:" -ForegroundColor Cyan
Write-Host "  aws-research-wizard --help                    # Show all available commands"
Write-Host "  aws-research-wizard config list              # List available research domains"
Write-Host "  aws-research-wizard deploy --help            # Deploy infrastructure help"
Write-Host "  aws-research-wizard monitor --help           # Monitoring dashboard help"
Write-Host ""
Write-Host "Quick Start Example:" -ForegroundColor Cyan
Write-Host "  aws-research-wizard config                   # Interactive domain configuration"
Write-Host "  aws-research-wizard deploy --domain genomics # Deploy a genomics environment"
Write-Host "  aws-research-wizard monitor                  # Launch monitoring dashboard"
Write-Host ""
Write-Host "Documentation:" -ForegroundColor Cyan
Write-Host "  https://github.com/aws-research-wizard/aws-research-wizard"
Write-Host ""
Write-Host "Note: AWS credentials are required for deployment and monitoring operations." -ForegroundColor Yellow
Write-Host "Configure with: aws configure" -ForegroundColor Yellow
