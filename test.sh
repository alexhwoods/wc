#!/bin/bash

# Build the ccwc binary
go build -o ccwc

# Compare counts and exit with code 1 if any difference is found
echo "Running: diff <(wc test.txt) <(./ccwc test.txt)" && diff <(wc test.txt) <(./ccwc test.txt) && \
echo "Running: diff <(wc -l test.txt) <(./ccwc -l test.txt)" && diff <(wc -l test.txt) <(./ccwc -l test.txt) && \
echo "Running: diff <(wc -w test.txt) <(./ccwc -w test.txt)" && diff <(wc -w test.txt) <(./ccwc -w test.txt) && \
echo "Running: diff <(wc -c test.txt) <(./ccwc -c test.txt)" && diff <(wc -c test.txt) <(./ccwc -c test.txt) && \
echo "Running: diff <(wc -l -w test.txt) <(./ccwc -l -w test.txt)" && diff <(wc -l -w test.txt) <(./ccwc -l -w test.txt) && \
echo "Running: diff <(wc -l -c test.txt) <(./ccwc -l -c test.txt)" && diff <(wc -l -c test.txt) <(./ccwc -l -c test.txt) && \
echo "Running: diff <(wc -w -c test.txt) <(./ccwc -w -c test.txt)" && diff <(wc -w -c test.txt) <(./ccwc -w -c test.txt) && \
# This is the case I'm choosing to differ from wc on.
# echo "Running: diff <(wc -cm test.txt) <(./ccwc -cm test.txt)" && diff <(wc -cm test.txt) <(./ccwc -cm test.txt) && \
echo "Running: diff <(wc -mc test.txt) <(./ccwc -mc test.txt)" && diff <(wc -mc test.txt) <(./ccwc -mc test.txt) && \
echo "Running: diff <(wc -l -w -c test.txt) <(./ccwc -l -w -c test.txt)" && diff <(wc -l -w -c test.txt) <(./ccwc -l -w -c test.txt) || \
exit 1

echo "All tests passed!"

# Clean up the ccwc binary
rm ccwc
