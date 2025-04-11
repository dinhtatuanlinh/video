#!/bin/bash
read -p "Enter the description: " DESCRIPTION
# Check if a description is provided
if [ -z "$DESCRIPTION" ]; then
  echo "Please provide a migration description."
  exit 1
fi


# Get the current timestamp
TIMESTAMP=$(date +"%Y%m%d%H%M%S")

# Format the description by replacing spaces with underscores
DESCRIPTION=$(echo "$DESCRIPTION" | tr " " "_")

# Define the migration directory
MIGRATION_DIR="./db/migration"

# Create the migration directory if it doesn't exist
mkdir -p "$MIGRATION_DIR"

# Create the file names
UP_FILE="${MIGRATION_DIR}/${TIMESTAMP}_${DESCRIPTION}.up.sql"
DOWN_FILE="${MIGRATION_DIR}/${TIMESTAMP}_${DESCRIPTION}.down.sql"

# Create the files
touch "$UP_FILE" "$DOWN_FILE"

echo "Created migration files:"
echo "$UP_FILE"
echo "$DOWN_FILE"