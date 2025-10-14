#!/bin/bash

# Name of the input HTML file
input_html="coverage/index-application.html"

# Check if the input file exists
if [ ! -f "$input_html" ]; then
    echo "Error: File '$input_html' does not exist."
    exit 1
fi

# Create output directory if it doesn't exist
output_dir="coverage_reports"
mkdir -p "$output_dir"

# Use htmlq and mapfile to read content and IDs from <option> tags
mapfile -t option_texts < <(htmlq --text 'option' < "$input_html")
mapfile -t pre_ids < <(htmlq --attribute value 'option' < "$input_html")

# Loop through the array indices
for i in "${!pre_ids[@]}"; do
  # Get the ID of the <pre> tag and the full text content
  id="${pre_ids[$i]}"
  full_text="${option_texts[$i]}"

  # Remove the coverage percentage at the end, example: ' (100.0%)'
  path_only=$(echo "$full_text" | sed 's/ *(.*)$//')

  # Use 'basename' to get only the filename from the path
  short_name=$(basename "$path_only")

  # Create the final filename with .html extension
  filename="${short_name}.html"

  # Use htmlq to extract the <pre> tag content and save to file
  htmlq "#${id}" < "$input_html" > "${output_dir}/${filename}"

done

echo "Completed! Files have been created in the '$output_dir' directory."