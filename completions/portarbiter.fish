complete -c portarbiter -l kill -d "Terminate the owner holding the port"
complete -c portarbiter -l force -d "Force termination"
complete -c portarbiter -l dry-run -d "Show what would be done"
complete -c portarbiter -l yes -s y -d "Auto-confirm dangerous actions"
complete -c portarbiter -l version -d "Show version information"
complete -c portarbiter -l help -d "Show help"

# Common ports
complete -c portarbiter -a "22 80 443 5432 6379 8080 9200"

