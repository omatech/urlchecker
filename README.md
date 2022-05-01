# Package urlchecker

## Description
Generate tokens or check for token validity


urls must comply with this format:
"/LVG/HEM/2010/01/09/file.pdf"
/{edition}/HEM/{year}/{month}/{day}/{file}

Additional timestamp must be passed to the generator function

func GenerateToken(url string, timestamp int64) string 

Pass the url, the timestamp (same as the generation) and the token to validate if it's correct.

func Check(url string, timestamp int64, token string) (bool, error) 

Timestamp is checked too, if the timestamp is too old an error is raised.

The validity time for urls is controled by const, for example:

```
MAX_SECONDS          = 48 * 60 * 60
```

To create Tokens from other systems use:

GetMD5Hash(edition + year + month + day + file + strconv.FormatInt(timestamp, 10) + url_checker_secret)

*IMPORTANT: The environment variable URL_CHECKER_SECRET must be set*

## test with 
```
export URL_CHECKER_SECRET="Test Secret String"
go test
```