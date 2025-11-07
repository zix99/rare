# Run this from the project root dir

Output docs/images/rare.gif

Require echo
Require rare
Require head

Set Shell "bash"
Set FontSize 12
Set Width 900
Set Height 600
Set WindowBar Colorful
Set Padding 16

Env PS1 "> "

Type "head -n 1 access.log | rare filter --match '(\d{3}) (\d+)'" Sleep 500ms Enter
Wait
Sleep 2s
Enter

Type "rare histo -m '(\d{3}) (\d+)' -e '{1}' -x access.log" Sleep 500ms Enter
Wait
Sleep 2s
Enter

Type `rare bars --stacked -m '\[(.+?)\].*" (\d+)' -e '{buckettime {1} month}' -e '{2}' access.log` Sleep 500ms Enter
Wait
Sleep 2s
Enter

Type `rare heatmap -m '\[(.+?)\].*" (\d+)' -e '{timeattr {time {1}} yearweek}' -e '{timeformat {time {1}} hour}' access.log` Sleep 1s Enter
Wait
Sleep 2s

Sleep 3s
