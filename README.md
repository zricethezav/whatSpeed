<p align="center">
  <img alt="whatSpeed" src="https://raw.githubusercontent.com/zricethezav/gifs/master/whatspeed.png" />
</p>

# whatSpeed tells you what speed
This is a tiny application written in Go that will report your download and upload speeds by running tests again speedtest.net by Ookla. The test works by downloading random images (literally np.random.random(500,500), like the image above) from the nearest* speedtest server and measuring the time it takes for the downloads to complete. Similar test for uploading.  

```
$ whatSpeed
whatSpeed v1.0, by zricethezav
45.1659259433093
118.98816376759761
131.59810886836965
195.33647016125454
186.12982979073337
143.58390724639452
150.6709725918009
170.66080382593154
144.5409862240571
169.5862135833997
avg: 145.63 mbps (download)
21.533985770260482
21.79260714188201
23.29998753674339
23.88618888507929
23.251959471046728
avg: 22.75 mbps (upload)
```

\* Nearest meaning smallest GREAT CIRCLE DISTANCE between your client machine and a speedtest.net server. Calculated using the haversine forumula: https://en.wikipedia.org/wiki/Haversine_formula

## Usage
```
Usage of ./whatSpeed:
  -download
        run download test only
  -upload
        run upload test only
  -version
        outputs version number
```

## Installing
```
go get -u github.com/zricethezav/whatSpeed
```
Or download from release binaries [here](https://github.com/zricethezav/whatSpeed/releases)
