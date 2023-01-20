package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/dghubble/oauth1"
)

var consumerKey string
var consumerSecret string
var accessToken string
var accessSecret string
var httpClient *http.Client

func init() {
	consumerKey = os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken = os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret = os.Getenv("TWITTER_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		log.Fatalln("TWITTER auth env vars not set")
	}

	// Authenticate with OAuth1 and HMAC-SHA1
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient = config.Client(oauth1.NoContext, token)
}

// Sends a tweet with a message string and optional media ID
func SendTweet(message string, mediaID *string) error {

	values := url.Values{}
	values.Set("status", message)
	if mediaID != nil {
		values.Set("media_ids", *mediaID)
	}
	resp, err := httpClient.PostForm("https://api.twitter.com/1.1/statuses/update.json", values)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Println(string(body))
		return fmt.Errorf("Error: %s", resp.Status)
	}

	return nil
}

// Uploads a media file to Twitter and returns the media ID
func UploadMediaFile(filename string) (*string, error) {
	// create body form
	b := &bytes.Buffer{}
	form := multipart.NewWriter(b)

	// create media paramater
	fw, err := form.CreateFormFile("media", filename)
	if err != nil {
		return nil, err
	}

	// open file
	opened, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// copy to form
	_, err = io.Copy(fw, opened)
	if err != nil {
		return nil, err
	}

	// close form
	form.Close()

	// upload media
	resp, err := httpClient.Post("https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image",
		form.FormDataContentType(), bytes.NewReader(b.Bytes()))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	log.Println(resp.Status)

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Println(string(body))
		return nil, fmt.Errorf("Error: %s", resp.Status)
	}

	// Decode response, only need media ID
	var mediaResp struct {
		MediaIDString string `json:"media_id_string"`
	}
	err = json.NewDecoder(resp.Body).Decode(&mediaResp)
	if err != nil {
		return nil, err
	}

	return &mediaResp.MediaIDString, nil
}

func UploadMediaImage(img *image.RGBA) (*string, error) {
	// create body form
	b := &bytes.Buffer{}
	form := multipart.NewWriter(b)

	// create media paramater
	fw, err := form.CreateFormFile("media", "fractal.png")
	if err != nil {
		return nil, err
	}

	// Encode image to PNG and into form
	err = png.Encode(fw, img)
	if err != nil {
		return nil, err
	}

	// close form
	form.Close()

	// upload media
	resp, err := httpClient.Post("https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image",
		form.FormDataContentType(), bytes.NewReader(b.Bytes()))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	log.Println(resp.Status)

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Println(string(body))
		return nil, fmt.Errorf("Error: %s", resp.Status)
	}

	// Decode response, only need media ID
	var mediaResp struct {
		MediaIDString string `json:"media_id_string"`
	}
	err = json.NewDecoder(resp.Body).Decode(&mediaResp)
	if err != nil {
		return nil, err
	}

	return &mediaResp.MediaIDString, nil
}
