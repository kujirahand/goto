#!/usr/bin/env python3
"""
Test for interactive mode without history file (alphabetical sorting).
This test verifies that directories are sorted alphabetically when no history exists.
"""

import os
import tempfile
import subprocess

def main():
    # Create temporary files
    temp_dir = tempfile.mkdtemp(prefix="goto_test_")
    config_file = os.path.join(temp_dir, "test_config.toml")

    try:
        # Create config with multiple destinations
        with open(config_file, 'w', encoding='utf-8') as f:
            f.write('''[zebra]
path = "~/zebra"
shortcut = "z"

[alpha]
path = "~/alpha"
shortcut = "a"

[beta]
path = "~/beta"
shortcut = "b"

[gamma]
path = "~/gamma"
shortcut = "g"
''')

        print("Created test config file without history")
        print(f"Config file: {config_file}")

        # Test interactive mode output by sending "0" (exit) as input
        goto_binary = "/Users/kujirahand/repos/goto/go/goto"
        cmd = [goto_binary, "--config", config_file, "-l"]

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
        entry_lines = [line for line in lines if line.strip() and (line.startswith('1.') or line.startswith('2.') or line.startswith('3.') or line.startswith('4.'))]
        
        print("\nExtracted entry lines:")
        for line in entry_lines:
            print(f"  {line}")

        # Expected order without history: alphabetical
        expected_labels = ["alpha", "beta", "gamma", "zebra"]
        
        print(f"\nExpected order (alphabetical): {expected_labels}")
        
        actual_order = []
        for line in entry_lines:
            for label in expected_labels:
                if label in line:
                    actual_order.append(label)
                    break
        
        print(f"Actual order: {actual_order}")
        
        if actual_order == expected_labels:
            print("✅ SUCCESS: Entries are correctly sorted alphabetically when no history exists")
            return True
        else:
            print("❌ FAILURE: Entries are not correctly sorted alphabetically")
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
