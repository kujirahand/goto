"""
Pytest configuration and shared fixtures.
"""
import json
import os
import shutil
import subprocess
import tempfile
from pathlib import Path
from typing import List

import pytest


class GotoTestHelper:
    """Helper class for testing goto command."""
    
    def __init__(self, goto_binary_path: str):
        self.goto_binary = goto_binary_path
        self.temp_dir = None
        self.config_file = None
        self.history_file = None
    
    def setup_temp_env(self) -> None:
        """Set up temporary test environment."""
        self.temp_dir = tempfile.mkdtemp(prefix="goto_test_")
        self.config_file = os.path.join(self.temp_dir, "test_config.toml")
        self.history_file = os.path.join(self.temp_dir, "test_history.json")
    
    def cleanup_temp_env(self) -> None:
        """Clean up temporary test environment."""
        if self.temp_dir and os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir)
    
    def create_config(self, config_content: str) -> str:
        """Create a test config file and return its path."""
        if not self.config_file:
            raise RuntimeError("Test environment not set up")
        with open(self.config_file, 'w', encoding='utf-8') as f:
            f.write(config_content)
        return self.config_file
    
    def create_history(self, history_entries) -> str:
        """Create a test history file and return its path."""
        if not self.history_file:
            raise RuntimeError("Test environment not set up")
        
        # Handle both old format (list) and new format (dict with entries)
        if isinstance(history_entries, list):
            history_data = {"entries": history_entries}
        elif isinstance(history_entries, dict) and "entries" in history_entries:
            history_data = history_entries
        else:
            history_data = {"entries": []}
            
        with open(self.history_file, 'w', encoding='utf-8') as f:
            json.dump(history_data, f, indent=2)
        return self.history_file
    
    def run_goto(self, args: List[str], timeout: int = 5):
        """
        Run goto command with given arguments.
        
        Returns:
            Tuple of (return_code, stdout, stderr)
        """
        cmd = [self.goto_binary]
        
        # Add config and history file options first if they exist
        if self.config_file and os.path.exists(self.config_file):
            cmd.extend(['--config', self.config_file])
        if self.history_file and os.path.exists(self.history_file):
            cmd.extend(['--history-file', self.history_file])
        
        # Add the user-provided arguments
        cmd.extend(args)
        
        try:
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                timeout=timeout,
                check=False
            )
            return result.returncode, result.stdout, result.stderr
        except subprocess.TimeoutExpired:
            return 124, "", "Command timed out"
        except FileNotFoundError:
            return 127, "", f"Command not found: {self.goto_binary}"
    
    def run_goto_interactive(self, args: List[str], input_text: str = "", timeout: int = 2):
        """
        Run goto command in interactive mode with input.
        
        Returns:
            Tuple of (return_code, stdout, stderr)
        """
        cmd = [self.goto_binary]
        
        # Add config and history file options first if they exist
        if self.config_file and os.path.exists(self.config_file):
            cmd.extend(['--config', self.config_file])
        if self.history_file and os.path.exists(self.history_file):
            cmd.extend(['--history-file', self.history_file])
        
        # Add the user-provided arguments
        cmd.extend(args)
        
        try:
            result = subprocess.run(
                cmd,
                input=input_text,
                capture_output=True,
                text=True,
                timeout=timeout,
                check=False
            )
            return result.returncode, result.stdout, result.stderr
        except subprocess.TimeoutExpired:
            return 124, "", "Interactive command timed out"
        except FileNotFoundError:
            return 127, "", f"Command not found: {self.goto_binary}"


def build_goto_binary(project_dir: Path) -> bool:
    """Build the goto binary if it doesn't exist or is outdated."""
    go_dir = project_dir / "go"
    goto_binary = go_dir / "goto"
    
    # Check if binary exists and is newer than source files
    if goto_binary.exists():
        binary_mtime = goto_binary.stat().st_mtime
        source_files = list(go_dir.glob("*.go"))
        if source_files and all(f.stat().st_mtime <= binary_mtime for f in source_files):
            return True
    
    # Build the binary
    try:
        result = subprocess.run(
            ["go", "build", "-o", "goto"],
            cwd=go_dir,
            capture_output=True,
            text=True,
            timeout=30,
            check=False
        )
        return result.returncode == 0
    except (subprocess.TimeoutExpired, FileNotFoundError):
        return False


@pytest.fixture(scope="session")
def ensure_goto_binary():
    """Ensure goto binary is built before running tests."""
    script_dir = Path(__file__).parent
    project_dir = script_dir.parent
    goto_binary = project_dir / "go" / "goto"
    
    if not goto_binary.exists():
        success = build_goto_binary(project_dir)
        if not success:
            pytest.fail("Failed to build goto binary")
    
    return str(goto_binary)


@pytest.fixture
def goto_helper(ensure_goto_binary):
    """Pytest fixture for GotoTestHelper."""
    helper = GotoTestHelper(ensure_goto_binary)
    helper.setup_temp_env()
    
    yield helper
    
    helper.cleanup_temp_env()


@pytest.fixture
def sample_config():
    """Sample TOML configuration for testing."""
    return '''# Test configuration for goto

[home]
path = "~/"
shortcut = "h"

[projects]
path = "~/Projects"
shortcut = "p"

[documents]
path = "~/Documents"

[website]
path = "https://example.com"
shortcut = "w"

[scripts]
path = "~/Scripts"
command = "ls -la"
'''


@pytest.fixture
def sample_history():
    """Sample history data for testing."""
    return [
        {
            "label": "home",
            "last_used": "2024-01-01T12:00:00Z"
        },
        {
            "label": "projects",
            "last_used": "2024-01-02T12:00:00Z"
        }
    ]
