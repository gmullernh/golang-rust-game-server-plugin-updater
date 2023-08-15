# rust-game-server-plugin-updater
An app written in Go to update your plugins for gameservers of "Rust".

What Rust are we talking about? This one: https://rust.facepunch.com/.

1. Plugins that work with Oxide or Carbon.
1. Plugins that are downloaded from [https://umod.org].

## What does it do?

1. Look for all files inside the `/plugins` directory.
1. Query all the plugins from `umod.org` API.
1. If a match is found, downloads the file to replace the original file inside the `/plugins` folder.

## How to use

1. Place your plugins `*.cs` files in the folder `/plugins`. 
1. Run the program.