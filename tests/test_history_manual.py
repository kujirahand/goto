import os
import tempfile
import subprocess
import json
from datetime import datetime

# Create temporary files
temp_dir = tempfile.mkdtemp(prefix="goto_test_")
config_file = os.path.join(temp_dir, "test_config.toml")
history_file = os.path.join(temp_dir, "test_history.json")

# Create config
with open(config_file, 'w') as f:
    f.write('''
[test_dest]
path = "~/test_destination"
''')

# Create history
history_data = [
    {
        "label": "test_dest",
        "last_used": datetime.now().isoformat() + "Z"
    }
]
with open(history_file, 'w') as f:
    json.dump(history_data, f)

print(f"Config file: {config_file}")
print(f"History file: {history_file}")

# Test goto command
goto_binary = "/Users/kujirahand/repos/goto/go/goto"
cmd = [goto_binary, "--config", config_file, "--history-file", history_file, "--history"]

print(f"Running: {' '.join(cmd)}")
result = subprocess.run(cmd, capture_output=True, text=True)

print(f"Return code: {result.returncode}")
print(f"STDOUT:\n{result.stdout}")
print(f"STDERR:\n{result.stderr}")

# Cleanup
import shutil
shutil.rmtree(temp_dir)
