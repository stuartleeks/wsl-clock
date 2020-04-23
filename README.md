# wsl-clock

**IMPORTANT: This is not an official solution and whilst it worked on my machine this is an unsupported workaround :-)**

There is an issue with WSL where the clock in WSl isn't updated when the host resumes from sleep/hibernate. E.g. [this issue](https://github.com/microsoft/WSL/issues/4245)

This repo has a workaround that creates a scheduled task that is triggered by Windows Events for resuming from sleep/hibernation. When the scheduled task executes it resets the clock in the WSL VM.

## Setup

To set up the scheduled task, clone the repo and run `add-wslclocktask.ps1`.

## Cleanup

To remove the scheduled task, run `remove-wslclocktask.ps1`

## Logs/Troubleshooting

The script invoked by the scheduled task logs output to `~/.wsl-clock.log`

## Known issues

### Fails if default distribution is configured for WSL 1

The `update-clock.ps1` script executes the `hwclock` against the default distro. If that distro is configured to run under WSL 1 then it will result in the following error as `hwclock` isn't supported on WSL 1 (and WSL 1 isn't executing in the lightweight VM for WSL 2 so this wouldn't have any effect anyway)

```log
hwclock: Cannot access the Hardware Clock via any known method.
hwclock: Use the --debug option to see the details of our search for an access method.
```

See <https://github.com/microsoft/wsl/issues/193>.

Possible workarounds:

* change your default distro to a distro that is configured for WSL 2
* change your default distro to run on WSL 2
* modify your copy of the script to specify a distro via `-d` parameter for `wsl` when it runs `hwclock`
