# bash completion for portarbiter

_portarbiter() {
    local cur prev
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    local opts="--kill --force --dry-run --yes -y --version --help"

    if [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
        return 0
    fi

    # complete port numbers
    if [[ ${COMP_CWORD} -eq 1 ]]; then
        COMPREPLY=( $(compgen -W "22 80 443 5432 6379 8080 9200" -- "${cur}") )
        return 0
    fi
}

complete -F _portarbiter portarbiter

