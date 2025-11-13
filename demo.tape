# https://github.com/charmbracelet/vhs

Output demo.gif
Require hcloud

Type "# Welcome, this is a quick demo of the hcloud CLI"
Enter 1
Sleep 50ms

Type "hcloud version"
Enter 1
Sleep 50ms

Type "# Note that we already configured a context attached to our Hetzner Cloud project"
Enter 1

Type "# First, let's upload our public ssh key"
Enter 1

Type "hcloud ssh-key create --name demo --public-key-from-file ~/.ssh/id_ed25519.pub"
Enter 1
Sleep 1s
Wait

Enter 1
