#!/bin/bash

base_dir="."

find "$base_dir/cmd" "$base_dir/utils" "$base_dir/common" "$base_dir/config" -type f | while read -r file; do
    echo "This is '$file' content:"
    echo
    cat "$file"
    echo 
done
