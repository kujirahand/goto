#!/bin/bash
# Bash completion script for goto command

_goto_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # Basic options
    opts="-h --help help --complete"

    if [[ ${cur} == -* ]]; then
        COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
        return 0
    fi

    # Get completion candidates from goto command
    if command -v goto >/dev/null 2>&1; then
        local candidates=$(goto --complete 2>/dev/null)
        COMPREPLY=($(compgen -W "${candidates}" -- ${cur}))
    fi

    return 0
}

# Register the completion function
complete -F _goto_completion goto
