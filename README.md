# pictree

pictree is a tool that organizes your pictures or videos automatically based on the metadata of the files.
It takes Live Photos into account since version 0.2.

## usage

`$ pictree -verbose -src "/Source/Folder" -dst "/Destination/Folder"`

```
$ pictree -h
Usage of pictree:
  -dst string
        Destination folder.
  -f string
        Name of the last level folder where the files will be stored.
  -r    Rename the file with the extracted date <Year>-<Month>-<Day>_<Hour>-<Minute>-<Second>_<Originale name>.jpg
  -src string
        Source folder that contains the files to process.
  -verbose
        Print the detailed log.
  -version
        Print the version number.
```