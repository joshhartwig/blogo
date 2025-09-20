---
title: "Creating a blog in Go part 1"
date: 2025-09-20
summary: "Let's create a simple blog from scratch using Go standard libarry and Markdown"
tags:
  - "#Go"
  - "#Blog"
  - "#Development"
---

I spent a decent part of this year and some of last trying to learn Go in my spare time (which is pratically non existent anymore). Go has
a fantastic standard libary for dealing with anything server related. Combine that with a built in templating features, gives you a nice platform to build a simple blog.

Through this 3 part series I will walk you through building a basic blog that serves up posts in Markdown format. Throughout this guide I will walk you through the following features

* Setting up a server
* Creating a router (mux)
* Setting up route handlers
* Hosting static files
* Reading Markdown files & parsing Frontmatter
* Using HTML templates
* Creating a search feature and pagination
* Settng up Tailwind CLI
* Implementing an RSS feed
* Dockerizing the app
* Deploying to a Raspberry Pi

I won't claim to be a Go expert, there are probably quite a few things that could be done better. I hope you learn something though.

## Why Use Templates?

Templates allow you to separate your HTML from your Go code, making your application easier to maintain and more secure.

## Basic Usage

Here's how to render a template in Go:

```go
tmpl, err := template.ParseFiles("template.html")
err = tmpl.Execute(w, data)
```

...
