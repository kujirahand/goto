"""
History sorting and ordering tests for goto command.
"""
import sys
from pathlib import Path
from datetime import datetime, timedelta

# Add parent directory to path to import conftest
sys.path.insert(0, str(Path(__file__).parent.parent))


class TestHistoryOrdering:
    """Test history ordering and sorting functionality."""
    
    def test_history_sorting_by_last_used(self, goto_helper):
        """Test that history is sorted by last_used time (most recent first)."""
        # Create config with multiple destinations
        config_content = '''
[home]
path = "~/"
shortcut = "h"

[projects]
path = "~/Projects"
shortcut = "p"

[documents]
path = "~/Documents" 
shortcut = "d"

[downloads]
path = "~/Downloads"
shortcut = "b"
'''
        goto_helper.create_config(config_content)
        
        # Create history with specific timestamps (newest to oldest)
        now = datetime.now()
        history_data = {
            "entries": [
                {
                    "label": "home",
                    "last_used": (now - timedelta(days=2)).strftime("%Y-%m-%dT%H:%M:%SZ")
                },
                {
                    "label": "projects", 
                    "last_used": now.strftime("%Y-%m-%dT%H:%M:%SZ")  # Most recent
                },
                {
                    "label": "documents",
                    "last_used": (now - timedelta(days=4)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Oldest
                },
                {
                    "label": "downloads",
                    "last_used": (now - timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")
                }
            ]
        }
        goto_helper.create_history(history_data)
        
        # Get history output
        return_code, stdout, _ = goto_helper.run_goto(["--history"])
        assert return_code == 0
        
        output_lines = stdout.strip().split('\n')
        
        # Find the entries in the output (skip header lines)
        entries = []
        for line in output_lines:
            if '. ' in line and ' → ' in line:
                # Extract number and label from lines like " 1. projects → /Users/..."
                parts = line.strip().split('. ', 1)
                if len(parts) == 2:
                    order_num = int(parts[0])
                    label = parts[1].split(' → ')[0]
                    entries.append((order_num, label))
        
        # Sort by order number to ensure correct sequence
        entries.sort()
        
        # Verify the order is correct (most recent first)
        expected_order = ["projects", "downloads", "home", "documents"]
        actual_order = [label for _, label in entries]
        
        assert actual_order == expected_order, f"Expected {expected_order}, got {actual_order}"
        
        # Verify timestamps are displayed correctly (not 0001-01-01)
        assert "0001-01-01" not in stdout, "Timestamps should not show as 0001-01-01"
        
        # Verify proper timestamp format is shown
        assert "📅" in stdout, "Timestamp emoji should be present"
        
    def test_history_with_invalid_timestamps(self, goto_helper):
        """Test history handling with invalid or missing timestamps."""
        config_content = '''
[test1]
path = "/tmp"
shortcut = "t"

[test2] 
path = "/var"
shortcut = "v"
'''
        goto_helper.create_config(config_content)
        
        # Create history with some invalid timestamps
        history_data = {
            "entries": [
                {
                    "label": "test1",
                    "last_used": "invalid-date"
                },
                {
                    "label": "test2",
                    "last_used": "2025-01-15T10:30:00Z"
                }
            ]
        }
        goto_helper.create_history(history_data)
        
        # Should still work without crashing
        return_code, stdout, _ = goto_helper.run_goto(["--history"])
        assert return_code == 0
        # Invalid timestamps may cause entries to be filtered out
        # Just verify the command doesn't crash and handles it gracefully
        assert "📈" in stdout or "使用履歴が見つかりません" in stdout
        
    def test_empty_history_file(self, goto_helper):
        """Test behavior with empty history file."""
        config_content = '''
[test]
path = "/tmp"
shortcut = "t"
'''
        goto_helper.create_config(config_content)
        
        # Create empty history
        goto_helper.create_history({"entries": []})
        
        return_code, stdout, _ = goto_helper.run_goto(["--history"])
        assert return_code == 0
        # Should show header but no entries
        assert "最近の使用履歴" in stdout or "使用履歴が見つかりません" in stdout
        
    def test_missing_history_file(self, goto_helper):
        """Test behavior when history file doesn't exist."""
        config_content = '''
[test]
path = "/tmp"  
shortcut = "t"
'''
        goto_helper.create_config(config_content)
        
        # Don't create history file
        return_code, stdout, _ = goto_helper.run_goto(["--history"])
        assert return_code == 0
        # Should handle gracefully
        assert "最近の使用履歴" in stdout or "使用履歴が見つかりません" in stdout

    def test_history_with_mixed_timestamps(self, goto_helper):
        """Test history with various timestamp formats."""
        config_content = '''
[recent]
path = "/tmp/recent"
shortcut = "r"

[old]
path = "/tmp/old"
shortcut = "o"

[middle]
path = "/tmp/middle"
shortcut = "m"
'''
        goto_helper.create_config(config_content)
        
        # Create history with different valid timestamp formats
        history_data = {
            "entries": [
                {
                    "label": "recent",
                    "last_used": "2025-07-20T15:30:45Z"
                },
                {
                    "label": "old", 
                    "last_used": "2025-07-15T08:15:30Z"
                },
                {
                    "label": "middle",
                    "last_used": "2025-07-18T12:45:15Z"
                }
            ]
        }
        goto_helper.create_history(history_data)
        
        return_code, stdout, _ = goto_helper.run_goto(["--history"])
        assert return_code == 0
        
        output_lines = stdout.strip().split('\n')
        
        # Find the entries in the output
        entries = []
        for line in output_lines:
            if '. ' in line and ' → ' in line:
                parts = line.strip().split('. ', 1)
                if len(parts) == 2:
                    order_num = int(parts[0])
                    label = parts[1].split(' → ')[0]
                    entries.append((order_num, label))
        
        entries.sort()
        expected_order = ["recent", "middle", "old"]  # Newest to oldest
        actual_order = [label for _, label in entries]
        
        assert actual_order == expected_order, f"Expected {expected_order}, got {actual_order}"
