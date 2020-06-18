# scanner for H1 Programs

## Dependencies.

Holy shit, this is long.

1. Golang
 -> Gin Framework
2. MongoDB ( `docker run -p 27017:27017 mongo:latest` will do its work for testing. )
3. Python3
 -> Check the requeriments.txt under tools/[TOOLNAME]/

## Define your own modules:

1. Add the scan logic in `scan/` you can user some of them as example
2. Add the template of the front end in `assets/html`. You can copy some of them and reeplace the first 2 lines
```
{{define "Dirsearch"}}
{{$name := "Dirsearch"}}
```

3. If theres any info you want to persist in de DB, create de model of it in `model` and declareit in `main.go`. This will try to auto migrate de DB and create the schema for it.
```
for _, models := range []interface{}{
		model.Program{},
		model.Target{},
	} 
```


# TODO

- Add Raw tool output to scan [DONE]
- crunchbase
- Add Tools
- Make a fucking docker image to handle all the dependencies.
- Auto import the Program Scope
- Fix Bugs
- Fix MORE BUGS

# Testing programs

https://hackerone.com/tron_btfs
https://hackerone.com/evernote