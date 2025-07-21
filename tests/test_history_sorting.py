import os
import tempfile
import subprocess
import json
from datetime import datetime, timedelta

# Create temporary files
temp_dir = tempfile.mkdtemp(prefix="goto_test_")
config_file = os.path.join(temp_dir, "test_config.toml")
history_file = os.path.join(temp_dir, "test_history.json")

# Create config with multiple destinations
with open(config_file, 'w') as f:
    f.write('''
[home]
path = "~/"

[projects]
path = "~/Projects"

[documents]
path = "~/Documents"

[downloads]
path = "~/Downloads"
''')

# Create history with different timestamps
base_time = datetime.now()
history_data = {
    "entries": [
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
}

with open(history_file, 'w') as f:
    json.dump(history_data, f, indent=2)

print("History content:")
with open(history_file, 'r') as f:
    print(f.read())

# Test goto command
goto_binary = "/Users/kujirahand/repos/goto/go/goto"
cmd = [goto_binary, "--config", config_file, "--history-file", history_file, "--history"]

print(f"\nRunning: {' '.join(cmd)}")
result = subprocess.run(cmd, capture_output=True, text=True)

print(f"Return code: {result.returncode}")
print(f"STDOUT:\n{result.stdout}")
print(f"STDERR:\n{result.stderr}")

# Cleanup
import shutil
shutil.rmtree(temp_dir)
