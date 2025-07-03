$ErrorActionPreference = 'Stop'

$packageName = 'aws-research-wizard'
$toolsDir = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$binaryPath = Join-Path $toolsDir "aws-research-wizard.exe"

Write-Host "Uninstalling AWS Research Wizard..." -ForegroundColor Yellow

# Remove the binary
if (Test-Path $binaryPath) {
    try {
        Remove-Item $binaryPath -Force
        Write-Host "✅ Removed binary: $binaryPath" -ForegroundColor Green
    } catch {
        Write-Warning "Could not remove binary: $binaryPath"
        Write-Warning "You may need to remove it manually"
    }
} else {
    Write-Host "Binary not found at expected location: $binaryPath" -ForegroundColor Yellow
}

# Clean up any additional files
$additionalFiles = @(
    "aws-research-wizard-windows-amd64.exe"
)

foreach ($file in $additionalFiles) {
    $filePath = Join-Path $toolsDir $file
    if (Test-Path $filePath) {
        try {
            Remove-Item $filePath -Force
            Write-Host "Removed: $file" -ForegroundColor Green
        } catch {
            Write-Warning "Could not remove: $file"
        }
    }
}

Write-Host ""
Write-Host "AWS Research Wizard has been uninstalled." -ForegroundColor Green
Write-Host ""
Write-Host "Note: Configuration files and AWS credentials are preserved." -ForegroundColor Cyan
Write-Host "If you want to remove all traces, you may also want to:" -ForegroundColor Cyan
Write-Host "  - Remove ~/.aws-research-wizard (if it exists)" -ForegroundColor Cyan
Write-Host "  - Remove AWS CLI configuration (if no longer needed)" -ForegroundColor Cyan
Write-Host ""
Write-Host "Thank you for using AWS Research Wizard!" -ForegroundColor Green
