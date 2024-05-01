#!/bin/bash

# Build the ccwc binary
go build -o ccwc

# Compare counts and exit with code 1 if any difference is found
diff <(wc -l test.txt) <(./ccwc -l test.txt) && \
diff <(wc -w test.txt) <(./ccwc -w test.txt) && \
diff <(wc -c test.txt) <(./ccwc -c test.txt) && \
diff <(wc -l -w test.txt) <(./ccwc -l -w test.txt) && \
diff <(wc -l -c test.txt) <(./ccwc -l -c test.txt) && \
diff <(wc -w -c test.txt) <(./ccwc -w -c test.txt) && \
diff <(wc -l -w -c test.txt) <(./ccwc -l -w -c test.txt) || \
exit 1

echo "All tests passed!"

# Clean up the ccwc binary
rm ccwc
