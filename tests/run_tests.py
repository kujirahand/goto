#!/usr/bin/env python3
"""
Main test runner for goto command tests.
Runs all tests and provides summary.
"""
import os
import sys
import subprocess
from pathlib import Path


def run_pytest():
    """Run pytest with appropriate arguments."""
    test_dir = Path(__file__).parent
    
    # Change to test directory
    os.chdir(test_dir)
    
    # Run pytest
    cmd = [
        sys.executable, "-m", "pytest",
        "-v",
        "--tb=short",
        "--cov=.",
        "--cov-report=term-missing",
        "tests/"
    ]
    
    print("🚀 Running goto command test suite...")
    print("=" * 60)
    
    try:
        result = subprocess.run(cmd, check=False)
        return result.returncode
    except FileNotFoundError:
        print("❌ pytest not found. Please install it with:")
        print("   pip install pytest pytest-cov pytest-mock")
        return 1


def build_goto_first():
    """Build goto binary before running tests."""
    project_dir = Path(__file__).parent.parent
    go_dir = project_dir / "go"
    
    print("🔨 Building goto binary...")
    
    try:
        result = subprocess.run(
            ["go", "build", "-o", "goto"],
            cwd=go_dir,
            check=False,
            capture_output=True,
            text=True
        )
        
        if result.returncode == 0:
            print("✅ goto binary built successfully")
            return True
        else:
            print("❌ Failed to build goto binary:")
            print(result.stderr)
            return False
    except FileNotFoundError:
        print("❌ Go compiler not found. Please install Go.")
        return False


def main():
    """Main test runner."""
    print("🎯 goto Command Test Suite")
    print("=" * 60)
    
    # Build goto binary first
    if not build_goto_first():
        return 1
    
    # Run pytest
    exit_code = run_pytest()
    
    print("=" * 60)
    if exit_code == 0:
        print("🎉 All tests passed!")
    else:
        print("❌ Some tests failed.")
    
    return exit_code


if __name__ == "__main__":
    sys.exit(main())
