---
title: "Serving Static Files in Go"
date: 2025-08-01
summary: "A quick guide to serving CSS, JavaScript, and images efficiently in your Go web applications."
tags:
  - "#Webdev"
  - "#Development"
---

Serving static files is essential for any web app. Go makes this easy with the http.FileServer function.

## Example

```go
mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))
```

This will serve your CSS, JS, and images from the /static/ path.

...
