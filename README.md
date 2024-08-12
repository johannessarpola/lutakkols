# lutakkols

A terminal UI to view the upcoming events on www.jelmu.net. Uses Bubbles and Bubbletea for UI layer. 
Can be used also in offline mode after syncing events from the site with `sync` command.

You can view the help with `lutakkols -help` but should run just calling `lutakkols` after you've installed it and added into $PATH. 

To install you can use `go install github.com/johannessarpola/lutakkols@latest` and it should work just fine. 

It can also be run with docker by building the `Dockerfile` with `docker build . -t <image>:<tag>` 
and then running it with ` docker run -it --rm <image>:<tag>`. Notice that the `-it` is required to capture TTY.

## Screenshots

List shows the upcoming events in which you cna use for example arrows on keyboard to navigate.
Help is printed in the lower section for the keybinds. Enter opens the event view (following screenshot)

![alt text](https://github.com/johannessarpola/lutakkols/blob/main/docs/imgs/lutakkols_1.png?raw=true)

Event view has the event details and event image converted into ascii art. G opens the default browser for the same
page on the Jelmu website. 

![alt text](https://github.com/johannessarpola/lutakkols/blob/main/docs/imgs/lutakkols_2.png?raw=true)
