package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/apex/log"
	lcli "github.com/apex/log/handlers/cli"
	"github.com/shahinam/cloudac-dl/client"
	"github.com/urfave/cli"

	"golang.org/x/crypto/ssh/terminal"
)

var version = "1.x-dev"

// CommandLineOptions Command line options.
type CommandLineOptions struct {
	userName   string
	passWord   string
	saveDir    string
	resolution string
	courseURL  string
	inputFile  string
}

// Init app.
func init() {
	log.SetHandler(lcli.Default)
	log.SetLevel(log.DebugLevel)
}

func main() {
	dir, _ := os.Getwd()

	app := cli.NewApp()
	app.Name = "cloudac-dl"
	app.Version = version
	app.Usage = `Downloads the video lectures for the given Cloud Academy course.
	 Homepage: https://github.com/shahinam/cloudac-dl`
	app.Authors = []cli.Author{
		{
			Name:  "Muhammad Inam",
			Email: "mohdinamshah@gmail.com",
		},
	}
	app.Action = func(c *cli.Context) error {
		return cli.ShowAppHelp(c)
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Usage: "The login email address for your Cloud Academy account.",
			Value: "",
		},
		cli.StringFlag{
			Name:  "pass, p",
			Usage: "The password for your Cloud Academy account.",
			Value: "",
		},
		cli.StringFlag{
			Name:  "out, o",
			Usage: "The directory where the videos are saved.",
			Value: dir,
		},
		cli.StringFlag{
			Name:  "res, r",
			Usage: "The required video resolution. Allowed values are 360, 720, and 1080.",
			Value: "720p",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Download URLs found in local or external FILE",
			Value: "",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "course",
			Usage: "Download a course.",
			Action: func(c *cli.Context) error {
				return download(c, "course")
			},
		},
		{
			Name:    "path",
			Aliases: []string{"learning-path"},
			Usage:   "Download all courses in learning path.",
			Action: func(c *cli.Context) error {
				return download(c, "path")
			},
		},
	}

	_ = app.Run(os.Args)
}

// Read the input file into URL array.
func readInputFile(args *CommandLineOptions) ([]string, error) {
	file, err := os.Open(args.inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.Trim(scanner.Text(), " ")
		if !strings.HasPrefix(text, "#") {
			lines = append(lines, scanner.Text())
		}
	}

	return lines, scanner.Err()
}

// Download.
func download(c *cli.Context, op string) error {
	args := parseCommandLineArgs(c)
	cl := getClient(c, args)
	links, err := readInputFile(args)
	if err != nil {
		return err
	}

	// Append the url if provided.
	if args.courseURL != "" {
		links = append(links, args.courseURL)
	}

	co := &client.Course{
		CourseURL:  args.courseURL,
		SaveDir:    args.saveDir,
		Resolution: args.resolution,
	}

	for _, link := range links {
		co.CourseURL = link

		// set an error if invalid op is provided.
		err = errors.New("invalid operation")

		if op == "course" {
			err = cl.DownloadCourse(co)
		} else if op == "path" {
			err = cl.DownloadLearningPath(co)
		}

		if err != nil {
			log.Error(err.Error())
		}
	}

	return nil
}

// Get client.
func getClient(c *cli.Context, args *CommandLineOptions) *client.Client {
	// Get the client & course.
	cl := client.New()
	cl.SetUserName(args.userName)
	cl.SetPassWord(args.passWord)

	// Login.
	err := cl.Login()
	if err != nil {
		log.Fatal("Failed to Login.")
	}

	return cl
}

// Parse command line arguments.
func parseCommandLineArgs(c *cli.Context) *CommandLineOptions {
	// Command line options.
	args := &CommandLineOptions{}
	args.userName = c.GlobalString("user")
	args.passWord = c.GlobalString("pass")
	args.saveDir = c.GlobalString("out")
	args.resolution = c.GlobalString("res")
	args.inputFile = c.GlobalString("file")
	args.courseURL = c.Args().First()

	// Validations.
	if args.userName == "" {
		_ = cli.ShowAppHelp(c)
		os.Exit(1)
	}

	if args.inputFile == "" && args.courseURL == "" {
		log.Fatalf("Please provide a URL to download or specify a URL list with --file flag.")
	}

	if args.courseURL != "" {
		_, err := url.ParseRequestURI(args.courseURL)
		if err != nil {
			log.Fatalf("The provided url %s is invalid.\n", args.courseURL)
		}
	}

	// If password is not provided - get it interactively.
	if args.passWord == "" {
		fmt.Print("Please enter password: ")
		password, _ := terminal.ReadPassword(int(syscall.Stdin))
		args.passWord = string(password)
	}

	return args
}
