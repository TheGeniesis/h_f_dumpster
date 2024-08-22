#!/bin/bash

# Get the list of user directories in the workspace
USER_DIRECTORIES=$(databricks workspace list --absolute /Users)

# Export notebooks for each user directory
for USER_DIRECTORY in $USER_DIRECTORIES; do
    echo "Exporting notebooks for user directory: $USER_DIRECTORY"
    OUTPUT_DIR="./imports/notebooks/$(basename $USER_DIRECTORY)"

    # Create the output directory
    mkdir -p $OUTPUT_DIR

    # Export notebooks for the user directory
    databricks workspace export_dir $USER_DIRECTORY $OUTPUT_DIR  -o
done
