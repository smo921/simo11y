package sources

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/ugorji/go/codec"

	"simo11y/internal/types"
)

// reads metrcis from addr and sends to returned channel
func Traces(done chan string, addr string) <-chan string {
	traceStream := make(chan string)
	go func() {
		defer close(traceStream)

		srv := &http.Server{Addr: addr, Handler: http.HandlerFunc(handle)}
		log.Printf("Serving on https://%s", addr)
		log.Fatal(srv.ListenAndServe())
	}()

	return traceStream
}

func handle(w http.ResponseWriter, r *http.Request) {
	// Log the request protocol
	defer r.Body.Close()

	traceCount, err := strconv.Atoi(r.Header["X-Datadog-Trace-Count"][0])
	if err != nil {
		fatalError(w, err)
		return
	} else if traceCount == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fatalError(w, err)
		return
	}

	var traces types.Traces
	var h codec.Handle = new(codec.MsgpackHandle)
	var dec *codec.Decoder = codec.NewDecoderBytes(contents, h)
	err = dec.Decode(&traces)
	if err != nil {
		log.Fatalf("Oops! Failed unmarshalling msgpack.\n %s", err)
		http.Error(w, err.Error(), 500)
	}

	w.WriteHeader(http.StatusOK)

	for traceNum, trace := range traces {
		for spanNum, span := range trace {
			fmt.Printf("Span %d - Trace: %d:\n%s\n", spanNum, traceNum, span.String())
		}
	}
}

func fatalError(w http.ResponseWriter, e error) {
	log.Fatalf("Oops! Failed reading body of the request.\n %s", e)
	http.Error(w, e.Error(), 500)
}
