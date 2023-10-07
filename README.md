# Not-so-hot browser reloading

Serves static content from a target directory. Then monitors the target directory for changes and reloads the browser on change. 

Very scuffed and probably not stable, but seems to generally get the job done for my small use cases.

Compile to bin with `make build` and use for basic static site dev stuff. Should output to a `build` directory.

Usage

`tepid -path=myRelativeDirectory`

Args

```
-path : relative path to directory to monitor for reloads (default current directory)
-port : port number for the reloader to serve your files on (default 3000)
-interval : how many MS to wait between checking for file updates (default 250)
-exclude-ext: comma-delimited list of file extensions to exclude (default empty / "") (example: "tmpl, exe")
-exclude-dir: comma-delimited list of directory names to exclude (default empty / "") (example: "node_modules, example_data")
-exclude-files: comma-delimited list of specific file names to exclude (default empty / "") (example: ".env, script.js")
```

----


Or import it into a Go project to use for a dev environment. Check out the `examples` folder for a general idea of usage in a Go project.