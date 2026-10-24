# mresource

Безплатни инструменти за ресурсни учители. Сайт: https://mresource.github.io/

Безплатно и без регистрация. Данните за децата никога не напускат устройството на учителя. Инструментите не са диагностични и не заменят професионалната оценка.

## Структура

```
site/                         статичният сайт (GitHub Pages)
  index.html                  начална страница
  practices/chertichki/       практика: едноминутно четене и правилото с четирите точки
  assets/logo.svg
go/                           логиката, на Go
  progress/                   целева линия и правилото с четирите точки, с тестове
  cmd/wasm/                   мост към браузъра (WebAssembly)
.github/workflows/pages.yml   тестове → WebAssembly → публикуване
```

## Как работи публикуването

При всеки push в `main` GitHub Actions:

1. проверява форматирането и пуска `go test ./...`;
2. компилира `go/cmd/wasm` до `site/assets/progress.wasm`;
3. публикува `site/` на GitHub Pages.

Ако тестовете не минат, сайтът остава на последната работеща версия.

## Локално

```
cd go
go test ./...
GOOS=js GOARCH=wasm go build -o ../site/assets/progress.wasm ./cmd/wasm
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" ../site/assets/
cd ../site && python3 -m http.server 8000
```

После отвори http://localhost:8000.

---

Прави го [Мария Пангарова](https://mpangarova.github.io/).
