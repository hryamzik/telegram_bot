#!/bin/bash
files=$(find testdata -name '*.txt')

messages[0]=$'word'
messages[1]=$'hello word'
messages[2]=$'hello\nline'

function performTest {
    result="$(curl -sw '\n%{http_code}' --data-binary "$1" 127.0.0.1:9037/${TELEGRAM_CHATID})"
    code=$(echo "$result"|tail -1)
    responce=$(echo "$result"|head -n -1 )
    test "$code" -eq 200 && msg="ok" || msg="not ok"
    echo "$msg $((++i)) -" ${1#@testdata/}
    echo $responce|jq -C . 2>/dev/null| sed -E 's/^(.)/# \1/'
}

./telegram_bot -d > bot.log 2>&1 &
sleep 3

echo "1..$(( $(echo "$files"|wc -l) + ${#messages[@]}))"

for msg in "${messages[@]}"
do
    performTest "$msg"
done

for file in $files
do
    performTest "@$file"
done

kill $! 2>/dev/null
wait $! 2>/dev/null

exit 0