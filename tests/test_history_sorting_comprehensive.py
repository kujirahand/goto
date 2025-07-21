#!/usr/bin/env python3
"""
Comprehensive test suite for goto history-based sorting functionality.
This test covers multiple scenarios:
1. Full history sorting (most recent first)
2. No history sorting (alphabetical)
3. Mixed history sorting (history first, then alphabetical)
"""

import os
import sys
import tempfile
import subprocess
import json
from datetime import datetime, timedelta

def create_test_case(name, config_content, history_data=None):
    """Create a test case with config and optional history files."""
    temp_dir = tempfile.mkdtemp(prefix=f"goto_test_{name}_")
    config_file = os.path.join(temp_dir, "test_config.toml")
    history_file = os.path.join(temp_dir, "test_history.json")
    
    with open(config_file, 'w', encoding='utf-8') as f:
        f.write(config_content)
    
    if history_data:
        with open(history_file, 'w', encoding='utf-8') as f:
            json.dump(history_data, f, indent=2)
        return temp_dir, config_file, history_file
    else:
        return temp_dir, config_file, None

def run_goto_test(config_file, history_file=None):
    """Run goto command and return the output."""
    goto_binary = "/Users/kujirahand/repos/goto/go/goto"
    cmd = [goto_binary, "--config", config_file, "-l"]
    
    if history_file:
        cmd.extend(["--history-file", history_file])
    
    result = subprocess.run(cmd, input="0\n", capture_output=True, text=True, check=False)
    
    # Extract entry lines
    lines = result.stdout.split('\n')
    entry_lines = [line for line in lines if line.strip() and 
                   any(line.startswith(f'{i}.') for i in range(1, 10))]
    
    # Extract labels from entry lines
    actual_order = []
    for line in entry_lines:
        # Extract label from format "1. label (shortcut) → path"
        parts = line.split(' ')
        if len(parts) >= 2:
            label_part = parts[1]
            # Remove shortcut part if exists
            if '(' in label_part:
                label = label_part.split('(')[0].strip()
            else:
                label = label_part.strip()
            actual_order.append(label)
    
    return actual_order

def test_full_history_sorting():
    """Test sorting when all entries have history."""
    print("🧪 Test 1: Full history sorting")
    
    config_content = '''[alpha]
path = "~/alpha"

[beta]
path = "~/beta"

[gamma]
path = "~/gamma"

[delta]
path = "~/delta"
'''
    
    base_time = datetime.now()
    history_data = {
        "entries": [
            {"label": "alpha", "last_used": (base_time - timedelta(days=2)).strftime("%Y-%m-%dT%H:%M:%SZ")},
            {"label": "beta", "last_used": (base_time - timedelta(minutes=30)).strftime("%Y-%m-%dT%H:%M:%SZ")},
            {"label": "gamma", "last_used": (base_time - timedelta(days=5)).strftime("%Y-%m-%dT%H:%M:%SZ")},
            {"label": "delta", "last_used": (base_time - timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")}
        ]
    }
    
    temp_dir, config_file, history_file = create_test_case("full_history", config_content, history_data)
    
    try:
        actual_order = run_goto_test(config_file, history_file)
        expected_order = ["beta", "delta", "alpha", "gamma"]  # Most recent first
        
        print(f"  Expected: {expected_order}")
        print(f"  Actual:   {actual_order}")
        
        if actual_order == expected_order:
            print("  ✅ PASS: Full history sorting works correctly")
            return True
        else:
            print("  ❌ FAIL: Full history sorting failed")
            return False
    finally:
        import shutil
        shutil.rmtree(temp_dir)

def test_no_history_sorting():
    """Test sorting when no history exists."""
    print("\n🧪 Test 2: No history sorting (alphabetical)")
    
    config_content = '''[zebra]
path = "~/zebra"

[alpha]
path = "~/alpha"

[gamma]
path = "~/gamma"

[beta]
path = "~/beta"
'''
    
    temp_dir, config_file, _ = create_test_case("no_history", config_content)
    
    try:
        actual_order = run_goto_test(config_file)
        expected_order = ["alpha", "beta", "gamma", "zebra"]  # Alphabetical
        
        print(f"  Expected: {expected_order}")
        print(f"  Actual:   {actual_order}")
        
        if actual_order == expected_order:
            print("  ✅ PASS: No history sorting works correctly")
            return True
        else:
            print("  ❌ FAIL: No history sorting failed")
            return False
    finally:
        import shutil
        shutil.rmtree(temp_dir)

def test_mixed_history_sorting():
    """Test sorting when some entries have history."""
    print("\n🧪 Test 3: Mixed history sorting")
    
    config_content = '''[alpha]
path = "~/alpha"

[beta]
path = "~/beta"

[gamma]
path = "~/gamma"

[delta]
path = "~/delta"

[epsilon]
path = "~/epsilon"
'''
    
    base_time = datetime.now()
    history_data = {
        "entries": [
            {"label": "gamma", "last_used": (base_time - timedelta(minutes=30)).strftime("%Y-%m-%dT%H:%M:%SZ")},
            {"label": "alpha", "last_used": (base_time - timedelta(days=1)).strftime("%Y-%m-%dT%H:%M:%SZ")}
        ]
    }
    
    temp_dir, config_file, history_file = create_test_case("mixed_history", config_content, history_data)
    
    try:
        actual_order = run_goto_test(config_file, history_file)
        expected_order = ["gamma", "alpha", "beta", "delta", "epsilon"]  # History first, then alphabetical
        
        print(f"  Expected: {expected_order}")
        print(f"  Actual:   {actual_order}")
        
        if actual_order == expected_order:
            print("  ✅ PASS: Mixed history sorting works correctly")
            return True
        else:
            print("  ❌ FAIL: Mixed history sorting failed")
            return False
    finally:
        import shutil
        shutil.rmtree(temp_dir)

def main():
    """Run all tests."""
    print("🚀 Running goto history sorting test suite...")
    print("=" * 60)
    
    results = []
    results.append(test_full_history_sorting())
    results.append(test_no_history_sorting())
    results.append(test_mixed_history_sorting())
    
    print("\n" + "=" * 60)
    print("📊 Test Results Summary:")
    
    passed = sum(results)
    total = len(results)
    
    if passed == total:
        print(f"🎉 All {total} tests PASSED!")
        print("\n✅ History-based sorting is working correctly:")
        print("   • Entries with history are sorted by most recent usage first")
        print("   • Entries without history are sorted alphabetically")
        print("   • Mixed scenarios work as expected")
        return True
    else:
        print(f"❌ {total - passed} out of {total} tests FAILED!")
        return False

if __name__ == "__main__":
    success = main()
    sys.exit(0 if success else 1)
