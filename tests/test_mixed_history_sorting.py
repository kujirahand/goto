"""
Test for interactive mode with partial history (mixed sorting).
This test verifies that directories with history are shown first (sorted by most recent),
followed by directories without history (sorted alphabetically).
"""

import sys
from datetime import datetime, timedelta
from pathlib import Path

# Add parent directory to path to import conftest
sys.path.insert(0, str(Path(__file__).parent))


class TestMixedHistorySorting:
    """Test mixed history sorting functionality."""

    def test_mixed_history_sorting(self, goto_helper):
        """Test that entries with history are shown first, then alphabetical entries without history."""
        # Create config with multiple destinations
        config_content = '''[alpha]
path = "~/alpha"
shortcut = "a"

[beta]
path = "~/beta"
shortcut = "b"

[gamma]
path = "~/gamma"
shortcut = "g"

[delta]
path = "~/delta"
shortcut = "d"

[epsilon]
path = "~/epsilon"
shortcut = "e"

[zeta]
path = "~/zeta"
shortcut = "z"
'''
        goto_helper.create_config(config_content)

        # Create history with only some entries
        # gamma (most recent), alpha (older)
        # beta, delta, epsilon, zeta should appear after in alphabetical order
        base_time = datetime.now()
        history_entries = [
            {
                "label": "alpha",
                "last_used": (base_time - timedelta(days=2)).strftime("%Y-%m-%dT%H:%M:%SZ")
            },
            {
                "label": "gamma",
                "last_used": (base_time - timedelta(minutes=30)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Most recent
            }
        ]
        goto_helper.create_history(history_entries)

        # Test interactive mode output by sending "0" (exit) as input
        returncode, stdout, stderr = goto_helper.run_goto_interactive(["-l"], input_text="0\n")

        # Verify command succeeded
        assert returncode == 0, f"Command failed with stderr: {stderr}"

        # Analyze the output to check order
        lines = stdout.split('\n')
        entry_lines = [line for line in lines if line.strip() and (
            line.startswith('1.') or line.startswith('2.') or line.startswith('3.') or 
            line.startswith('4.') or line.startswith('5.') or line.startswith('6.'))]

        # Expected order: gamma (history, most recent), alpha (history, older), 
        # then beta, delta, epsilon, zeta (alphabetical, no history)
        expected_labels = ["gamma", "alpha", "beta", "delta", "epsilon", "zeta"]
        
        actual_order = []
        for line in entry_lines:
            for label in expected_labels:
                if label in line:
                    actual_order.append(label)
                    break
        
        assert actual_order == expected_labels, f"Expected order: {expected_labels}, but got: {actual_order}"

    def test_all_entries_have_history(self, goto_helper):
        """Test sorting when all entries have history."""
        # Create config with multiple destinations
        config_content = '''[first]
path = "~/first"

[second]
path = "~/second"

[third]
path = "~/third"
'''
        goto_helper.create_config(config_content)

        # Create history for all entries with different timestamps
        base_time = datetime.now()
        history_entries = [
            {
                "label": "second",
                "last_used": (base_time - timedelta(minutes=10)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Most recent
            },
            {
                "label": "third",
                "last_used": (base_time - timedelta(hours=1)).strftime("%Y-%m-%dT%H:%M:%SZ")
            },
            {
                "label": "first",
                "last_used": (base_time - timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Oldest
            }
        ]
        goto_helper.create_history(history_entries)

        # Test interactive mode (simulate exit with '0')
        _, stdout, _ = goto_helper.run_goto_interactive([], input_text="0\n")

        # Verify sorting by history (most recent first)
        lines = stdout.split('\n')
        
        # Remove ANSI escape sequences and find entry lines
        import re
        clean_lines = [re.sub(r'\x1b\[[0-9;]*m', '', line) for line in lines]
        entry_lines = [line for line in clean_lines if line.strip() and (
            line.startswith('1.') or line.startswith('2.') or line.startswith('3.'))]

        expected_order = ["second", "third", "first"]  # Most recent to oldest
        actual_order = []
        for line in entry_lines:
            for label in expected_order:
                if label in line:
                    actual_order.append(label)
                    break

        assert actual_order == expected_order, f"Expected: {expected_order}, but got: {actual_order}"
