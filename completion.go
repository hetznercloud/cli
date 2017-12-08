package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

const (
	bashCompletionFunc = `
	__hcloud_server_ids() {
		local ctl_output out
		if ctl_output=$(hcloud server list 2>/dev/null); then
				COMPREPLY=($(echo "${ctl_output}" | grep -v '^ID' | awk '{print $1}'))
		fi
	}

	__custom_func() {
		case ${last_command} in
			hcloud_server_delete | hcloud_server_describe )
				__hcloud_server_ids
				return
				;;
			*)
				;;
		esac
	}
	`

	completionShortDescription = "Output shell completion code for the specified shell (bash or zsh)"
	completionLongDescription  = completionShortDescription + `

Note: this requires the bash-completion framework, which is not installed by default on Mac. This can be installed by using homebrew:

	$ brew install bash-completion

Once installed, bash completion must be evaluated. This can be done by adding the following line to the .bash profile:

	$ source $(brew --prefix)/etc/bash_completion

Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2

Examples:
	# Load the hcloud completion code for bash into the current shell
	source <(hcloud completion bash)

	# Load the hcloud completion code for zsh into the current shell
	source <(hcloud completion zsh)`
)

var (
	completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash": runCompletionBash,
		"zsh":  runCompletionZsh,
	}
)

func newCompletionCommand(cli *CLI) *cobra.Command {
	shells := []string{}
	for s := range completionShells {
		shells = append(shells, s)
	}

	cmd := &cobra.Command{
		Use:       "completion SHELL",
		Short:     "Output shell completion code for the specified shell (bash or zsh)",
		Long:      completionLongDescription,
		RunE:      cli.wrap(runCompletion),
		Args:      cobra.ExactArgs(1),
		ValidArgs: shells,
	}
	return cmd
}

func runCompletion(cli *CLI, cmd *cobra.Command, args []string) error {
	run, found := completionShells[args[0]]
	if !found {
		return fmt.Errorf("unsupported shell type %q", args[0])
	}

	return run(os.Stdout, cmd.Parent())
}

func runCompletionBash(out io.Writer, cmd *cobra.Command) error {
	return cmd.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, cmd *cobra.Command) error {
	zshInitialization := `
__hcloud_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__hcloud_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__hcloud_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__hcloud_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__hcloud_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__hcloud_declare() {
	if [ "$1" == "-F" ]; then
		whence -w "$@"
	else
		builtin declare "$@"
	fi
}
__hcloud_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__hcloud_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__hcloud_filedir() {
	local RET OLD_IFS w qw
	__debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__hcloud_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__hcloud_quote() {
	if [[ $1 == \'* || $1 == \"* ]]; then
		# Leave out first character
		printf %q "${1:1}"
	else
		printf %q "$1"
	fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__hcloud_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__hcloud_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__hcloud_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__hcloud_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__hcloud_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__hcloud_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/__hcloud_declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__hcloud_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	out.Write([]byte(zshInitialization))

	buf := new(bytes.Buffer)
	cmd.Root().GenBashCompletion(buf)
	out.Write(buf.Bytes())

	zshTail := `
BASH_COMPLETION_EOF
}
__hcloud_bash_source <(__hcloud_convert_bash_to_zsh)
`
	out.Write([]byte(zshTail))
	return nil
}
