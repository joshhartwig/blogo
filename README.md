# Blogo

Blogo is a lightweight, fast, and customizable blog engine written in Go, designed to run with minimal dependencies. It reads content directly from Markdown files, supports flexible theming, and leverages Go’s standard library for reliable performance. Themes are easy to create and manage, combining templates, CSS, and JavaScript for a unique look and feel. Blogo is ideal for anyone who wants a simple, efficient blogging platform that’s easy to extend and maintain.

![Screenshot of Blogo](blog.png)

## Run

```bash
blogo serve --config ./config.toml --addr :8080
```

## Installation

```bash
# Clone the repository
git clone https://github.com/joshhartwig/blogo.git
cd blogo

# Build and Start the container
docker compose up --build
```

## Features

- No database needed, reads markdown files direct from content directory
- Site configuration via external `.toml`
- Theme support via external `/theme` folder but also contains embedded default theme
- Leverages built in standard libraries for almost everything short of Markdown & Frontmatter parsing.
- Lightweight

## Contributing

Go away!

## AI usage

The original application was built almost entirely without AI, the most recent refactor leveraged claude code to help with the refactor and to add quite a bit of unit tests. For areas where I need to polish my skillset, I try my best to not use AI, or at very least to ask what the AI would do and implement myself. For things that I have done frequently, I am fine with AI doing some of it.

## License

This project is licensed under the [MIT License](LICENSE).
