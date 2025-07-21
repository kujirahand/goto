"""
Debug script to check current history functionality.
"""
import sys
import tempfile
import os
import json
from datetime import datetime, timedelta
from pathlib import Path

# Add parent directory to path
sys.path.insert(0, str(Path(__file__).parent.parent))
from conftest import GotoTestHelper, build_goto_binary


def debug_history():
    """Debug the current history functionality."""
    # Find goto binary
    script_dir = Path(__file__).parent.parent
    project_dir = script_dir.parent
    goto_binary = project_dir / "go" / "goto"
    
    if not goto_binary.exists():
        print("Building goto binary...")
        if not build_goto_binary(project_dir):
            print("Failed to build goto binary")
            return
    
    # Create helper
    helper = GotoTestHelper(str(goto_binary))
    helper.setup_temp_env()
    
    try:
        # Create test config
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
'''
        config_path = helper.create_config(config_content)
        print(f"Created config at: {config_path}")
        
        # Create history with timestamps
        base_time = datetime.now()
        history_entries = [
            {
                "label": "home",
                "last_used": (base_time - timedelta(days=3)).isoformat() + "Z"
            },
            {
                "label": "projects", 
                "last_used": (base_time - timedelta(days=1)).isoformat() + "Z"  # Most recent
            },
            {
                "label": "documents",
                "last_used": (base_time - timedelta(days=5)).isoformat() + "Z"  # Oldest
            },
            {
                "label": "downloads",
                "last_used": (base_time - timedelta(days=2)).isoformat() + "Z"
            }
        ]
        history_path = helper.create_history(history_entries)
        print(f"Created history at: {history_path}")
        
        # Show the created files
        print("\n=== Config File ===")
        with open(config_path, 'r') as f:
            print(f.read())
        
        print("\n=== History File ===")
        with open(history_path, 'r') as f:
            print(f.read())
        
        # Test history command
        print("\n=== History Command Output ===")
        returncode, stdout, stderr = helper.run_goto(["--history"])
        print(f"Return code: {returncode}")
        print(f"STDOUT:\n{stdout}")
        print(f"STDERR:\n{stderr}")
        
        # Test completion command
        print("\n=== Completion Command Output ===")
        returncode, stdout, stderr = helper.run_goto(["--complete"])
        print(f"Return code: {returncode}")
        print(f"STDOUT:\n{stdout}")
        print(f"STDERR:\n{stderr}")
        
        # Test help to see available options
        print("\n=== Help Output ===")
        returncode, stdout, stderr = helper.run_goto(["-h"])
        print(f"Return code: {returncode}")
        print(f"STDOUT:\n{stdout}")
        print(f"STDERR:\n{stderr}")
        
    finally:
        helper.cleanup_temp_env()


if __name__ == "__main__":
    debug_history()
