rups2html
=========

A silly go program that runs "rups" and uses [flot](https://github.com/flot) to create an html page showing some cpu usage history as graphs

=========

"rups" is expected to be defined as a short script that looks like this:

```bash
date
rup SERVER1 SERVER2 SERVER3 ....
```

**date** gives the timestamp to mark the data points and the loads across the supplied servers are parsed out of **rup**'s output

=========

I also find that doing something like this is quite helpful:
```bash
function checkAndRunrups2html() {
    if [ "$(pidof rups2html)" ]; then
        echo "rups2html is already running"
    else
        echo "starting rups2html"
	(cd ~/gopath/src/github.com/skiesel/rups2html; ~/gopath/bin/rups2html -scpremote="SERVER:path/to/index.html" &)
    fi
}

alias rups2html="checkAndRunrups2html"
```
=========
