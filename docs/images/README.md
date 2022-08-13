# Recording Script

## Commands

```bash
nvm use ---lts
npm install -g terminalizer
terminalizer record -k output.yml
terminalizer render -o temp.gif output.yml
gifsicle -O3 --colors 128 -i temp.gif -o output.gif
```

Note on environment; Make sure bashrc when terminalizer starts is set by changing `command:` in config yaml
```bash
export PS1="$ "
export PATH="./:$PATH"
```

### Recording

```bash
### Filter

rare filter -n 5 --match "(\d{3}) (\d+)" access.log

rare filter -n 5 --match "(\d{3}) (\d+)" -e "{1} {2}" access.log

head -n 5 access.log | rare filter -m "(\d{3}) (\d+)"

### Output json for further analysis
rare filter -n 5 --match "(?P<code>\d{3}) (?P<bytes>\d+)" -e '{.}' access.log

rare filter -n 5 --match "(?P<code>\d{3}) (?P<bytes>\d+)" -e '{.}' access.log | jq

### Histogram

rare histo --percentage simple.log

rare histo -m '(\d{3}) (\d+)' -e '{1}' -x access.log

### Bars

rare bars -s -m '\[(.+?)\].*" (\d+)' -e '{buckettime {1} year}' -e '{2}' access.log

### Table

rare table -m '\[(.+?)\].*" (\d+)' -e '{buckettime {1} year}' -e '{2}' access.log

### Heatmap

rare heatmap -m '\[(.+?)\].*" (\d+)' -e "{timeattr {time {1}} yearweek}" -e "{2}" access.log

### Analyze bytes sent, only looking at 200's

rare analyze -m '(\d{3}) (\d+)' -e '{2}' -i '{neq {1} 200}' access.log

```