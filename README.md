# godisk

A utility to display disk usage, similar to `du`, but it works with network file systems.

I implemented it to see where my Google Drive storage was going. It aggregates the total
disk usage per folder (including all of its children), and displays children sorted by
descending usage. It only displays folders with their totals.

In other words, if your directory structure and files look like this:
```
dirA/
  file1 (1 mb)
  file2 (2 mb)
  subdir1/
    file3 (3 mb)
    file4 (4 mb)
  subdir2/
    file5 (5 mb)
```

Then running `godisk dirA` will output something like:
```
`-- 15.0 MB dirA (self: 3.0 MB)
   |-- 7.0 MB subdir1 (self: 7.0 MB)
   |-- 5.0 MB subdir2 (self: 5.0 MB)
```
