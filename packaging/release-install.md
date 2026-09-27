## Install

macOS:

```bash
brew tap yeixio/yggdrasil https://github.com/yeixio/yggdrasil-core
brew install yggdrasil
```

Debian and Ubuntu:

```bash
echo "deb [trusted=yes] https://raw.githubusercontent.com/yeixio/yggdrasil-core/apt stable main" | sudo tee /etc/apt/sources.list.d/yggdrasil.list
sudo apt-get update
sudo apt-get install yggdrasil
```

RPM packages for x86_64 and aarch64 are attached to this release. Install one with `sudo rpm -i` or `sudo dnf install`.
