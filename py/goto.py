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

    print("\n➕ [+] Add current directory")
    print("\nPlease enter the number, shortcut key, label name, or [+] to add current directory:")
    
    try:
        choice = input("Enter number, shortcut key, label name, or [+]: ").strip()
    except (EOFError, KeyboardInterrupt):
        print("\nOperation cancelled.")
        return None, None, None

    # Check if user wants to add current directory
    if choice == "+":
        return "ADD_CURRENT", None, None

    # Determine input and get corresponding entry
    index = None
    if choice.isdigit():
        index = int(choice)
    elif choice in shortcut_map:
        index = shortcut_map[choice]
    else:
        # Check if it's a label name (case-insensitive)
        for i, (label, values) in enumerate(entries, start=1):
            if label.lower() == choice.lower():
                index = i
                break

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


def add_current_path_to_config(toml_file):
    """Add current directory to the TOML configuration file"""
    current_dir = os.getcwd()
    
    # Get label name from user
    print(f"📍 Current directory: {current_dir}")
    try:
        label = input("Enter a label for this directory: ").strip()
        if not label:
            print("❌ Label cannot be empty.")
            return False
            
        shortcut = input("Enter a shortcut key (optional, press Enter to skip): ").strip()
        
        # Generate TOML entry
        toml_entry = f"\n[{label}]\n"
        toml_entry += f'path = "{current_dir}"\n'
        if shortcut:
            toml_entry += f'shortcut = "{shortcut}"\n'
        
        # Append to TOML file
        with open(toml_file, "a", encoding="utf-8") as f:
            f.write(toml_entry)
        
        print(f"✅ Added '{label}' → {current_dir}")
        if shortcut:
            print(f"🔑 Shortcut: {shortcut}")
        
        return True
        
    except (EOFError, KeyboardInterrupt):
        print("\n❌ Operation cancelled.")
        return False
    except (OSError, IOError) as e:
        print(f"❌ Error adding path: {e}")
        return False


def find_destination_by_arg(arg, entries, shortcut_map):
    """Find destination by command line argument (number, label or shortcut)"""
    # Check if it's a number
    if arg.isdigit():
        index = int(arg)
        if 1 <= index <= len(entries):
            label, values = entries[index - 1]
            path = os.path.expanduser(values["path"])
            command = values.get("command")
            return path, command, label
        else:
            return None, None, None
    
    # Check if it's a shortcut
    if arg in shortcut_map:
        index = shortcut_map[arg]
        label, values = entries[index - 1]
        path = os.path.expanduser(values["path"])
        command = values.get("command")
        return path, command, label
    
    # Check if it's a label (case-insensitive)
    for label, values in entries:
        if label.lower() == arg.lower():
            path = os.path.expanduser(values["path"])
            command = values.get("command")
            return path, command, label
    
    return None, None, None


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

    # Check for command line argument
    if len(sys.argv) > 1:
        arg = sys.argv[1]
        
        # Handle help option
        if arg in ['-h', '--help', 'help']:
            print("🚀 goto - Navigate directories quickly")
            print("\nUsage:")
            print("  goto                 Show interactive menu")
            print("  goto <number>        Go to destination by number (e.g., goto 1)")
            print("  goto <label>         Go to destination by label name")
            print("  goto <shortcut>      Go to destination by shortcut key")
            print("  goto -h, --help      Show this help message")
            print("\nExamples:")
            print("  goto 1              # Navigate to 1st destination")
            print("  goto Home           # Navigate to 'Home' destination")
            print("  goto h              # Navigate using shortcut 'h'")
            print("  goto                # Show interactive menu")
            sys.exit(0)
        
        # Find destination by argument
        target_dir, command, label = find_destination_by_arg(arg, entries, shortcut_map)
        
        if target_dir is None:
            print(f"❌ Destination '{arg}' not found.")
            print("\n📋 Available destinations:")
            for i, (label, values) in enumerate(entries, start=1):
                shortcut = values.get("shortcut", "")
                path = os.path.expanduser(values.get("path", ""))
                shortcut_str = f" (shortcut: {shortcut})" if shortcut else ""
                print(f"  • {label}{shortcut_str} → {path}")
            sys.exit(1)
        
        print(f"🎯 Found destination: {label}")
        # Open new shell in the found directory
        success = open_new_shell(target_dir, command, label)
        sys.exit(0 if success else 1)

    # Get user choice
    target_dir, command, label = get_user_choice(entries, shortcut_map)
    
    # Handle adding current directory
    if target_dir == "ADD_CURRENT":
        success = add_current_path_to_config(toml_file)
        sys.exit(0 if success else 1)
    
    if target_dir is None:
        print("ℹ️  No directory selected or operation cancelled.")
        sys.exit(0)

    # Open new shell in the selected directory
    success = open_new_shell(target_dir, command, label)
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
