# Cloud Academy Downloader

This utility can be used to download [Cloud Academy](https://cloudacademy.com) lecture videos for offline viewing.

## Installation
Download the binary from the latest release.

## Usage
To download a course
```
cloudac-dl --user=<login email> --pass=<password> course <couse url>
```
If you omit `--pass` the app will ask for password.

To download all the courses in a learning path.
```
cloudac-dl --user=<login email> --pass=<password> path <learning path url>
```

You can specify output directory and video resolution by `--out` and `--res` options respectively.

## Commands and options
```
USAGE:
   cloudac-dl [global options] command [command options] [arguments...]
COMMANDS:
     course               Download a course.
     path, learning-path  Download all courses in learning path.
     help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --user value, -u value  The login email address for your Cloud Academy account.
   --pass value, -p value  The password for your Cloud Academy account.
   --out value, -o value   The directory where the videos are saved.
   --res value, -r value   The required video resolution. Allowed values are 360, 720, and 1080. (default: "720p")
   --file value, -f value  Download URLs found in local or external FILE
   --help, -h              show help
   --version, -v           print the version
```


[![Release](https://img.shields.io/github/release/shahinam/cloudac-dl.svg?style=flat-square)](https://github.com/shahinam/cloudac-dl/releases/latest)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://raw.githubusercontent.com/shahinam/cloudac-dl/master/LICENSE)
