# Use local helm charts

You can use your locally developed charts.

## From file system

If you use a file path (relative to the DSF, or absolute) for the ```chart``` attribute
helmsman will try to resolve that chart from the local file system. The chart on the
local file system must have a version matching the version specified in the DSF.


