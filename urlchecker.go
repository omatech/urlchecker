package go_urlchecker

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var url_checker_secret string
var timezone = "UTC"
var max_seconds int64 = 48 * 60 * 60
var token_mismatch_error = "Token mismatch"
var url_too_old_error = "Url too old"
var source_url_base = "https://www.google.com"

func init() {
	found := false
	_, found = os.LookupEnv("URL_CHECKER_SECRET")
	if found {
		url_checker_secret = os.Getenv("URL_CHECKER_SECRET")
	} else {
		fmt.Println("Critical: URL_CHECKER_SECRET Environment variable not set!")
		os.Exit(1)
	}

	_, found = os.LookupEnv("TIMEZONE")
	if found {
		timezone = os.Getenv("TIMEZONE")
	}

	_, found = os.LookupEnv("MAX_SECONDS")
	if found {
		err := errors.New("Dummy Error")
		max_seconds, err = strconv.ParseInt(os.Getenv("MAX_SECONDS"), 10, 64)
		if err != nil {
			fmt.Println("Critical: MAX_SECONDS Environment variable cannot be converted to int64!")
			os.Exit(1)
		}
	}

	_, found = os.LookupEnv("TOKEN_MISMATCH_ERROR")
	if found {
		token_mismatch_error = os.Getenv("TOKEN_MISMATCH_ERROR")
	}

	_, found = os.LookupEnv("URL_TOO_OLD_ERROR")
	if found {
		url_too_old_error = os.Getenv("URL_TOO_OLD_ERROR")
	}

	_, found = os.LookupEnv("SOURCE_URL_BASE")
	if found {
		source_url_base = os.Getenv("SOURCE_URL_BASE")
	}
}

func GetTimestamp() int64 {
	date := time.Now()
	year, month, day := date.Date()
	loc, _ := time.LoadLocation(timezone)
	return time.Date(year, month, day, 0, 0, 0, 0, loc).Unix()
}

func Check(url string, timestamp int64, token string) error {
	if GenerateToken(url, timestamp) != token {
		return errors.New(token_mismatch_error)
	}

	limitTime := time.Now().Unix() - max_seconds
	if timestamp < limitTime {
		return errors.New(url_too_old_error)
	}

	return nil
}

func GenerateToken(url string, timestamp int64) string {

	edition, year, month, day, file := getTokenComponentsFromUrl(url)

	return GetMD5Hash(edition + year + month + day + file + strconv.FormatInt(timestamp, 10) + url_checker_secret)
}

func getTokenComponentsFromUrl(url string) (string, string, string, string, string) {
	s := strings.Split(url, "/")
	edition := s[1]
	year := s[3]
	month := s[4]
	day := s[5]
	file := s[6]

	return edition, year, month, day, file
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetRemoteFileInBase64(path string) string {
	res, err := http.Get(source_url_base + path)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	encodedContent := base64.StdEncoding.EncodeToString(content)
	res.Body.Close()
	return encodedContent
}

func Debug(url string, timestamp int64, token string) string {
	edition, year, month, day, file := getTokenComponentsFromUrl(url)

	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("URL=%v \n", url)
	message += fmt.Sprintf("token=%v \n", token)
	message += fmt.Sprintf("md5=%v \n", GenerateToken(url, timestamp))
	message += fmt.Sprintf("timestamp=%v \n", timestamp)
	message += fmt.Sprintf("edition=%v \n", edition)
	message += fmt.Sprintf("year=%v \n", year)
	message += fmt.Sprintf("month=%v \n", month)
	message += fmt.Sprintf("day=%v \n", day)
	message += fmt.Sprintf("file=%v \n", file)

	return message
}
