# The hashing algorithms supported by this module.
public type Algorithm "SHA1"|"SHA256"|"MD5";

# The `SHA1` hashing algorithm
@final public Algorithm SHA1 = "SHA1";
# The `SHA256` hashing algorithm
@final public Algorithm SHA256 = "SHA256";
# The `MD5` hashing algorithm
@final public Algorithm MD5 = "MD5";

# Returns the hash of the given file using the specified algorithm.
#
# + path - The path of the file to be hashed
# + algorithm - The hashing algorithm to be used
# + return - The hashed string
public extern function hash(string path, Algorithm algorithm) returns (string|error);