# neatfs
a file server create by golang, used to manage file on web server.

It is nano, lite, beta, and even with bugs. 

What's feature here?
 - list files
 - download file
 - upload files
 - copy a file to another path on server
 - move a file to another path on server
 - delete a file
 - rename a file
 - create a dir

How to use it?
 - download a released zip file to you server
 - unzip the zip file
 - edit the config.json
 - change the value of initial from false to true, and change the port and root dir if need
 - finally, run neatfs with "./neatfs" or "./neatfs.exe" on windows
 - maybe "nohup ./neatfs > log &" is better