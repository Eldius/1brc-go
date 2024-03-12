# 1BRC #

## links ##

- [1BRC Site](https://mrkaran.dev/posts/1brc/)
- [1BRC Github Repository](https://github.com/gunnarmorling/1brc)


## code snippets ##

```shell
# count unique locations
sed 's/;[0-9.\-]*//g' internal/parser/sample_data/measurements_1m.txt  | sort | uniq | wc -l

# count location ocurrences
sed 's/;[0-9.\-]*//g' internal/parser/sample_data/measurements_1m.txt  | sort | uniq -c
```
