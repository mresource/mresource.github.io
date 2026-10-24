# mresource

Free tools for resource teachers in Bulgarian schools. Site: https://mresource.github.io/

Free and with no sign-up. Data about the children never leaves the teacher's device. The tools are not diagnostic and do not replace a professional assessment.

## Structure

```
site/                         the static site (GitHub Pages)
  index.html                  home page
  practices/chertichki/       practice: one-minute reading and the four-point rule
  assets/logo.svg
go/                           the logic, in Go
  progress/                   aim line and the four-point rule, with tests
  cmd/wasm/                   bridge to the browser (WebAssembly)
.github/workflows/pages.yml   tests → WebAssembly → publish
```

## How publishing works

On every push to `main`, GitHub Actions:

1. checks formatting, then runs `go vet ./...` and `go test ./...`;
2. compiles `go/cmd/wasm` to `site/assets/progress.wasm`;
3. publishes `site/` to GitHub Pages.

If the tests fail, the site stays on the last working version.

## Running locally

```
cd go
go test ./...
GOOS=js GOARCH=wasm go build -o ../site/assets/progress.wasm ./cmd/wasm
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" ../site/assets/
cd ../site && python3 -m http.server 8000
```

Then open http://localhost:8000.

---

mpangarova](https://mpangarova.github.io/).
