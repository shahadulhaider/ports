Source demo/tapes/common.tape

Output demo/out/diff.gif

Set Height 320

Hide
Type "docker run --rm -it ports-demo"
Enter
Wait+Screen@60s /❯/
Sleep 2s
Show

Sleep 800ms
Type "ports --diff"
Sleep 300ms
Enter
Sleep 2500ms

Type "node -tcp4 0.0.0.0:4321 &"
Sleep 300ms
Enter
Sleep 1500ms

Type "pkill redis-server"
Sleep 300ms
Enter
Sleep 1500ms

Type "ports --diff"
Sleep 300ms
Enter
Sleep 4s
Screenshot demo/out/diff.png
Sleep 1s

Hide
Type "exit"
Enter
Sleep 1s
