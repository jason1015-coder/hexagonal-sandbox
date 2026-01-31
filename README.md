# TesselBox

A 2D sandbox adventure game inspired by *Terraria*, but built on a **hexagonal grid**.

Explore worlds, mine resources, build structures, craft items, fight enemies, and survive — all in beautiful hex tiles.


## Migration Notice

**Sorry!**  
The game is currently migrating from a pure-Python prototype to a **Golang ** language  
This refactor improves performance (Go for core engine/world/physics) 

Expect some features to be temporarily broken or incomplete during this phase.  
Thanks for your patience — progress is ongoing!



## Installation

### Via Releases (Recommended once binaries are ready)

1. Go to the **[Releases page](https://github.com/tesselstudio/TesselBox-game/releases)**.
2. Select the latest release (e.g. "Testing_stage1" or higher).
3. Under **Assets**, download the file for your platform:
   - **Windows**: `TesselBox.exe` (or similar)
   - **Linux**: `TesselBox.AppImage` → run with `./TesselBox.AppImage` (make executable first: `chmod +x TesselBox.AppImage`)
   - **macOS**: `TesselBox.app.zip` → unzip and open the .app
4. Launch the file — no installation needed!



### Temporary: Download Latest Build from GitHub Actions (Nightly-ish Builds)

While official release binaries are pending:

1. Visit the **[Actions tab](https://github.com/tesselstudio/TesselBox-game/actions)**.
2. Find the most recent **successful** workflow run (green checkmark, e.g. "Build Executables" or similar).
3. Scroll to the bottom → **Artifacts** section.
4. Download the zip for your OS (e.g. `tesselbox-windows-exe`, `tesselbox-linux-appimage`, `tesselbox-macos-app`).
5. Unzip and run the executable.

This gives you the freshest compiled version directly from recent commits.





