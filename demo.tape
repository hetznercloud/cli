Output demo.gif

Set Height 1080
Set Width  1920
Set FontSize 26

Set Margin 24
Set MarginFill "#d50c2d"
Set BorderRadius 24

Set Padding 24

Set TypingSpeed 150ms # 50ms default
Set PlaybackSpeed 3 # to make real-time actions quicker

Type "# Welcome, this is a quick demo of the hcloud CLI"
Enter
Type "# Note that we already configured a context with a token for our Hetzner project"
Enter
Enter

Type "# First, let's upload our public ssh key"
Enter
Type "hcloud ssh-key create --name demo --public-key-from-file ~/.ssh/id_ed25519.pub"
Enter
Wait
Sleep 4s
Enter

Type "# We can now configure a default ssh key that will be used when creating new servers"
Enter
Type "hcloud config set default-ssh-keys demo"
Enter
Wait
Sleep 4s
Enter

Type "# Now, we create a new private network"
Enter
Type "hcloud network create --name demo --ip-range 10.0.0.0/16"
Enter
Wait
Type "hcloud network add-subnet demo --type cloud --network-zone eu-central --ip-range 10.0.0.0/24"
Enter
Wait
Sleep 8s
Enter

Type "# Then, we create a new server and attach it the private network"
Enter
Type "hcloud server create --name demo --image debian-12 --type cpx22 --network demo"
Enter
Wait@30s
Sleep 8s
Enter

Type "# We can also list the resources in the project"
Enter
Type "hcloud all list"
Enter
Wait
Sleep 20s