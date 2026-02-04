cmd directory respresents main entry point in to your application. each folder is a separate
appliction (executable program/ cli tool)

internal packages is not accessible from outside of it's parent directory

any other packages inside folder name other than internal are accessible, but for apis typical 
folder name used is pkg.

internal/pkg is used for packages that are used in multiple programs in cmd.
