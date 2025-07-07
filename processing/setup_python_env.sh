#!/bin/bash

# fail on error
set -e

ENV_DIR="venv"

# checking if python is installed
if ! command -v python3 &> /dev/null; then
    echo "Python3 is required but not found. Please install Python3 first."
    exit 1
fi

# create virtual environment if it doesn't exist
if [ ! -d "$ENV_DIR" ]; then
    echo "Creating Python virtual environment..."
    python3 -m venv "$ENV_DIR"
else
    echo "Virtual environment already exists."
fi


echo "Installing dependencies..."
source "$ENV_DIR/bin/activate"
pip install --upgrade pip
pip install -r requirements.txt

echo "Python environment setup complete. To activate, run: source $ENV_DIR/bin/activate"