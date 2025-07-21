# goto Command Test Suite

This directory contains a comprehensive test suite for the goto command, implemented using Python and pytest.

## Overview

The test suite has been completely rewritten from bash scripts to Python for better maintainability, reliability, and extensibility. It uses pytest as the testing framework with additional plugins for coverage reporting.

## Test Structure

```
test/
├── conftest.py              # Pytest configuration and shared fixtures
├── run_tests.py             # Main test runner script
├── pyproject.toml           # Python project configuration
├── requirements.txt         # Python dependencies
└── tests/
    ├── __init__.py          # Package initialization
    ├── test_basic.py        # Basic functionality tests
    ├── test_advanced.py     # Advanced features and edge cases
    └── test_utils.py        # Legacy utilities (not actively used)
```

## Features

### Test Categories

1. **Basic Functionality** (`test_basic.py`)
   - Binary existence and execution
   - Help and version output
   - TOML configuration parsing
   - History functionality
   - Interactive modes (-c, -l)

2. **Advanced Features** (`test_advanced.py`)
   - Custom config and history file paths
   - Path expansion (tilde, relative, absolute)
   - Command execution configurations
   - URL destinations
   - Edge cases and error handling
   - Performance tests

### Key Features

- **Isolated Testing**: Each test uses temporary config and history files
- **Custom File Paths**: Tests the new `--config` and `--history-file` options
- **Comprehensive Coverage**: Tests normal operation, edge cases, and error conditions
- **Performance Testing**: Validates performance with large configurations
- **Interactive Mode Testing**: Tests cursor mode (-c) and label mode (-l)

## Usage

### Quick Start

```bash
# Run all tests
cd test
python3 run_tests.py
```

### Using pytest directly

```bash
# Install dependencies first
pip install pytest pytest-cov pytest-mock

# Run all tests with verbose output
pytest -v

# Run specific test file
pytest tests/test_basic.py -v

# Run specific test class
pytest tests/test_basic.py::TestConfigParsing -v

# Run with coverage report
pytest --cov=. --cov-report=html

# Run tests matching a pattern
pytest -k "config" -v
```

### Dependencies

The test suite requires:
- Python 3.8+
- pytest >= 7.0.0
- pytest-cov >= 4.0.0
- pytest-mock >= 3.10.0

Install dependencies:
```bash
cd test
pip install -r requirements.txt
```

## Test Environment

### Fixtures

- `goto_helper`: Main test helper with methods for running goto commands
- `sample_config`: Standard TOML configuration for testing
- `sample_history`: Sample history data for testing
- `ensure_goto_binary`: Ensures the goto binary is built before tests run

### GotoTestHelper Methods

```python
# Setup and teardown
helper.setup_temp_env()
helper.cleanup_temp_env()

# File creation
helper.create_config(toml_content)
helper.create_history(history_list)

# Command execution
returncode, stdout, stderr = helper.run_goto(args, timeout=5)
returncode, stdout, stderr = helper.run_goto_interactive(args, input_text="", timeout=2)
```

## Test Coverage

The test suite covers:

✅ **Configuration Parsing**
- Valid TOML configurations
- Invalid/malformed configurations
- Empty configurations
- Large configurations (50+ entries)
- Special characters in labels
- Custom config file paths

✅ **History Management**
- History file creation and reading
- Custom history file paths
- Large history files (1000+ entries)
- Empty history handling

✅ **Interactive Modes**
- Cursor mode (-c) option
- Label mode (-l) option
- Interactive input handling

✅ **Path Handling**
- Tilde expansion (~/)
- Relative paths (./path)
- Absolute paths (/path)
- URL destinations (https://)

✅ **Command Options**
- Help output (-h)
- Version output (--version)
- Completion output (--complete)
- History display (--history)
- Custom file options (--config, --history-file)

✅ **Error Handling**
- Non-existent config files
- Malformed TOML syntax
- Invalid command arguments
- Graceful failure modes

✅ **Performance**
- Startup time with large configs
- Memory usage with large history
- Timeout handling

## Benefits of Python Test Suite

### Compared to the previous bash scripts:

1. **Better Error Handling**: More precise error detection and reporting
2. **Isolation**: Each test runs in a clean environment with temporary files
3. **Maintainability**: Cleaner code structure with classes and fixtures
4. **Extensibility**: Easy to add new test cases and modify existing ones
5. **Coverage**: Built-in coverage reporting to identify untested code
6. **Cross-platform**: Works consistently across different operating systems
7. **IDE Support**: Better integration with IDEs and debugging tools
8. **Parallel Execution**: Can run tests in parallel for faster execution

### Testing the new features:

- **Custom Config Files**: Tests `--config FILE` option thoroughly
- **Custom History Files**: Tests `--history-file FILE` option thoroughly
- **Test Isolation**: No longer modifies user's actual ~/.goto.toml file
- **Interactive Modes**: Properly tests -c and -l options with timeout handling

## Running Specific Tests

```bash
# Test only configuration parsing
pytest tests/test_basic.py::TestConfigParsing -v

# Test only interactive modes
pytest tests/test_basic.py::TestInteractiveModes -v

# Test only advanced features
pytest tests/test_advanced.py::TestAdvancedFeatures -v

# Test only performance
pytest tests/test_advanced.py::TestPerformance -v

# Test with pattern matching
pytest -k "config" -v          # All config-related tests
pytest -k "history" -v         # All history-related tests
pytest -k "interactive" -v     # All interactive mode tests
```

## Continuous Integration

The test suite is designed to work well in CI environments:

```bash
# Example CI command
python3 run_tests.py || exit 1
```

All tests are non-interactive and complete within a reasonable time frame (typically under 1 second total).

## Legacy Bash Scripts

The original bash test scripts are still available for reference:
- `test.sh` - Basic functionality tests
- `functional_test.sh` - Functional tests  
- `performance_test.sh` - Performance tests
- `run_tests.sh` - Main test runner

However, the Python test suite is now the recommended approach for testing.
