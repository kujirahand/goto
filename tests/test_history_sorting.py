"""
Test history sorting functionality of goto command.
"""
import sys
from datetime import datetime, timedelta
from pathlib import Path

# Add parent directory to path to import conftest
sys.path.insert(0, str(Path(__file__).parent))


class TestHistorySorting:
    """Test history-based sorting functionality."""

    def test_history_display_sorting(self, goto_helper):
        """Test that history command displays entries in correct order."""
        # Create config with multiple destinations
        config_content = '''
[home]
path = "~/"

[projects]
path = "~/Projects"

[documents]
path = "~/Documents"

[downloads]
path = "~/Downloads"
'''
        goto_helper.create_config(config_content)

        # Create history with different timestamps
        base_time = datetime.now()
        history_entries = [
            {
                "label": "home",
                "last_used": (base_time - timedelta(days=3)).strftime("%Y-%m-%dT%H:%M:%SZ")
            },
            {
                "label": "projects",
                "last_used": (base_time - timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Most recent
            },
            {
                "label": "documents",
                "last_used": (base_time - timedelta(days=5)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Oldest
            },
            {
                "label": "downloads",
                "last_used": (base_time - timedelta(days=2)).strftime("%Y-%m-%dT%H:%M:%SZ")
            }
        ]
        goto_helper.create_history(history_entries)

        # Test goto command with history flag
        returncode, stdout, stderr = goto_helper.run_goto(["--history"])

        # Verify the command succeeded
        assert returncode == 0, f"Command failed with stderr: {stderr}"
        
        # Verify that projects (most recent) appears first
        lines = stdout.strip().split('\n')
        assert any("projects" in line and "1." in line for line in lines), "Projects should be listed as #1"
        assert any("downloads" in line and "2." in line for line in lines), "Downloads should be listed as #2"
        assert any("home" in line and "3." in line for line in lines), "Home should be listed as #3"
        assert any("documents" in line and "4." in line for line in lines), "Documents should be listed as #4"

    def test_interactive_mode_sorting(self, goto_helper):
        """Test that interactive mode displays entries in history order."""
        # Create config with multiple destinations
        config_content = '''
[home]
path = "~/"

[projects]
path = "~/Projects"

[documents]
path = "~/Documents"

[downloads]
path = "~/Downloads"
'''
        goto_helper.create_config(config_content)

        # Create history with different timestamps
        base_time = datetime.now()
        history_entries = [
            {
                "label": "home",
                "last_used": (base_time - timedelta(days=3)).strftime("%Y-%m-%dT%H:%M:%SZ")
            },
            {
                "label": "projects",
                "last_used": (base_time - timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Most recent
            },
            {
                "label": "documents",
                "last_used": (base_time - timedelta(days=5)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Oldest
            },
            {
                "label": "downloads",
                "last_used": (base_time - timedelta(days=2)).strftime("%Y-%m-%dT%H:%M:%SZ")
            }
        ]
        goto_helper.create_history(history_entries)

        # Test goto command in interactive mode (simulate pressing 'q' to quit)
        _, stdout, _ = goto_helper.run_goto_interactive([], input_text="q\n")

        # Verify that projects (most recent) appears as option 1
        lines = stdout.strip().split('\n')
        project_line = None
        for line in lines:
            if "1." in line and "projects" in line:
                project_line = line
                break
        
        assert project_line is not None, f"Projects should be option 1. Stdout: {stdout}"
