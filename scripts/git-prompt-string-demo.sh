#!/usr/bin/env bash

# Utility script to execute demo steps via Kitty

# gif generated with
# ffmpeg -i git-prompt-string-demo.mov -r 5 frame%04d.png
# gifski --quality 99 --motion-quality 99 --lossy-quality 99 --width 1200 --height 627 -o git-prompt-string.gif frame*.png

win="$1" # kitty window ID
char_delay='0.2'
feed_delay='2.5'

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
		sleep "$char_delay"
	done
}

function feed() {
	feed_no_newline "$@"
	feed_str '\n'
	sleep "$feed_delay"
}

function feed_str() {
	str="$1"
	kitty_send_text "$str"
	sleep "$char_delay"
}

feed clear
feed git reset --hard 7e47962
feed 'sed -i "16s/$/TODO: add demo/" README.md'
feed 'touch new_file.txt'
feed 'rm -f new_file.txt'
feed 'git commit -am "chore: add TODO message"'
feed git merge
feed 'git checkout --theirs README.md && git add .'
feed 'git merge --abort'
feed git rebase
feed git rebase --abort
feed 'git reset --hard @{u}'
feed cd .git
feed cd -
feed git bisect start HEAD~3
feed git bisect reset
feed git checkout -b demo_branch
feed git push -u
feed git push origin :demo_branch
feed git checkout -
feed git branch -d demo_branch
feed clear
