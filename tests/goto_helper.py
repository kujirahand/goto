import os
import sys
import json
import subprocess

DIR_TESTS = os.path.dirname(os.path.abspath(__file__))
DIR_ROOT = os.path.dirname(DIR_TESTS)
FILE_GOTO = os.path.join(DIR_ROOT, "go", "goto")

FILE_CONFIG = "/tmp/goto/goto.toml"
FILE_HISTORY = "/tmp/goto/history.json"

def run(args, input_text=None):
    """Run the goto command with given arguments and optional input."""
    command = [FILE_GOTO] + args
    # Use with statement for resource management
    with subprocess.Popen(
        command,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        stdin=subprocess.PIPE,
        text=True
    ) as process:
        stdout, stderr = process.communicate(input=input_text)
        return process.returncode, stdout, stderr

def create_config(path, toml_str):
    """Create a TOML configuration file."""
    with open(path, 'w', encoding='utf-8') as f:
        f.write(toml_str)

def create_history(path, history_entries):
    """Create a history file with given entries."""
    history = {"entries": history_entries}
    with open(path, 'w', encoding='utf-8') as f:
        json.dump(history, f, indent=4, ensure_ascii=False)

def prepare_test():
    """Prepare the test environment."""
    tmp_goto = os.path.dirname(FILE_CONFIG)
    os.makedirs(tmp_goto, exist_ok=True)
    os.makedirs("/tmp/goto/dir1", exist_ok=True)
    os.makedirs("/tmp/goto/dir2", exist_ok=True)
    os.makedirs("/tmp/goto/dir3", exist_ok=True)
    config = """
[dir1]
path = "/tmp/goto/dir1"
[dir2]
path = "/tmp/goto/dir2"
[dir3]
path = "/tmp/goto/dir3"
"""
    history = {
        "entries": [
            {"label": "dir1", "last_used": "2025-01-01T12:00:03Z"},  # Most recent
            {"label": "dir2", "last_used": "2025-01-01T12:00:02Z"},
            {"label": "dir3", "last_used": "2025-01-01T12:00:01Z"}  # Oldest
        ]
    }
    # save
    create_config(FILE_CONFIG, config)
    create_history(FILE_HISTORY, history["entries"])

def load_config_org():
    """Load the original configuration."""
    path_config = os.path.expanduser("~/.goto.toml")
    if not os.path.exists(path_config):
        raise FileNotFoundError(f"Configuration file not found: {path_config}")
    with open(path_config, "r", encoding="utf-8") as f:
        config_content = f.read()
    return config_content

def load_history_org():
    """Load the original configuration."""
    path_config = os.path.expanduser("~/.goto.history.json")
    if not os.path.exists(path_config):
        raise FileNotFoundError(f"History file not found: {path_config}")
    with open(path_config, "r", encoding="utf-8") as f:
        config_content = f.read()
    return config_content

# 必ず実行してテストの準備をする
prepare_test()


if __name__ == "__main__":
    code, out, err = run(["--version"])
    if code != 0:
        print(f"Error running goto: {err}", file=sys.stderr)
        sys.exit(code)
    if out:
        print(f"goto version: {out.strip()}")
    # input_text
    code, out, err = run([], input_text="0\n")
    if code != 0:
        print(f"Error running goto: {err}", file=sys.stderr)
        sys.exit(code)
    if out:
        print(f"goto version: {out.strip()}")
