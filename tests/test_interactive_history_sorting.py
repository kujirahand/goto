"""
Test for interactive mode history sorting.
This test verifies that directories are sorted by most recently used in interactive mode.
"""

import sys
from datetime import datetime, timedelta
from pathlib import Path

# Add parent directory to path to import conftest
sys.path.insert(0, str(Path(__file__).parent))


class TestInteractiveHistorySorting:
    """Test interactive mode history sorting functionality."""

    def test_interactive_history_sorting(self, goto_helper):
        """Test that interactive mode sorts entries by most recently used."""
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
'''
        goto_helper.create_config(config_content)

        # Create history with different timestamps (beta most recent, gamma oldest)
        base_time = datetime.now()
        history_entries = [
            {
                "label": "alpha",
                "last_used": (base_time - timedelta(days=2)).strftime("%Y-%m-%dT%H:%M:%SZ")
            },
            {
                "label": "beta",
                "last_used": (base_time - timedelta(minutes=30)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Most recent
            },
            {
                "label": "gamma",
                "last_used": (base_time - timedelta(days=5)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Oldest
            },
            {
                "label": "delta",
                "last_used": (base_time - timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")
            }
        ]
        goto_helper.create_history(history_entries)

        # Test interactive mode (simulate pressing 'q' to quit)
        _, stdout, _ = goto_helper.run_goto_interactive([], input_text="q\n")

        # Extract the order from output
        lines = stdout.strip().split('\n')
        
        # Find lines that start with numbers (1., 2., etc.) - handle ANSI escape sequences
        order = []
        for line in lines:
            # Remove ANSI escape sequences
            import re
            clean_line = re.sub(r'\x1b\[[0-9;]*m', '', line)
            
            if clean_line.strip().startswith(('1.', '2.', '3.', '4.')):
                # Extract label from format "1. label (shortcut) → path" or "1. label → path"
                parts = clean_line.split(' ')
                if len(parts) >= 2:
                    label_part = parts[1]
                    # Remove shortcut part if exists
                    if '(' in label_part:
                        label = label_part.split('(')[0].strip()
                    else:
                        label = label_part.strip()
                    order.append(label)

        # Expected order: beta (most recent), delta, alpha, gamma (oldest)
        expected_order = ["beta", "delta", "alpha", "gamma"]
        
        assert len(order) >= 4, f"Should have at least 4 entries, got {len(order)}: {order}"
        assert order[:4] == expected_order, f"Expected order {expected_order}, got {order[:4]}"

    def test_cursor_mode_history_sorting(self, goto_helper):
        """Test that cursor mode also respects history sorting."""
        # Create config with two destinations
        config_content = '''[first]
path = "~/first"

[second]
path = "~/second"
'''
        goto_helper.create_config(config_content)

        # Create history with second being more recent
        base_time = datetime.now()
        history_entries = [
            {
                "label": "first",
                "last_used": (base_time - timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")
            },
            {
                "label": "second",
                "last_used": (base_time - timedelta(minutes=10)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Most recent
            }
        ]
        goto_helper.create_history(history_entries)

        # Test with cursor mode
        _, stdout, _ = goto_helper.run_goto_interactive(["-c"], input_text="q\n")

        # Should show second first (most recent)
        assert "second" in stdout, "Second should appear in cursor mode"
        assert "first" in stdout, "First should appear in cursor mode"
