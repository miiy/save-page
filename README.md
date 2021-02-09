# save-page

Save page ad html or pdf from link

## Build

```bash
docker build . -t miiy/save-page:latest
```

## Save as html

```bash
docker run --rm -it -v "$PWD"/data:/app/data miiy/save-page -t html http://test.com/index.html
```

## Save as pdf

```bash
docker run --rm -it -v "$PWD"/data:/app/data miiy/save-page -t pdf https://www.test.com test.pdf
```

## Debug

```bash
docker run --rm -it -v "$PWD":/go/build -v "$PWD"/data:/app/data miiy/save-page bash
cd /go/build
```