# scanner for H1 Programs

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
Add Raw tool output to scan
crunchbase

# Testing programs

https://hackerone.com/tron_btfs
https://hackerone.com/evernote