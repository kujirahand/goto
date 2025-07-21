#!/usr/bin/env python3
"""
Test for interactive mode with partial history (mixed sorting).
This test verifies that directories with history are shown first (sorted by most recent),
followed by directories without history (sorted alphabetically).
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
        with open(config_file, 'w', encoding='utf-8') as f:
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

[epsilon]
path = "~/epsilon"
shortcut = "e"

[zeta]
path = "~/zeta"
shortcut = "z"
''')

        # Create history with only some entries
        # gamma (most recent), alpha (older)
        # beta, delta, epsilon, zeta should appear after in alphabetical order
        base_time = datetime.now()
        history_data = {
            "entries": [
                {
                    "label": "alpha",
                    "last_used": (base_time - timedelta(days=2)).strftime("%Y-%m-%dT%H:%M:%SZ")
                },
                {
                    "label": "gamma",
                    "last_used": (base_time - timedelta(minutes=30)).strftime("%Y-%m-%dT%H:%M:%SZ")  # Most recent
                }
            ]
        }

        with open(history_file, 'w', encoding='utf-8') as f:
            json.dump(history_data, f, indent=2)

        print("Created test config and partial history files")
        print(f"Config file: {config_file}")
        print(f"History file: {history_file}")
        
        print("\nHistory data:")
        with open(history_file, 'r', encoding='utf-8') as f:
            print(f.read())

        # Test interactive mode output by sending "0" (exit) as input
        goto_binary = "/Users/kujirahand/repos/goto/go/goto"
        cmd = [goto_binary, "--config", config_file, "--history-file", history_file, "-l"]

        print(f"\nRunning: {' '.join(cmd)}")
        print("Sending '0' to exit interactive mode...")
        
        # Send "0" as input to exit the interactive mode
        result = subprocess.run(cmd, input="0\n", capture_output=True, text=True, check=False)

        print(f"Return code: {result.returncode}")
        print(f"STDOUT:\n{result.stdout}")
        if result.stderr:
            print(f"STDERR:\n{result.stderr}")

        # Analyze the output to check order
        lines = result.stdout.split('\n')
        entry_lines = [line for line in lines if line.strip() and (
            line.startswith('1.') or line.startswith('2.') or line.startswith('3.') or 
            line.startswith('4.') or line.startswith('5.') or line.startswith('6.'))]
        
        print("\nExtracted entry lines:")
        for line in entry_lines:
            print(f"  {line}")

        # Expected order: gamma (history, most recent), alpha (history, older), 
        # then beta, delta, epsilon, zeta (alphabetical, no history)
        expected_labels = ["gamma", "alpha", "beta", "delta", "epsilon", "zeta"]
        
        print(f"\nExpected order: {expected_labels}")
        
        actual_order = []
        for line in entry_lines:
            for label in expected_labels:
                if label in line:
                    actual_order.append(label)
                    break
        
        print(f"Actual order: {actual_order}")
        
        if actual_order == expected_labels:
            print("✅ SUCCESS: Entries are correctly sorted (history first, then alphabetical)")
            return True
        else:
            print("❌ FAILURE: Entries are not correctly sorted")
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
