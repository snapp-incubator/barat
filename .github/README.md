![Logo](https://user-images.githubusercontent.com/49960770/214579375-0a8e1cee-d9cd-4cc5-9453-42aff2c26357.png)

# Barat
Barat is a linter for localization in Golang projects. Barat check your Toml files for missing keys and duplicate keys.
Also, Barat can check your code for missing localization keys.

## Installation
```bash
go get github.com/snapp-incubator/barat
```

or

```bash
go install github.com/snapp-incubator/barat@latest
```

You can check all flags and helps by using this command:

```bash
$ barat checker --help 

Usage: barat checker <path> ...

Check all toml file and make sure that all of them are valid.

Arguments:
  <path> ...    paths to load toml files.

Flags:
  -h, --help                                          Show context-sensitive help.

  -e, --exclude-key-regex=EXCLUDE-KEY-REGEX,...       exclude keys that match the given regex.
      --map-function-names-to-arg-no=KEY=VALUE;...    it's map of the function's name that returns the message by i18n To number of MessageID in arguments.
      --project-path=STRING                           paths to project for check all files.
      --exclude-folders=EXCLUDE-FOLDERS,...           list of exclude folders for check localization.
```

## How to use

### Check your toml files
For check your toml files, you can use this command:

```bash
barat checker <path to toml files> <path to toml files> ...
```

You can exclude some keys from checking by adding them to the `--exclude-key-regex` flag in your command.
```bash 
barat checker <path to toml files> <path to toml files> ... \
      --exclude-key-regex="ExcludeKey","key*"
```

Barat support simple regex for exclude keys. When you use `*` in part of your regex, it will be replaced with `(.*?)`.

## Check your code
If you want to check your code for missing localization keys, you can use this command:

```bash
barat checker <path to toml files> <path to toml files> ... 
      --exclude-key-regex="ExcludeKey","key*" \ 
      --project-path <path to your project> \
      --exclude-folders <folder name>,<folder name>,... \
      --map-function-names-to-arg-no "GetMessages=1;getMessage=0"
```

We describe the flags in the following table:

| Flag                             | Description                                                   |
|----------------------------------|---------------------------------------------------------------|
| `--project-path`                 | Path to your project. Barat search recursive for `.go` files. |
| `--exclude-folders`              | Folders that you want to exclude from checking.               |
| `--map-function-names-to-arg-no` | Map function names to argument number of MessageID.           |

### More details about `--map-function-names-to-arg-no`

In this flag, you can map function names to argument number of MessageID. For example, if you have a function like this:

```go
// GetMessages is function for getting internationalized messages.
func GetMessages(lang string, messageID string, templateData interface{}) (string, error) {
loc := i18n.NewLocalizer(Bundle, lang)
return loc.Localize(&i18n.LocalizeConfig{MessageID: messageID, TemplateData: templateData})
}
```

You can map this function to argument number 1 by adding this to your command:

```bash
--map-function-names-to-arg-no "GetMessages=1"
```

It's imported to mention that the number of arguments in your function starts from 0.

# Help Us
You can contribute to improving this tool by sending pull requests or issues on GitHub.  
Please send us your feedback. Thanks!
