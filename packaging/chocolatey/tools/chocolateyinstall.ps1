$ErrorActionPreference = 'Stop'

$packageName = 'aws-research-wizard'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$version = '0.3.1'

# Package parameters
$packageArgs = @{
  packageName    = $packageName
  unzipLocation  = $toolsDir
  fileType       = 'exe'
  url64bit       = "https://github.com/scttfrdmn/aws-research-wizard/releases/download/v$version/aws-research-wizard-windows-amd64.exe"
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
Write-Host "  1. Configure AWS credentials: aws configure"
Write-Host "  2. See available research domains: aws-research-wizard config list"
Write-Host "  3. Deploy a research environment: aws-research-wizard deploy --domain genomics"
Write-Host ""
Write-Host "Available Research Domains:" -ForegroundColor Cyan
Write-Host "  • Genomics & DNA Analysis      • Climate Modeling"
Write-Host "  • Machine Learning & AI        • Materials Science"
Write-Host "  • Neuroscience                 • Astronomy & Astrophysics"
Write-Host "  • Cybersecurity Research       • Digital Humanities"
Write-Host "  • Drug Discovery               • Economics & Finance"
Write-Host "  • And 17 more domains..."
Write-Host ""
Write-Host "Quick Start Guide:" -ForegroundColor Cyan
Write-Host "  https://researchwizard.app/getting-started/"
Write-Host ""
Write-Host "Full Documentation:" -ForegroundColor Cyan
Write-Host "  https://researchwizard.app"
Write-Host ""
Write-Host "Note: Requires valid AWS account with appropriate permissions." -ForegroundColor Yellow
