![Logo](https://user-images.githubusercontent.com/49960770/214579375-0a8e1cee-d9cd-4cc5-9453-42aff2c26357.png)

# Barat
Barat is a linter for [localization in Golang](https://github.com/nicksnyder/go-i18n) projects. Barat check 
your Toml files for missing keys and duplicate keys. Also, Barat can check your code for missing localization keys.

## Installation

```bash
go install github.com/snapp-incubator/barat@latest
```

You can check all flags and helps by using this command:

```bash
$ barat lint --help 

Usage: barat lint

Check all toml file for validation and find message keys in code.

Flags:
  -h, --help                                          Show context-sensitive help.

      --config-path=STRING                            Path to config file.
      --toml-paths=TOML-PATHS,...                     paths to load toml files.
  -e, --exclude-regex-keys=EXCLUDE-REGEX-KEYS,...     exclude keys that match the given regex.
      --map-function-names-to-arg-no=KEY=VALUE;...    it's map of the function's name that returns the message by i18n To number of MessageID in arguments.
      --project-path=STRING                           paths to project for check all files.
      --exclude-folders=EXCLUDE-FOLDERS,...           list of exclude folders for check localization.
```

## How to use Barat

### Check your toml files
For check your toml files, you should specify the path to your toml files by `--toml-paths` flag.

```bash
barat lint --toml-paths <path to toml files>, <path to toml files>, ...
```

You can exclude some keys from checking by adding them to the `--exclude-regex-keys` flag in your command.
```bash 
barat lint --toml-paths <path to toml files>, <path to toml files>, ... \
      --exclude-regex-keys="ExcludeKey","key*"
```

Barat support simple regex for exclude keys. When you use `*` in part of your regex, it will be replaced with `(.*?)`.

## Check your code
If you want to check your code for missing localization keys, you can use this command:

```bash
barat lint --toml-paths <path to toml files>, <path to toml files>, ... \
      --exclude-regex-keys="ExcludeKey","key*" \ 
      --project-path <path to your project> \
      --exclude-folders <folder name>, <folder name>, ... \
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

### Run Barat With Config File

You can run Barat with config file. You can create a file named `barat.yaml` in your project and add your config to it.
After creating your config file, you can run Barat with this command:

```bash
barat lint --config-path <path to your config file>
```

#### Config File

You can use `barat.yml` in [Barat source code](https://github.com/snapp-incubator/barat/blob/master/barat.yml) as your
config file template.

```yaml
project-path: "./project-path"

toml-path:
  - "./project-path/toml/en"
  - "./project-path/toml/fa"
  - "./project-path/toml/ru"

exclude:
  folders:
    - "vendor"
    - "assets"
  regex-keys:
    - "ToBeExcludedKey1"
    - "KeyToBe*"

message-functions:
  - name: "GetMessage"
    message-id-no: 1

options:
  enable:
    toml-check: true
    description-check: true
    other-key-check: true
    code-check: true
```

Most part of the config file was described before. But we describe the `options.enable` part in the following:
* toml-check
    * if you set this to `false`, Barat will not check your toml files.
* description-check
    * if you set this to `false`, Barat will not check your toml files for missing descriptions.
* other-key-check
    * if you set this to `false`, Barat will not check your toml files for missing other keys.
* code-check
    * if you set this to `false`, Barat will not check your code for missing localization keys.

#### Note: If you use config file, you can't use other flags. If you use other flags, Barat will ignore them.


## Help Us
You can contribute to improving this tool by sending pull requests or issues on GitHub.  
Please send us your feedback. Thanks!
