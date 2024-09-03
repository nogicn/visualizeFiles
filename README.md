
## How to Use

### Build it yourself

1. **Install Go >= 1.19**

2. **Install dependencies**
   - Install dependencies for your system:
     - Debian / Ubuntu: `sudo apt-get install gcc libgl1-mesa-dev xorg-dev`
     - Fedora: `sudo dnf install gcc libXcursor-devel libXrandr-devel mesa-libGL-devel libXi-devel libXinerama-devel libXxf86vm-devel`
     - Arch Linux: `sudo pacman -S xorg-server-devel libxcursor libxrandr libxinerama libxi`
     - Solus: `sudo eopkg it -c system.devel mesalib-devel libxrandr-devel libxcursor-devel libxi-devel libxinerama-devel`
     - openSUSE: `sudo zypper install gcc libXcursor-devel libXrandr-devel Mesa-libGL-devel libXi-devel libXinerama-devel libXxf86vm-devel`
     - Void Linux: `sudo xbps-install -S base-devel xorg-server-devel libXrandr-devel libXcursor-devel libXinerama-devel libXxf86vm-devel`
     - Alpine Linux: `sudo apk add gcc libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev linux-headers mesa-dev`
     - NixOS: `nix-shell -p libGL pkg-config xorg.libX11.dev xorg.libXcursor xorg.libXi xorg.libXinerama xorg.libXrandr xorg.libXxf86vm`

3. **Run the Program**:
   ```bash
   go run main.go
   ```

4. **Review Output**:
   - Drag and drop your file into the application and move the slider to see the heatmap change

