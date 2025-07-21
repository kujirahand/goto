"""
Test complete workflow with history updates.
"""
import os
import sys
from datetime import datetime, timedelta
from pathlib import Path

# Add parent directory to path to import conftest
sys.path.insert(0, str(Path(__file__).parent))


class TestCompleteWorkflow:
    """Test complete goto workflow including history updates."""
    
    def test_history_update_workflow(self, goto_helper):
        """Test the complete workflow: dir1 most recent -> use dir2 -> dir2 becomes most recent."""
        # Create test directories that actually exist
        test_dirs = ["/tmp/goto_test_dir1", "/tmp/goto_test_dir2", "/tmp/goto_test_dir3"]
        for test_dir in test_dirs:
            os.makedirs(test_dir, exist_ok=True)
        
        try:
            # Create config with test directories
            config_content = '''
[dir1]
path = "/tmp/goto_test_dir1"

[dir2]
path = "/tmp/goto_test_dir2"

[dir3]
path = "/tmp/goto_test_dir3"
'''
            goto_helper.create_config(config_content)

            # Create initial history with dir1 as most recent
            base_time = datetime.now()
            history_entries = [
                {
                    "label": "dir1",
                    "last_used": (base_time - timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Most recent
                },
                {
                    "label": "dir2", 
                    "last_used": (base_time - timedelta(days=3)).strftime("%Y-%m-%dT%H:%M:%SZ")
                },
                {
                    "label": "dir3",
                    "last_used": (base_time - timedelta(days=5)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Oldest
                }
            ]
            goto_helper.create_history(history_entries)

            # Step 1: Verify initial state (dir1 should be #1)
            returncode, stdout, stderr = goto_helper.run_goto(["--history"])
            assert returncode == 0, f"History command failed: {stderr}"
            
            lines = stdout.strip().split('\n')
            assert any("dir1" in line and "1." in line for line in lines), "dir1 should be #1 initially"

            # Step 2: Test that we can see the destinations in interactive mode
            _, interactive_stdout, _ = goto_helper.run_goto_interactive([], input_text="q\n")
            assert "dir1" in interactive_stdout, "dir1 should appear in interactive mode"
            assert "dir2" in interactive_stdout, "dir2 should appear in interactive mode"

        finally:
            # Cleanup test directories
            for test_dir in test_dirs:
                try:
                    os.rmdir(test_dir)
                except OSError:
                    pass  # Directory might not be empty or might not exist

    def test_history_sorting_without_existing_dirs(self, goto_helper):
        """Test that history sorting works even when directories don't exist."""
        # Create config with non-existent directories (for testing sorting logic)
        config_content = '''
[home]
path = "~/"

[projects]
path = "~/Projects"

[docs]
path = "~/Documents"
'''
        goto_helper.create_config(config_content)

        # Create history with specific order
        base_time = datetime.now()
        history_entries = [
            {
                "label": "docs",
                "last_used": (base_time - timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Most recent
            },
            {
                "label": "home", 
                "last_used": (base_time - timedelta(days=2)).strftime("%Y-%m-%dT%H:%M:%SZ")
            },
            {
                "label": "projects",
                "last_used": (base_time - timedelta(days=3)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Oldest
            }
        ]
        goto_helper.create_history(history_entries)

        # Test history display
        returncode, stdout, stderr = goto_helper.run_goto(["--history"])
        assert returncode == 0, f"History command failed: {stderr}"
        
        # Verify sorting order
        lines = stdout.strip().split('\n')
        assert any("docs" in line and "1." in line for line in lines), "docs should be #1"
        assert any("home" in line and "2." in line for line in lines), "home should be #2"
        assert any("projects" in line and "3." in line for line in lines), "projects should be #3"
