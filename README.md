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

### Directory

The directory is the path where torrent files should be stored. Many torrent clients have a feature to automatically start downloading torrents that appear in a certain folder. Use this to automate your workflow. The example config uses the current directory as the destination folder.

```json
"Directory":"./",
```

### Quality

Use the quality parameter to refine the search for torrent. If there is no version matching the given quality, no torrent will be downloaded.

```json
"Quality": "720p WEB",
 ```
 
### Series

This is where you specify the TV shows you follow. Add or remove entries as you desire.

```json
"Series":[
  "Angie Tribeca",
  "Mercy Street",
  "Billions",
  "Better Call Saul",
  "Brooklyn 99",
  "The Walking Dead",
  "Shameless US",
  "American Dad!",
  "New Girl",
  "Family Guy",
  "Castle"
]
```
