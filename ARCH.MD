# Applications
## app.moon
- sets layout to views/layout.moon
- renders views/index.moon at "/"
- renders views/login.moon
- handles login POST requests
- imports applications/protected.moon
## applications/protected.moon
- contains the auth middleware
- renders views/dashboard.moon
- imports applications/blog.moon
## applications/blog.moon
- handles /posts/add and /posts/mod
# Views
## views/layout.moon
- is the main layout every route uses
- has a two item navbar
  - Homepage: "/"
  - Dashboard: "/dashboard"
## views/index.moon
- is rendered directly inside views/layout.moon
- lists posts in a grid pattern
## views/login.moon
- is rendered directly inside views/layout.moon
- provides the login form
## views/dashboard.moon
- is rendered directly inside views/layout.moon
- renders a second navbar, that points to child routes like
  - views/postform.moon with no arguments
## views/postform.moon
- is rendered wherever necessary
- accepts one optional argument
  - a post to modify
- 