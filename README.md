# go-scan

## Intro

1. Folder scanner - Scans a given folder to look for stable files
2. Stable is defined as files not being written.
3. In linux looks at modified time and waits a set time before considering it stable.
4. On Windows tries to rename the file and if file is locked by a writing process will wait.
5. Returns the filepaths which can be process through outChannel
6. Takes a white list of sub folders under the scan path and only scans those sub folders

## More to do

1. Tests to be written for scanner
2. Allow to scan all sub folders if white list is blank
3. Add state management to use at processing.
