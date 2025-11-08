jwt = require "resty.jwt"
bcrypt = require "bcrypt"
import Users from require "models"

class Auth
  -- Get JWT secret from config or use default (should be overridden in production)
  @get_secret: =>
    secret = os.getenv("JWT_SECRET")
    if not secret
      secret = "default-insecure-secret-change-in-production"
      print "WARNING: Using default JWT secret. Set JWT_SECRET environment variable!"
    secret

  -- Generate a JWT token for a user
  @generate_token: (username) =>
    secret = @get_secret!
    jwt_token = jwt:sign(secret, {
      header: { typ: "JWT", alg: "HS256" }
      payload: {
        username: username
        exp: os.time! + (7 * 24 * 60 * 60) -- 7 days
      }
    })
    jwt_token

  -- Verify a JWT token and return the username if valid
  @verify_token: (token) =>
    return nil unless token
    secret = @get_secret!
    jwt_obj = jwt:verify(secret, token)
    return nil unless jwt_obj and jwt_obj.verified
    jwt_obj.payload.username

  -- Hash a password using bcrypt
  @hash_password: (password) =>
    bcrypt.digest(password, 12)

  -- Verify a password against a hash
  @verify_password: (password, hash) =>
    bcrypt.verify(password, hash)

  -- Authenticate user with username and password
  @authenticate: (username, password) =>
    user = Users\find username: username
    return nil unless user
    return nil unless @verify_password(password, user.passhash)
    user

  -- Create a new user
  @create_user: (username, password) =>
    -- Check if user already exists
    existing = Users\find username: username
    return nil, "User already exists" if existing
    
    -- Hash password and create user
    hash = @hash_password(password)
    user = Users\create {
      username: username
      passhash: hash
      xp: 0
    }
    user

{ :Auth }
