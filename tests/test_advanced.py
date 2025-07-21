"""
Advanced functionality tests for goto command.
"""
import os
import sys
import tempfile
from pathlib import Path

# Add parent directory to path to import conftest
sys.path.insert(0, str(Path(__file__).parent.parent))
from conftest import GotoTestHelper, build_goto_binary


class TestAdvancedFeatures:
    """Test advanced goto features."""
    
    def test_custom_config_and_history_files(self, goto_helper):
        """Test using custom config and history file paths."""
        # Create custom config
        config_content = '''
[test_dest]
path = "~/test_location"
shortcut = "t"
'''
        config_path = goto_helper.create_config(config_content)
        history_path = goto_helper.create_history([])
        
        # Run goto with custom files
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        
        assert returncode == 0, f"Custom config/history failed: {stderr}"
        assert "test_dest" in stdout.lower(), "Custom config not loaded"
    
    def test_path_expansion(self, goto_helper):
        """Test tilde and relative path expansion."""
        config_content = '''
[home]
path = "~/"

[relative]
path = "./relative_path"

[absolute]
path = "/tmp"
'''
        goto_helper.create_config(config_content)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        
        assert returncode == 0, f"Path expansion test failed: {stderr}"
        output = stdout.lower()
        assert "home" in output, "Tilde path not handled"
        assert "relative" in output, "Relative path not handled"
        assert "absolute" in output, "Absolute path not handled"
    
    def test_command_execution_config(self, goto_helper):
        """Test destinations with command configurations."""
        config_content = '''
[scripts]
path = "~/Scripts"
command = "ls -la"

[logs]
path = "/var/log"
command = "tail -f"
'''
        goto_helper.create_config(config_content)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        
        assert returncode == 0, f"Command config test failed: {stderr}"
        output = stdout.lower()
        assert "scripts" in output, "Command destination not found"
        assert "logs" in output, "Command destination not found"
    
    def test_url_destinations(self, goto_helper):
        """Test URL-type destinations."""
        config_content = '''
[github]
path = "https://github.com"
shortcut = "gh"

[docs]
path = "https://docs.example.com"

[localhost]
path = "http://localhost:3000"
'''
        goto_helper.create_config(config_content)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        
        assert returncode == 0, f"URL destinations test failed: {stderr}"
        output = stdout.lower()
        assert "github" in output, "GitHub URL not found"
        assert "docs" in output, "Docs URL not found" 
        assert "localhost" in output, "Localhost URL not found"


class TestEdgeCases:
    """Test edge cases and error conditions."""
    
    def test_empty_config_file(self, goto_helper):
        """Test handling of empty config file."""
        goto_helper.create_config("")
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        
        # Empty config should be handled gracefully - may succeed or fail
        assert returncode in [0, 1], f"Empty config handling failed unexpectedly: {stderr}"
    
    def test_config_with_only_comments(self, goto_helper):
        """Test config file with only comments."""
        config_content = '''
# This is a comment
# Another comment
# No actual configuration
'''
        goto_helper.create_config(config_content)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        
        # Comment-only config should be handled gracefully - may succeed or fail
        assert returncode in [0, 1], f"Comment-only config failed unexpectedly: {stderr}"
    
    def test_config_with_special_characters(self, goto_helper):
        """Test config with special characters in paths and labels."""
        config_content = '''
[test-with-dashes]
path = "~/test-path"

[test_with_underscores]
path = "~/test_path"

[test.with.dots]
path = "~/test.path"

[test123]
path = "~/test123"
'''
        goto_helper.create_config(config_content)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        
        assert returncode == 0, f"Special characters test failed: {stderr}"
        output = stdout.lower()
        assert "test-with-dashes" in output, "Dashes in label not handled"
        assert "test_with_underscores" in output, "Underscores in label not handled"
    
    def test_long_config_file(self, goto_helper):
        """Test handling of config file with many entries."""
        config_lines = ["# Large config file test"]
        
        # Generate 50 destinations
        for i in range(50):
            config_lines.extend([
                f"[dest{i:02d}]",
                f'path = "~/destination_{i}"',
                f'shortcut = "d{i}"' if i < 26 else "",  # Only first 26 get shortcuts
                ""
            ])
        
        config_content = "\n".join(config_lines)
        goto_helper.create_config(config_content)
        
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
        
        assert returncode == 0, f"Large config test failed: {stderr}"
        output = stdout.lower()
        assert "dest00" in output, "First destination not found"
        assert "dest49" in output, "Last destination not found"
    
    def test_malformed_toml_sections(self, goto_helper):
        """Test various malformed TOML sections."""
        malformed_configs = [
            # Missing closing bracket
            '''
[incomplete
path = "~/test"
''',
            # Invalid characters in section name
            '''
[section with spaces]
path = "~/test"
''',
            # Missing quotes
            '''
[test]
path = ~/unquoted
''',
            # Invalid escape sequence
            '''
[test]
path = "~/test\\invalid"
'''
        ]
        
        for i, config in enumerate(malformed_configs):
            goto_helper.create_config(config)
            returncode, stdout, stderr = goto_helper.run_goto(["--complete"])
            
            # Malformed configs should fail gracefully
            # We don't require specific return codes, just that it doesn't crash
            assert returncode in [0, 1, 2], f"Malformed config {i} caused unexpected failure: {stderr}"


class TestPerformance:
    """Test performance characteristics."""
    
    def test_startup_time(self, goto_helper):
        """Test that goto starts up quickly even with large configs."""
        # Create large config
        config_lines = []
        for i in range(100):
            config_lines.extend([
                f"[perf_test_{i:03d}]",
                f'path = "~/perf_test_{i}"',
                f'shortcut = "p{i}"' if i < 10 else "",
                ""
            ])
        
        config_content = "\n".join(config_lines)
        goto_helper.create_config(config_content)
        
        # Test completion with timeout (should complete quickly)
        returncode, stdout, stderr = goto_helper.run_goto(["--complete"], timeout=3)
        
        assert returncode == 0, f"Performance test failed: {stderr}"
        assert len(stdout) > 0, "No completion output generated"
    
    def test_memory_usage_with_large_history(self, goto_helper):
        """Test behavior with large history file."""
        # Create large history
        large_history = []
        for i in range(1000):
            large_history.append({
                "label": f"hist_entry_{i:04d}",
                "last_used": f"2024-01-{(i % 30) + 1:02d}T12:00:00Z"
            })
        
        goto_helper.create_config('''
[test]
path = "~/test"
''')
        goto_helper.create_history(large_history)
        
        # Test history display
        returncode, stdout, stderr = goto_helper.run_goto(["--history"], timeout=5)
        
        assert returncode == 0, f"Large history test failed: {stderr}"
        # Should handle large history without issues
