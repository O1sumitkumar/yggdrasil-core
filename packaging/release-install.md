## Install

macOS. Homebrew installs Yggdrasil Core from this repository. It does not install Yggdrasil Desktop. The short command `brew install yggdrasil` is a different Homebrew cask.

```bash
brew tap yeixio/yggdrasil https://github.com/yeixio/yggdrasil-core
brew install yeixio/yggdrasil/yggdrasil
yggdrasil-daemon
```

Or use the headless archive attached to this release:

```bash
tar -xzf yggdrasil-*-darwin-*-headless.tar.gz
cd yggdrasil-*-darwin-*-headless
./yggdrasil-daemon
```

Open `http://127.0.0.1:7331`.

Debian and Ubuntu:

```bash
echo "deb [trusted=yes] https://raw.githubusercontent.com/yeixio/yggdrasil-core/apt stable main" | sudo tee /etc/apt/sources.list.d/yggdrasil.list
sudo apt-get update
sudo apt-get install yggdrasil
```

RPM packages for x86_64 and aarch64 are attached to this release. Install one with `sudo rpm -i` or `sudo dnf install`.

Windows amd64. The headless archive attached to this release is not code-signed.

```bash
tar -xzf yggdrasil-*-windows-amd64-headless.tar.gz
cd yggdrasil-*-windows-amd64-headless
./yggdrasil-daemon.exe
```

Open `http://127.0.0.1:7331`.
