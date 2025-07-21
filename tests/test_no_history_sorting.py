"""
Test for interactive mode without history file (alphabetical sorting).
This test verifies that directories are sorted alphabetically when no history exists.
"""

import sys
from pathlib import Path

# Add parent directory to path to import conftest
sys.path.insert(0, str(Path(__file__).parent))


class TestNoHistorySorting:
    """Test alphabetical sorting when no history exists."""

    def test_alphabetical_sorting_no_history(self, goto_helper):
        """Test that entries are sorted alphabetically when no history file exists."""
        # Create config with multiple destinations in non-alphabetical order
        config_content = '''[zebra]
path = "~/zebra"
shortcut = "z"

[alpha]
path = "~/alpha"
shortcut = "a"

[beta]
path = "~/beta"
shortcut = "b"

[gamma]
path = "~/gamma"
shortcut = "g"
'''
        goto_helper.create_config(config_content)

        # Don't create any history file
        # This will trigger alphabetical sorting

        # Test interactive mode output by sending "0" (exit) as input
        returncode, stdout, stderr = goto_helper.run_goto_interactive(["-l"], input_text="0\n")

        # Verify command succeeded
        assert returncode == 0, f"Command failed with stderr: {stderr}"

        # Analyze the output to check order
        lines = stdout.split('\n')
        entry_lines = [line for line in lines if line.strip() and (
            line.startswith('1.') or line.startswith('2.') or 
            line.startswith('3.') or line.startswith('4.'))]

        # Expected order without history: alphabetical
        expected_labels = ["alpha", "beta", "gamma", "zebra"]
        
        actual_order = []
        for line in entry_lines:
            for label in expected_labels:
                if label in line:
                    actual_order.append(label)
                    break
        
        assert actual_order == expected_labels, f"Expected alphabetical order: {expected_labels}, but got: {actual_order}"

    def test_empty_history_file_sorting(self, goto_helper):
        """Test alphabetical sorting when history file exists but is empty."""
        # Create config
        config_content = '''[charlie]
path = "~/charlie"

[alice]
path = "~/alice"

[bob]
path = "~/bob"
'''
        goto_helper.create_config(config_content)

        # Create empty history
        goto_helper.create_history([])

        # Test interactive mode (simulate exit with '0')
        _, stdout, _ = goto_helper.run_goto_interactive([], input_text="0\n")

        # Should be sorted alphabetically
        lines = stdout.split('\n')
        
        # Remove ANSI escape sequences and find entry lines
        import re
        clean_lines = [re.sub(r'\x1b\[[0-9;]*m', '', line) for line in lines]
        entry_lines = [line for line in clean_lines if line.strip() and (
            line.startswith('1.') or line.startswith('2.') or line.startswith('3.'))]

        expected_order = ["alice", "bob", "charlie"]
        actual_order = []
        for line in entry_lines:
            for label in expected_order:
                if label in line:
                    actual_order.append(label)
                    break

        assert actual_order == expected_order, f"Expected: {expected_order}, but got: {actual_order}"
