# gotee
It's like "tee" but uses regex and a print server. Basic usage:
```
gotee "REGEXP"
```

For example:
```
ping localhost | gotee "(?P<bytes>\d+) bytes from (?P<host>\w+)"
```
Will print out a JSON stream, with groups as keys and the whole match under the key "match".

Templates can also be used:
```
ping localhost | gotee -f "{{.host}} - {{.bytes}} bytes" "(?P<bytes>\d+) bytes from (?P<host>\w+)"
```
