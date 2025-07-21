# test for history sorting in goto command
import os
import goto_helper as helper

def test_version():
    """Test that the version command works."""
    ret, out, err = helper.run(["--version"])
    assert ret == 0, f"Command failed with error: {err}"
    assert "goto version" in out, f"Expected version output but got: {out.strip()}"

def test_history():
    """Test that history entries are sorted by last used time."""
    ret, out, err = helper.run([
        "--config-file", helper.FILE_CONFIG, 
        "--history-file", helper.FILE_HISTORY,
        "--list-label"
    ])
    assert ret == 0, f"Command failed with error: {err}"
    assert out.strip() == "dir1\ndir2\ndir3", f"Expected order: dir1, dir2, dir3 but got: {out.strip()}"

def test_exec_change_history():
    """Test that executing a directory updates its history."""
    # Create config with multiple destinations
    helper.run([
        "--config-file", helper.FILE_CONFIG, 
        "--history-file", helper.FILE_HISTORY,
        "dir3",
    ])
    code, out, err = helper.run([
        "--config-file", helper.FILE_CONFIG,
        "--history-file", helper.FILE_HISTORY,
        "--list-label",
    ])
    out_list = out.strip().split('\n')
    assert code == 0, f"Command failed with error: {err}"
    assert out_list[0] == "dir3", f"Expected dir3 to be most recent but got: {out_list[0]}"
    assert out_list[1] == "dir1", f"Expected dir1 to be second but got: {out_list[1]}"
    assert out_list[2] == "dir2", f"Expected dir2 to be third but got: {out_list[2]}"
    

def test_history_org():
    """Test that history entries are sorted by last used time."""
    ret, out, err = helper.run([
        "--list-label"
    ])
    config_text = helper.load_history_org()
    out_list = out.strip().split('\n')
    # exec 3
    ret, out, err = helper.run([
        "3"
    ])
    # check history
    ret, out, err = helper.run([
        "--list-label"
    ])
    config_text2 = helper.load_history_org()
    out2_list = out.strip().split('\n')
    assert out2_list[0] == out_list[2], f"Expected {out_list[1]} to be most recent but got: {out2_list[0]}"
    assert config_text != config_text2, "Expected history file to be updated after executing a directory"

def test_history_org2():
    """Test that history entries are sorted by last used time."""
    ret, out, err = helper.run([
        "--list-label"
    ])
    config_text = helper.load_history_org()
    out_list = out.strip().split('\n')
    # exec 3
    ret, out, err = helper.run([
        "-l"
    ], input_text="3\n")
    # check history
    ret, out, err = helper.run([
        "--list-label"
    ])
    config_text2 = helper.load_history_org()
    out2_list = out.strip().split('\n')
    assert out2_list[0] == out_list[2], f"Expected {out_list[1]} to be most recent but got: {out2_list[0]}"
    assert config_text != config_text2, "Expected history file to be updated after executing a directory"
