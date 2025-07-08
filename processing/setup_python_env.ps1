$ErrorActionPreference = "Stop"

# Check Python
if (-not (Get-Command python -ErrorAction SilentlyContinue)) {
    Write-Host "Python not found. Please install Python 3.8+" -ForegroundColor Red
    exit 1
}

# Create venv if missing
if (-not (Test-Path "venv")) {
    python -m venv venv
}

# Install requirements
.\venv\Scripts\Activate.ps1
python -m pip install --upgrade pip
python -m pip install -r requirements.txt

Write-Host "✅ Setup complete! Activate with:" -ForegroundColor Green
Write-Host "    .\venv\Scripts\Activate.ps1" -ForegroundColor Cyan