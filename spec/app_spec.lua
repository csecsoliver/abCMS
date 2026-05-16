local request = require("lapis.spec.request")
require("app")

local Models = require("models")
local Users = Models.Users
local Posts = Models.Posts
local bcrypt = require("bcrypt")
local spec_helper = require("spec.spec_helper")

local function header_value(headers, name)
  return headers[name] or headers[name:lower()] or headers[name:upper()]
end

local function create_user(username, password)
  return Users:create({
    username = username,
    passhash = bcrypt.digest(password, 12)
  })
end

local function auth_session(user)
  return {
    user = user.username,
    expiry = os.time() + 3600
  }
end

describe("abCMS routes", function()
  before_each(function()
    spec_helper.truncate_tables()
  end)

  it("renders the login page", function()
    local status, body = request("/login")
    assert.are.equal(200, status)
    assert.matches("Login", body)
  end)

  it("redirects authenticated users away from login", function()
    local user = create_user("alice", "secret")
    local status, _, headers = request("/login", {
      session = auth_session(user)
    })

    assert.are.equal(302, status)
    assert.are.equal("/dashboard", header_value(headers, "Location"))
  end)

  it("signs up new users", function()
    local status, _, headers = request("/login", {
      method = "POST",
      params = {
        option = "signup",
        username = "alice",
        password = "secret"
      }
    })

    assert.are.equal(302, status)
    assert.are.equal("/dashboard", header_value(headers, "Location"))
    assert.is_not_nil(Users:find({ username = "alice" }))
  end)

  it("rejects duplicate signups", function()
    create_user("alice", "secret")

    local status, body = request("/login", {
      method = "POST",
      params = {
        option = "signup",
        username = "alice",
        password = "secret"
      }
    })

    assert.are.equal(200, status)
    assert.matches("Username already exists", body)
  end)

  it("reports invalid login credentials", function()
    create_user("alice", "secret")

    local status, body = request("/login", {
      method = "POST",
      params = {
        option = "login",
        username = "alice",
        password = "wrong"
      }
    })

    assert.are.equal(200, status)
    assert.matches("Invalid username or password", body)
  end)

  it("logs in valid users", function()
    create_user("alice", "secret")

    local status, _, headers = request("/login", {
      method = "POST",
      params = {
        option = "login",
        username = "alice",
        password = "secret"
      }
    })

    assert.are.equal(302, status)
    assert.are.equal("/dashboard", header_value(headers, "Location"))
  end)

  it("redirects unauthenticated dashboard access", function()
    local status, _, headers = request("/dashboard")
    assert.are.equal(302, status)
    assert.are.equal("/login", header_value(headers, "Location"))
  end)

  it("redirects expired sessions to login", function()
    local user = create_user("alice", "secret")

    local status, _, headers = request("/dashboard", {
      session = {
        user = user.username,
        expiry = os.time() - 10
      }
    })

    assert.are.equal(302, status)
    assert.are.equal("/login", header_value(headers, "Location"))
  end)

  it("rejects empty posts without an image", function()
    local user = create_user("alice", "secret")

    local status, _, headers = request("/formapi/posts/add", {
      method = "POST",
      session = auth_session(user),
      params = {
        title = "",
        content = "",
        image = {
          filename = ""
        }
      }
    })

    assert.are.equal(302, status)
    assert.are.equal("/dashboard/posts/add?error_message=At least one of title, content, or image must be provided.",
      header_value(headers, "Location"))
  end)

  it("creates posts for authenticated users", function()
    local user = create_user("alice", "secret")

    local status, _, headers = request("/formapi/posts/add", {
      method = "POST",
      session = auth_session(user),
      params = {
        title = "Hello",
        content = "World"
      }
    })

    assert.are.equal(302, status)
    assert.are.equal("/dashboard/posts", header_value(headers, "Location"))
    assert.is_not_nil(Posts:find({ title = "Hello" }))
  end)

  it("handles missing posts on delete", function()
    local user = create_user("alice", "secret")

    local status, _, headers = request("/formapi/posts/delete/999", {
      session = auth_session(user)
    })

    assert.are.equal(302, status)
    assert.are.equal("/dashboard/posts?error_message=Post not found.", header_value(headers, "Location"))
  end)

  it("blocks unauthorized deletes", function()
    local owner = create_user("owner", "secret")
    local intruder = create_user("intruder", "secret")
    local post = Posts:create({
      user_id = owner.id,
      title = "Secret",
      content = "Hidden",
      created_at = os.time(),
      has_image = 0,
      path = "",
      thumbnail_path = ""
    })

    local status, _, headers = request("/formapi/posts/delete/" .. post.id, {
      session = auth_session(intruder)
    })

    assert.are.equal(302, status)
    assert.are.equal("/dashboard/posts?error_message=You are not authorized to delete this post.",
      header_value(headers, "Location"))
    assert.is_not_nil(Posts:find({ id = post.id }))
  end)

  it("deletes owned posts", function()
    local owner = create_user("owner", "secret")
    local post = Posts:create({
      user_id = owner.id,
      title = "Disposable",
      content = "Bye",
      created_at = os.time(),
      has_image = 0,
      path = "",
      thumbnail_path = ""
    })

    local status, _, headers = request("/formapi/posts/delete/" .. post.id, {
      session = auth_session(owner)
    })

    assert.are.equal(302, status)
    assert.are.equal("/dashboard/posts?error_message=Post deleted successfully.",
      header_value(headers, "Location"))
    assert.is_nil(Posts:find({ id = post.id }))
  end)

  it("renders posts and handles missing ids", function()
    local user = create_user("alice", "secret")
    local post = Posts:create({
      user_id = user.id,
      title = "Visible",
      content = "Sample",
      created_at = os.time(),
      has_image = 1,
      path = "/static/images/1.png",
      thumbnail_path = "/static/images/2.png"
    })

    local status, body = request("/posts/" .. post.id)
    assert.are.equal(200, status)
    assert.matches("Visible", body)

    local missing_status, missing_body = request("/posts/999")
    assert.are.equal(404, missing_status)
    assert.matches("Post not found", missing_body)
  end)
end)
