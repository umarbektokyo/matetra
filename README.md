# matetra
The ultimate math game made by self-claimed "The best future mathematicians on this tiny planet."
```
                         $$\                $$\                        
                         $$ |               $$ |                       
$$$$$$\$$$$\   $$$$$$\ $$$$$$\    $$$$$$\ $$$$$$\    $$$$$$\  $$$$$$\  
$$  _$$  _$$\  \____$$\\_$$  _|  $$  __$$\\_$$  _|  $$  __$$\ \____$$\ 
$$ / $$ / $$ | $$$$$$$ | $$ |    $$$$$$$$ | $$ |    $$ |  \__|$$$$$$$ |
$$ | $$ | $$ |$$  __$$ | $$ |$$\ $$   ____| $$ |$$\ $$ |     $$  __$$ |
$$ | $$ | $$ |\$$$$$$$ | \$$$$  |\$$$$$$$\  \$$$$  |$$ |     \$$$$$$$ |
\__| \__| \__| \_______|  \____/  \_______|  \____/ \__|      \_______|

```

For the latest versions and documentation on the matetra-tui, please check [matetra-engine repository](https://github.com/umarbektokyo/matetra-engine)

# Installation
```bash
git clone https://github.com/umarbektokyo/matetra
cd matetra-engine

# For the client
go install cmd/matetra-client
# For the server
go install .cmd/matetra-server
```
# Running
```bash
# For the client
matetra-client "ip-address":1729
# ex: matetra-client localhost:1729
# For the server
matetra-server start "game-name"
# ex: matetra-server start WonderfulGame
```

## Progress
 - [x] Game Server: Functining and running
 - [x] CLI-Client: Functining and running
 - [x] TUI-Client: In development
 - [ ] Godot-Client: In development
 - [ ] Physical Game: Design Stage

Approximately overall: ~12%

## Welcome the crew!
- Flush! - Esia
- Noga L.
- Pikachu - Meharwan
- Umarbek
