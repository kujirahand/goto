"""
Simplified history test to debug the issue.
"""
import os
import sys
import tempfile
import json
from datetime import datetime, timedelta
from pathlib import Path

# Add parent directory to path to import conftest  
sys.path.insert(0, str(Path(__file__).parent.parent))
from conftest import GotoTestHelper, build_goto_binary


def test_simple_history():
    """Simple test to debug history functionality."""
    # Find goto binary
    script_dir = Path(__file__).parent.parent
    project_dir = script_dir.parent
    goto_binary = project_dir / "go" / "goto"
    
    if not goto_binary.exists():
        print("Building goto binary...")
        if not build_goto_binary(project_dir):
            print("Failed to build goto binary")
            return False
    
    # Create helper
    helper = GotoTestHelper(str(goto_binary))
    helper.setup_temp_env()
    
    try:
        # Create simple config
        config_content = '''
[test_dest]
path = "~/test_destination"
'''
        config_path = helper.create_config(config_content)
        print(f"Created config at: {config_path}")
        
        # Create simple history
        history_entries = [
            {
                "label": "test_dest",
                "last_used": datetime.now().isoformat() + "Z"
            }
        ]
        history_path = helper.create_history(history_entries)
        print(f"Created history at: {history_path}")
        
        # Check files exist
        print(f"Config exists: {os.path.exists(config_path)}")
        print(f"History exists: {os.path.exists(history_path)}")
        
        # Read files
        with open(config_path, 'r', encoding='utf-8') as f:
            print(f"Config content:\n{f.read()}")
        
        with open(history_path, 'r', encoding='utf-8') as f:
            print(f"History content:\n{f.read()}")
        
        # Test command
        print("\n=== Testing History Command ===")
        returncode, stdout, stderr = helper.run_goto(["--history"])
        
        print(f"Return code: {returncode}")
        print(f"STDOUT:\n{stdout}")
        print(f"STDERR:\n{stderr}")
        
        return returncode == 0 and "test_dest" in stdout
        
    finally:
        helper.cleanup_temp_env()


if __name__ == "__main__":
    success = test_simple_history()
    print(f"\nTest result: {'PASS' if success else 'FAIL'}")
    sys.exit(0 if success else 1)
