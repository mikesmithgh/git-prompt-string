# Better Git Prompt String
A better Bash prompt for Git.

![demo](screenshots/demo.gif)

## Installation
1. Copy bgps file to a bin directory on your machine. Choose one of the following options.
   - Option A) Your local bin directory:
      Check if local bin directory exists. If it does not exist, then create the directory. 
      ```bash
      ! [[ -d "~/bin" ]] && mkdir ~/bin
      ```
      Download bgps to your local bin directory.
      ```bash
      wget -O ~/bin/bgps https://raw.githubusercontent.com/mjsmith1028/bgps/master/bgps 
      ```
   - Option B) Shared bin directory:
      Download bgps to the shared bin directory.
      ```bash
      sudo wget -O /usr/local/bin/bgps https://raw.githubusercontent.com/mjsmith1028/bgps/master/bgps
      ```
2. Copy or create your desired bgps configuration at `~/.bgps_config`.
   For example, copy my configuration:
   ```bash
   wget -O ~/.bgps_config https://raw.githubusercontent.com/mjsmith1028/bgps/master/examples/mine
   ```
3. Modify `PROMPT_COMMAND` environment variable in startup file.
   I prefer to source my `.bashrc` file from my `.bash_profile` file. Then modify the `PROMPT_COMMAND` in the `.bashrc` file.

   For example:

   Add the following to `~/.bash_profile`.
   ```bash
   if [ -f "${HOME}/.bashrc" ] ; then
     source "${HOME}/.bashrc"
   fi
   ```
   Add the following to `~/.bashrc`.
   ```bash
   PROMPT_COMMAND="source bgps"
   ```
4. Source `.bashrc` or open a new terminal window. You should now see the BGPS prompt.
   ```bash
   source ~/.bashrc
   ```
5. Optional. Add git symbol.
   1. Install Powerline symbols font. This package provides the git symbol referenced in [my](examples/mine) example configuration.
     ```bash
     sudo apt update && sudo apt install fonts-powerline
     ```
   2. Open a new terminal window. If you are using [my](examples/mine) configuration, then you should now see the git symbol.

## Configuration
TODO

## FAQ
### My prompt is a shrug face `¯\_(ツ)_/¯`. What is going on?
You are missing the BGPS configuration file `~/.bgps_config`. Make sure `~/bgps-config` exists and that it is properly configured. 

### Why do I have a weird symbol `` for my git icon?
Powerline fonts is used to display the git icon. If you are using [my](examples/mine) configuration and the icon is incorrect, then you need to install the `fonts-powerline` package.

## Compatibilty
BGPS has been tested with the following:
- Operating System: Linux 4.10.0-35-generic #39~16.04.1-Ubuntu
- Bash: 4.3.48(1)-release
- Git: 2.7.4
