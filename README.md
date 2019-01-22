# Overview

`todocli` is a small utility I wrote to help search my to do list for upcoming items that are due. The format of my to do list is just basic markdown combined with a little bit of go style tagging, e.g.:

```
* [ ] Item 1 `due:"01-01-2019"`
* [ ] Item 2 `tag:"important"`
* [ ] Item 3 `created_at:"01-02-2019"`
* [ ] Item 3 `created_at:"01-02-2019" modified_at:"01-03-2019" due:"02-01-2019"`
```

The utility is built using GoLang + [Cobra](https://github.com/spf13/cobra) and all the command are discoverable by running `todocli --help`.
