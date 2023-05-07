# go-docs-go
Generating scrappy docs from source code with Go

## Usage

```
./go-docs-go [-o outputDir] source1, source2, ...
```

If specified by `-o` or `--outputDir`, the documents generated are written to the `{outputDir}` directory. Otherwise,
they are written to `./docs`.

The last arguments are assumed to be source file paths to analyze and document.
