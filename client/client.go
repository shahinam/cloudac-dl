package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	lcli "github.com/apex/log/handlers/cli"

	"github.com/PuerkitoBio/goquery"
	"github.com/gosimple/slug"
)

const (
	baseURL  = "https://cloudacademy.com"
	loginURL = "https://cloudacademy.com/login/"
)

// Course Options.
type Course struct {
	CourseURL  string
	SaveDir    string
	Resolution string
}

// Client HTTP client.
type Client struct {
	*http.Client
	*goquery.Document
	loginURL string
	userName string
	passWord string
}

// Link represents html links.
type Link struct {
	Title string
	URL   string
}

func init() {
	log.SetHandler(lcli.Default)
	log.SetLevel(log.DebugLevel)
}

// New returns a  Client.
func New() *Client {
	client := &Client{}
	client.Client = &http.Client{}
	client.Document = &goquery.Document{}
	client.loginURL = loginURL
	client.userName = ""
	client.passWord = ""

	return client
}

// DownloadLearningPath downloads all courses in a learning path.
func (c *Client) DownloadLearningPath(co *Course) error {
	log.Info("Downloading Learning path: " + co.CourseURL)

	d, err := c.GetDocument(co.CourseURL)
	if err != nil {
		return err
	}

	links := []Link{}
	d.Find("article[data-type='course'] a").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title, _ := s.Attr("title")

		link := Link{title, href}
		links = append(links, link)
	})

	if len(links) == 0 {
		return errors.New("No courses found")
	}

	// Download courses.
	for _, link := range links {
		co.CourseURL = link.URL
		_ = c.DownloadCourse(co)
	}

	return nil
}

// DownloadCourse retrives all videos of a give course.
func (c *Client) DownloadCourse(co *Course) error {
	log.Info("Downloading course: " + co.CourseURL)

	u, _ := url.Parse(co.CourseURL)
	coursePath := u.Path

	// Configure the target directory.
	dir := filepath.Join(co.SaveDir, coursePath)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	m, err := c.CourseContents(co.CourseURL)
	if err != nil {
		return err
	}

	if len(m) == 0 {
		return errors.New("No videos found")
	}

	i := 1
	for _, link := range m {
		fileName := fmt.Sprintf("%02d-%s", i, link.Title)
		fileName = slug.Make(fileName)
		fileName += ".mp4"

		filePath := filepath.Join(co.SaveDir, coursePath, fileName)
		i++

		log.Infof("Downloading Video: %s", link.Title)

		videoURL, err := c.GetVideoURL(link.URL, co)
		if err != nil {
			log.Error("Unable to grab video file")
		} else {
			err = c.DownloadFile(videoURL, filePath)
			if err != nil {
				log.Error("Unable to download")
			} else {
				log.Info("Done")
			}
		}
	}

	return nil
}

// CourseContents Get video Urls.
func (c *Client) CourseContents(url string) ([]Link, error) {
	d, err := c.GetDocument(url)
	if err != nil {
		return nil, err
	}

	links := []Link{}
	d.Find("#course-contents a").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title, _ := s.Attr("title")

		if href == "javascript:void(0);" {
			link := Link{title, url}
			links = append(links, link)
		} else if strings.HasSuffix(href, ".html") {
			link := Link{title, baseURL + href}
			links = append(links, link)
		}
	})

	return links, nil
}

// GetDocument return goquery Document.
func (c *Client) GetDocument(url string) (*goquery.Document, error) {
	res, e := c.Get(url)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()

	return goquery.NewDocumentFromResponse(res)
}

// GetVideoURL - Get video URL from page.
func (c *Client) GetVideoURL(link string, co *Course) (string, error) {
	d, err := c.GetDocument(link)
	if err != nil {
		return "", err
	}

	url := ""
	d.Find("source[type='video/mp4']").Each(func(_ int, s *goquery.Selection) {
		t, _ := s.Attr("type")
		r, _ := s.Attr("data-res")
		if t == "video/mp4" && r == co.Resolution {
			url, _ = s.Attr("src")
		}
	})

	if url != "" {
		return url, nil
	}

	return "", errors.New("could not find the video file")
}

// DownloadFile Download the video.
func (c *Client) DownloadFile(url string, filePath string) error {
	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

// Login into the application.
func (c *Client) Login() error {
	if c.loginURL == "" {
		return errors.New("login url is not set")
	}

	data := url.Values{
		"email":    {c.userName},
		"password": {c.passWord},
	}
	cookieJar, _ := cookiejar.New(nil)
	c.Jar = cookieJar

	resp, err := c.PostForm(c.loginURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d, err := c.GetDocument(baseURL)
	if err != nil {
		return err
	}

	if d.Find(".user-image").Length() == 0 {
		return errors.New("failed to login, please verify credentials")
	}

	return nil
}

// Setters

// SetUserName property.
func (c *Client) SetUserName(name string) {
	c.userName = name
}

// SetPassWord property.
func (c *Client) SetPassWord(pass string) {
	c.passWord = pass
}
