package main


import (
    "net/http"
    "fmt"
    "io"
    "sync"
    "time"
    "math/rand"
    "os"
    "strings"
    "flag"
)


func getUrl(url string) {
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(resp.StatusCode, len(body), url)
}


func parseUrlFile(urlsFile *string) *[]string {
    data, err := os.ReadFile(*urlsFile)
    if err != nil {
        panic(err)
    }
    urls := strings.Split(string(data), "\n")
    return &urls
}


func makeUrlList(rawUrls *[]string, rps *int, duration *int) []string {
    size := *rps * *duration
    urls := make([]string, 0, size)
    rUrls := *rawUrls
    for i := 0; i < size; i++ {
        urls = append(urls, rUrls[i % len(rUrls)])
    }
    return urls
}


func UrlRequester(wg *sync.WaitGroup, urlsChannel <-chan string) {
    defer wg.Done()
    for url := range urlsChannel {
        getUrl(url)
        time.Sleep(time.Duration(900 + rand.Intn(201)) * time.Millisecond)
    }
    fmt.Println("Stopping worker")
}


func UrlProducer(wg *sync.WaitGroup, urls *[]string, urlsChannel chan<- string) {
    defer wg.Done()
    defer close(urlsChannel)
    for _, url := range *urls {
        urlsChannel <- url
    }
    fmt.Println("Done with URLS")
}


func main() {
    rps := flag.Int("rps", 20, "Requests per second")
    duration := flag.Int("duration", 10, "Time in seconds to run")
    urlsFile := flag.String("urlsFile", "urls.tsv", "Newline-delimited list of URLs to GET")
    flag.Parse()
    fmt.Printf("Using %d requests/second for %d seconds\n", *rps, *duration)

    var wg sync.WaitGroup
    urlsChannel := make(chan string, 1000)

    rawUrls := parseUrlFile(urlsFile)
    urls := makeUrlList(rawUrls, rps, duration)
    fmt.Println("Got URLS", len(urls))

    fmt.Println("Starting URL producer")
    wg.Add(1)
    go UrlProducer(&wg, &urls, urlsChannel)

    fmt.Println("Starting workers")
    start := time.Now()
    for i := 0; i < *rps; i++ {
        fmt.Println("Worker", i)
        wg.Add(1)
        time.Sleep(time.Duration(rand.Intn(301)) * time.Millisecond)
        go UrlRequester(&wg, urlsChannel)
    }

    wg.Wait()

    took := time.Since(start)

    fmt.Println("Total time for workers to process URLS:", took)
}
