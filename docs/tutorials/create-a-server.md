# Creating a server

This tutorial covers the process of creating a server using the Hetzner Cloud CLI. It includes creating an SSH Key, 
uploading it, creating a Server and then connecting to it via SSH.

## Prerequisites

- A functioning installation of the [hcloud CLI](setup-hcloud-cli.md), with a valid active context.

## 1. Create an SSH Key

### 1.1 Generate an SSH Key

While an SSH key is not strictly required to create a server, it is highly recommended for secure access.
If you don't have an SSH key yet, you can generate one using the following command:

```bash
ssh-keygen -t ed25519 -f ~/.ssh/hcloud
```

Your private key will now be located at `~/.ssh/hcloud`. **Do not share your private key with anyone!**

### 1.2 Upload the SSH Key

You can upload your SSH key to Hetzner Cloud using the following command:

```bash
hcloud ssh-key create --name my-ssh-key --public-key-from-file ~/.ssh/hcloud.pub
```

> [!TIP]
> You can set this SSH key as the default SSH key for your context using the following command:
> ```bash
> hcloud config set default-ssh-keys my-ssh-key
> ```

## 2. Create a Server

### 2.1 Pick a Server Type

Before creating a server, you need to choose a server type. You can list all available server types using the following command:

```bash
hcloud server-type list
```

For this example we will use the `cpx11` server type.
You can view the details of this server type using the following command:

```bash
hcloud server-type describe cpx11
```

### 2.2 Pick an Image

You need to choose an image for your server. You can list all available images using the following command:

```bash
hcloud image list
```

There are many images available, including various Linux distributions and pre-configured app images.
For this example we will use the `ubuntu-24.04` image.

### 2.3 Pick a Location (Optional)

You can choose a location for your server. You can list all available locations using the following command:

```bash
hcloud location list
```

If you don't specify a location, one will be chosen for you. This is what we will do in this example.

### 2.4 Create the Server

Now you can create the server using the following command:

```bash
hcloud server create --name my-server --type cpx11 --image ubuntu-24.04 --ssh-key my-ssh-key
```

If you set the SSH key as the default SSH key for your context, you can omit the `--ssh-key` flag.
After the server was created, you will see information about the server, including its IP address.

## 3. Connect to the Server

You can connect to the server using SSH. The CLI contains a built-in utility to do this:

```bash
hcloud server ssh my-server -i ~/.ssh/hcloud
```

This command will open an SSH connection to the server using the private key you generated earlier.

## 4. Clean Up

After you are done with the server, you can delete it using the following command:

```bash
hcloud server delete my-server
```

You can also delete the SSH key using the following command:

```bash
hcloud ssh-key delete my-ssh-key
```
