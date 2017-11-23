## Octogon

### An automatic file upload utility for Octoprint/Octopi

Octogon is a simple command line utility which will monitor a single folder on your local file system for the addition/modification of ```.stl``` and ```.gcode``` extension files and automatically copy them to your Octopi so they appear in your files list.

Octogon uses SSH secure copy to move the files rather than Octoprint's REST Api and is very fast. Files are copied to the ```watched``` folder so they will automatically appear in your file list with no browser refresh required.

A folder can specified which must exist on the Octopi. Uploads which target a folder are sent the ```uploads``` folder directly, a browser refresh may be required if you are already viewing the contents of that folder.

### Usage

By default Octogon connects with the user "pi" on octoprint.local, port 22, in most cases only a password -p is required i.e

```
octogon -p mypassword
```

Operation can be tailored by passing additional flags, you can see these with ```octogon -h```

```
Usage of octogon:

  -d	bool	Delete the file after sending
  -f	string	Absolute path to local folder to monitor. Default is current folder
  -hp	string	Remote Hostname and Port to connect on (default "octopi.local:22")
  -ip	string	Remote IP address and Port to connect on
  -p	string	Password, required
  -r	string	Remote folder to send files to. Default is none. Must exist.
  -u	string	User account to connect with (default "pi")

```


### Binaries

Built binaries for Linux, Mac & Windows can be found [here]. (https://github.com/olliephillips/octogon/releases/)

### License

Ocotogon is made available under MIT license.
