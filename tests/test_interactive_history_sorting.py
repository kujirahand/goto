#!/usr/bin/env python3
"""
Test for interactive mode history sorting.
This test verifies that directories are sorted by most recently used in interactive mode.
"""

import os
import tempfile
import subprocess
import json
from datetime import datetime, timedelta

def main():
    # Create temporary files
    temp_dir = tempfile.mkdtemp(prefix="goto_test_")
    config_file = os.path.join(temp_dir, "test_config.toml")
    history_file = os.path.join(temp_dir, "test_history.json")

    try:
        # Create config with multiple destinations
        with open(config_file, 'w') as f:
            f.write('''[alpha]
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
''')

        # Create history with different timestamps
        # Expected order: beta (most recent), delta, alpha, gamma (oldest)
        base_time = datetime.now()
        history_data = {
            "entries": [
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
        }

        with open(history_file, 'w') as f:
            json.dump(history_data, f, indent=2)

        print("Created test config and history files")
        print(f"Config file: {config_file}")
        print(f"History file: {history_file}")
        
        print("\nHistory data:")
        with open(history_file, 'r') as f:
            print(f.read())

        # Test interactive mode output by sending "0" (exit) as input
        goto_binary = "/Users/kujirahand/repos/goto/go/goto"
        cmd = [goto_binary, "--config", config_file, "--history-file", history_file, "-l"]

        print(f"\nRunning: {' '.join(cmd)}")
        print("Sending '0' to exit interactive mode...")
        
        # Send "0" as input to exit the interactive mode
        result = subprocess.run(cmd, input="0\n", capture_output=True, text=True)

        print(f"Return code: {result.returncode}")
        print(f"STDOUT:\n{result.stdout}")
        if result.stderr:
            print(f"STDERR:\n{result.stderr}")

        # Analyze the output to check order
        lines = result.stdout.split('\n')
        entry_lines = [line for line in lines if line.strip() and (line.startswith('1.') or line.startswith('2.') or line.startswith('3.') or line.startswith('4.'))]
        
        print("\nExtracted entry lines:")
        for line in entry_lines:
            print(f"  {line}")

        # Expected order based on history: beta, delta, alpha, gamma
        expected_labels = ["beta", "delta", "alpha", "gamma"]
        
        print(f"\nExpected order: {expected_labels}")
        
        actual_order = []
        for line in entry_lines:
            for label in expected_labels:
                if label in line:
                    actual_order.append(label)
                    break
        
        print(f"Actual order: {actual_order}")
        
        if actual_order == expected_labels:
            print("✅ SUCCESS: Entries are correctly sorted by history (most recent first)")
            return True
        else:
            print("❌ FAILURE: Entries are not correctly sorted by history")
            print(f"Expected: {expected_labels}")
            print(f"Actual: {actual_order}")
            return False

    finally:
        # Cleanup
        import shutil
        shutil.rmtree(temp_dir)

if __name__ == "__main__":
    success = main()
    exit(0 if success else 1)
