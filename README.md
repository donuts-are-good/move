![donuts-are-good's followers](https://img.shields.io/github/followers/donuts-are-good?&color=555&style=for-the-badge&label=followers) ![donuts-are-good's stars](https://img.shields.io/github/stars/donuts-are-good?affiliations=OWNER%2CCOLLABORATOR&color=555&style=for-the-badge) ![donuts-are-good's visitors](https://komarev.com/ghpvc/?username=donuts-are-good&color=555555&style=for-the-badge&label=visitors)

# move

move is like `mv` but with a progress indicator and visual output


## usage
here's how you use `move`:

```
move /path/to/source /path/to/destination
```
`source` is the file or directory you want to move, and `destination` is where you want to move it to.

if the `source` is a directory, all its contents including subdirectories will be moved to the `destination` directory.

## examples

### moving a single file
to move a single file, use the path to the file as the `source` and the path to the `destination` directory (or the full path including new file name) as the `destination`.

```
move /home/user/documents/file.txt /home/user/desktop
```
this will move `file.txt` from the `documents` directory to the `desktop` directory.

### moving a directory
to move an entire directory, use the path to the directory as the source and the path to the parent of the destination directory as the destination.

```
move /home/user/documents /home/user/desktop
```
this will move the documents directory and all of its contents to the desktop directory.

## license

MIT License 2023 donuts-are-good, for more info see license.md
