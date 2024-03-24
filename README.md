# 📍 git-prompt-string

> [!WARNING]\
> 03/20/2024: git-prompt-string (previously bgps) is actively undergoing a major rewrite in ![go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)
>
> This is a breaking change that will simplify and improve maintainability of git-prompt-string
>
> If you prefer to keep using legacy bgps, then use the tag [v0.0.1](https://github.com/mikesmithgh/bgps/releases/tag/v0.0.1)

A better bash prompt for git. bgps provides a convenient way to customize
the PS1 prompt and to determine information about the current git branch. 
bgps can indicate if the branch has a clean or dirty working tree, whether
or not it is tracking a remote branch, and the number of commits the local
branch is ahead or behind the remote branch.

![demo](screenshots/demo.gif)
Note: Demo is currently outdated and will be updated.

## Installation
### 1. Download bgps
You can choose to install bgps only for yourself or share it with other users 
on the system.
#### 1a. Only for yourself, add bgps to your local bin directory.
1. Check if local bin directory exists. If it does not exist, then create the 
directory. 
```bash
! [[ -d "~/bin" ]] && mkdir ~/bin
```
2. Download bgps to your local bin directory.
```bash
wget -O ~/bin/bgps https://raw.githubusercontent.com/mjsmith1028/bgps/master/bgps 
```
#### 1b. All users, add bgps to the shared bin directory.
1. Download bgps to the shared bin directory.
```bash
sudo wget -O /usr/local/bin/bgps https://raw.githubusercontent.com/mjsmith1028/bgps/master/bgps
```
### 2. Configure bgps
#### 2a. Only for yourself, add user configuration to your home directory.
Create `~/.bgps_config`, if it does not exist. Modify it according to your 
preferences. See the [configuration](#configuration) section for more
information.

If you like what you see in the demo, copy [my configuration](examples/mine):
```bash
wget -O ~/.bgps_config https://raw.githubusercontent.com/mjsmith1028/bgps/master/examples/mine
```
#### 2b. All users, add global configuration to the etc directory.
Create `/etc/bgps_config`, if it does not exist. Modify it according to your 
preferences. See the [configuration](#configuration) section for more 
information.

Note: configurations in the user configuration file, `~/.bgps_config`, will take
precedence over the configurations in the global configuration file,
`/etc/bgps_config`.  

### 3. Update `${PROMPT_COMMAND}`
Modify `PROMPT_COMMAND` environment variable in startup file.  I prefer to 
source my `~/.bashrc` file from my `~/.bash_profile` file. Then, modify the
`PROMPT_COMMAND` in the `~/.bashrc` file.

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
### 4. Reopen terminal
Source the file containing the updates to the `PROMPT_COMMAND` environment
variable. If you followed my suggestion, then you will be sourcing `~/.bashrc`.
Alternativly, just open a new terminal window. You should now see the bgps 
prompt.
```bash
source ~/.bashrc
```
### 5. (Optional) Add git icon
1. Install the powerline symbols font. This package provides the git icon 
referenced in [my configuration](examples/mine).
```bash
sudo apt update && sudo apt install fonts-powerline
```
2. Open a new terminal window. If you are using [my configuration](examples/mine),
then you should now see the git icon.

## Configuration
- `PS1_FORMAT`
  - Format of PS1. `%s` is replaced with the contents of the git section. `%(` and `%)` provide additional formatting when `%s` is within the parentheses. Formatting within the parentheses is only applied to the PS1 when the current working directory is in a git repository.
  - **Default Value:** `%(%s %)\[\033[0m\]\$ `
- `GIT_AHEAD`
  - Format of section that indicates the local branch is ahead of the remote branch. `%s` is replaced with the number of commits that the local branch has not yet merged into the remote branch.
  - **Default Value:** `[ahead %s]`
- `GIT_BEHIND`
  - Format of section that indicates the local branch is behind the remote branch. `%s` is replaced with the number of commits that the remote branch has not yet merged into the local branch.
  - **Default Value:** `[behind %s]`
- `GIT_DIVERGED`
  - Format of section that indicates the local branch has diverged from the remote branch. `%a` is replaced with the number of commits that the local branch has not yet merged into the remote branch. `%b` is replaced with the number of commits that the remote branch has not yet merged into the local branch.
  - **Default Value:** `[ahead %a] [behind %b]`
- `GIT_POSTFIX`
  - Optional section that is displayed before git section.                          
- `GIT_PREFIX`
  - Optional section that is displayed after git section.                           
- `GIT_COLOR_ENABLED`
  - If `true`, then enable color for git section.                                   
  - **Default Value:** `false`
- `GIT_COLOR_CLEAN`
  - Color of git section when the branch has a clean working tree.                  
  - **Default Value:** `\[\033[2;92m\]` (green)
- `GIT_COLOR_CONFLICT`
  - Color of git section when the local branch is ahead or behind the remote branch.
  - **Default Value:** `\[\033[0;33m\]` (yellow)
- `GIT_COLOR_DIRTY`
  - Color of git section when the branch has a dirty working tree.                  
  - **Default Value:** `\[\033[2;91m\]` (red)
- `GIT_COLOR_UNTRACKED`
  - Color of git section when the branch has untracked files.                       
  - **Default Value:** `\[\033[2;95m\]` (purple)
- `GIT_COLOR_NO_UPSTREAM`
  - Color of git section when the local branch is not tracking a remote branch.     
  - **Default Value:** `\[\033[2;40m\]` (gray)
- `TEXT_COLOR`
  - Color of text after PS1.                                                        
  - **Default Value:** `\[\033[0m\]` (reset)

## Environment Variables
- `BGPS_GLOBAL_CONFIG`
  - Global configuration file that is shared between all users.
  - **Default Value:** `/etc/bgps_config`
- `BGPS_USER_CONFIG`
  - User-specific configuration file.                          
  - **Default Value:** `~/.bgps_config`

## Command-line Options
- `--config-file`
  - Path of the desired configuration file. Specifying the configuration file will override and ignore all configurations referenced in `BGPS_GLOBAL_CONFIG` and `BGPS_USER_CONFIG`.
  - **Example:** `source bgps --config-file "/path/to/config/file"`
- `--ls-config`
  - List the current configuration and environment variables.
  - **Example:** `source bgps --ls-config`

## FAQ
### My prompt is a shrug face `¯\_(ツ)_/¯`. What is going on?
You are missing the bgps configuration file `~/.bgps_config`. Make sure 
`~/.bgps_config` exists and that it is properly configured. See the 
[configuration](#configuration) section for more information.

### Why do I have a weird symbol `` for my git icon?
Powerline fonts is used to display the git icon. If you are using 
[my configuration](examples/mine) and the icon is incorrect, then you need to
install the `fonts-powerline` package. See the
[add git icon](#5-optional-add-git-icon) section for instructions.

## Compatibilty
bgps has been tested with the following:
- Operating System: 20.6.0 Darwin Kernel Version 20.6.0
- Bash: 5.1.16(1)-release (x86_64-apple-darwin20.6.0)
- Git: 2.37.0

