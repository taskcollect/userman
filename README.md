# userman
Manages everything related to user state in the backend.

Exposes a HTTP api for other taskcollect microservices to connect to.

See the HTTP specification below.

## HTTP Spec

* POST /v1/register
    * JSON Payload:
        * `user` (string) The person's username.
        * `secret` (string) The person's password, in plaintext. Used to generate the DB hash.
        * `creds` (json) Various tokens for other platforms. Should never include plaintext passwords.
    * Returns:
        * 200 OK - User successfully registered
        * 409 Conflict - User already registered
        * 400 Bad Request - JSON Payload was invalid

        <br>

        **Example request:** (comments not present in real request)
        ```jsonc
        {
            "username": "someuser",
            "secret": "PlaintextToBeHashed123",
            "creds": {
                // API tokens, never passwords
                // don't rely on these keys, this is an example
                "google": "TOKEN123"
            },
            "prefs": {
                // you can omit default settings here if you want
                "time24h": true
            }
        }
        ```
* GET /v1/get
    * JSON Payload:
        * `user` (string) the person's username
        * `secret` (string) the person's password
        * `creds` (boolean) send the person's credentials back?
        * `prefs` (boolean) send the person's preferences back? 
    * Returns:
        * 200 OK - Data can be sent back
        * 403 Forbidden - Credential verification failed; no access to data
        * 404 Not Found - Requested user does not exist
        * 400 Bad Request - JSON Payload was invalid

        <br>
        
        **Example request:**
        ```jsonc
        {
            "username": "someuser",
            "secret": "plaintext123",
            "prefs": true,
            "creds": true
        }
        ```
        **Example response:** (comments not present in real response)
        ```jsonc
        {
            "prefs": {
                // don't rely on these keys, this is an example
                "time24h": true,
                "accentColor": "blue"
            },
            "creds": {
                // API tokens, never passwords
                // again, don't rely on these keys
                "google": "TOKEN123"
            }
        }
        ```
