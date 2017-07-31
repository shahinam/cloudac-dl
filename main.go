package main

import (
	"fmt"
	"net/url"
	"os"
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
}

func main() {
	dir, _ := os.Getwd()

	log.SetHandler(lcli.Default)
	log.SetLevel(log.DebugLevel)

	app := cli.NewApp()
	app.Name = "cloudac-dl"
	app.Version = version
	app.Usage = `Downloads the video lectures for the given Cloud Academy course.
	 Homepage: https://github.com/shahinam/cloudac-dl`
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Muhammad Inam",
			Email: "mohdinamshah@gmail.com",
		},
	}
	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
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
	}

	app.Commands = []cli.Command{
		{
			Name:    "course",
			Aliases: []string{"i"},
			Usage:   "Download a course.",
			Action:  downloadCourse,
		},
		{
			Name:    "path",
			Aliases: []string{"i"},
			Usage:   "Download all courses in learning path.",
			Action:  downloadLearningPath,
		},
	}

	app.Run(os.Args)
}

// Download learning path
func downloadLearningPath(c *cli.Context) error {
	args := parseCommandLineArgs(c)

	co := &client.Course{
		CourseURL:  args.courseURL,
		SaveDir:    args.saveDir,
		Resolution: args.resolution,
	}

	cl := getClient(c, args)
	err := cl.DownloadCourse(co)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// Download course
func downloadCourse(c *cli.Context) error {
	args := parseCommandLineArgs(c)

	co := &client.Course{
		CourseURL:  args.courseURL,
		SaveDir:    args.saveDir,
		Resolution: args.resolution,
	}

	cl := getClient(c, args)
	err := cl.DownloadLearningPath(co)

	if err != nil {
		log.Error(err.Error())
		return err
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
	args := &CommandLineOptions{}
	args.userName = c.GlobalString("user")
	args.passWord = c.GlobalString("pass")
	args.saveDir = c.GlobalString("out")
	args.resolution = c.GlobalString("res")
	args.courseURL = c.Args().First()

	// Command line options.
	if args.userName == "" || c.NArg() == 0 {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	_, err := url.ParseRequestURI(args.courseURL)
	if err != nil {
		log.Fatalf("The provided url %s is invalid.\n", args.courseURL)
	}

	// If password is not provided - get it interactively.
	if args.passWord == "" {
		fmt.Print("Please enter password: ")
		password, _ := terminal.ReadPassword(int(syscall.Stdin))
		args.passWord = string(password)
	}

	return args
}
