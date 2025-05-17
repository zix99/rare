# Recording Script

## Commands

```bash
nvm use --lts
npm install -g terminalizer
terminalizer record -k output.yml
# Do any yaml cleanup/delays
terminalizer render -o temp.gif output.yml
gifsicle -O3 --colors 128 -i temp.gif -o output.gif
```

Note on environment; Make sure bashrc when terminalizer starts is set by changing `command: bash --rcfile ~/terminalizer/bashrc` in config yaml
```bash
export PS1="$ "
export PATH="./:$PATH"
```

### Main Image Script

```bash
head -n 1 access.log | rare filter --match '(\d{3}) (\d+)'

rare histo -m '(\d{3}) (\d+)' -e '{1}' -x access.log

rare bars --stacked -m '\[(.+?)\].*" (\d+)' -e '{buckettime {1} month}' -e '{2}' access.log

rare heatmap -m '\[(.+?)\].*" (\d+)' -e "{timeattr {time {1}} yearweek}" -e "{timeformat {time {1}} hour}" access.log
```

### Recording

```bash
### Filter

rare filter -n 5 --match "(\d{3}) (\d+)" access.log

rare filter -n 5 --match "(\d{3}) (\d+)" -e "{1} {2}" access.log

head -n 5 access.log | rare filter -m "(\d{3}) (\d+)"

### Search

rare search --include "*.go" test

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

### Sparkline

rare spark -m '\[(.+?)\].*" (\d+)' -e "{timeattr {time {1}} yearweek}" -e "{2}" access.log

### Analyze bytes sent, only looking at 200's

rare analyze -m '(\d{3}) (\d+)' -e '{2}' -i '{neq {1} 200}' access.log

### Reduce http data
rare reduce -m "(\d{3}) (\d+)" -g "http={1}" -a "total={sumi {.} {2}}" -a "count={sumi {.} 1}" -a "avg={divi {total} {count}}" --sort="-{avg}" access.log

```