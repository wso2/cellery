# Error struct for encapsulating artifact handling errors.
# + message - Error message.
# + code - The error code. This is not an http status code.
# + regError - Inner registry error.
# + cause -  The error(s) that cause the registry error.
type RegistryError object {
    public string message;
    public string code = "";
    public RegistryError? regError;
    public error? cause;
    new (message, code, regError = (), cause = ()) {}

    public function toString() returns @untainted (string) {
        string toStringValue = string `message: {{message}}. code: {{code}} `;
        match regError {
            () => {}
            RegistryError innerRegError => {
                toStringValue = string `{{toStringValue}}
                    - {{innerRegError.toString()}}`;
            }
        }
        match self.cause {
            () => {}
            error bError => {
                toStringValue = string `{{toStringValue}}
                        - {{errorToString(bError)}}`;
            }
        }
        return toStringValue;
    }
};

public function errorToString(error bError) returns @untainted (string) {
    string toStringValue = bError.message;
    match bError.cause {
        error innerBError => {
            toStringValue = string `{{toStringValue}}
                - {{errorToString(innerBError)}}`;
        }
        () => {}
    }
    return toStringValue;
}
