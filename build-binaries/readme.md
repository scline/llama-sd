Folder for building the golang dropbox binaries.

```
docker build . -t llama
docker run -v $(pwd)/bin:/go/bin llama
```