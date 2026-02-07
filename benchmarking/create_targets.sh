OUTPUT="targets.txt"
> "$OUTPUT"

N="${1:-10000}"

i=1
while [ "$i" -le "$N" ]; do
    echo "GET http://localhost:8080/" >> "$OUTPUT"
    echo "X-Forwarded-For: 10.0.$((i/256)).$((i%256))" >> "$OUTPUT"
    echo "" >> "$OUTPUT"
    i=$((i+1))
done

echo "Generated $OUTPUT with $N unique IPs"
