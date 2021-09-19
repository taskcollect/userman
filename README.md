# userman
Manages everything related to user state.

## HTTP Spec

* POST /v1/register
    * JSON Payload:
        * `name` (string) The person's username.
        * `secret` (argon2 hash) The person's hashed password. Used to verify identity to the server.
        * `creds` (json) Various tokens for other platforms. Should never include plaintext passwords.
    * Returns:
        * 200 OK - User successfully registered
        * 409 Conflict - User already registered
        * 400 Bad Request - JSON Payload was invalid (includes reason)