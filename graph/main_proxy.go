package graph

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//--------------------------------------------------------------------------------------|

const (
	TARGET    = "https://programforyou.ru/graph-redactor"
	SERV_IP   = "127.0.0.1"
	SERV_PORT = uint16(8080)
)

//--------------------------------------------------------------------------------------|

func proxy(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse(TARGET)
	u.Path = "/graph-redactor" + r.URL.Path
	u.RawQuery = r.URL.RawQuery

	req, err := http.NewRequest(r.Method, u.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read response body", http.StatusInternalServerError)
		return
	}

	isGzip := strings.EqualFold(resp.Header.Get("Content-Encoding"), "gzip")
	if isGzip {
		body, err = modifyGzipBody(body)
		if err != nil {
			http.Error(w, "failed to modify gzip body", http.StatusInternalServerError)
			return
		}
	} else {
		body = replaceSaveFile(body) // если без сжатия
	}

	for k, v := range resp.Header {
		if strings.EqualFold(k, "Content-Length") {
			w.Header().Set("Content-Length", fmt.Sprint(len(body)))
		} else {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

//--------------------------------------------------------------------------------------|

func modifyGzipBody(gzipData []byte) ([]byte, error) {
	// Разархивируем
	gzReader, err := gzip.NewReader(bytes.NewReader(gzipData))
	if err != nil {
		return nil, err
	}
	decompressed, err := io.ReadAll(gzReader)
	gzReader.Close()
	if err != nil {
		return nil, err
	}

	// Вносим корректировки
	modified := replaceSaveFile(decompressed)

	// Повторно сжимаем
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	if _, err := gzWriter.Write(modified); err != nil {
		gzWriter.Close()
		return nil, err
	}
	gzWriter.Close()

	return buf.Bytes(), nil
}

//--------------------------------------------------------------------------------------|

func replaceSaveFile(body []byte) []byte {
	old := []byte(`GraphRedactor.prototype.SaveFile = function() {
    let graph = {
        x0: this.x0,
        y0: this.y0,
        vertices: this.vertices.ToJSON(),
        edges: this.edges.map((edge) => edge.ToJSON(this.vertices)),
        texts: this.texts.map((text) => text.ToJSON())
    }

    this.SaveObject(new Blob([JSON.stringify(graph)], { type: 'application/octet-stream' }), SAVE_FILE_NAME)
}`)

	newCode := []byte(`// lignigno
GraphRedactor.prototype.SaveFile = function() {
    const vertices = this.vertices.ToJSON();
    const lines = vertices.map(v => ` + "`${v.name} ${v.x} ${v.y}`" + `);

    const edges = this.edges.map(e => {
		tmpEdges = e.ToJSON(this.vertices)
        const from = vertices[tmpEdges.vertex1].name;
        const to = vertices[tmpEdges.vertex2].name;
        return ` + "`${from}-${to}`" + `;
    });

    const content = [...lines, '', ...edges].join('\n');
    this.SaveObject(new Blob([content], { type: 'text/plain' }), "graph.txt");
}`)

	newBody := bytes.Replace(body, old, newCode, 1)
	if len(newBody) == len(body) {
		return append(body, []byte("\n\033[38;2;255;0;128m<!-- Unicorn!!! -->\033[m\n")...)
	}

	return bytes.Replace(newBody, old, old, 1)
}

//--------------------------------------------------------------------------------------|

func main() {
	http.HandleFunc("/", proxy)
	addr := fmt.Sprintf("%s:%d", SERV_IP, SERV_PORT)

	fmt.Printf("server running {\033[38;2;255;0;128mhttp://%s\033[m}\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
