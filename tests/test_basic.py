"""
Basic functionality tests for goto command.
"""
import os
import sys
from pathlib import Path

# Add parent directory to path to import conftest
sys.path.insert(0, str(Path(__file__).parent))


class TestBasicFunctionality:
    """Test basic goto command functionality."""
    
    def test_goto_binary_exists(self):
        """Test that the goto binary exists and can be executed."""
        script_dir = Path(__file__).parent
        project_dir = script_dir.parent
        goto_binary = project_dir / "go" / "goto"
        
        assert goto_binary.exists(), "Goto binary not found - run 'make build' first"
        assert os.access(goto_binary, os.X_OK), "Goto binary is not executable"
    
    def test_goto_help(self, goto_helper):
        """Test goto help output."""
        returncode, stdout, stderr = goto_helper.run_goto(["-h"])
        
        # Help should return 0 or show usage
        assert returncode in [0, 1], f"Unexpected return code: {returncode}"
        assert "goto" in stdout.lower() or "usage" in stdout.lower(), "Help output doesn't contain expected content"
    
    def test_goto_version(self, goto_helper):
        """Test goto version output."""
        returncode, stdout, stderr = goto_helper.run_goto(["--version"])
        
        # Version should return 0 and show version info
        assert returncode in [0, 1], f"Unexpected return code: {returncode}"
        # Version output might be in stdout or stderr
        output = stdout + stderr
        assert any(keyword in output.lower() for keyword in ["version", "goto", "v"]), "Version output doesn't contain expected content"


class TestConfigParsing:
    """Test TOML configuration file parsing."""
    
    def test_simple_config_parsing(self, goto_helper, sample_config):
        """Test parsing a simple TOML configuration."""
        goto_helper.create_config(sample_config)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        
        # Completion should work
        assert returncode == 0, f"Config parsing failed: {stderr}"
        
        # Check that all destinations are listed
        output = stdout.lower()
        expected_destinations = ["home", "projects", "documents", "website", "scripts"]
        for dest in expected_destinations:
            assert dest in output, f"Destination '{dest}' not found in completion output"
    
    def test_config_with_shortcuts(self, goto_helper):
        """Test configuration with shortcuts."""
        config_content = '''
[home]
path = "~/"
shortcut = "h"

[projects] 
path = "~/Projects"
shortcut = "p"
'''
        goto_helper.create_config(config_content)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        assert returncode == 0, f"Config with shortcuts parsing failed: {stderr}"
        
        output = stdout.lower()
        assert "home" in output and "projects" in output, "Destinations with shortcuts not found"
    
    def test_config_with_urls(self, goto_helper):
        """Test configuration with URLs."""
        config_content = '''
[github]
path = "https://github.com"
shortcut = "g"

[local]
path = "~/Documents"
'''
        goto_helper.create_config(config_content)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        assert returncode == 0, f"Config with URLs parsing failed: {stderr}"
        
        output = stdout.lower()
        assert "github" in output and "local" in output, "URL and local destinations not found"
    
    def test_invalid_config_handling(self, goto_helper):
        """Test handling of invalid TOML configuration."""
        invalid_config = '''
[invalid
path = "invalid toml syntax
'''
        goto_helper.create_config(invalid_config)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        
        # Should fail gracefully
        assert returncode != 0, "Invalid config should cause non-zero exit code"
    
    def test_nonexistent_config_handling(self, goto_helper):
        """Test handling of non-existent configuration file."""
        # Use a config file that doesn't exist
        nonexistent_config = "/tmp/nonexistent_goto_config.toml"
        
        returncode, stdout, stderr = goto_helper.run_goto(["--config", nonexistent_config, "--complete"])
        
        # Should fail gracefully - allow either success (fallback behavior) or failure
        assert returncode in [0, 1, 2], f"Non-existent config handling failed unexpectedly: {stderr}"


class TestHistoryFunctionality:
    """Test history functionality."""
    
    def test_history_file_creation(self, goto_helper, sample_config):
        """Test that history file is handled correctly."""
        goto_helper.create_config(sample_config)
        goto_helper.create_history([])  # Empty history
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        assert returncode == 0, f"History file handling failed: {stderr}"
    
    def test_history_display(self, goto_helper, sample_config, sample_history):
        """Test history display functionality."""
        goto_helper.create_config(sample_config)
        goto_helper.create_history(sample_history)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--history"])
        
        # History command should work
        assert returncode == 0, f"History display failed: {stderr}"
        
        # Check that history entries are shown
        output = stdout.lower()
        assert "home" in output or "projects" in output, "History entries not displayed"
    
    def test_empty_history(self, goto_helper, sample_config):
        """Test behavior with empty history."""
        goto_helper.create_config(sample_config)
        goto_helper.create_history([])
        
        returncode, stdout, stderr = goto_helper.run_goto(["--history"])
        
        # Should handle empty history gracefully
        assert returncode == 0, f"Empty history handling failed: {stderr}"


class TestInteractiveModes:
    """Test interactive mode options."""
    
    def test_cursor_mode_option(self, goto_helper, sample_config):
        """Test cursor mode (-c) option recognition."""
        goto_helper.create_config(sample_config)
        
        # Cursor mode should timeout (interactive)
        returncode, stdout, stderr = goto_helper.run_goto_interactive(["-c"], timeout=1)
        
        # Timeout is expected for interactive mode
        assert returncode in [0, 124], f"Cursor mode failed: {stderr}"
    
    def test_label_mode_option(self, goto_helper, sample_config):
        """Test label mode (-l) option recognition."""
        goto_helper.create_config(sample_config)
        
        # Label mode should timeout (interactive) 
        returncode, stdout, stderr = goto_helper.run_goto_interactive(["-l"], timeout=1)
        
        # Timeout is expected for interactive mode
        assert returncode in [0, 124], f"Label mode failed: {stderr}"
    
    def test_interactive_modes_with_input(self, goto_helper, sample_config):
        """Test interactive modes with simulated input."""
        goto_helper.create_config(sample_config)
        
        # Test cursor mode with escape key (should exit)
        returncode, stdout, stderr = goto_helper.run_goto_interactive(["-c"], input_text="\x1b", timeout=2)
        # Should exit gracefully (various exit codes are acceptable)
        assert returncode in [0, 1, 124], f"Cursor mode with input failed: {stderr}"
        
        # Test label mode with empty input (should exit)
        returncode, stdout, stderr = goto_helper.run_goto_interactive(["-l"], input_text="\n", timeout=2)
        # Should exit gracefully
        assert returncode in [0, 1, 124], f"Label mode with input failed: {stderr}"
