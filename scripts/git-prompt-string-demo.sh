#!/usr/bin/env bash

# Utility script to execute demo steps via Kitty

win= # replace with kitty window ID
delay='0.2'

if [[ -z "$win" ]]; then
	printf "Please povide kitty window ID"
	exit 1
fi

function kitty_send_text() {
	kitty @ send-text "--match=id:$win" "$@"
}

function feed_no_newline() {
	str=""
	for word in "$@"; do
		str="$str $word"
	done
	str="${str:1}" # trim leading space
	for ((i = 0; i < ${#str}; i++)); do
		char="${str:i:1}"
		kitty_send_text "$char"
		sleep "$delay"
	done
}

function feed() {
	feed_no_newline "$@"
	feed_str '\n'
}

function feed_str() {
	str="$1"
	kitty_send_text "$str"
	sleep "$delay"
}

feed clear
sleep 3
feed git reset --hard 7e47962
feed vi README.md
feed '/Installation'
feed 'OTODO: add demo'
feed_str '\x1b'
feed ':wq'
feed 'touch new_file.txt'
feed 'rm -f new_file.txt'
feed git add .
feed 'git commit -m "chore: add TODO message"'
feed git merge
feed git mergetool
feed ':%diffg REMOTE'
feed ':wqa'
feed 'git merge --abort'
feed git rebase
feed git rebase --abort
feed 'git reset --hard @{u}'
feed cd .git
feed cd -
feed 'rm -rf tmp/bare_repo && mkdir -p tmp/bare_repo && cd tmp/bare_repo'
feed git init --bare
feed cd -
feed git bisect HEAD~3
feed Y
feed git bisect reset
feed git checkout -b demo_branch
feed git push -u
sleep 2.5
feed git push origin :demo_branch
sleep 2.5
feed git checkout -
feed git branch -d demo_branch
feed clear
