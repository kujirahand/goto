#!/usr/bin/env python3
"""goto command - Python version with new shell functionality"""
import os
import sys
import subprocess
import tempfile

try:
    import tomli
except ImportError:
    print("❌ Error: tomli module is required. Please install it with: pip3 install tomli")
    sys.exit(1)

from goto_template import DEFAULT_CONFIG


def create_default_config(toml_file):
    """Create default TOML configuration file if it doesn't exist"""
    with open(toml_file, "w", encoding="utf-8") as f:
        f.write(DEFAULT_CONFIG)
    print(f"Created default configuration file: {toml_file}")


def get_user_choice(entries, shortcut_map):
    """Get user's choice for destination"""
    print("👉 Available destinations:")
    for i, (label, values) in enumerate(entries, start=1):
        shortcut = values.get("shortcut", "")
        path = os.path.expanduser(values.get("path", ""))
        print(f"{i}. {label} → {path} (shortcut: {shortcut})")

    print("\nPlease enter the number or shortcut key of the destination you want to go to:")
    
    try:
        choice = input("Enter number or shortcut key: ").strip()
    except (EOFError, KeyboardInterrupt):
        print("\nOperation cancelled.")
        return None, None, None

    # Determine input and get corresponding entry
    index = None
    if choice.isdigit():
        index = int(choice)
    elif choice in shortcut_map:
        index = shortcut_map[choice]

    if index and 1 <= index <= len(entries):
        label, values = entries[index - 1]
        path = os.path.expanduser(values["path"])
        command = values.get("command")
        return path, command, label
    else:
        print("Invalid input.")
        return None, None, None


def open_new_shell(target_dir, command=None, label=None):
    """Open a new shell in the target directory"""
    if not os.path.isdir(target_dir):
        print(f"❌ Directory does not exist: {target_dir}")
        return False

    print(f"🚀 Opening new shell in: {target_dir}")
    if label:
        print(f"📍 Destination: {label}")
    
    # Get user's preferred shell
    user_shell = os.environ.get('SHELL', '/bin/sh')
    
    try:
        if command:
            print(f"⚡ Will execute: {command}")
            print("=" * 50)
            
            # Create a temporary startup script
            with tempfile.NamedTemporaryFile(mode='w', suffix='.sh', delete=False) as temp_file:
                temp_script = temp_file.name
                temp_file.write(f"""#!/bin/sh
cd "{target_dir}"
echo "📍 Current directory: $(pwd)"
echo "⚡ Executing: {command}"
echo "{'-' * 40}"
{command}
echo "{'-' * 40}"
echo "✅ Command completed. You are now in: $(pwd)"
echo "💡 Type 'exit' to return to previous shell"
exec "{user_shell}"
""")
            
            # Make the script executable
            os.chmod(temp_script, 0o755)
            
            try:
                # Execute the temporary script
                subprocess.run(['/bin/sh', temp_script], cwd=target_dir, check=False)
            finally:
                # Clean up the temporary script
                if os.path.exists(temp_script):
                    os.unlink(temp_script)
        else:
            # Simply open shell in the target directory
            print("💡 Type 'exit' to return to previous shell")
            print("=" * 50)
            os.chdir(target_dir)
            print(f"✅ You are now in: {os.getcwd()}")
            subprocess.run([user_shell], check=False)
        
        return True
        
    except KeyboardInterrupt:
        print("\n🛑 Shell session interrupted.")
        return True
    except (FileNotFoundError, PermissionError, OSError) as e:
        print(f"❌ Error opening shell: {e}")
        return False


def main():
    """Main function to handle directory navigation with new shell"""
    # Load TOML file
    toml_file = os.path.expanduser("~/.goto.toml")

    # Create default TOML file if it doesn't exist
    if not os.path.exists(toml_file):
        create_default_config(toml_file)

    try:
        with open(toml_file, "rb") as f:
            config = tomli.load(f)
    except (FileNotFoundError, PermissionError, tomli.TOMLDecodeError) as e:
        print(f"❌ Error reading configuration file: {e}")
        sys.exit(1)

    # Get list of entries
    entries = list(config.items())
    if not entries:
        print("⚠️  No destinations configured in ~/.goto.toml")
        sys.exit(1)

    # Build shortcut map
    shortcut_map = {}
    for i, (label, values) in enumerate(entries, start=1):
        shortcut = values.get("shortcut", "")
        if shortcut:
            shortcut_map[shortcut] = i

    # Get user choice
    target_dir, command, label = get_user_choice(entries, shortcut_map)
    
    if target_dir is None:
        print("ℹ️  No directory selected or operation cancelled.")
        sys.exit(0)

    # Open new shell in the selected directory
    success = open_new_shell(target_dir, command, label)
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
