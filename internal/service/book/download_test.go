package book

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/tp86/legimi-go/internal/model"
)

type noopDownloadPresenter struct{}

func (noopDownloadPresenter) Start(model.BookMetadata) {}
func (noopDownloadPresenter) Part(model.BookMetadata)  {}
func (noopDownloadPresenter) End(model.BookMetadata)   {}
func (noopDownloadPresenter) Wait(model.BookMetadata)  {}

const testBookId = 1234

func testData(size int) []byte {
	data := make([]byte, size)
	rand.New(rand.NewSource(1)).Read(data)
	return data
}

// parseRange parses "bytes=start-end" request header and clamps it to data size
func parseRange(t *testing.T, r *http.Request, size int) (int, int) {
	var start, end int
	if _, err := fmt.Sscanf(r.Header.Get("range"), "bytes=%d-%d", &start, &end); err != nil {
		t.Errorf("invalid range header %q: %v", r.Header.Get("range"), err)
	}
	if end >= size {
		end = size - 1
	}
	return start, end
}

func writePartial(w http.ResponseWriter, data []byte, start, end int) {
	if start >= len(data) || start > end {
		w.Header().Set("content-range", fmt.Sprintf("bytes */%d", len(data)))
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(data)))
	w.Header().Set("content-length", fmt.Sprint(end-start+1))
	w.WriteHeader(http.StatusPartialContent)
	w.Write(data[start : end+1])
}

// writeEmpty sends successful response to range request without any data
func writeEmpty(w http.ResponseWriter, data []byte, start, end int) {
	w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(data)))
	w.Header().Set("content-length", "0")
	w.WriteHeader(http.StatusPartialContent)
}

// rangeHandler honours range requests, handlers below wrap it to misbehave
func rangeHandler(t *testing.T, data []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start, end := parseRange(t, r, len(data))
		writePartial(w, data, start, end)
	}
}

// nthRequest calls special for n-th request (counted from 0) and regular for all others
func nthRequest(n int, special, regular http.HandlerFunc) http.HandlerFunc {
	var mu sync.Mutex
	count := 0
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		current := count
		count++
		mu.Unlock()
		if current == n {
			special(w, r)
		} else {
			regular(w, r)
		}
	}
}

// runDownload downloads book served by handler into new temporary working directory
func runDownload(t *testing.T, handler http.Handler, size uint64) error {
	t.Helper()
	chdirToTempDir(t)
	return download(t, handler, size)
}

func chdirToTempDir(t *testing.T) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(wd) })
}

func download(t *testing.T, handler http.Handler, size uint64) error {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	oldDelay := downloadRetryDelay
	downloadRetryDelay = time.Millisecond
	t.Cleanup(func() { downloadRetryDelay = oldDelay })

	bs := defaultBookService{downloadPresenter: noopDownloadPresenter{}}
	done := make(chan error, 1)
	go func() {
		done <- bs.download(
			model.BookMetadata{Id: testBookId},
			model.BookDownloadDetails{Url: server.URL + "/book.mobi", Size: size},
		)
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(20 * time.Second):
		t.Fatal("download didn't finish in time")
		return nil
	}
}

func bookFileName() string {
	return fmt.Sprintf("%d.mobi", testBookId)
}

func assertDownloaded(t *testing.T, err error, data []byte) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content, err := os.ReadFile(bookFileName())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(content, data) {
		t.Errorf("downloaded file differs from served data: got %d bytes, expected %d", len(content), len(data))
	}
	assertNoFile(t, bookFileName()+".part")
}

func assertNoFile(t *testing.T, name string) {
	t.Helper()
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		t.Errorf("file %s should not exist, stat error: %v", name, err)
	}
}

func assertFailed(t *testing.T, err error, expectedInError string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), expectedInError) {
		t.Errorf("expected error containing %q, got: %v", expectedInError, err)
	}
	assertNoFile(t, bookFileName())
	assertNoFile(t, bookFileName()+".part")
}

// not a multiple of chunk size, so that last chunk is shorter
var testDataSize = 3*int(downloadChunkSize) + 12345

func TestDownload(t *testing.T) {
	data := testData(testDataSize)
	var requestedRanges []string
	handler := func(w http.ResponseWriter, r *http.Request) {
		requestedRanges = append(requestedRanges, r.Header.Get("range"))
		rangeHandler(t, data)(w, r)
	}

	err := runDownload(t, http.HandlerFunc(handler), uint64(len(data)))

	assertDownloaded(t, err, data)
	expectedRanges := []string{
		"bytes=0-81919",
		"bytes=81920-163839",
		"bytes=163840-245759",
		fmt.Sprintf("bytes=245760-%d", len(data)-1),
	}
	if fmt.Sprint(requestedRanges) != fmt.Sprint(expectedRanges) {
		t.Errorf("requested ranges %v, expected %v", requestedRanges, expectedRanges)
	}
}

func TestDownloadContinuesAfterShortResponse(t *testing.T) {
	data := testData(testDataSize)
	short := func(w http.ResponseWriter, r *http.Request) {
		start, end := parseRange(t, r, len(data))
		writePartial(w, data, start, start+(end-start)/2)
	}

	err := runDownload(t, nthRequest(1, short, rangeHandler(t, data)), uint64(len(data)))

	assertDownloaded(t, err, data)
}

func TestDownloadContinuesAfterConnectionCut(t *testing.T) {
	data := testData(testDataSize)
	cut := func(w http.ResponseWriter, r *http.Request) {
		start, end := parseRange(t, r, len(data))
		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(data)))
		w.Header().Set("content-length", fmt.Sprint(end-start+1))
		w.WriteHeader(http.StatusPartialContent)
		w.Write(data[start : start+1000])
		w.(http.Flusher).Flush()
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			t.Errorf("couldn't hijack connection: %v", err)
			return
		}
		conn.Close()
	}

	err := runDownload(t, nthRequest(1, cut, rangeHandler(t, data)), uint64(len(data)))

	assertDownloaded(t, err, data)
}

func TestDownloadRetriesTransientErrors(t *testing.T) {
	data := testData(testDataSize)
	unavailable := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	empty := func(w http.ResponseWriter, r *http.Request) {
		start, end := parseRange(t, r, len(data))
		writeEmpty(w, data, start, end)
	}
	handler := nthRequest(1, unavailable, nthRequest(2, empty, rangeHandler(t, data)))

	err := runDownload(t, handler, uint64(len(data)))

	assertDownloaded(t, err, data)
}

func TestDownloadWhenRangeIgnored(t *testing.T) {
	data := testData(testDataSize)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}

	err := runDownload(t, http.HandlerFunc(handler), uint64(len(data)))

	assertDownloaded(t, err, data)
}

func TestDownloadFailsOnErrorStatus(t *testing.T) {
	err := runDownload(t, http.NotFoundHandler(), uint64(testDataSize))

	assertFailed(t, err, "404")
}

func TestDownloadFailsWhenNoProgress(t *testing.T) {
	data := testData(testDataSize)
	requests := 0
	empty := func(w http.ResponseWriter, r *http.Request) {
		requests++
		start, end := parseRange(t, r, len(data))
		writeEmpty(w, data, start, end)
	}

	err := runDownload(t, nthRequest(0, rangeHandler(t, data), empty), uint64(len(data)))

	assertFailed(t, err, fmt.Sprintf("stalled at %d of %d bytes", downloadChunkSize, len(data)))
	if requests != maxDownloadAttemptsWithoutProgress {
		t.Errorf("made %d requests without progress, expected %d", requests, maxDownloadAttemptsWithoutProgress)
	}
}

func TestDownloadFailsOnRangeNotSatisfiable(t *testing.T) {
	data := testData(testDataSize)
	notSatisfiable := func(w http.ResponseWriter, r *http.Request) {
		writePartial(w, data, len(data), len(data))
	}

	err := runDownload(t, nthRequest(0, rangeHandler(t, data), notSatisfiable), uint64(len(data)))

	assertFailed(t, err, "416")
}

func TestDownloadFailsWhenServerHasLessData(t *testing.T) {
	data := testData(testDataSize)

	err := runDownload(t, rangeHandler(t, data), uint64(len(data)+100))

	assertFailed(t, err, "total size")
}

func TestDownloadFailsWhenServerHasMoreData(t *testing.T) {
	data := testData(testDataSize)

	err := runDownload(t, rangeHandler(t, data), uint64(len(data)-100))

	assertFailed(t, err, "total size")
}

func TestDownloadFailsWhenWholeFileHasWrongSize(t *testing.T) {
	data := testData(testDataSize)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}

	err := runDownload(t, http.HandlerFunc(handler), uint64(len(data)-100))

	assertFailed(t, err, fmt.Sprintf("downloaded %d bytes", len(data)))
}

func TestDownloadFailsWhenRangeIgnoredAfterStart(t *testing.T) {
	data := testData(testDataSize)
	whole := func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}

	err := runDownload(t, nthRequest(0, rangeHandler(t, data), whole), uint64(len(data)))

	assertFailed(t, err, "200 OK")
}

func TestDownloadFailsOnWrongContentRange(t *testing.T) {
	data := testData(testDataSize)
	fromStart := func(w http.ResponseWriter, r *http.Request) {
		writePartial(w, data, 0, int(downloadChunkSize)-1)
	}

	err := runDownload(t, http.HandlerFunc(fromStart), uint64(len(data)))

	assertFailed(t, err, "starts at 0, requested 81920")
}

func TestFailedDownloadKeepsExistingBook(t *testing.T) {
	chdirToTempDir(t)
	if err := os.WriteFile(bookFileName(), []byte("existing book"), 0644); err != nil {
		t.Fatal(err)
	}

	err := download(t, http.NotFoundHandler(), uint64(testDataSize))

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	content, err := os.ReadFile(bookFileName())
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "existing book" {
		t.Errorf("existing book was overwritten with %q", content)
	}
	assertNoFile(t, bookFileName()+".part")
}
