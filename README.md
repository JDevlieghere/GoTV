# GoTV

Automaticall find new episodes from your favorite TV shows and download the torrent file. All you need to do is specify the shows you're interested in. For each show the app will check whether there's a new episode available and download the respective torrent file. 

## Usage

Running the application is as easy as executing the command below:

```
$ ./GoTV run
```

## Configuration

The configuration is stored in the `config.json` file. This file contains the path to the download directory, as well as a list of tv shows. An example configuration is included in the release or can be found [here](https://github.com/JDevlieghere/GoTV/blob/master/config.json). It's imporant this file is in the same folder as the binary. To display the configuration, run the command below.

```
$ ./GoTV info
```
