#!/usr/bin/env bash
#
# Copyright (C) 2017 Michael Smith <nvimmike@gmail.com>
# Copyright (C) 2006,2007 Shawn O. Pearce <spearce@spearce.org>
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
#
# A better bash prompt for git. bgps provides a convenient way to customize
# the PS1 prompt and to determine information about the current git branch.
# bgps can indicate if the branch has a clean or dirty working tree, whether
# or not it is tracking a remote branch, and the number of commits the local
# branch is ahead or behind the remote branch.
#
# Parts of this program were copied and modified from
# <https://github.com/git/git/blob/master/contrib/completion/git-prompt.sh>
#

#######################################
# Read contents of file to variables
# Arguments:
#   $filepath the path of the file to read
#   $variable_names... the name of the variables to store the contents
#######################################
_eread() {
	local filepath="${1}"
	shift
	[[ -r "${filepath}" ]] && read -r "$@" <"${filepath}"
}

#######################################
# Get current git branch and additional repository information
# Returns: 0 on success
# Returns: 2 on success and in a detached head state
#######################################
_branch_info() {

	local git_dir="$1"
	local inside_gitdir="$2"
	local bare_repo="$3"
	local inside_worktree="$4"
	local short_sha="$5"

	local is_detached="false"
	local merge_status=""
	local sparse=""
	local branch=""
	local step=""
	local total=""
	if [[ -d "${git_dir}/rebase-merge" ]]; then

		_eread "${git_dir}/rebase-merge/head-name" branch
		_eread "${git_dir}/rebase-merge/msgnum" step
		_eread "${git_dir}/rebase-merge/end" total

		if [[ -f "${git_dir}/rebase-merge/interactive" ]]; then
			merge_status="|REBASE-i"
		else
			merge_status="|REBASE-m"
		fi

	else

		if [[ -d "${git_dir}/rebase-apply" ]]; then

			_eread "${git_dir}/rebase-apply/next" step
			_eread "${git_dir}/rebase-apply/last" total

			if [[ -f "${git_dir}/rebase-apply/rebasing" ]]; then
				_eread "${git_dir}/rebase-apply/head-name" branch
				merge_status="|REBASE"
			elif [[ -f "${git_dir}/rebase-apply/applying" ]]; then
				merge_status="|AM"
			else
				merge_status="|AM/REBASE"
			fi

		elif [[ -f "${git_dir}/MERGE_HEAD" ]]; then
			merge_status="|MERGING"
		elif [[ -f "${git_dir}/CHERRY_PICK_HEAD" ]]; then
			merge_status="|CHERRY-PICKING"
		elif [[ -f "${git_dir}/REVERT_HEAD" ]]; then
			merge_status="|REVERTING"
		elif [[ -f "${git_dir}/BISECT_LOG" ]]; then
			merge_status="|BISECTING"
		fi

		# TODO stopped here

		if [[ "${branch}" ]]; then
			:
		elif [[ -L "${git_dir}/HEAD" ]]; then
			# symlink symbolic ref
			branch="$(git symbolic-ref HEAD 2>/dev/null)"
		else
			local head=""
			if ! _eread "${git_dir}/HEAD" head; then
				return 1
			fi
			branch="${head#ref: }"
			local upstream_info
			upstream_info="$(git rev-parse --abbrev-ref '@{upstream}' 2>/dev/null)"
			if [[ "${head}" == "${branch}" ]]; then
				is_detached="true"
				branch="$(git describe --tags --exact-match HEAD 2>/dev/null)" || branch="${head:0:7}..." # was short
				branch="(${branch})"
			elif [[ "$upstream_info" == '@{upstream}' ]] || [[ "$upstream_info" == '' ]]; then
				is_detached="true"
			fi
		fi
	fi

	if [[ -n "${step}" ]] && [[ -n "${total}" ]]; then
		merge_status="${merge_status} ${step}/${total}"
	fi

	if [[ "$merge_status" ]] && [[ $(git ls-files --unmerged 2>/dev/null) ]]; then
		merge_status="$merge_status|CONFLICT"
	fi

	local prefix=""

	if [[ "${inside_gitdir}" == "true" ]]; then
		if [[ "${bare_repo}" == "true" ]]; then
			prefix="BARE:"
		else
			branch="GIT_DIR!"
		fi
	fi

	branch="${branch##refs/heads/}"

	if [[ "$(git config --bool core.sparseCheckout)" == "true" ]]; then
		sparse="|SPARSE"
	fi
	printf -- "%s%s%s%s" "${prefix}" "${branch}" "${sparse}" "${merge_status}"

	if [[ "$is_detached" == "true" ]]; then
		return 2
	fi

	return 0
}

#######################################
# Get formatted git status
# Globals:
#   $stop_color
# Returns: 0 on success
#######################################
_bgps_git_status() {
	local color_clean='\x1b[38;2;167;192;128m'

	local color_no_upstream='\x1b[38;2;146;131;116m'
	local color_untracked='\x1b[38;2;211;134;155m'
	local color_dirty='\x1b[38;2;255;105;97m'
	local color_conflict='\x1b[38;2;219;188;95m'
	local prefix='  '
	local postfix=''
	local git_ahead='↑[%s]'
	local git_behind='↓[%s]'
	local git_diverged='↕ ↑[%a] ↓[%b]'
	local git_bare='󱣻'
	local git_color="${stop_color}"

	local git_symbol=""

	local repo_info
	repo_info="$(git rev-parse --git-dir --is-inside-git-dir --is-bare-repository --is-inside-work-tree --short '@{upstream}' 2>/dev/null)"

	local repo_info_arr
	readarray -t repo_info_arr <<<"$repo_info"
	local git_dir="${repo_info_arr[0]}"
	local inside_gitdir="${repo_info_arr[1]}"
	local bare_repo="${repo_info_arr[2]}"
	local inside_worktree="${repo_info_arr[3]}"
	local short_sha="${repo_info_arr[4]}"

	local is_detached="false"
	local git_branch_exit
	local git_branch
	git_branch=$(_branch_info "$git_dir" "$inside_gitdir" "$bare_repo" "$inside_worktree" "$short_sha")
	git_branch_exit="$?"

	if (("$git_branch_exit" == 2)); then
		is_detached="true"
	elif (("$git_branch_exit")); then
		return 1
	fi

	if [[ "$bare_repo" == "true" ]]; then
		git_symbol="${git_bare}"
	fi
	if [[ "$bare_repo" == "true" ]] || [[ "$inside_gitdir" == "true" ]]; then
		git_color="${color_no_upstream}"
		printf -- "%s" "${git_color}"
		printf -- "%s%s%s%s" "${prefix:+${prefix}}" "${git_branch}" "${git_symbol:+ ${git_symbol}}" "${postfix:+${postfix}}"
		return 0
	fi

	local diff_error
	diff_error=$(git diff --no-ext-diff --quiet HEAD 2>&1)
	local dirty_exit_code="${?}" # code == 0 clean working tree, code == 1 dirty working tree

	local empty_git_error="ambiguous argument 'HEAD'"
	local no_upstream_git_error="no upstream configured"
	local no_such_branch_git_error="no such branch"
	local no_upstream=0
	if [[ "$short_sha" == "" ]]; then
		no_upstream=1
	fi
	if [[ "${diff_error}" == *"${empty_git_error}"* ]] || [[ "${diff_error}" == *"${no_upstream_git_error}"* ]] || [[ "${diff_error}" == *"${no_such_branch_git_error}"* ]]; then
		no_upstream=1
		# there is no upstream so compare against staging area
		diff_error=$(git diff --cached --no-ext-diff --quiet 2>&1)
		dirty_exit_code="${?}" # code == 0 clean working tree, code == 1 dirty working tree
	fi

	local untracked_exit_code="1" # code == 0 untracked files exist, code > 0 no untracked files
	# allow no upstream error to passthrough to apply coloring and formatting
	if (("${dirty_exit_code}" == 0)) || (("${dirty_exit_code}" == 1)) || (("${no_upstream}")); then

		git ls-files --others --exclude-standard --directory --no-empty-directory --error-unmatch -- ':/*' >/dev/null 2>/dev/null
		untracked_exit_code="${?}" # code == 0 untracked files exist, code > 0 no untracked files

		local commit_counts
		local commit_counts_exit_code=1
		if [[ "$is_detached" == "true" ]]; then
			commit_counts=(0 0)
		elif ((no_upstream)); then
			ahead_count="$(git rev-list --count HEAD 2>/dev/null)"
			commit_counts=("$ahead_count" 0)
		else
			IFS=$'\t' read -r -a commit_counts <<<"$(git rev-list --left-right --count ...'@{upstream}' 2>/dev/null)"
			commit_counts_exit_code="$?"
		fi

		if ((!dirty_exit_code)); then
			git_color="${color_clean}"
		fi

		if ((commit_counts[0])); then
			git_color="${color_conflict}"
			git_symbol="${git_ahead/\%s/${commit_counts[0]}}"
		fi

		if ((commit_counts[1])); then
			git_color="${color_conflict}"
			git_symbol="${git_behind/\%s/${commit_counts[1]}}"
		fi

		if ((commit_counts[0] && commit_counts[1])); then
			git_symbol="${git_diverged/\%a/${commit_counts[0]}}"
			git_symbol="${git_symbol/\%b/${commit_counts[1]}}"
		fi

		# continue to check for untracked and dirty because
		# it is still possible even without an upstream

		if ((no_upstream)) || ((commit_counts_exit_code)); then
			# no upstream configured for branch
			git_color="${color_no_upstream}"
		fi

		if ((!untracked_exit_code)); then
			git_color="${color_untracked}"
			git_symbol="*${git_symbol#\*}"
		fi

		if ((dirty_exit_code == 1)) && ((untracked_exit_code)); then
			git_color="${color_dirty}"
			git_symbol="*${git_symbol#\*}"
		fi

		printf -- "%s" "${git_color}"
		printf -- "%s%s%s%s" "${prefix:+${prefix}}" "${git_branch}" "${git_symbol:+ ${git_symbol}}" "${postfix:+${postfix}}"
	else
		printf -- "%s" "${color_dirty}"
		printf -- "ERROR(bgps): %s" "${diff_error}"
	fi
	printf -- ''

	return 0
}

stop_color="\033[0m"
# shellcheck disable=SC2059
printf -- "$(_bgps_git_status)$stop_color"

# unset variables and functions
unset stop_color
unset -f _bgps_git_status
unset -f _branch_info
unset -f _eread
