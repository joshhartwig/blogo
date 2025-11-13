---
title: "Building a blog in Go Part 1"
date: 2025-11-03
summary: "Building a simple Blog in 3 parts"
draft: false
tags:
  - "#Go"
  - "#Development"
---

Let's build a blog!

Why a blog? So boring and simple. When the rubber meets the road, most things are boring and (sometimes) simple. Complex systems are usually many simple systems coaxed together. Also, I am not the best programmer on the planet. I learn by creating stuff I would actually use.

So what will you learn building a blog? A suprising amount actually. Just off the top of my head...

* HTTP servers
* Routing
* HTML templates
* Serving Markdown

Now time to get going!

1. open the terminal, create a new directory, cd into the directory and go mod init it

```bash
mkdir example_blog && cd example_blog && go mod init github.com/yourname/exampleblog && touch main.go
```

2.open the folder in vscode or your code editor of choice
3.open the main.go and add the following

```go
package main

import "fmt"

func main() {
 fmt.Println("hello world")
}
```

open the integrated terminal in the code editor or bring up a terminal in the same folder and run:

```bash
go run .
```

You should see `hello world` print in the terminal. Lets move a little quicker and lets get basic server going.

```go
import (
 "fmt"
 "net/http"
)

func main() {
 // register a handle function for the route default "/" and write the string "hi" in a byte slice
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("hello world"))
 })

 // use the default serve mux to start a server on port 5040 and check for errors
 if err := http.ListenAndServe(":5040", nil); err != nil {
  fmt.Printf("server crashed with error %v", err)
 }
}
```

We now have a very basic http server that writes `hello world` on port 5040. Go ahead and test it in your browser or via terminal with curl

```bash
curl http://localhost:5040/
# hello world
```

Let's think of the things we need to build the starts of a blog. We can deliver some content to a browser at this point but it is not very useful. Let's work on reading the markdown files.

Create a `content` folder in the terminal or in the code editor and add this file and it's content.

```mark
<!-- firstpost.md -->
# Hello World

This is a first post! There are many others but this is mine.
```

Now would be a good time to talk about the strucutre of this project. For our initial iteration we are going to serve very basic html at the default route `/` and list out our posts that are read from the content directory. Each post will render a link and clicking the link will take you to `/posts/{slug}` with the postname being the slug.

## Reading Markdown Content

At the start of the application we need to read all the markdown files into memory and serve them up to our router. Before we start down that path let's create a few structs, one to store server related configuration data and the other that will store data for our actual markdown post. At the top of the `main.go` file create the following.

```go
// stores server related configuration (port no, content directory etc..)
type config struct {
 port int
 path string
}

// stores post data
type post struct {
 title   string
 content []byte
}
```

Let's populate our configuration struct with some useful data like our port and content paths. Ideally we don't want to hardcode this type of data into our main function.

Add these lines to the top of our main function

```go
var cfg config

flag.IntVar(&cfg.port, "port", 5040, "5040")
flag.StringVar(&cfg.path, "content path", "./content", "/content")
```

We will create a variable for a config struct and assign values to them at the start of the main function. The `flag` package provides some useful functions to parse command line options. Both `intVar` and `stringVar` functions parse command line options at runtime into our struct fields. If you do not pass any command line options, the defaults will be used. This actually presents a little bit of a problem with our port variable that we need to fix.

Take a look at this line

```go
// use the default serve mux to start a server on port 5040 and check for errors
 if err := http.ListenAndServe(":5040", nil); err != nil {
  fmt.Printf("server crashed with error %v", err)
 }
```

Ideally we would just add `cfg.port` to the first parameter of `ListenAndServe`. Can you spot the issue? Let's try... Make the following change.

```go
if err := http.ListenAndServe(cfg.port, nil); err != nil {
  fmt.Printf("server crashed with error %v", err)
 }
```

Your editor should give you a squigly line or error out when you try to run the program. The issue is that you are assigning an integer to a function that expects a string. Even converting this to a string though won't fix the issue entirely as it expects a string like this `:port`. Let's convert the port to a string using `fmt.Sprintf` function. Make the following change to the function call.****

```go
http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), nil)
```

The `Sprintf` function will format the variable for us and return a string. We are basically telling the function to return us a string with the colon at the start and replace the `%d` with the integar variable.

Let's parse those markdown files. Create a function with this name and signature

```go
readMarkdown(fSys fs.FS) (map[string][]byte, error)
```

 right under the main function.

`readMarkdown` just takes a single parameter, the `fs.FS` interface. Let's build this bitch...

We will start by creating a "cache" `cache := make(map[string][]byte)` this creates a map that will hold our slugs by name and our markown data. Using `fs.FS` provides for a nice file system abstraction that makes it easier to unit test.

We have a cache, now let's get our files. Add the following lines

```go
files, err := fs.Glob(fSys, "*.md")
 if err != nil {
  return nil, err
}
```

We are using fs.Glob to return a slice of strings of the files that match the extension passed in to the function `*.md`. Now let's iterate through these files and read them, and store them in our cache.

```go
for _, f := range files {
  data, err := io.ReadAll(strings.NewReader(f))
  if err != nil {
   return nil, err
  }
  cache[f] = data
}
```

Here is the order of operations for the code above

* Range over each string entry in the slice (ex helloworld.md)
* Use `io.ReadAll` to read a reader, which we will pass using `strings.NewReader` passing in our file.
* Lastly assign the filename as the key in our cache along with the data

This is what main should look like now.

```go
package main

import (
 "flag"
 "fmt"
 "io"
 "io/fs"
 "net/http"
 "strings"
)

type config struct {
 port int
 path string
}

type post struct {
 title   string
 content []byte
}

func main() {
 var cfg config

 flag.IntVar(&cfg.port, "port", 5040, "5040")
 flag.StringVar(&cfg.path, "content path", "./content", "/content")

 // register a handle function for the route default "/" and write the string "hi" in a byte slice
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("<html><body><h1>hi</h1></body></html>"))
 })

 // use the default serve mux to start a server on port 5040 and check for errors
 if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), nil); err != nil {
  fmt.Printf("server crashed with error %v", err)
 }
}

func readMarkdown(fSys fs.FS) (map[string][]byte, error) {
 cache := make(map[string][]byte)
 files, err := fs.Glob(fSys, "*.md")
 if err != nil {
  return nil, err
 }

 for _, f := range files {
  data, err := io.ReadAll(strings.NewReader(f))
  if err != nil {
   return nil, err
  }

  cache[f] = data
 }

 return cache, nil
}
```

Time to put this in action and see if it works as expected. Go to your `main.go` file and add the following right after our calls to `flag`.

```go
cache, err := readMarkdown(os.DirFS(cfg.path))
 if err != nil {
  return
 }
 fmt.Println("cache", cache)
```

We are passing in our path via `os.DirFS` and printing the results. Assuming the markdown we created is in that directory, you should see something like this when you run the program with `go run .`

```bash
cache map[firstpost.md:[102 105 114 115 116 112 111 115 116 46 109 100]]
```

Jackpot... It found our post and it's contents are in the map. In Part 2 we will work on rendering the content and building a router.
