# 📍 git-prompt-string

git-prompt-string is a shell agnostic git prompt written in Go. git-prompt-string provides
information about the current git branch and is inspired by
[git-prompt.sh](https://github.com/git/git/blob/master/contrib/completion/git-prompt.sh).

[![go](https://img.shields.io/static/v1?style=flat-square&label=&message=v1.22.0&logo=go&labelColor=282828&logoColor=9dbad4&color=9dbad4)](https://go.dev/)
[![semantic-release: angular](https://img.shields.io/static/v1?style=flat-square&label=semantic-release&message=angular&logo=semantic-release&labelColor=282828&logoColor=d8869b&color=8f3f71)](https://github.com/semantic-release/semantic-release)

> [!WARNING]\
> 03/25/2024: git-prompt-string (previously bgps) is actively undergoing a major rewrite 
>
> This is a breaking change that will simplify and improve maintainability of git-prompt-string
>
> If you prefer to keep using legacy bgps, then use the tag [v0.0.1](https://github.com/mikesmithgh/git-prompt-string/tree/v0.0.1)

## 📦 Installation

### homebrew tap
```sh
brew install mikesmithgh/homebrew-git-prompt-string/git-prompt-string
```

### go install

```sh
go install github.com/mikesmithgh/git-prompt-string@latest 
```

## 🛠️ Setup

Add git-prompt-string to your prompt. For example,

### bash
```sh
PROMPT_COMMAND='PS1="\[\n \e[0;33m\w\e[0m$(git-prompt-string)\n \e[0;32m\u@local \e[0;36m\$\e[0m \]"'
```

## 📌 Alternatives
- [git-prompt.sh](https://github.com/git/git/blob/master/contrib/completion/git-prompt.sh) - bash/zsh git prompt support
- [bash-git-prompt](https://github.com/magicmonty/bash-git-prompt) - An informative and fancy bash prompt for Git users
- [zsh-git-prompt](https://github.com/olivierverdier/zsh-git-prompt) - Informative git prompt for zsh
- [starship](https://starship.rs/) - The minimal, blazing-fast, and infinitely customizable prompt for any shell!

