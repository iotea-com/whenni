#!/bin/sh

yarn prisma generate

# Specify the original and new file names
original_file="./db/query-engine-linux-arm64-openssl-3.0.x_gen.go"
new_file="./db/query-engine-linux-static-arm64_gen.go"

# Check if the file exists
if [ -f "$original_file" ]; then
    # Rename the file
    mv "$original_file" "$new_file"

    # Modify a specific line in the file using awk
    awk '{
      if ($0 ~ /unpack.Unpack\(data, "linux-arm64-openssl-3.0.x"/) {
          sub("linux-arm64-openssl-3.0.x", "linux-static-arm64", $0)
      }
      print
    }' "$new_file" > temp_file && mv temp_file "$new_file"

    # Check if the original file still exists for any reason and remove it
    if [ -f "$original_file" ]; then
        rm "$original_file"
    fi

    echo "File has been successfully renamed and modified."

else
    echo "File does not exist."
fi

# Final check to ensure the original file is indeed removed
if [ ! -f "$original_file" ]; then
    echo "Original file successfully removed."
else
    echo "Error: Original file still exists."
fi
