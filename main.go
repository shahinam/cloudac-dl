package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shahinam/cloudac-dl/client"
)

var version = "1.0"

// CommandLineOptions Command line options.
type CommandLineOptions struct {
	userName   string
	passWord   string
	saveDir    string
	resolution string
	courseURL  string
}

func main() {
	// Command line options.
	args := parseCommandLineArgs()

	// Get the client & course.
	c := client.New()
	c.SetUserName(args.userName)
	c.SetPassWord(args.passWord)

	course := &client.Course{args.courseURL, args.saveDir, args.resolution}
	// Login.
	err := c.Login()
	if err != nil {
		fmt.Printf("Failed to Login.")
		os.Exit(1)
	}

	// Download course.
	c.DownloadCourse(course)
}

// Parse command line arguments.
func parseCommandLineArgs() *CommandLineOptions {
	args := &CommandLineOptions{}
	dir, _ := os.Getwd()

	flag.StringVar(&args.userName, "user", "", "The login email address for your CloudAcademy account.")
	flag.StringVar(&args.passWord, "password", "", "The password for your CloudAcademy account.")
	flag.StringVar(&args.saveDir, "out", dir, "The directory where the videos are saved.")
	flag.StringVar(&args.resolution, "res", "720p", "The required video resolution. Allowed values are 360, 720, and 1080.")

	flag.Usage = func() {
		fmt.Printf("cloudac-dl version 1.0\n")
		fmt.Printf("  Downloads the video lectures for the given CloudAcademy course.\n")
		fmt.Printf("  https://github.com/shahinam/cloudac-dl\n\n")

		fmt.Printf("Usage\n")
		fmt.Printf("  cloudac-dl [OPTIONS] course-url \n")
		fmt.Printf("  cloudac-dl -user=user -password=password https://cloudacademy.com/amazon-web-services/aws-security-fundamentals-course\n\n")
		fmt.Printf("Options\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() == 0 || args.userName == "" || args.passWord == "" {
		flag.Usage()
		os.Exit(1)
	}

	args.courseURL = flag.Arg(0)

	return args
}
