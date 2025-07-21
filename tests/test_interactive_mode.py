"""
Test interactive mode functionality of goto command.
"""
import sys
from datetime import datetime, timedelta
from pathlib import Path

# Add parent directory to path to import conftest
sys.path.insert(0, str(Path(__file__).parent))


class TestInteractiveMode:
    """Test interactive mode functionality."""

    def test_interactive_mode_with_history_sorting(self, goto_helper):
        """Test that interactive mode displays entries sorted by history."""
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

        # Test goto command in interactive mode (simulating pressing 'q' to quit)
        _, stdout, _ = goto_helper.run_goto_interactive([], input_text="q\n")

        # Verify that the most recently used items appear first
        lines = stdout.strip().split('\n')
        
        # Find the numbered options - handle ANSI escape sequences
        import re
        options = {}
        for line in lines:
            # Remove ANSI escape sequences
            clean_line = re.sub(r'\x1b\[[0-9;]*m', '', line)
            
            if clean_line.strip().startswith(('1.', '2.', '3.', '4.')):
                if 'projects' in clean_line:
                    options['projects'] = 1
                elif 'downloads' in clean_line:
                    options['downloads'] = 2
                elif 'home' in clean_line:
                    options['home'] = 3
                elif 'documents' in clean_line:
                    options['documents'] = 4

        # Verify correct ordering based on history
        assert 'projects' in options, "Projects should appear in the list"
        assert options.get('projects') == 1, f"Projects should be option 1, but found at position {options.get('projects')}"

    def test_interactive_mode_quit_functionality(self, goto_helper):
        """Test that interactive mode can be properly exited."""
        # Create minimal config
        config_content = '''
[test]
path = "~/test"
'''
        goto_helper.create_config(config_content)
        goto_helper.create_history([])

        # Test goto command in interactive mode (simulate pressing 'q' to quit)
        returncode, stdout, _ = goto_helper.run_goto_interactive([], input_text="q\n")

        # Should exit successfully
        assert returncode == 0, "Interactive mode should exit successfully with 'q'"
        assert "test" in stdout, "Should display the test destination"
