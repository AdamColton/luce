## Compile Templates
Takes one arg, the directory to read from. Ouputs a file in the directory from
which the command is invoked.

There should be one file named config.json with five values:
```json
{
    "TemplateName": "name to give template in call to template.New",
	"Package": "package name",
	"Var": "variable name",
	"FileName": "file name, if not set, will write to stdout",
	"Path": "base path of templates relative to running directory",
	"Comment": "adds a comment to the output",
	"SkipImport": "if true, will not output the import statement"
}
```