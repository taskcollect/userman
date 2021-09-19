# userman
Manages everything related to user state.

## HTTP Spec

* POST /v1/register
    * JSON Payload:
        * `user` (string) The person's username.
        * `secret` (string) The person's password. Used to verify identity to the server.
        * `creds` (json) Various tokens for other platforms. Should never include plaintext passwords.
    * Returns:
        * 200 OK - User successfully registered
        * 409 Conflict - User already registered
        * 400 Bad Request - JSON Payload was invalid

* GET /v1/get
    * JSON Payload:
        * `user` (string) the person's username
        * `secret` (string) the person's password OR authentication token
        * `fields` (* see below) what fields are required to be sent back?
    * Returns:
        * 200 OK - Data can be sent back
        * 403 Forbidden - Credential verification failed; no access to data
        * 404 Not Found - Requested user does not exist
        * 400 Bad Request - JSON Payload was invalid (includes reason)

\* Here is an example of a field selector
```json
{
    // I want this user's credentials for the Google API & Stile
    "creds": ["google", "stile"], 
    // Send back the preferences of this user, with the defaults filled in
    "prefs": true 
}
```