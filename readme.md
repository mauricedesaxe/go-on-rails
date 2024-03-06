# Go on Rails / Rails for Go

A starter kit for Go, not a framework, resembling patterns found in Rails and other MVC frameworks.

What I mean to do here is speed up web app development in Go, by providing a set of tools and patterns that are familiar to developers coming from other languages and frameworks. 

I am also highly opinionated against thick clients.

If you like [Go](https://golang.org/), [Docker](https://www.docker.com/), [HTMX](https://htmx.org/), [Templ](https://templ.guide/), and [Tailwind](https://tailwindcss.com/), then you might enjoy this project and find value in it.

## How to use

1. Clone this repository
2. Run `make dev` to start the development server
3. Open your browser at `http://localhost:3000`

## How to deploy

1. Clone this repository
2. Build the Docker image with `docker build -t yourapp .`
3. Run the Docker container with `docker run -p 8080:8080 yourapp`

Alternatively, you can write a `docker-compose.yml` file and use it to deploy your app.
This is not included here because you may likely use one for multiple services, not just this one.