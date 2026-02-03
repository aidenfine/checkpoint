OUTPUT="targets.txt"
> $OUTPUT

for i in {1..10000}; do
    echo "GET http://localhost:8080/" >> $OUTPUT
    echo "X-Forwarded-For: 10.0.$((i/256)).$((i%256))" >> $OUTPUT
    echo "" >> $OUTPUT
done

echo "Generated $OUTPUT with 10000 unique IPs"