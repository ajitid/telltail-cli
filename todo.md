(another soln. from https://github.com/golang/go/issues/29202#issuecomment-1233042513)

to write code for windows, open cli and write

```
setenv "GOOS" "windows"
```

before opening nvim. See https://github.com/golang/go/issues/29202.

If it errors out with exit code 1, its fine. Doing `echo $GOOS` will give you the right result.

The equivalent for fish is `set -x GOOS windows`. ref: https://fishshell.com/docs/current/cmds/set.html

---

check if a local file can be downloaded

make sure that services are disabled and removed in a failsafe manner before re-enabling them by re-installing
.. so say like doing systemctl --user stop/disable telltail-center << and it fails because the service doesn't exist

correct tailnet name in docs and other places. Tailscale considers the whole `jasklfj.ts.net` to be tailnet name. And there are many type of names available for net https://tailscale.com/kb/1136/tailnet/

// put commands into categories

think about using fmt vs log

main issues (probably put them in tldr?):

- dns-fight
- manually going and removing telltail node because telltail-1 appears. Check if with tailscale cli all telltail[-num] nodes can be removed via telltail cli without sudo
