# go-load-test

Make 50 requests/second for 60 seconds to URLs specified in `urls.tsv`:

```bash
go run glt.go -rps=50 -duration=60 -urlsFile=urls.tsv
```