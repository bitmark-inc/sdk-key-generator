# Bitmark SDK's key generator
A tool for generating Bitmark SDK's client id and client secret

## Prequisites
- Go 1.12
- Go module

## Installation
```bash
go install
```

## Usage
### Generate new sdk token:
Generate a Client ID and Client Certificate.  
Example:
```bash
sdk-key-generator generate --file key.pem
```
After running this, it will generate `key.pem` as Client Certificate file and a `Client ID`.  
Use the Client Certificate file as input of authentication server to issue new Bitmark SDK token.  
Use the generated Client ID to submit to Bitmark Inc in order to issue and use Bitmark SDK token with Bitmark System.

### Test issuing new sdk token:
Use a generated Client Certificate file to issue a Bitmark SDK token, for testing purpose.  
Example:
```bash
sdk-key-generator issuetoken eLt16cSbeLdZLtpjcqgXfw5jbp5r5Z1KHCtQ5s3XY7DGaPQ3NZ --file key.pem
```
It will return a SDK token. This token will be valid for 1 hour from the time issuing.  
We can use this token to input to Bitmark SDK and use.